# Agentic Network RPC Protocol
# Support for ZMQ, Cap'n Proto, and gRPC wire protocols with PQ privacy

using Cxx = import "/capnp/c++.capnp"

# Agentic network message types
struct AgenticMessage
  # Message metadata
  messageId UInt64
  timestamp UInt64
  sourceAgent AgentId
  targetAgent AgentId

  # Post-quantum encryption
  encryption PQEncryption

  # Message payload (union for different types)
  payload
    union
      valueTransfer ValueTransfer
      computation ComputationRequest
      intelligence IntelligenceShare
      coordination CoordinationMessage
      discovery DiscoveryMessage

  # Message routing and delivery
  routing RoutingInfo

  # Zero-copy shared data reference
  sharedDataRef SharedDataReference

# Agent identity with post-quantum keys
struct AgentId
  publicKey PQPublicKey
  address Data
  capabilities AgentCapabilities
  reputation Float64
  stake UInt64

# Post-quantum cryptographic primitives
struct PQEncryption
  algorithm PQAlgorithm
  encryptedData Data
  keyExchange PQKeyExchange
  signature PQSignature

enum PQAlgorithm
  kyber1024
  ntruHrss701
  classicMcEliece
  sphincsHaraka128f
  dilithium5
  falcon1024

struct PQKeyExchange
  algorithm PQAlgorithm
  publicKey Data
  encapsulatedSecret Data
  sessionKey Data

struct PQSignature
  algorithm PQAlgorithm
  signature Data
  publicKey Data

struct PQPublicKey
  algorithm PQAlgorithm
  keyData Data
  keyId UInt64

# Agent capabilities for intelligent routing
struct AgentCapabilities
  computeTypes List(ComputeType)
  maxMemory UInt64
  maxCompute UInt64
  specializations List(Text)
  supportedProtocols List(Protocol)
  trustLevel TrustLevel

enum ComputeType
  cpu
  gpu
  tpu
  fpga
  quantum
  neuromorphic

enum Protocol
  zmq
  capnp
  grpc
  libp2p
  websocket

enum TrustLevel
  untrusted
  basic
  verified
  premium
  governance

# Value transfer between agents
struct ValueTransfer
  amount UInt64
  currency CurrencyType
  sourceAccount Data
  targetAccount Data

  # Privacy features
  zkProof ZKProof
  mixerData MixerData

  # Smart contract execution
  contractCall ContractCall

  # Fees and incentives
  fee UInt64
  tip UInt64

enum CurrencyType
  lux
  hanzo
  compute
  storage
  bandwidth
  intelligence

# Zero-knowledge privacy proofs
struct ZKProof
  proofType ZKProofType
  proof Data
  publicInputs Data
  circuitId UInt64

enum ZKProofType
  plonk
  groth16
  stark
  bulletproofs
  nova

# Privacy mixing for transaction unlinkability
struct MixerData
  mixerSet UInt64
  commitment Data
  nullifier Data
  merkleProof MerkleProof

struct MerkleProof
  path List(Data)
  indices List(Bool)
  root Data

# Smart contract execution
struct ContractCall
  contractAddress Data
  method Text
  parameters Data
  gasLimit UInt64
  gasPrice UInt64

# Computation requests between agents
struct ComputationRequest
  computationType ComputationType
  inputData ComputationInput
  requirements ComputeRequirements
  payment ComputePayment
  privacy ComputePrivacy

enum ComputationType
  mlInference
  mlTraining
  zkProofGeneration
  cryptographicOp
  dataProcessing
  simulation

struct ComputationInput
  dataType DataType
  data Data
  sharedMemoryRef UInt64
  metadata Data

enum DataType
  tensor
  image
  text
  audio
  video
  structured
  encrypted

struct ComputeRequirements
  minMemory UInt64
  minCompute UInt64
  maxLatency UInt32
  specialized List(Text)
  trustRequired TrustLevel

struct ComputePayment
  paymentType PaymentType
  amount UInt64
  currency CurrencyType
  escrow EscrowTerms

enum PaymentType
  prepaid
  postpaid
  escrow
  reputation

struct EscrowTerms
  releaseConditions List(Text)
  timeoutBlocks UInt64
  arbitrator AgentId

# Privacy-preserving computation
struct ComputePrivacy
  inputPrivacy PrivacyLevel
  outputPrivacy PrivacyLevel
  fhe FHEConfig
  mpc MPCConfig
  tee TEEConfig

enum PrivacyLevel
  none
  encrypted
  fhe
  mpc
  tee
  zksnark

struct FHEConfig
  scheme FHEScheme
  publicKey Data
  parameters Data

