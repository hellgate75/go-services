package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/hellgate75/go-services/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mongoDriver struct {
}

func (md *mongoDriver) Connect(config database.DbConfig) (database.Connection, error) {
	var client *mongo.Client
	var err error
	defer func(){
		if r:=recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Driver::Connect %v", r))
		}
	}()
	if config.Name != "" && config.Password != "" {
		client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%v", config.Name, config.Password, config.Host, config.Port) ))
	} else {
		client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Host, config.Port) ))
	}
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	var conn = mongoConnection{
		Valid: (client != nil),
		Client: client,
		Context: &ctx,
		Cancel: cancel,
		Configuration: config,
	}
	return &conn, err
}

func GetMongoDriver() database.Driver{
	return &mongoDriver{
	}
}