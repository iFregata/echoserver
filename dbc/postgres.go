package dbc

import (
	"context"
	"echoserver/gorux"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDSN struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	DbName   string `json:"dbname"`
	PoolSize int    `json:"pool_size"`
}

var once sync.Once

func ConnectPostgres() *pgxpool.Pool {
	var db *pgxpool.Pool
	once.Do(func() {
		cnf := loadPostgresDSN()

		//user=jack password=secret host=pg.example.com port=5432 dbname=mydb pool_max_conns=10
		dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=%d",
			cnf.User, cnf.Pass, cnf.Host, cnf.Port, cnf.DbName, cnf.PoolSize)
		ctx := context.Background()
		var err error
		db, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping(ctx)
		if err != nil {
			log.Fatalf("Error: Can not connect to Postgres! cause by %v", err)
		}
	})
	return db
}
func loadPostgresDSN() *PostgresDSN {
	return gorux.LoadConfigFile("config/postgres.json", &PostgresDSN{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Pass:     "pass4pass",
		DbName:   "gorux",
		PoolSize: 10,
	}).(*PostgresDSN)
}
