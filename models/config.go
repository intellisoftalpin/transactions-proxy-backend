package models

type Config struct {
	ServerPort   string `json:"server_port"`
	CNodeAddress string `json:"cnode_host"`
	HostName     string `json:"host_name"`

	CertPath string `json:"cert_path"`

	DB    DBConfig `json:"db"`
	Pools []string `json:"pools"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
