package utilities

type ConfigStruct struct {
	Address      string
	AdminToken   string
	TokenKeyFile string
	Secure       bool
	ApiCertFile  string
	ApiKeyFile   string
	Interval     int

	MeshHostname string
	MeshUsername string
	MeshPassword string
	PyVenvName   string
}

type InfoResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type TokenCreateDetails struct {
	Name string `json:"name"`
}

type TokenCreateBody struct {
	AuthToken string             `json:"authtoken"`
	Details   TokenCreateDetails `json:"details"`
}

type TokenListBody struct {
	AuthToken string `json:"authtoken"`
}

type TaskData struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Nodeids  []string `json:"nodeids"`
	Creation string   `json:"creation"`
	Status   string   `json:"status"`
}

type TaskBody struct {
	AuthToken string   `json:"authtoken"`
	Details   TaskData `json:"details"`
}

// Python wrapper objects.

type Device struct {
	Name   string `json:"name"`
	NodeID string `json:"nodeid"`
}

type PyOnlineDevices struct {
	OnlineDevices []Device `json:"online_devices"`
	TotalDevices  int      `json:"total_devices"`
}
