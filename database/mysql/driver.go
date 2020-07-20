package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hellgate75/go-services/database"
)

type mySQLDriver struct {
}

func (d *mySQLDriver) Connect(config database.DbConfig) (database.Connection, error) {
	connStr := config.Url
	//db, err := sql.Open("mysql", "<username>:<pw>@tcp(<HOST>:<port>)/<dbname>")
	if connStr == "" {
		if config.Database.Database == "" {
			if config.Host != "" && config.Port > 0 {
				connStr = fmt.Sprintf("%s:%s@tcp(%s:%v)", config.Name, config.Password, config.Host, config.Port)
			} else {
				connStr = fmt.Sprintf("%s:%s@", config.Name, config.Password)
			}
		} else {
			if config.Host != "" && config.Port > 0 {
				connStr = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", config.Name, config.Password, config.Host, config.Port, config.Database.Database)
			} else {
				connStr = fmt.Sprintf("%s:%s@/%s", config.Name, config.Password, config.Database.Database)
			}
		}
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &mySqlConnection{
		Configuration: config,
		DB:            db,
	}, nil
}

func GetMySqlDriver() database.Driver {
	return &mySQLDriver{}
}

