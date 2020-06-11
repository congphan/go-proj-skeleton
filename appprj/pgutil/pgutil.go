package pgutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// This is needed to have the driver to acess PG server
	"github.com/go-pg/pg/v9"
	_ "github.com/lib/pq"

	"github.com/golang/glog"

	"go-prj-skeleton/appprj/setting"
)

// Configuration struct
type Configuration struct {
	URL             string
	Host            string
	Port            string
	Database        string
	User            string
	Password        string
	MaxConnections  string
	ApplicationName string
}

var dbSession *pg.DB

// ConnectPG is for accessing PG server
func ConnectPG(query string) (*sql.Rows, error) {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		setting.ProjectEnvSettings.PostgreHost,
		setting.ProjectEnvSettings.PostgrePort,
		setting.ProjectEnvSettings.PostgreUser,
		setting.ProjectEnvSettings.PostgrePassword,
		setting.ProjectEnvSettings.PostgreDatabaseName)

	db, err := sql.Open("postgres", dbinfo)
	defer db.Close()
	if err != nil {
		glog.Errorf("Cannot connect to DB, error: %s", err)
		return nil, err
	}
	rows, err := db.Query(query)
	if err != nil {
		glog.Errorf("Cannot query, error: %s", err)
		return nil, err
	}
	return rows, nil
}

// DB : return connection to DB
func DB() *pg.DB {
	return dbSession
}

// StartUp ...
func StartUp(config Configuration) {
	pg.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	if config.URL != "" {
		options, err := pg.ParseURL(config.URL)
		if err != nil {
			panic(err)
		}
		options.ApplicationName = config.ApplicationName
		dbSession = pg.Connect(options)
		return
	}

	dbSession = pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%s", config.Host, config.Port),
		User:            config.User,
		Password:        config.Password,
		Database:        config.Database,
		TLSConfig:       nil,
		ApplicationName: config.ApplicationName,
	})
}

// Shutdown ...
func Shutdown() {
	dbSession.Close()
}