enum FHEScheme
  bfv
  ckks
  bgv
  tfhe

struct MPCConfig
  protocol MPCProtocol
  parties List(AgentId)
  threshold UInt32

enum MPCProtocol
  shamir
  bgw
  gmw
  aby

struct TEEConfig
  teeType TEEType
  attestation Data
  enclave Data

enum TEEType
  sgx
  sev
  armTrustzone
  risc5Keystone

# Intelligence sharing between agents
struct IntelligenceShare
  shareType IntelligenceType
  data IntelligenceData
  rights SharingRights
  provenance DataProvenance
  incentives IncentiveStructure

enum IntelligenceType
  model
  dataset
  insights
  patterns
  predictions

struct IntelligenceData
  format DataFormat
  size UInt64
  hash Data
  encryption PQEncryption
  sharedMemRef UInt64

enum DataFormat
  onnx
  pytorch
  tensorflow
  hanzoMl
  raw
  compressed

struct SharingRights
  usage List(UsageRight)
  restrictions List(Text)
  expiration UInt64
  transferable Bool

enum UsageRight
  view
  inference
  training
  modification
  redistribution
  commercial

struct DataProvenance
  creator AgentId
  contributors List(AgentId)
  derivedFrom List(Data)
  timestamp UInt64
  lineage List(ProvenanceStep)

struct ProvenanceStep
  operation Text
  operator AgentId
  timestamp UInt64
  inputs List(Data)
  outputs List(Data)

struct IncentiveStructure
  paymentModel PaymentModel
  revenue RevenueSharing
  attribution AttributionRewards

enum PaymentModel
  oneTime
  subscription
  usage
  revenue
  reputation

struct RevenueSharing
  shares List(RevenueShare)
  duration UInt64
  triggers List(Text)

struct RevenueShare
  agent AgentId
  percentage Float32
  minimum UInt64
  maximum UInt64

struct AttributionRewards
  citationReward UInt64
  usageReward UInt64
  derivationReward UInt64

# Agent coordination and consensus
struct CoordinationMessage
  coordinationType CoordinationType
  proposal CoordinationProposal
  consensus ConsensusData
  execution ExecutionPlan

enum CoordinationType
  taskAllocation
  resourceSharing
  consensusBuilding
  conflictResolution
  emergencyResponse

struct CoordinationProposal
  proposalId UInt64
  proposer AgentId
  description Text
  participants List(AgentId)
  deadline UInt64
  votingMechanism VotingMechanism

enum VotingMechanism
  simple
  weighted
  quadratic
  futarchy
  liquid

struct ConsensusData
  votes List(Vote)
  result ConsensusResult
  finalized Bool
  proof ConsensusProof

struct Vote
  voter AgentId
  choice VoteChoice
  weight UInt64
  reasoning Text
  signature PQSignature

struct VoteChoice
  union
    binary Bool
    multiple UInt32
    ranked List(UInt32)
    continuous Float64

struct ConsensusResult
  outcome VoteChoice
  confidence Float64
  participation Float64
  finalizationTime UInt64

struct ConsensusProof
  proofType ProofType
  proof Data
  validators List(AgentId)

enum ProofType
  aggregateSignature
  merkleTree
  zkProof
  threshold

struct ExecutionPlan
  tasks List(Task)
  dependencies List(Dependency)
  timeline Timeline
  monitoring MonitoringPlan

struct Task
  taskId UInt64
  assignee AgentId
  description Text
  requirements TaskRequirements
  deliverables List(Deliverable)
  deadline UInt64

struct TaskRequirements
  compute ComputeRequirements
  data List(DataRequirement)
  permissions List(Permission)
  dependencies List(UInt64)

struct DataRequirement
  dataType DataType
  source AgentId
  access AccessLevel
  privacy PrivacyLevel

enum AccessLevel
  read
  write
  execute
  admin

struct Permission
  resource Text
  action Text
  grantor AgentId
  expiration UInt64

struct Deliverable
  deliverableId UInt64
  description Text
  format DataFormat
  verification VerificationMethod
  quality QualityMetrics

enum VerificationMethod
  manual
  automated
  consensus
  oracle

struct QualityMetrics
  accuracy Float64
  completeness Float64
  timeliness Float64
  relevance Float64

struct Dependency
  fromTask UInt64
  toTask UInt64
  dependencyType DependencyType
  data List(DataDependency)

enum DependencyType
  sequential
  parallel
  conditional
  resource

struct DataDependency
  dataId UInt64
  source AgentId
  format DataFormat
  size UInt64

struct Timeline
  start UInt64
  end UInt64
  milestones List(Milestone)
  criticalPath List(UInt64)

