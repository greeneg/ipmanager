package globals

type Config struct {
	TcpPort int    `json:"tcpPort"`
	DbPath  string `json:"dbPath"`
}
