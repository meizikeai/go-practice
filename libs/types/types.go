package types

type ConfMySQL struct {
	Master   []string `json:"master"`
	Slave    []string `json:"slave"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Database string   `json:"database"`
}

type ConfRedis struct {
	Master   []string `json:"master"`
	Password string   `json:"password"`
	Db       int      `json:"db"`
}
