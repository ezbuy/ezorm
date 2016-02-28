package cmd

import (
	"fmt"
	"testing"
	"github.com/ezbuy/ezorm/mysql"
	"github.com/ezbuy/ezorm/page"
)

func TestGetRefIf(t *testing.T) {
	conf := new(mysql.MySQLConfig)
	conf.MySQLDB = "root:zhangpei@/test_db"
	mysql.Setup(conf)
	result := page.PageMgr.Query("SELECT * FROM mail")

	for key, value := range result {
		fmt.Printf("%s : %s\n", key, value)
	}
}
