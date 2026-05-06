# ZAP — Zero-Copy App Proto

## Overview

ZAP is Lux's binary RPC wire protocol. Two implementations exist; both speak compatible binary wire formats and carry the same name. Pick the right one for the job.

| Implementation | Path | Use |
|---|---|---|
| `luxfi/api/zap` | this package | **Hot-path Lux RPC**: VM ↔ Node, sender, warp signing, netrunner control. Plain TCP framing, message-type-per-RPC. Optimised for low-latency in-process or local-cluster RPC. |
| `zap-protocol/zap-go` | `~/work/zap/zap-go` | **Canonical cross-language ZAP**: Cap'n Proto-based, multi-transport (TCP/Unix/WS/UDP/HTTP-SSE/Stdio), MCP gateway, agent consensus, full PQ crypto. Used for AI agent communication and cross-language schemas. |

Both are valid; they target different workloads. The `luxfi/api/zap` package is the tight Go-only wire format used inside the Lux node ecosystem. `zap-protocol` is the multi-language schema-driven implementation.

## Wire format (`luxfi/api/zap`)

Per-message frame:

```
[4 bytes: payload length, big-endian][1 byte: message type][payload bytes...]
```

`MaxMessageSize` = 16 MiB. Header size = 5 bytes.

### Message type byte layout

```
bit 7 (0x80) — MsgResponseFlag : set on every response
bit 6 (0x40) — MsgErrorFlag    : response payload is an error string
bits 5..0    — RPC opcode      : application-defined message type (1..63)
```

Reserved opcode ranges (see `wire.go` constants):
- `1..31` — VM service methods (Initialize, BuildBlock, ParseBlock, …)
- `40..43` — p2p.Sender methods (SendRequest, SendResponse, SendError, SendGossip)
- `50..52` — Warp signing (WarpSign, WarpGetPublicKey, WarpBatchSign)
- `60..63` — Reserved for future hot-path RPCs (do not redefine — pick a fresh range above)

When you add a new opcode, **stay below `0x40`** so the response and error flags remain free.

### Request/response encoding

The application is responsible for embedding a request ID at the start of payloads when call/response correlation is needed. The `Conn.Call` helper in `transport.go` does this automatically:

```
client → server:  [hdr] [reqID:uint32] [args bytes...]
server → client:  [hdr|MsgResponseFlag]               [reqID:uint32] [result bytes...]
                  [hdr|MsgResponseFlag|MsgErrorFlag]  [reqID:uint32] [error string]
```

The 4-byte request ID is mandatory in every payload sent through `Call`.

### Buffer + Reader

`Buffer` is a pooled, growable, big-endian-LE writer (`WriteUint8/16/32/64`, `WriteInt32/64`, `WriteBool`, `WriteBytes` length-prefixed, `WriteString`). `Reader` is the symmetric zero-copy reader. `BufferPool` is a `sync.Pool`; call `GetBuffer` / `PutBuffer` for reuse.

`ReadBytes` returns a slice into the original buffer (zero-copy). The caller must finish reading before recycling the underlying frame.

## Transport (`transport.go`)

- `Listen(addr, *Config) (*Listener, error)` — TCP listener.
- `Dial(ctx, addr, *Config) (*Conn, error)` — client dialer.
- `NewServer(*Listener, Handler)` + `Server.Serve(ctx)` — server loop. `Handler` is `func(ctx, msgType, payload) (msgType, payload, error)`.
- `Conn.Call(ctx, msgType, payload) (msgType, payload, error)` — sync request/response. Demultiplexes by request ID.

Multiple in-flight calls per connection are supported (each gets a fresh request ID). The server runs handlers in goroutines; writes are serialised by `writeMu`.

## Post-Quantum (PQ ZAP)

`luxfi/api/zap` is a wire format only — it does not implement PQ key exchange itself. PQ for ZAP-flavoured connections in Lux is delivered at two layers:

### 1. RNS transport (canonical hybrid PQ for Lux node-to-node)

For Lux node-to-node connections, the canonical PQ layer is `~/work/lux/node/network/dialer/rns_*` (Reticulum Network Stack with hybrid PQ identity). It tunnels app-layer protocols (including ZAP) inside an authenticated, PQ-secure session:

| Purpose | Classical | Post-quantum | NIST level |
|---|---|---|---|
| Identity signing | Ed25519 | ML-DSA-65 (FIPS 204) | 3 |
| Key exchange | X25519 | ML-KEM-768 (FIPS 203) | 3 |
| Session encryption | AES-256-GCM | — | — |

Hybrid combination: `combined_secret = HKDF-SHA256(X25519_shared || ML_KEM_shared)`. Forward secrecy via fresh ephemeral keys per session, zeroed after handshake. Optional `requirePostQuantum: true` rejects classical-only peers; default allows fallback for backward compatibility. See `~/work/lux/node/network/dialer/rns_link.go` and `rns_identity_pq.go` for the full implementation; the spec is LP-9701.

