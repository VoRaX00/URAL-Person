package postgres

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string
	Database string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl_mode"`
}
