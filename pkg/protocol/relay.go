package protocol

type NodeStats struct {
	Connections int    `json:"connections"`
	Ingress     uint64 `json:"ingress_bytes"`
	Egress      uint64 `json:"egress_bytes"`
	Dropped     uint64 `json:"dropped_packets"`
}
type NodeLimits struct {
	MaxConnections int   `json:"max_connections"`
	BytesPerSecond int64 `json:"bytes_per_second"`
}
