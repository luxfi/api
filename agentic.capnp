# Agentic Network RPC Protocol
# Support for ZMQ, Cap'n Proto, and gRPC wire protocols with PQ privacy

@0xdeadbeef12345678;

using Cxx = import "/capnp/c++.capnp";
$Cxx.namespace("lux::agentic");

# Agentic network message types
struct AgenticMessage {
  # Message metadata
  messageId @0: UInt64;
  timestamp @1: UInt64;
  sourceAgent @2: AgentId;
  targetAgent @3: AgentId;
  
  # Post-quantum encryption
  encryption @4: PQEncryption;
  
  # Message payload (union for different types)  
  payload :union {
    valueTransfer @5: ValueTransfer;
    computation @6: ComputationRequest;
    intelligence @7: IntelligenceShare;
    coordination @8: CoordinationMessage;
    discovery @9: DiscoveryMessage;
  }
  
  # Message routing and delivery
  routing @10: RoutingInfo;
  
  # Zero-copy shared data reference
  sharedDataRef @11: SharedDataReference;
}

# Agent identity with post-quantum keys
struct AgentId {
  publicKey @0: PQPublicKey;        # Post-quantum public key
  address @1: Data;                 # Network address
  capabilities @2: AgentCapabilities;
  reputation @3: Float64;           # Trust score
  stake @4: UInt64;                 # Staked value for security
}

# Post-quantum cryptographic primitives
struct PQEncryption {
  algorithm @0: PQAlgorithm;
  encryptedData @1: Data;           # Encrypted message payload
  keyExchange @2: PQKeyExchange;    # Key exchange for session
  signature @3: PQSignature;        # Post-quantum signature
}

enum PQAlgorithm {
  kyber1024 @0;                     # NIST PQC standard
  ntruHrss701 @1;                   # NTRU-based
  classicMcEliece @2;               # Code-based
  sphincsHaraka128f @3;             # Hash-based signatures
  dilithium5 @4;                    # Lattice-based signatures
  falcon1024 @5;                    # NTRU-based signatures
}

struct PQKeyExchange {
  algorithm @0: PQAlgorithm;
  publicKey @1: Data;               # Ephemeral public key
  encapsulatedSecret @2: Data;      # KEM encapsulated secret
  sessionKey @3: Data;              # Derived session key (encrypted)
}

struct PQSignature {
  algorithm @0: PQAlgorithm;
  signature @1: Data;               # Post-quantum signature
  publicKey @2: Data;               # Signing public key
}

struct PQPublicKey {
  algorithm @0: PQAlgorithm;
  keyData @1: Data;                 # Public key bytes
  keyId @2: UInt64;                 # Unique key identifier
}

# Agent capabilities for intelligent routing
struct AgentCapabilities {
  computeTypes @0: List(ComputeType);
  maxMemory @1: UInt64;             # Available memory (bytes)
  maxCompute @2: UInt64;            # Compute units available
  specializations @3: List(Text);   # AI/ML specializations
  supportedProtocols @4: List(Protocol);
  trustLevel @5: TrustLevel;
}

enum ComputeType {
  cpu @0;
  gpu @1;
  tpu @2;
  fpga @3;
  quantum @4;
  neuromorphic @5;
}

enum Protocol {
  zmq @0;
  capnp @1;
  grpc @2;
  libp2p @3;
  websocket @4;
}

enum TrustLevel {
  untrusted @0;
  basic @1;
  verified @2;
  premium @3;
  governance @4;
}

# Value transfer between agents
struct ValueTransfer {
  amount @0: UInt64;                # Transfer amount
  currency @1: CurrencyType;        # Currency/token type
  sourceAccount @2: Data;           # Source account/address
  targetAccount @3: Data;           # Target account/address
  
  # Privacy features
  zkProof @4: ZKProof;              # Zero-knowledge proof
  mixerData @5: MixerData;          # Privacy mixing
  
  # Smart contract execution
  contractCall @6: ContractCall;    # Optional contract interaction
  
  # Fees and incentives
  fee @7: UInt64;                   # Network fee
  tip @8: UInt64;                   # Agent tip/incentive
}