The X25519+ML-KEM-768 hybrid pattern matches the IETF **X-Wing** KEM construction (`draft-connolly-cfrg-xwing-kem`). Lux's RNS hybrid uses the same component primitives but its own KDF wrapping; consumers expecting the literal X-Wing wire format should not assume drop-in interop with RNS without checking the KDF labels. New code that wants standards-conformant interop should use the X-Wing path in `zap-protocol/zap-go` (or its sibling Rust crate) once it lands.

### 2. zap-protocol (canonical PQ for AI-agent ZAP)

For the multi-language ZAP at `~/work/zap`, the Rust core (`zap/src/crypto.rs`) provides the same PQ primitives directly in the protocol layer:

- `PQKeyExchange` — ML-KEM-768 (`encapsulate` / `decapsulate`)
- `PQSignature` — ML-DSA-65 (`sign` / `verify`)
- `HybridHandshake` — X25519 + ML-KEM-768 with HKDF combine. This is the **X-Wing** construction (`draft-connolly-cfrg-xwing-kem`); the Rust core implements the IETF spec's KDF labels exactly so consumers across language bindings interop on the wire.

Backend chain (Rust):
1. `luxcrypto-sys` (libluxcrypto FFI) — preferred, FIPS 203/204 via Cloudflare CIRCL
2. `pqcrypto` Rust crate — fallback

Sizes: `MLKEM_PUBLIC_KEY_SIZE`=1184, `MLKEM_CIPHERTEXT_SIZE`=1088, `MLDSA_PUBLIC_KEY_SIZE`=1952, `MLDSA_SIGNATURE_SIZE`=3309.

### Why two hybrids?

| | RNS hybrid | X-Wing |
|---|---|---|
| Purpose | Lux node-to-node session | AI-agent ZAP, cross-language |
| Components | X25519 + ML-KEM-768 + Ed25519/ML-DSA identity | X25519 + ML-KEM-768 |
| KDF | HKDF-SHA256 with Lux-specific labels | IETF X-Wing combiner |
| Standard | LP-9701 | draft-connolly-cfrg-xwing-kem |
| Identity tied to handshake? | Yes (Ed25519 + ML-DSA cert chain) | Separate (apps layer their own auth) |

For brand-new code: prefer X-Wing where standards-conformant interop matters (open clients, cross-org boundaries). Keep RNS hybrid where the handshake must bind validator identity to the PQ session, because RNS rolls identity into the construction.

### When to use which

- **Hot-path Lux RPC over TCP inside a trusted network or already inside an RNS session** → plain `luxfi/api/zap`. PQ already terminated by RNS (or trust boundary makes it unnecessary).
- **AI agent comms, cross-language clients, end-to-end PQ at the protocol layer, MCP gateway aggregation** → `zap-protocol/zap-go` (or sibling language binding) which speaks ZAP-on-Cap'nProto with PQ in the handshake.
- **New Lux internal services that need PQ at the message layer** → wrap `luxfi/api/zap` connections in `node/network/dialer/rns_link.RNSLink` rather than reimplementing PQ inside the wire format.

## Common pitfalls

1. **Don't OR `MsgResponseFlag` and `MsgErrorFlag` into a request opcode.** Keep all defined message types in `1..63`. Flags `0x40` and `0x80` are wire-protocol bits.
2. **Echo handlers must still get the response flag** — the server applies it automatically before writing. Don't pre-OR in your handler.
3. **Always include a request ID** in payloads sent via `Call`. Bare `WriteMessage` skips this; use it only for fire-and-forget gossip.
4. **`Buffer` data is reused** — don't retain `buf.Bytes()` across `PutBuffer`. Copy or hand the buffer downstream.
5. **`Reader.ReadBytes` returns a slice into the original frame** — same lifetime caveat. Copy if you need to outlive the frame.

## Adding a new RPC

1. Pick a fresh opcode in `wire.go`. Stay below `0x40`. Document the range (e.g. `// netrunner control: 60-99`).
2. Define a hand-written Go struct for request and response. Implement `Encode(*Buffer) []byte` and `Decode(*Reader) error`.
3. Server: register a `Handler` that switches on `msgType`, decodes the request, runs the call, encodes the response.
4. Client: implement a typed wrapper that calls `Conn.Call(ctx, opcode, request.Encode(...))` and decodes the result.
5. No protobuf, no `.proto` files. Hand-written Go types are the source of truth.

## Testing

```bash
go test ./...
```

Critical tests:
- `TestTransport` — round-trip request/response with response-flag discrimination.
- `TestConcurrentCalls` — many in-flight calls per connection demultiplex by request ID.
- `BenchmarkRoundTrip` / `BenchmarkRoundTripLargeBlock` — perf gates.

## See also

- `~/work/zap/zap/src/crypto.rs` — canonical PQ implementation (Rust)
- `~/work/lux/node/network/dialer/rns_link.go` — RNS hybrid-PQ wrapper for Lux RPC
- `~/work/lux/node/vms/rpcchainvm/zap/` — VM plugin handshake using `luxfi/api/zap`
- `~/work/lux/node/vms/platformvm/warp/zwarp/` — Warp signing client/server pattern
