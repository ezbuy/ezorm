package db

type Config struct {
	DB *DB
}

type DB struct {
	SqlServerConn string
}
