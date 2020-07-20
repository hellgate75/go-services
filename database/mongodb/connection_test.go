package mongodb

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hellgate75/go-services/database"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestMongoConnection(t *testing.T) {
	driver := GetMongoDriver()
	conn, err := driver.Connect(database.DbConfig{
		Host:     "localhost",
		Port:     27017,
		Name:     "root",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("Connection error occured: %v", err)
	}
	config := database.DataRef{
		Database:  "test",
		Namespace: "sample",
	}
	err = conn.CreateDb(config)
	if err != nil {
		t.Fatalf("Database creation error occured: %v", err)
	}
	err = conn.Create(config, []database.Field{})
	if err != nil {
		t.Fatalf("Database collection creation error occured: %v", err)
	}
	var values = make([]database.Value, 0)
	values = append(values, database.Value{
		Type: "struct",
		Value: Data{
			Code:    uuid.New().String(),
			Name:    "Fabrizio",
			Surname: "Torelli",
			Age:     45,
			Role:    "System Architect",
		},
	},
		database.Value{
			Type: "struct",
			Value: Data{
				Code:    uuid.New().String(),
				Name:    "Francesco",
				Surname: "Torelli",
				Age:     42,
				Role:    "Software Developer",
			},
		})
	err = conn.Insert(config, []database.Field{}, values)
	if err != nil {
		t.Fatalf("Database collection creation error occured: %v", err)
	}
	conditions := make([]database.Condition, 0)
	conditions = append(conditions, database.Condition{
		Field: "name",
		Value: database.Value{
			Type:  "struct",
			Value: "Francesco",
		},
	})
	conditions = append(conditions, database.Condition{
		Field: "surname",
		Value: database.Value{
			Type:  "struct",
			Value: "Torelli",
		},
	})
	var rs database.ResultSet
	rs, err = conn.Query(config, []string{}, conditions, false)
	if err != nil {
		t.Fatalf("Database collection querying error occured: %v\n", err)
	}
	if len(rs.Records) == 0 {
		t.Fatal("Database collection querying: no results")
	}
	if len(rs.Records) > 1 {
		t.Fatal("Database collection querying: too many results")
	}
	var record = rs.Records[0]
	fmt.Printf("Value: %v\n", record)
	fmt.Printf("Value set: %v\n", record.Values)
	var count int64
	var newValue = database.Value{
		Type: "struct",
		Value: bson.D{{
			"$set",
			bson.D{{
				"role",
				"Hardware Specialist",
			}},
		}},
	}
	count, err = conn.Update(config, conditions, []database.Field{}, []database.Value{newValue}, true)
	fmt.Printf("Update: %v\n", count)
	if err != nil {
		t.Fatalf("Database collection update error occured: %v\n", err)
	}
	if count != 1 {
		t.Fatalf("Wrong number of updated records %v\n", count)
	}
	conditions = make([]database.Condition, 0)
	conditions = append(conditions, database.Condition{
		Field: "name",
		Value: database.Value{
			Type:  "struct",
			Value: "Fabrizio",
		},
	})
	conditions = append(conditions, database.Condition{
		Field: "surname",
		Value: database.Value{
			Type:  "struct",
			Value: "Torelli",
		},
	})
	count, err = conn.Delete(config, conditions, true)
	fmt.Printf("Delete: %v\n", count)
	if err != nil {
		t.Fatalf("Database collection delete record error occured: %v\n", err)
	}
	if count != 1 {
		t.Fatalf("Wrong number of deleted records %v\n", count)
	}
	count, err = conn.Purge(config)
	fmt.Printf("Purge: %v\n", count)
	if err != nil {
		t.Fatalf("Database collection purge error occured: %v\n", err)
	}
	if count != 1 {
		t.Fatalf("Wrong number of purged records %v\n", count)
	}
	err = conn.Drop(config)
	if err != nil {
		t.Fatalf("Database collection drop error occured: %v\n", err)
	}
}
