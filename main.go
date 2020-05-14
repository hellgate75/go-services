package go_services

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-services/database"
	"github.com/hellgate75/go-services/database/mongodb"
)

func GetDatabaseDriver(dType database.DriverType) (database.Driver, error) {
	switch dType {
	case database.MongoDbDriver:
		return mongodb.GetMongoDriver(), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown Database Driver type: %v", dType))
	}
}
