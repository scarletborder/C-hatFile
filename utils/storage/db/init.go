package db

import (
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

func init() {
	// 启动数据库服务init

	// 读取所有user-dbname
	ReadConfig()

	// 连接数据库推迟到组件初始化时
}

func ReadConfig() {
	DBReadConfigWG.Add(1)
	go func() {
		const config_path = `utils\storage\db\config.toml`
		// 异步加载文件配置
		defer func() {
			DBReadConfigWG.Done()
			// 写回文件
			writeConfigWithBackup(config_path, &DBCfg)
		}()
		if _, err := toml.DecodeFile(config_path, &DBCfg); err != nil {
			logrus.Warnln("无法加载DB配置文件，即将初始化")
		}

		// Use reflection to iterate over DatabaseConfig fields
		dbConfigValue := reflect.ValueOf(&DBCfg.SubDBCfgs).Elem()
		dbConfigType := dbConfigValue.Type()

		for i := 0; i < dbConfigValue.NumField(); i++ {
			field := dbConfigValue.Field(i)
			fieldName := dbConfigType.Field(i).Tag.Get("toml")
			if fieldName == "" {
				fieldName = dbConfigType.Field(i).Name
			}

			if field.Kind() == reflect.Struct {
				// Update the User and Password fields
				DB_nameField := field.FieldByName("DB_name")
				passwordField := field.FieldByName("Password")
				if DB_nameField.IsValid() && DB_nameField.String() == "" {
					// 有相关条目但是没有初始化
					pwd, err := generateRandomString(24)
					if err != nil {
						logrus.Errorln("无法为db生成随机密码，使用固定密码代替")
						pwd = "sdf!SsoOjkb5j49FFXtrh05!"
					}
					DB_nameField.SetString(fieldName)
					passwordField.SetString(pwd)
				}
			}
		}
	}()
}