enum CurrencyType {
  lux @0;                           # Native LUX token
  hanzo @1;                         # HANZO token
  compute @2;                       # Compute credits
  storage @3;                       # Storage credits
  bandwidth @4;                     # Bandwidth credits
  intelligence @5;                  # AI/intelligence credits
}

# Zero-knowledge privacy proofs
struct ZKProof {
  proofType @0: ZKProofType;
  proof @1: Data;                   # ZK proof bytes
  publicInputs @2: Data;            # Public verification inputs
  circuitId @3: UInt64;             # Circuit identifier
}

enum ZKProofType {
  plonk @0;
  groth16 @1;
  stark @2;
  bulletproofs @3;
  nova @4;
}

# Privacy mixing for transaction unlinkability
struct MixerData {
  mixerSet @0: UInt64;              # Anonymity set ID
  commitment @1: Data;              # Pedersen commitment
  nullifier @2: Data;               # Spend nullifier
  merkleProof @3: MerkleProof;      # Membership proof
}

struct MerkleProof {
  path @0: List(Data);              # Merkle path hashes
  indices @1: List(Bool);           # Left/right path indicators
  root @2: Data;                    # Merkle root
}

# Smart contract execution
struct ContractCall {
  contractAddress @0: Data;         # Contract address
  method @1: Text;                  # Method name
  parameters @2: Data;              # Encoded parameters
  gasLimit @3: UInt64;              # Gas limit
  gasPrice @4: UInt64;              # Gas price
}

# Computation requests between agents
struct ComputationRequest {
  computationType @0: ComputationType;
  inputData @1: ComputationInput;
  requirements @2: ComputeRequirements;
  payment @3: ComputePayment;
  privacy @4: ComputePrivacy;
}

enum ComputationType {
  mlInference @0;
  mlTraining @1;
  zkProofGeneration @2;
  cryptographicOp @3;
  dataProcessing @4;
  simulation @5;
}

struct ComputationInput {
  dataType @0: DataType;
  data @1: Data;                    # Input data (may be encrypted)
  sharedMemoryRef @2: UInt64;       # Zero-copy shared memory reference
  metadata @3: Data;                # Additional metadata
}

enum DataType {
  tensor @0;
  image @1;
  text @2;
  audio @3;
  video @4;
  structured @5;
  encrypted @6;
}

struct ComputeRequirements {
  minMemory @0: UInt64;             # Minimum memory required
  minCompute @1: UInt64;            # Minimum compute required
  maxLatency @2: UInt32;            # Maximum acceptable latency (ms)
  specialized @3: List(Text);       # Specialized hardware requirements
  trustRequired @4: TrustLevel;     # Required trust level
}

struct ComputePayment {
  paymentType @0: PaymentType;
  amount @1: UInt64;                # Payment amount
  currency @2: CurrencyType;        # Payment currency
  escrow @3: EscrowTerms;           # Escrow conditions
}

enum PaymentType {
  prepaid @0;
  postpaid @1;
  escrow @2;
  reputation @3;
}

struct EscrowTerms {
  releaseConditions @0: List(Text); # Conditions for payment release
  timeoutBlocks @1: UInt64;         # Timeout in blocks
  arbitrator @2: AgentId;           # Optional arbitrator
}

# Privacy-preserving computation
struct ComputePrivacy {
  inputPrivacy @0: PrivacyLevel;
  outputPrivacy @1: PrivacyLevel;
  fhe @2: FHEConfig;                # Fully homomorphic encryption
  mpc @3: MPCConfig;                # Multi-party computation
  tee @4: TEEConfig;                # Trusted execution environment
}

enum PrivacyLevel {
  none @0;
  encrypted @1;
  fhe @2;                           # Fully homomorphic
  mpc @3;                           # Multi-party
  tee @4;                           # Trusted execution
  zksnark @5;                       # Zero-knowledge
}

struct FHEConfig {
  scheme @0: FHEScheme;
  publicKey @1: Data;               # FHE public key
  parameters @2: Data;              # Scheme parameters
}

enum FHEScheme {
  bfv @0;
  ckks @1;
  bgv @2;
  tfhe @3;
}

struct MPCConfig {
  protocol @0: MPCProtocol;
  parties @1: List(AgentId);        # MPC participants
  threshold @2: UInt32;             # Threshold for reconstruction
}

