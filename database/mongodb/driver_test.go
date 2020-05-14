package mongodb

import (
	"github.com/hellgate75/go-services/database"
	"testing"
)

type Data struct {
	Code	string				`json:"$code, omitempty"`
	Name	string 				`json:"name, omitempty"`
	Surname	string 				`json:"surname, omitempty"`
	Age 	int 				`json:"age, omitempty"`
	Role    string 				`json:"role, omitempty"`
}

func TestMongoDriver(t *testing.T) {
	driver := GetMongoDriver()
	conn, err := driver.Connect(database.DbConfig{
		Host: "localhost",
		Port: 27017,
		Name: "root",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("Connection error occured: %v", err)
	}
	config := database.DataRef{
		Database: "test",
		Namespace: "sample",
	}
	err = conn.CreateDb(config)
	if err != nil {
		t.Fatalf("Database creation error occured: %v", err)
	}
	err = conn.DropDb(config)
	if err != nil {
		t.Fatalf("Database drop error occured: %v\n", err)
	}
}
