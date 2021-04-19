package echolet

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var once sync.Once

func ConnectMySQL() *sql.DB {
	var db *sql.DB
	once.Do(func() {
		cnf := LoadMySqlDSN()
		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", cnf.User, cnf.Pass, cnf.Host, cnf.Port, cnf.DbName)
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		// SetConnMaxLifetime() is required to ensure connections are closed by the
		// driver safely before connection is closed by MySQL server, OS, or other middlewares.
		// Since some middlewares close idle connections by 5 minutes, we recommend timeout
		// shorter than 5 minutes. This setting helps load balancing and changing system
		// variables too.
		db.SetConnMaxLifetime(time.Minute * 3)
		//SetMaxOpenConns() is highly recommended to limit the number of connection used by the
		// application. There is no recommended limit number because it depends on application
		// and MySQL server.
		db.SetMaxOpenConns(cnf.PoolSize)
		// SetMaxIdleConns() is recommended to be set same to (or greater than)
		// db.SetMaxOpenConns(). When it is smaller than SetMaxOpenConns(), connections can be
		// opened and closed very frequently than you expect. Idle connections can be closed by
		// the db.SetConnMaxLifetime(). If you want to close idle connections more rapidly,
		// you can use db.SetConnMaxIdleTime() since Go 1.15.
		db.SetMaxIdleConns(cnf.PoolSize)

		err = db.Ping()
		if err != nil {
			log.Fatalf("Error: Can not connect to mysql! cause by %v", err)
		}
	})
	return db
}
