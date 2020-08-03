package config

type Config struct {
	RootPath string
	Port     uint16
	Auth     LocalAuth
}

type LocalAuth struct {
	Username string
	Password string
}