enum MPCProtocol {
  shamir @0;                        # Shamir secret sharing
  bgw @1;                           # BGW protocol
  gmw @2;                           # GMW protocol
  aby @3;                           # ABY framework
}

struct TEEConfig {
  teeType @0: TEEType;
  attestation @1: Data;             # Remote attestation
  enclave @2: Data;                 # Enclave configuration
}

enum TEEType {
  sgx @0;                           # Intel SGX
  sev @1;                           # AMD SEV
  armTrustzone @2;                  # ARM TrustZone
  risc5Keystone @3;                 # RISC-V Keystone
}

# Intelligence sharing between agents
struct IntelligenceShare {
  shareType @0: IntelligenceType;
  data @1: IntelligenceData;
  rights @2: SharingRights;
  provenance @3: DataProvenance;
  incentives @4: IncentiveStructure;
}

enum IntelligenceType {
  model @0;                         # ML model
  dataset @1;                       # Training data
  insights @2;                      # Derived insights
  patterns @3;                      # Discovered patterns
  predictions @4;                   # Prediction results
}

struct IntelligenceData {
  format @0: DataFormat;
  size @1: UInt64;                  # Data size in bytes
  hash @2: Data;                    # Content hash for integrity
  encryption @3: PQEncryption;      # Encrypted intelligence data
  sharedMemRef @4: UInt64;          # Zero-copy reference
}

enum DataFormat {
  onnx @0;
  pytorch @1;
  tensorflow @2;
  hanzoMl @3;
  raw @4;
  compressed @5;
}

struct SharingRights {
  usage @0: List(UsageRight);
  restrictions @1: List(Text);      # Usage restrictions
  expiration @2: UInt64;            # Rights expiration timestamp
  transferable @3: Bool;            # Can rights be transferred
}

enum UsageRight {
  view @0;
  inference @1;
  training @2;
  modification @3;
  redistribution @4;
  commercial @5;
}

struct DataProvenance {
  creator @0: AgentId;              # Original creator
  contributors @1: List(AgentId);   # Contributors
  derivedFrom @2: List(Data);       # Source data hashes
  timestamp @3: UInt64;             # Creation timestamp
  lineage @4: List(ProvenanceStep); # Transformation lineage
}

struct ProvenanceStep {
  operation @0: Text;               # Operation performed
  operator @1: AgentId;             # Who performed it
  timestamp @2: UInt64;             # When performed
  inputs @3: List(Data);            # Input data hashes
  outputs @4: List(Data);           # Output data hashes
}

struct IncentiveStructure {
  paymentModel @0: PaymentModel;
  revenue @1: RevenueSharing;
  attribution @2: AttributionRewards;
}

enum PaymentModel {
  oneTime @0;
  subscription @1;
  usage @2;
  revenue @3;
  reputation @4;
}

struct RevenueSharing {
  shares @0: List(RevenueShare);
  duration @1: UInt64;              # Sharing duration
  triggers @2: List(Text);          # Revenue triggers
}

struct RevenueShare {
  agent @0: AgentId;
  percentage @1: Float32;           # Revenue percentage
  minimum @2: UInt64;               # Minimum payment
  maximum @3: UInt64;               # Maximum payment
}

struct AttributionRewards {
  citationReward @0: UInt64;        # Reward per citation
  usageReward @1: UInt64;           # Reward per usage
  derivationReward @2: UInt64;      # Reward for derived works
}

# Agent coordination and consensus
struct CoordinationMessage {
  coordinationType @0: CoordinationType;
  proposal @1: CoordinationProposal;
  consensus @2: ConsensusData;
  execution @3: ExecutionPlan;
}

enum CoordinationType {
  taskAllocation @0;
  resourceSharing @1;
  consensusBuilding @2;
  conflictResolution @3;
  emergencyResponse @4;
}

struct CoordinationProposal {
  proposalId @0: UInt64;
  proposer @1: AgentId;
  description @2: Text;
  participants @3: List(AgentId);
  deadline @4: UInt64;              # Proposal deadline
  votingMechanism @5: VotingMechanism;
}

