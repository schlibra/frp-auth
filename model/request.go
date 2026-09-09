package model

type requestBodyContentMetas struct {
	Token string `json:"token"`
}
type requestBodyContent struct {
	Arch          string                  `json:"arch"`
	ClientAddress string                  `json:"client_address"`
	ClientSpec    any                     `json:"client_spec"`
	Hostname      string                  `json:"hostname"`
	Metas         requestBodyContentMetas `json:"metas"`
	OS            string                  `json:"os"`
	PoolCount     int                     `json:"pool_count"`
	PrivilegeKey  string                  `json:"privilege_key"`
	RunID         string                  `json:"run_id"`
	Timestamp     int                     `json:"timestamp"`
	User          string                  `json:"user"`
	Version       string                  `json:"version"`
}
type RequestBody struct {
	Version string             `json:"version"`
	Op      string             `json:"op"`
	Content requestBodyContent `json:"content"`
}
