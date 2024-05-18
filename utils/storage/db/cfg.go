package db

import (
	"sync"
)

// host
type DBConfig struct {
	SubDBCfgs SubDBConfigs `toml:"database"` // db名-验证用户，仅管理非客户

	Addr         string
	RootUser     string `toml:"root_user"`
	RootPassword string `toml:"root_password"`
}

type SubDBConfigs struct {
	AuthDB       SubDBConfig `toml:"auth_db"`
	FileReaderDB SubDBConfig `toml:"file_reader"`
	FileWriterDB SubDBConfig `toml:"file_writer"`
}

type SubDBConfig struct {
	DB_name  string `toml:"DB_name"`
	Password string
}

var (
	DBCfg          DBConfig
	DBReadConfigWG sync.WaitGroup
)
