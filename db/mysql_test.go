package db

import (
	"testing"
)

func TestMysqlConfigConvert(t *testing.T) {
	cases := []struct {
		cfg *MysqlFieldConfig
		dsn string
	}{
		{
			cfg: &MysqlFieldConfig{
				Addr:     "localhost:3306",
				UserName: "root",
				Password: "u1s1",
				Database: "spk",
				Options: map[string]string{
					"parseTime": "true",
				},
			},
			dsn: "root:u1s1@tcp(localhost:3306)/spk?charset=utf8mb4&parseTime=true",
		},
		{
			cfg: &MysqlFieldConfig{
				Addr:     "10.24.32.12:8806",
				UserName: "blog_user",
				Password: "",
				Database: "blog",
			},
			dsn: "blog_user@tcp(10.24.32.12:8806)/blog?charset=utf8mb4",
		},
		{
			cfg: &MysqlFieldConfig{
				Addr:     "192.168.12.3",
				Database: "test",
			},
			dsn: "@tcp(192.168.12.3)/test?charset=utf8mb4",
		},
		{
			cfg: &MysqlFieldConfig{
				Addr:     "192.168.12.68",
				UserName: "product",
				Password: "product_123",

				Options: map[string]string{
					"parseTime":  "true",
					"autocommit": "true",
				},
			},
			dsn: "product:product_123@tcp(192.168.12.68)/?charset=utf8mb4&parseTime=true&autocommit=true",
		},
		{
			cfg: &MysqlFieldConfig{
				Addr:     "10.12.32.1:3306,10.12.32.2:8806,10.12.32.3:8921",
				UserName: "cluster",
				Password: "cluster_pass",
			},
			dsn: "cluster:cluster_pass@tcp(10.12.32.1:3306,10.12.32.2:8806,10.12.32.3:8921)/?charset=utf8mb4",
		},
	}

	for _, cs := range cases {
		cfg := cs.cfg.Convert()
		if cfg.DataSource != cs.dsn {
			t.Fatalf("invalid dsn: %s", cfg.DataSource)
		}
	}
}
