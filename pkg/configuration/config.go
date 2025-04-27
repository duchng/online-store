package configuration

type DbConfig struct {
	User        string `koanf:"user"`
	Password    string `koanf:"password"`
	DbName      string `koanf:"dbName"`
	Port        string `koanf:"port"`
	Host        string `koanf:"host"`
	EnableSsl   bool   `koanf:"enableSsl"`
	AutoMigrate bool   `koanf:"autoMigrate"`
}

type RedisConfig struct {
	Host     string `koanf:"host,omitempty" `
	Password string `koanf:"password,omitempty" `
	Database int    `koanf:"database,omitempty" `
}