struct Milestone
  milestoneId UInt64
  name Text
  deadline UInt64
  criteria List(Text)
  rewards List(Reward)

struct Reward
  recipient AgentId
  amount UInt64
  currency CurrencyType
  condition Text

struct MonitoringPlan
  metrics List(MonitoringMetric)
  alerts List(AlertRule)
  reports List(ReportSchedule)

struct MonitoringMetric
  name Text
  source AgentId
  frequency UInt32
  threshold List(Threshold)

struct Threshold
  level AlertLevel
  value Float64
  action Text

enum AlertLevel
  info
  warning
  error
  critical

struct AlertRule
  condition Text
  level AlertLevel
  recipients List(AgentId)
  cooldown UInt32

struct ReportSchedule
  name Text
  frequency ReportFrequency
  recipients List(AgentId)
  template Text

enum ReportFrequency
  realtime
  hourly
  daily
  weekly
  monthly

# Network discovery and routing
struct DiscoveryMessage
  discoveryType DiscoveryType
  query DiscoveryQuery
  response DiscoveryResponse
  routing RoutingUpdate

enum DiscoveryType
  agentDiscovery
  serviceDiscovery
  resourceDiscovery
  routingUpdate
  healthCheck

struct DiscoveryQuery
  queryId UInt64
  requester AgentId
  criteria SearchCriteria
  maxResults UInt32
  timeout UInt32

struct SearchCriteria
  capabilities List(ComputeType)
  specializations List(Text)
  location GeographicConstraints
  trustLevel TrustLevel
  priceRange PriceRange

struct GeographicConstraints
  regions List(Text)
  maxLatency UInt32
  dataResidency List(Text)

struct PriceRange
  minPrice UInt64
  maxPrice UInt64
  currency CurrencyType
  pricingModel PricingModel

enum PricingModel
  fixed
  auction
  dynamic
  tiered

struct DiscoveryResponse
  queryId UInt64
  responder AgentId
  matches List(AgentMatch)
  total UInt32

struct AgentMatch
  agent AgentId
  score Float64
  availability Availability
  pricing PricingInfo
  testimonials List(Testimonial)

struct Availability
  status AvailabilityStatus
  nextAvailable UInt64
  capacity CapacityInfo
  schedule List(ScheduleSlot)

enum AvailabilityStatus
  available
  busy
  maintenance
  offline

struct CapacityInfo
  currentLoad Float64
  maxCapacity UInt64
  reservedCapacity UInt64
  availableCapacity UInt64

struct ScheduleSlot
  start UInt64
  end UInt64
  capacity UInt64
  price UInt64

struct PricingInfo
  model PricingModel
  basePrice UInt64
  currency CurrencyType
  discounts List(Discount)
  premiums List(Premium)

struct Discount
  condition Text
  percentage Float32
  maximum UInt64

struct Premium
  condition Text
  percentage Float32
  minimum UInt64

struct Testimonial
  client AgentId
  rating Float32
  review Text
  timestamp UInt64
  verified Bool

# Routing information for message delivery
struct RoutingInfo
  path List(AgentId)
  hops UInt32
  latency UInt32
  cost UInt64
  priority MessagePriority
  qos QoSRequirements

enum MessagePriority
  low
  normal
  high
  critical
  realtime

struct QoSRequirements
  maxLatency UInt32
  minBandwidth UInt64
  reliability Float32
  security SecurityLevel

enum SecurityLevel
  none
  basic
  encrypted
  pqEncrypted
  anonymous

struct RoutingUpdate
  updateType UpdateType
  routes List(RouteInfo)
  topology NetworkTopology
  metrics NetworkMetrics

enum UpdateType
  add
  remove
  update
  full

struct RouteInfo
  destination AgentId
  nextHop AgentId
  cost UInt64
  latency UInt32
  reliability Float32
  lastUpdate UInt64

struct NetworkTopology
  nodes List(NetworkNode)
  edges List(NetworkEdge)
  clusters List(NetworkCluster)

struct NetworkNode
  agent AgentId
  position NetworkPosition
  connections UInt32
  centrality Float64

struct NetworkPosition
  x Float64
  y Float64
  z Float64

struct NetworkEdge
  from AgentId
  to AgentId
  weight Float64
  bandwidth UInt64
  latency UInt32

struct NetworkCluster
  clusterId UInt64
  members List(AgentId)
  representative AgentId
  properties ClusterProperties

struct ClusterProperties
  specialization Text
  trustLevel TrustLevel
  capacity UInt64
  load Float64

struct NetworkMetrics
  totalNodes UInt64
  totalEdges UInt64
  averageLatency Float64
  networkThroughput UInt64
  reliability Float64
  timestamp UInt64

