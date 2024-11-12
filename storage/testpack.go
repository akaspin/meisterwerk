package storage

import (
	"math/rand/v2"
	"os"
	"time"
)

func TestPackDBConfig(database string) *DBConfig {
	host := os.Getenv("TEST_DB_CONFIG_HOST")
	if host == "" {
		host = "localhost"
	}
	return &DBConfig{
		Host:         host,
		Port:         5432,
		Database:     database,
		User:         "postgres",
		Options:      "sslmode=disable TimeZone=UTC",
		MaxConn:      50,
		ConnLifetime: time.Minute * 5,
		Password:     "postgres",
	}
}

var r = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

func TestPackNextID() int32 {
	return r.Int32()
}
