# agentic.zap moved

The agentic-network RPC ZAP schema previously at

    github.com/luxfi/api/agentic.zap

now lives at

    github.com/luxfi/proto/schemas/api/agentic.zap

Generated `*_zap.go` siblings will land colocated with their
consuming Go package once a `//go:generate zapgen` directive is
authored. The legacy `agentic.capnp` (Cap'n Proto schema, same
shapes) remains in this directory unchanged.
