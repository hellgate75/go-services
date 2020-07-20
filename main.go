package go_services

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-services/database"
	"github.com/hellgate75/go-services/database/mongodb"
	"github.com/hellgate75/go-services/database/mysql"
	"strings"
)

func GetDatabaseRef(dType database.DriverType, databaseName string, databaseEntity string) (database.DataRef, error) {
	switch dType {
	case database.MongoDbDriver:
		return database.DataRef{
			Database: databaseName,
			Namespace: databaseEntity,
		}, nil
	default:
		return database.DataRef{}, errors.New(fmt.Sprintf("Unknown Database Driver type: %v", dType))
	}

}

func GetDatabaseDriver(dType database.DriverType) (database.Driver, error) {
	switch dType {
	case database.MongoDbDriver:
		return mongodb.GetMongoDriver(), nil
	case database.MySQLDriver:
		return mysql.GetMySqlDriver(), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown Database Driver type: %v", dType))
	}
}

// You can request database using the driver name: mysql, mongodb
func GetDatabaseDriverByName(driverName string) (database.Driver, error) {
	var dType = database.DriverToType(driverName)
	switch dType {
	case database.MongoDbDriver:
		return mongodb.GetMongoDriver(), nil
	case database.MySQLDriver:
		return mysql.GetMySqlDriver(), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown Database Driver type: %v", dType))
	}
}

type ServiceType byte
const(
	UnknownService		ServiceType = 0
	DatabaseService		ServiceType = iota + 1
)

type ServiceDiscovery interface {
	GetServiceByName(name string) (ServiceType, error)
	GetDatabaseDriver(dType database.DriverType) (database.Driver, error)
	GetDatabaseRef(dType database.DriverType, databaseName string, databaseEntity string) (database.DataRef, error)
}

type _serviceDiscovery struct{}

func (sd *_serviceDiscovery) GetServiceByName(name string) (ServiceType, error) {
	var nmLwr = strings.ToLower(name)
	switch nmLwr  {
	case "mongodb":
		return DatabaseService, nil
	default:
		return UnknownService, errors.New(fmt.Sprintf("Unable to discover service type: %s", name))
	}
}

func (sd *_serviceDiscovery) GetDatabaseDriver(dType database.DriverType) (database.Driver, error) {
	return GetDatabaseDriver(dType)
}

func (sd *_serviceDiscovery) GetDatabaseRef(dType database.DriverType, databaseName string, databaseEntity string) (database.DataRef, error) {
	return GetDatabaseRef(dType, databaseName, databaseEntity)
}

var Discovery ServiceDiscovery = &_serviceDiscovery{}