enum VotingMechanism {
  simple @0;                        # Simple majority
  weighted @1;                      # Stake-weighted
  quadratic @2;                     # Quadratic voting
  futarchy @3;                      # Prediction markets
  liquid @4;                        # Liquid democracy
}

struct ConsensusData {
  votes @0: List(Vote);
  result @1: ConsensusResult;
  finalized @2: Bool;
  proof @3: ConsensusProof;
}

struct Vote {
  voter @0: AgentId;
  choice @1: VoteChoice;
  weight @2: UInt64;                # Vote weight
  reasoning @3: Text;               # Optional reasoning
  signature @4: PQSignature;        # Cryptographic signature
}

union VoteChoice {
  binary @0: Bool;                  # Yes/No
  multiple @1: UInt32;              # Multiple choice index
  ranked @2: List(UInt32);          # Ranked preferences
  continuous @3: Float64;           # Continuous value
}

struct ConsensusResult {
  outcome @0: VoteChoice;
  confidence @1: Float64;           # Confidence level
  participation @2: Float64;        # Participation rate
  finalizationTime @3: UInt64;      # When finalized
}

struct ConsensusProof {
  proofType @0: ProofType;
  proof @1: Data;                   # Cryptographic proof
  validators @2: List(AgentId);     # Validating agents
}

enum ProofType {
  aggregateSignature @0;
  merkleTree @1;
  zkProof @2;
  threshold @3;
}

struct ExecutionPlan {
  tasks @0: List(Task);
  dependencies @1: List(Dependency);
  timeline @2: Timeline;
  monitoring @3: MonitoringPlan;
}

struct Task {
  taskId @0: UInt64;
  assignee @1: AgentId;
  description @2: Text;
  requirements @3: TaskRequirements;
  deliverables @4: List(Deliverable);
  deadline @5: UInt64;
}

struct TaskRequirements {
  compute @0: ComputeRequirements;
  data @1: List(DataRequirement);
  permissions @2: List(Permission);
  dependencies @3: List(UInt64);   # Task IDs this depends on
}

struct DataRequirement {
  dataType @0: DataType;
  source @1: AgentId;
  access @2: AccessLevel;
  privacy @3: PrivacyLevel;
}

enum AccessLevel {
  read @0;
  write @1;
  execute @2;
  admin @3;
}

struct Permission {
  resource @0: Text;                # Resource identifier
  action @1: Text;                  # Permitted action
  grantor @2: AgentId;              # Who granted permission
  expiration @3: UInt64;            # Permission expiration
}

struct Deliverable {
  deliverableId @0: UInt64;
  description @1: Text;
  format @2: DataFormat;
  verification @3: VerificationMethod;
  quality @4: QualityMetrics;
}

enum VerificationMethod {
  manual @0;
  automated @1;
  consensus @2;
  oracle @3;
}

struct QualityMetrics {
  accuracy @0: Float64;
  completeness @1: Float64;
  timeliness @2: Float64;
  relevance @3: Float64;
}

struct Dependency {
  fromTask @0: UInt64;              # Source task
  toTask @1: UInt64;                # Dependent task
  dependencyType @2: DependencyType;
  data @3: List(DataDependency);
}

enum DependencyType {
  sequential @0;                    # Must complete before
  parallel @1;                      # Can run in parallel
  conditional @2;                   # Depends on outcome
  resource @3;                      # Shared resource constraint
}

struct DataDependency {
  dataId @0: UInt64;
  source @1: AgentId;
  format @2: DataFormat;
  size @3: UInt64;
}

struct Timeline {
  start @0: UInt64;                 # Project start time
  end @1: UInt64;                   # Project end time
  milestones @2: List(Milestone);
  criticalPath @3: List(UInt64);    # Critical path task IDs
}

struct Milestone {
  milestoneId @0: UInt64;
  name @1: Text;
  deadline @2: UInt64;
  criteria @3: List(Text);          # Success criteria
  rewards @4: List(Reward);
}

struct Reward {
  recipient @0: AgentId;
  amount @1: UInt64;
  currency @2: CurrencyType;
  condition @3: Text;               # Reward condition
}

struct MonitoringPlan {
  metrics @0: List(MonitoringMetric);
  alerts @1: List(AlertRule);
  reports @2: List(ReportSchedule);
}

