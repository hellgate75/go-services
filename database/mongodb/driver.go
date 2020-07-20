package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/hellgate75/go-services/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type mongoDriver struct {
}

func (md *mongoDriver) Connect(config database.DbConfig) (database.Connection, error) {
	var client *mongo.Client
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Driver::Connect %v", r))
		}
	}()
	if config.Name != "" && config.Password != "" {
		if strings.Index(config.Host, ":") > 0 {
			if strings.Index(config.Host, ",") > 0 {
				//Cluster compressed configuration
				var clusterURI = ""
				for _, srv := range strings.Split(config.Host, ",") {
					lst := strings.Split(srv, ":")
					xhost, xport := lst[0], lst[1]
					var uri = fmt.Sprintf("%s:%s@%s:%s", config.Name, config.Password, xhost, xport)
					var sep = ""
					if len(clusterURI) > 0 {
						sep = ","
					}
					clusterURI = fmt.Sprintf("%s%s%s", clusterURI, sep, uri)
				}
				clusterURI = fmt.Sprintf("mongodb://%s", clusterURI)
				client, err = mongo.NewClient(options.Client().ApplyURI(clusterURI))
			} else {
				//Single server compressed configuration
				lst := strings.Split(config.Host, ":")
				xhost, xport := lst[0], lst[1]
				var uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Name, config.Password, xhost, xport)
				client, err = mongo.NewClient(options.Client().ApplyURI(uri))
			}
		} else {
			//No compressed configuration
			client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%v", config.Name, config.Password, config.Host, config.Port)))
		}
	} else {
		if strings.Index(config.Host, ":") > 0 {
			if strings.Index(config.Host, ",") > 0 {
				//Cluster compressed configuration
				var clusterURI = ""
				for _, srv := range strings.Split(config.Host, ",") {
					lst := strings.Split(srv, ":")
					xhost, xport := lst[0], lst[1]
					var uri = fmt.Sprintf("%s:%s", xhost, xport)
					var sep = ""
					if len(clusterURI) > 0 {
						sep = ","
					}
					clusterURI = fmt.Sprintf("%s%s%s", clusterURI, sep, uri)
				}
				clusterURI = fmt.Sprintf("mongodb://%s", clusterURI)
				client, err = mongo.NewClient(options.Client().ApplyURI(clusterURI))
			} else {
				//Single server compressed configuration
				lst := strings.Split(config.Host, ":")
				xhost, xport := lst[0], lst[1]
				var uri = fmt.Sprintf("mongodb://%s:%s", xhost, xport)
				client, err = mongo.NewClient(options.Client().ApplyURI(uri))
			}
		} else {
			//No compressed configuration
			client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Host, config.Port)))
		}
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
		Valid:         client != nil,
		Client:        client,
		Context:       &ctx,
		Cancel:        cancel,
		Configuration: config,
	}
	return &conn, err
}

func GetMongoDriver() database.Driver {
	return &mongoDriver{}
}

func GetMongoDbConfig(host string, port int, username string, password string) database.DbConfig {
	return database.DbConfig{
		Host:     host,
		Port:     port,
		Name:     username,
		Password: password,
	}
}
