package model

type FrpsClientDataItem struct {
	ClientID         string `json:"clientID"`
	ClientIP         string `json:"clientIP"`
	FirstConnectedAt int64  `json:"firstConnectedAt"`
	Hostname         string `json:"hostname"`
	Key              string `json:"key"`
	LastConnectedAt  int64  `json:"lastConnectedAt"`
	Online           bool   `json:"online"`
	RunID            string `json:"runID"`
	User             string `json:"user"`
	Version          string `json:"version"`
	WireProtocol     string `json:"wireProtocol"`
}
type FrpsClientData struct {
	Items    []FrpsClientDataItem `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int                  `json:"total"`
}

type FrpsClient struct {
	Code int            `json:"code"`
	Data FrpsClientData `json:"data"`
	Msg  string         `json:"msg"`
}
type FrpsProxyDataItemSpecTypeTransport struct {
	UseEncryption      bool   `json:"useEncryption"`
	UseCompression     bool   `json:"useCompression"`
	BandWidthLimit     string `json:"bandWidthLimit"`
	BandWidthLimitMode string `json:"bandWidthLimitMode"`
}
type FrpsProxyDataItemSpecTypeLoadBalancer struct {
	Group string `json:"group"`
}
type FrpsProxyDataItemSpecType struct {
	Transport    FrpsProxyDataItemSpecTypeTransport    `json:"transport"`
	LoadBalancer FrpsProxyDataItemSpecTypeLoadBalancer `json:"loadBalancer"`
	RemotePort   int                                   `json:"remotePort"`
}
type FrpsProxyDataItemSpec struct {
	Type string                    `json:"type"`
	Tcp  FrpsProxyDataItemSpecType `json:"tcp"`
	Udp  FrpsProxyDataItemSpecType `json:"udp"`
}
type FrpsProxyDataItemStatus struct {
	Phase           string `json:"phase"`
	TodayTrafficIn  int64  `json:"todayTrafficIn"`
	TodayTrafficOut int64  `json:"todayTrafficOut"`
	CurConns        int    `json:"curConns"`
	LastStartAt     int64  `json:"lastStartAt"`
	LastCloseAt     int64  `json:"lastCloseAt"`
}
type FrpsProxyDataItem struct {
	Name     string                  `json:"name"`
	User     string                  `json:"user"`
	ClientID string                  `json:"clientID"`
	Spec     FrpsProxyDataItemSpec   `json:"spec"`
	Status   FrpsProxyDataItemStatus `json:"status"`
}
type FrpsProxyData struct {
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	Items    []FrpsProxyDataItem `json:"items"`
}
type FrpsProxy struct {
	Code int           `json:"code"`
	Data FrpsProxyData `json:"data"`
	Msg  string        `json:"msg"`
}