struct MonitoringMetric {
  name @0: Text;
  source @1: AgentId;
  frequency @2: UInt32;             # Collection frequency (seconds)
  threshold @3: List(Threshold);
}

struct Threshold {
  level @0: AlertLevel;
  value @1: Float64;
  action @2: Text;                  # Action to take
}

enum AlertLevel {
  info @0;
  warning @1;
  error @2;
  critical @3;
}

struct AlertRule {
  condition @0: Text;               # Alert condition
  level @1: AlertLevel;
  recipients @2: List(AgentId);
  cooldown @3: UInt32;              # Cooldown period (seconds)
}

struct ReportSchedule {
  name @0: Text;
  frequency @1: ReportFrequency;
  recipients @2: List(AgentId);
  template @3: Text;                # Report template
}

enum ReportFrequency {
  realtime @0;
  hourly @1;
  daily @2;
  weekly @3;
  monthly @4;
}

# Network discovery and routing
struct DiscoveryMessage {
  discoveryType @0: DiscoveryType;
  query @1: DiscoveryQuery;
  response @2: DiscoveryResponse;
  routing @3: RoutingUpdate;
}

enum DiscoveryType {
  agentDiscovery @0;
  serviceDiscovery @1;
  resourceDiscovery @2;
  routingUpdate @3;
  healthCheck @4;
}

struct DiscoveryQuery {
  queryId @0: UInt64;
  requester @1: AgentId;
  criteria @2: SearchCriteria;
  maxResults @3: UInt32;
  timeout @4: UInt32;               # Query timeout (seconds)
}

struct SearchCriteria {
  capabilities @0: List(ComputeType);
  specializations @1: List(Text);
  location @2: GeographicConstraints;
  trustLevel @3: TrustLevel;
  priceRange @4: PriceRange;
}

struct GeographicConstraints {
  regions @0: List(Text);           # Allowed regions
  maxLatency @1: UInt32;            # Maximum network latency
  dataResidency @2: List(Text);     # Data residency requirements
}

struct PriceRange {
  minPrice @0: UInt64;
  maxPrice @1: UInt64;
  currency @2: CurrencyType;
  pricingModel @3: PricingModel;
}

enum PricingModel {
  fixed @0;
  auction @1;
  dynamic @2;
  tiered @3;
}

struct DiscoveryResponse {
  queryId @0: UInt64;
  responder @1: AgentId;
  matches @2: List(AgentMatch);
  total @3: UInt32;                 # Total matches available
}

struct AgentMatch {
  agent @0: AgentId;
  score @1: Float64;                # Match score
  availability @2: Availability;
  pricing @3: PricingInfo;
  testimonials @4: List(Testimonial);
}

struct Availability {
  status @0: AvailabilityStatus;
  nextAvailable @1: UInt64;         # Next availability timestamp
  capacity @2: CapacityInfo;
  schedule @3: List(ScheduleSlot);
}

enum AvailabilityStatus {
  available @0;
  busy @1;
  maintenance @2;
  offline @3;
}

struct CapacityInfo {
  currentLoad @0: Float64;          # Current utilization (0-1)
  maxCapacity @1: UInt64;           # Maximum capacity
  reservedCapacity @2: UInt64;      # Reserved/committed capacity
  availableCapacity @3: UInt64;     # Available capacity
}

struct ScheduleSlot {
  start @0: UInt64;                 # Slot start time
  end @1: UInt64;                   # Slot end time
  capacity @2: UInt64;              # Available capacity in slot
  price @3: UInt64;                 # Price for slot
}

struct PricingInfo {
  model @0: PricingModel;
  basePrice @1: UInt64;             # Base price
  currency @2: CurrencyType;
  discounts @3: List(Discount);
  premiums @4: List(Premium);
}

struct Discount {
  condition @0: Text;               # Discount condition
  percentage @1: Float32;           # Discount percentage
  maximum @2: UInt64;               # Maximum discount amount
}

struct Premium {
  condition @0: Text;               # Premium condition  
  percentage @1: Float32;           # Premium percentage
  minimum @2: UInt64;               # Minimum premium amount
}

