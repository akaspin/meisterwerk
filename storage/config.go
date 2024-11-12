package storage

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/pflag"
)

type DBConfig struct {
	Host         string
	Port         int
	Database     string
	User         string
	Password     string
	Options      string
	MaxConn      int
	ConnLifetime time.Duration
}

func (c *DBConfig) FS(fs *pflag.FlagSet, dbName string) {
	fs.StringVar(&c.Host, "db-host", "localhost", "database host")
	fs.IntVar(&c.Port, "db-port", 5432, "database port")
	fs.StringVar(&c.User, "db-user", "postgres", "database user")
	fs.StringVar(&c.Password, "db-pass", "postgres", "database password")
	fs.StringVar(&c.Database, "db-name", dbName, "database name")
	fs.StringVar(&c.Options, "db-options", "", "database connection options")
	fs.IntVar(&c.MaxConn, "db-max-conn", 25, "maximum database connections")
	fs.DurationVar(&c.ConnLifetime, "db-conn-lifetime", time.Minute*5, "database connection lifetime")
}

func (c DBConfig) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf(
		"%s@%s:%d/%s",
		c.User, c.Host, c.Port, c.Database))
}
