package db_test

import (
	"chatFileBackend/utils/storage/db"
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	db.ReadConfig()
	fmt.Println(db.DBCfg)
}
