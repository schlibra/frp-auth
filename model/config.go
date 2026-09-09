package model

type mysqlConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"user"`
	Password string `toml:"pass"`
	Database string `toml:"name"`
}

type serverConfig struct {
	Host  string `toml:"host"`
	Port  int    `toml:"port"`
	Debug bool   `toml:"debug"`
}

type jwtConfig struct {
	Key string `toml:"key"`
}

type frpsConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type frpsWebConfig struct {
	Host  string `toml:"host"`
	Port  int    `toml:"port"`
	Token string `toml:"token"`
}

type dangerousConfig struct {
	AllowClean bool `toml:"allowClean"`
}

type Config struct {
	MySQL     mysqlConfig     `toml:"mysql"`
	Server    serverConfig    `toml:"server"`
	Jwt       jwtConfig       `toml:"jwt"`
	Frps      frpsConfig      `toml:"frps"`
	FrpsWeb   frpsWebConfig   `toml:"frps-web"`
	Dangerous dangerousConfig `toml:"dangerous"`
}
