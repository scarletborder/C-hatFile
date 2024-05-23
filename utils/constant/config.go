package constant

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

const Db_config_path = "./db_config.toml"     // 数据库
const S3_config_path = "./s3_config.json"     // s3节点
const Misc_config_path = "./misc_config.toml" // 杂项

const db_config = `Addr = "127.0.0.1:3306"
root_user = "root"
root_password = ""

[database]
  [database.auth_db]
    Username = "auth_db"
    Password = ""
  [database.file_reader]
    Username = "file_reader"
    Password = ""
  [database.file_writer]
    Username = "file_writer"
    Password = ""`

const s3_config = `{
    "s3_points": [
        {
            "your-s3-endpoint-com": "127.0.0.1:9000",
            "your-access-key": "admin",
            "your-secret-key": "password",
            "chunksize": 0,
            "usessl": false
        }
    ],
    "bucket_name": "cffiles"
}`

const misc_config = `[blog]
directory_path = "./blogs"`

func init() {
	first_time := false
	if !fileExists(Db_config_path) {
		first_time = true
	}
	// 创建配置文件
	// 创建并写入第一个配置文件

	if first_time {
		err := os.WriteFile(Db_config_path, []byte(db_config), 0644)
		if err != nil {
			logrus.Errorf("创建%s失败: %v\n", Db_config_path, err)
		}
		err = os.WriteFile(S3_config_path, []byte(s3_config), 0644)
		if err != nil {
			logrus.Errorf("创建%s失败: %v\n", S3_config_path, err)
		}
		err = os.WriteFile(Misc_config_path, []byte(misc_config), 0644)
		if err != nil {
			logrus.Errorf("创建%s失败: %v\n", Misc_config_path, err)
		}

		logrus.Infoln("This is the first time to launch of program, fill config file and start again")
		os.Exit(0)
	}

	InitMiscConfig()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func InitMiscConfig() {
	if _, err := toml.DecodeFile(Misc_config_path, &MiscCfg); err != nil {
		logrus.Warnln("无法加载misc配置文件")
	}
}
