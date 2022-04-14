package constants

const ServerHost string = "localhost"
const ServerPort string = "8080"
const ProjectName string = "cosmos-opt-api"
const DefaultLoggingLevel string = "info"
const PrometheusUpdateTime int = 5
const LatestBlock int = -1
const CosmosApiUrl string = "http://localhost:26657"
const BlockPath string = "/block"
const BlockchainPath string = "/blockchain"
const BlockHeightParam string = "?height="
const MinHeightParam string = "?minHeight="
const MaxHeightParam string = "&maxHeight="
const AverageBlockWindow int = 20
const CosmostationApi string = "https://api.cosmostation.io/v1/staking/validator/"

const AverageBlocktimeGaugeName string = "average_blocktime"
const AverageBlocktimeGaugeHelp string = "Average block time in seconds"

const ValidatorUptimeGaugeName string = "validator_uptime"
const ValidatorUptimeGaugeHelp string = "Validator uptime in seconds"
