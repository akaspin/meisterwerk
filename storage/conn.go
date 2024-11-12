package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Conn struct {
	GORM *gorm.DB
	SQL  *sql.DB

	isClosed  int32
	closeChan chan struct{}
}

func Connect(c *DBConfig) (*Conn, error) {
	dialector := postgres.Open(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s %s TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.Database, c.Options))

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stderr, "\r\n", log.Lshortfile),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
			},
		),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(c.MaxConn)
	sqlDB.SetMaxIdleConns(c.MaxConn)
	sqlDB.SetConnMaxLifetime(c.ConnLifetime)

	conn := &Conn{
		GORM:      db,
		SQL:       sqlDB,
		closeChan: make(chan struct{}),
	}

	go func() { // watchdog
		<-conn.closeChan
		_ = conn.SQL.Close()
	}()
	return conn, nil
}

func (c *Conn) Close() error {
	if atomic.CompareAndSwapInt32(&c.isClosed, 0, 1) {
		close(c.closeChan)
	}
	return nil
}