struct Testimonial {
  client @0: AgentId;
  rating @1: Float32;               # Rating (1-5)
  review @2: Text;                  # Review text
  timestamp @3: UInt64;             # Review timestamp
  verified @4: Bool;                # Verified transaction
}

# Routing information for message delivery
struct RoutingInfo {
  path @0: List(AgentId);           # Route path
  hops @1: UInt32;                  # Number of hops
  latency @2: UInt32;               # Expected latency (ms)
  cost @3: UInt64;                  # Routing cost
  priority @4: MessagePriority;
  qos @5: QoSRequirements;
}

enum MessagePriority {
  low @0;
  normal @1;
  high @2;
  critical @3;
  realtime @4;
}

struct QoSRequirements {
  maxLatency @0: UInt32;            # Maximum acceptable latency
  minBandwidth @1: UInt64;          # Minimum bandwidth required
  reliability @2: Float32;          # Required reliability (0-1)
  security @3: SecurityLevel;
}

enum SecurityLevel {
  none @0;
  basic @1;
  encrypted @2;
  pqEncrypted @3;                   # Post-quantum encrypted
  anonymous @4;                     # Anonymous routing
}

struct RoutingUpdate {
  updateType @0: UpdateType;
  routes @1: List(RouteInfo);
  topology @2: NetworkTopology;
  metrics @3: NetworkMetrics;
}

enum UpdateType {
  add @0;
  remove @1;
  update @2;
  full @3;                          # Full routing table
}

struct RouteInfo {
  destination @0: AgentId;
  nextHop @1: AgentId;
  cost @2: UInt64;
  latency @3: UInt32;
  reliability @4: Float32;
  lastUpdate @5: UInt64;
}

struct NetworkTopology {
  nodes @0: List(NetworkNode);
  edges @1: List(NetworkEdge);
  clusters @2: List(NetworkCluster);
}

struct NetworkNode {
  agent @0: AgentId;
  position @1: NetworkPosition;
  connections @2: UInt32;           # Number of connections
  centrality @3: Float64;           # Network centrality measure
}

struct NetworkPosition {
  x @0: Float64;                    # X coordinate in network space
  y @1: Float64;                    # Y coordinate in network space
  z @2: Float64;                    # Z coordinate for 3D networks
}

struct NetworkEdge {
  from @0: AgentId;
  to @1: AgentId;
  weight @2: Float64;               # Edge weight/cost
  bandwidth @3: UInt64;             # Available bandwidth
  latency @4: UInt32;               # Edge latency
}

struct NetworkCluster {
  clusterId @0: UInt64;
  members @1: List(AgentId);
  representative @2: AgentId;       # Cluster representative
  properties @3: ClusterProperties;
}

struct ClusterProperties {
  specialization @0: Text;          # Cluster specialization
  trustLevel @1: TrustLevel;        # Average trust level
  capacity @2: UInt64;              # Total cluster capacity
  load @3: Float64;                 # Current load
}

struct NetworkMetrics {
  totalNodes @0: UInt64;
  totalEdges @1: UInt64;
  averageLatency @2: Float64;
  networkThroughput @3: UInt64;
  reliability @4: Float64;
  timestamp @5: UInt64;
}

# Zero-copy shared data reference
struct SharedDataReference {
  sharedId @0: UInt64;              # Shared memory segment ID
  offset @1: UInt64;                # Offset within segment
  size @2: UInt64;                  # Data size
  checksum @3: UInt64;              # Data integrity checksum
  permissions @4: List(Text);       # Access permissions
}

# RPC protocol configuration
struct RPCConfig {
  protocol @0: RPCProtocol;
  endpoint @1: Text;                # Protocol endpoint
  options @2: ProtocolOptions;
  security @3: SecurityConfig;
}

enum RPCProtocol {
  zmq @0;
  capnp @1;
  grpc @2;
  websocket @3;
  custom @4;
}

union ProtocolOptions {
  zmq @0: ZMQOptions;
  capnp @1: CapnpOptions;
  grpc @2: GRPCOptions;
  websocket @3: WebSocketOptions;
}

struct ZMQOptions {
  socketType @0: ZMQSocketType;
  highWaterMark @1: UInt64;         # High water mark for queuing
  linger @2: UInt32;                # Linger time on close
  reconnectInterval @3: UInt32;     # Reconnection interval
  maxReconnects @4: UInt32;         # Maximum reconnection attempts
}

