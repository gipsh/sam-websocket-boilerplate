package common

type ConnectionItem struct {
	// ID  string `json:"connectionID"`
	ConnectionID string   `json:"ConnectionID"`
	Created      string   `json:"Created"`
	SrcIP        string   `json:"SrcIP"`
	Messages     []string `json:"Messages"`
}
