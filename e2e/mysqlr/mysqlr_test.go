package mysqlr

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	port, err := strconv.ParseInt(os.Getenv("MYSQL_PORT"), 10, 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse MYSQL_PORT: %s", err))
	}
	MySQLSetup(&MySQLConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     int(port),
		UserName: "ezbuy",
		Password: "ezbuyisthebest",
		Database: "test",
	})
	sql, err := os.ReadFile("gen.script.mysqlr.blog.sql")
	if err != nil {
		panic(fmt.Errorf("failed to read gen.script.mysqlr.blog.sql: %s", err))
	}
	if _, err := MySQL().Exec(ctx, "CREATE DATABASE IF NOT EXISTS test"); err != nil {
		panic(fmt.Errorf("failed to create database: %s", err))
	}
	for _, q := range strings.Split(string(sql), ";") {
		if len(strings.TrimSpace(q)) == 0 {
			continue
		}
		if _, err := MySQL().Exec(ctx, q); err != nil {
			panic(fmt.Errorf("failed to create table: %s", err))
		}
	}
	os.Exit(m.Run())
}

func TestBlogs(t *testing.T) {

}

func TestBlogsTx(t *testing.T) {}