enum ZMQSocketType {
  req @0;                           # Request socket
  rep @1;                           # Reply socket
  dealer @2;                        # Dealer socket
  router @3;                        # Router socket
  pub @4;                           # Publisher socket
  sub @5;                           # Subscriber socket
  xpub @6;                          # Extended publisher
  xsub @7;                          # Extended subscriber
  push @8;                          # Push socket
  pull @9;                          # Pull socket
  pair @10;                         # Pair socket
  stream @11;                       # Stream socket
}

struct CapnpOptions {
  compression @0: CompressionType;
  packedEncoding @1: Bool;          # Use packed encoding
  traversalLimit @2: UInt64;        # Message traversal limit
  nestingLimit @3: UInt32;          # Nesting limit
}

enum CompressionType {
  none @0;
  gzip @1;
  lz4 @2;
  zstd @3;
  snappy @4;
}

struct GRPCOptions {
  maxReceiveMessageSize @0: UInt64;
  maxSendMessageSize @1: UInt64;
  keepaliveTime @2: UInt32;
  keepaliveTimeout @3: UInt32;
  maxConnectionIdle @4: UInt32;
  compression @5: GRPCCompression;
}

enum GRPCCompression {
  identity @0;                      # No compression
  gzip @1;
  deflate @2;
}

struct WebSocketOptions {
  maxFrameSize @0: UInt64;
  maxMessageSize @1: UInt64;
  compression @2: Bool;
  pingInterval @3: UInt32;
  pongTimeout @4: UInt32;
}

struct SecurityConfig {
  tls @0: TLSConfig;
  authentication @1: AuthConfig;
  authorization @2: AuthzConfig;
  rateLimit @3: RateLimitConfig;
}

struct TLSConfig {
  enabled @0: Bool;
  certificate @1: Data;             # TLS certificate
  privateKey @2: Data;              # TLS private key
  cipherSuites @3: List(Text);      # Allowed cipher suites
  minVersion @4: Text;              # Minimum TLS version
}

struct AuthConfig {
  method @0: AuthMethod;
  credentials @1: AuthCredentials;
  tokenExpiry @2: UInt64;           # Token expiry time
  refreshEnabled @3: Bool;          # Enable token refresh
}

enum AuthMethod {
  none @0;
  apiKey @1;
  jwt @2;
  oauth2 @3;
  mutual @4;                        # Mutual authentication
  pqc @5;                           # Post-quantum authentication
}

union AuthCredentials {
  apiKey @0: Data;
  jwt @1: JWTCredentials;
  oauth2 @2: OAuth2Credentials;
  pqc @3: PQCredentials;
}

struct JWTCredentials {
  token @0: Text;                   # JWT token
  publicKey @1: Data;               # Verification public key
  algorithm @2: Text;               # Signature algorithm
}

struct OAuth2Credentials {
  clientId @0: Text;
  clientSecret @1: Text;
  accessToken @2: Text;
  refreshToken @3: Text;
  scope @4: List(Text);
}

struct PQCredentials {
  algorithm @0: PQAlgorithm;
  publicKey @1: Data;
  signature @2: Data;
  certificate @3: Data;
}

struct AuthzConfig {
  enabled @0: Bool;
  policies @1: List(AuthzPolicy);
  defaultAction @2: AuthzAction;
}

struct AuthzPolicy {
  resource @0: Text;                # Resource pattern
  action @1: Text;                  # Action pattern
  principal @2: Text;               # Principal pattern
  effect @3: AuthzAction;           # Allow or deny
  conditions @4: List(Text);        # Additional conditions
}

enum AuthzAction {
  allow @0;
  deny @1;
}

struct RateLimitConfig {
  enabled @0: Bool;
  requests @1: UInt64;              # Requests per window
  window @2: UInt32;                # Time window (seconds)
  burst @3: UInt64;                 # Burst allowance
  strategy @4: RateLimitStrategy;
}

enum RateLimitStrategy {
  fixedWindow @0;
  slidingWindow @1;
  tokenBucket @2;
  leakyBucket @3;
}