package constant

type Misc_Config struct {
	BlogCfg Blog_Config `toml:"blog"`
}

type Blog_Config struct {
	Dir_path string `toml:"directory_path"`
}

var MiscCfg Misc_Config
