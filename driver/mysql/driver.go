package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jayvib/golog"
	"gophr.v2/config"
	"net/url"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

var _testInitCount int // Used for unit testing

var (
	mu sync.Mutex
	dbPool map[string]*sql.DB
)

func init() {
	dbPool = make(map[string]*sql.DB)
}

func New(conf *config.Config, dbName string) (db *sql.DB, err error) {
	mu.Lock()
	defer mu.Unlock()

	db, ok := dbPool[dbName]
	if ok {
		return db, nil
	}

	// Initialize for the first time
	for _, c := range conf.MySQL {
		// Initialize only the matched database configuration
		if c.Database == dbName {
			db, err = Initialize(c)
			if err != nil {
				return
			}
			dbPool[c.Database] = db
			_testInitCount++
			return
		}
	}

	return nil, errors.New("database configuration is not define")
}

func Initialize(s config.MySQL) (*sql.DB, error) {
	var err error

	// Using Singleton Pattern
		var e error
		format := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			s.User, s.Password, s.Host, s.Port, s.Database)
		golog.Info(format)
		val := url.Values{}
		val.Add("parseTime", "1")
		val.Add("loc", "Asia/Manila")
		dsn := fmt.Sprintf("%s?%s", format, val.Encode())
		golog.Debug("DSN:", dsn)
		db, e = sql.Open("mysql", dsn)
		if e != nil {
			err = e
		}
		if e := db.Ping(); e != nil {
			err = e
		}

	if err != nil {
		return nil, err
	}

	return db, nil
}
