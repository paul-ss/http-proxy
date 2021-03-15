package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	config "github.com/paul-ss/http-proxy/configs"
	"log"
	"sync"
)

var (
	once     = &sync.Once{}
	instance *pgx.Conn
	callErr  error
	opened   bool
)


func NewPostgresConn() (*pgx.Conn, error) {
	once.Do(func() {
		conn, err := pgx.Connect(
			context.Background(),
			fmt.Sprintf(
				"user=%s password=%s host=%s port=%s dbname=%s",
				config.C.Db.User,
				config.C.Db.Password,
				config.C.Db.Host,
				config.C.Db.Port,
				config.C.Db.DbName,
				),
			)

		if err != nil {
			callErr = err
			return
		}

		opened = true
		instance = conn
	})

	if callErr != nil {
		return nil, callErr
	}

	if !opened {
		return nil, fmt.Errorf("database is not opened")
	}

	return instance, nil
}

func Close() {
	if opened {
		if err := instance.Close(context.Background()); err != nil {
			log.Println("Postgres conn close error: " + err.Error())
		} else {
			log.Println("Postgres conn closed")
		}
	}
}
