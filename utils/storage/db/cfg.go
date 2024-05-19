package db

import (
	"sync"
)

// host
type DBConfig struct {
	SubDBCfgs SubUserConfigs `toml:"database"` // db名-验证用户，仅管理非客户

	Addr         string
	RootUser     string `toml:"root_user"`
	RootPassword string `toml:"root_password"`
}

type SubUserConfigs struct {
	AuthDB       SubUserConfig `toml:"auth_db"`
	FileReaderDB SubUserConfig `toml:"file_reader"`
	FileWriterDB SubUserConfig `toml:"file_writer"`
}

type SubUserConfig struct {
	Username string `toml:"Username"`
	Password string
}

var (
	DBCfg          DBConfig
	DBReadConfigWG sync.WaitGroup
)