# Zero-copy shared data reference
struct SharedDataReference
  sharedId UInt64
  offset UInt64
  size UInt64
  checksum UInt64
  permissions List(Text)

# RPC protocol configuration
struct RPCConfig
  protocol RPCProtocol
  endpoint Text
  options ProtocolOptions
  security SecurityConfig

enum RPCProtocol
  zmq
  capnp
  grpc
  websocket
  custom

struct ProtocolOptions
  union
    zmq ZMQOptions
    capnp CapnpOptions
    grpc GRPCOptions
    websocket WebSocketOptions

struct ZMQOptions
  socketType ZMQSocketType
  highWaterMark UInt64
  linger UInt32
  reconnectInterval UInt32
  maxReconnects UInt32

enum ZMQSocketType
  req
  rep
  dealer
  router
  pub
  sub
  xpub
  xsub
  push
  pull
  pair
  stream

struct CapnpOptions
  compression CompressionType
  packedEncoding Bool
  traversalLimit UInt64
  nestingLimit UInt32

enum CompressionType
  none
  gzip
  lz4
  zstd
  snappy

struct GRPCOptions
  maxReceiveMessageSize UInt64
  maxSendMessageSize UInt64
  keepaliveTime UInt32
  keepaliveTimeout UInt32
  maxConnectionIdle UInt32
  compression GRPCCompression

enum GRPCCompression
  identity
  gzip
  deflate

struct WebSocketOptions
  maxFrameSize UInt64
  maxMessageSize UInt64
  compression Bool
  pingInterval UInt32
  pongTimeout UInt32

struct SecurityConfig
  tls TLSConfig
  authentication AuthConfig
  authorization AuthzConfig
  rateLimit RateLimitConfig

struct TLSConfig
  enabled Bool
  certificate Data
  privateKey Data
  cipherSuites List(Text)
  minVersion Text

struct AuthConfig
  method AuthMethod
  credentials AuthCredentials
  tokenExpiry UInt64
  refreshEnabled Bool

enum AuthMethod
  none
  apiKey
  jwt
  oauth2
  mutual
  pqc

struct AuthCredentials
  union
    apiKey Data
    jwt JWTCredentials
    oauth2 OAuth2Credentials
    pqc PQCredentials

struct JWTCredentials
  token Text
  publicKey Data
  algorithm Text

struct OAuth2Credentials
  clientId Text
  clientSecret Text
  accessToken Text
  refreshToken Text
  scope List(Text)

struct PQCredentials
  algorithm PQAlgorithm
  publicKey Data
  signature Data
  certificate Data

struct AuthzConfig
  enabled Bool
  policies List(AuthzPolicy)
  defaultAction AuthzAction

struct AuthzPolicy
  resource Text
  action Text
  principal Text
  effect AuthzAction
  conditions List(Text)

enum AuthzAction
  allow
  deny

struct RateLimitConfig
  enabled Bool
  requests UInt64
  window UInt32
  burst UInt64
  strategy RateLimitStrategy

enum RateLimitStrategy
  fixedWindow
  slidingWindow
  tokenBucket
  leakyBucket

# Agentic network service interface
interface AgenticNetwork
  # Send a message to another agent
  sendMessage (message AgenticMessage) -> (result SendResult)

  # Discover agents matching criteria
  discover (query DiscoveryQuery) -> (response DiscoveryResponse)

  # Request computation from another agent
  requestComputation (request ComputationRequest) -> (result ComputationResult)

  # Share intelligence data
  shareIntelligence (share IntelligenceShare) -> (result ShareResult)

  # Coordinate with other agents
  coordinate (message CoordinationMessage) -> (result CoordinationResult)

  # Subscribe to network events
  subscribe (filter EventFilter) -> stream (event NetworkEvent)

struct SendResult
  success Bool
  messageId UInt64
  deliveryTime UInt64
  error Text

struct ComputationResult
  success Bool
  output Data
  metrics ComputeMetrics
  proof Data
  error Text

struct ComputeMetrics
  executionTime UInt64
  memoryUsed UInt64
  computeUsed UInt64

struct ShareResult
  success Bool
  shareId UInt64
  accessGrants List(AgentId)
  error Text

struct CoordinationResult
  success Bool
  consensus ConsensusResult
  execution ExecutionPlan
  error Text

struct EventFilter
  eventTypes List(EventType)
  agents List(AgentId)
  topics List(Text)

enum EventType
  messageReceived
  computationComplete
  intelligenceShared
  consensusReached
  routingChanged

struct NetworkEvent
  eventType EventType
  timestamp UInt64
  source AgentId
  data Data
