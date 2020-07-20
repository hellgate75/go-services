package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/hellgate75/go-services/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoConnection struct {
	Configuration database.DbConfig
	Client        *mongo.Client
	Context       *context.Context
	Valid         bool
	Cancel        context.CancelFunc
	err           error
}

func (conn *mongoConnection) Query(dbRef database.DataRef, fields []string, conditions []database.Condition, withAnd bool) (database.ResultSet, error) {
	if !conn.Valid || conn.Client == nil {
		return database.ResultSet{}, errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Query %v", r))
			conn.err = err
		}
	}()
	var resultSet = database.ResultSet{
		MetaData: database.MetaData{
			Columns:   make([]database.Column, 0),
			EntityRef: dbRef,
		},
		Records: make([]database.Result, 0),
		Lines:   int64(0),
	}
	if conn.Context == nil {
		err = errors.New(fmt.Sprint("Mongo Context unavailable"))
	} else {
		var cursor *mongo.Cursor
		var filter = make(bson.D, 0)
		for _, cond := range conditions {
			filter = append(filter, bson.E{
				Key:   cond.Field,
				Value: cond.Value.Value,
			})
		}
		cursor, err = conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).Find(*conn.Context, filter)
		records := 0
		recordSet := make([]database.Result, 0)
		for cursor.Next(*conn.Context) {
			raw := cursor.Current
			if err = cursor.Err(); err != nil {
				return database.ResultSet{}, err
			}
			values, err := raw.Values()
			if err == nil {
				records += 1
				res := database.Result{
					Document: raw,
					Columns:  int64(len(values)),
					Values:   make([]interface{}, 0),
				}
				for _, value := range values {
					val := convertRawValue(value)
					res.Values = append(res.Values, val)
				}
				recordSet = append(recordSet, res)
			} else {
				return database.ResultSet{}, err
			}
			resultSet.Lines = int64(len(recordSet))
			resultSet.Records = recordSet
			return resultSet, nil
		}
	}
	return resultSet, err
}

func convertRawValue(value bson.RawValue) interface{} {
	switch value.Type {
	case bsontype.Array:
		arrRaw := value.Array()
		vals, errV := arrRaw.Values()
		if errV != nil {
			return nil
		} else {
			var dat = make([]interface{}, 0)
			for _, val := range vals {
				dat = append(dat, convertRawValue(val))
			}
			return dat
		}
	case bsontype.Binary:
		_, d, ok := value.BinaryOK()
		if ok {
			return d
		} else {
			return nil
		}
	case bsontype.String:
		return value.String()
	case bsontype.Timestamp:
		t, _ := value.Timestamp()
		return t
	case bsontype.Int64:
		return value.Int64()
	case bsontype.Boolean:
		return value.Boolean()
	case bsontype.CodeWithScope:
		_, raw := value.CodeWithScope()
		vals, errV := raw.Values()
		if errV != nil {
			return nil
		} else {
			var dat = make([]interface{}, 0)
			for _, val := range vals {
				dat = append(dat, convertRawValue(val))
			}
			return dat
		}
	case bsontype.DateTime:
		return value.DateTime()
	case bsontype.DBPointer:
		//TODO Improve DBPointer type
		_, poid := value.DBPointer()
		return poid.String()
	case bsontype.Decimal128:
		//TODO Improve Decimal type
		return value.Decimal128().String()
	case bsontype.Double:
		return value.Double()
	case bsontype.EmbeddedDocument:
		raw := value.Document()
		vals, errV := raw.Values()
		if errV != nil {
			return nil
		} else {
			var dat = make([]interface{}, 0)
			for _, val := range vals {
				dat = append(dat, convertRawValue(val))
			}
			return dat
		}
	case bsontype.Int32:
		return value.Int32()
	case bsontype.JavaScript:
		return value.JavaScript()
	case bsontype.Null:
		return nil
	case bsontype.ObjectID:
		//TODO Improve type ObjectID
		return value.ObjectID()
	case bsontype.Regex:
		p, _ := value.Regex()
		return p
	case bsontype.Symbol:
		return value.Symbol()
	default:
		return nil
	}

}

func (conn *mongoConnection) Insert(dbRef database.DataRef, fields []database.Field, values []database.Value, withAnd bool) error {
	if !conn.Valid || conn.Client == nil {
		return errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Insert %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		var valMany = make([]interface{}, 0)
		for _, v := range values {
			valMany = append(valMany, v.Value)
		}
		_, err = conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).InsertMany(*conn.Context, valMany)
	}
	return err
}

func (conn *mongoConnection) Update(dbRef database.DataRef, conditions []database.Condition, fields []database.Field, values []database.Value, withAnd bool) (int64, error) {
	if !conn.Valid || conn.Client == nil {
		return 0, errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Update %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		var filter = make(bson.D, 0)
		for _, cond := range conditions {
			filter = append(filter, bson.E{
				Key:   cond.Field,
				Value: cond.Value.Value,
			})
		}
		var res *mongo.UpdateResult
		for _, v := range values {
			res, err = conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).UpdateMany(*conn.Context, filter, v.Value)
			if err != nil {
				return 0, err
			}
		}
		if err == nil {
			return res.ModifiedCount, nil
		}
	}
	return int64(len(values)), err
}

func (conn *mongoConnection) Delete(dbRef database.DataRef, conditions []database.Condition) (int64, error) {
	if !conn.Valid || conn.Client == nil {
		return 0, errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Delete %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		var filter = make(bson.D, 0)
		for _, cond := range conditions {
			filter = append(filter, bson.E{
				Key:   cond.Field,
				Value: cond.Value.Value,
			})
		}
		var res *mongo.DeleteResult
		res, err = conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).DeleteMany(*conn.Context, filter)
		if err == nil {
			return res.DeletedCount, nil
		}
	}
	return 0, err
}

func (conn *mongoConnection) Purge(dbRef database.DataRef) (int64, error) {
	if !conn.Valid || conn.Client == nil {
		return 0, errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Purge %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		var res *mongo.DeleteResult
		res, err = conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).DeleteMany(*conn.Context, bson.D{})
		if err == nil {
			return res.DeletedCount, nil
		}
	}
	return 0, err
}

func (conn *mongoConnection) Create(dbRef database.DataRef, fields []database.Field) error {
	if !conn.Valid || conn.Client == nil {
		return errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Create %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		name := conn.Client.Database(dbRef.Database).Name()
		db := conn.Client.Database(dbRef.Database)
		if db == nil {
			return errors.New("Errors retriving databse")
		} else {
			//			err = db.CreateCollection(*conn.Context, dbRef.Namespace)
			_ = db.Collection(dbRef.Namespace)
			if err != nil {
				collName := conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).Name()
				fmt.Printf("Created database: %s collection: %s\n", name, collName)
			}
		}
	}
	return err
}
func (conn *mongoConnection) CreateDb(dbRef database.DataRef) error {
	if !conn.Valid || conn.Client == nil {
		return errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::CreateDb %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		name := conn.Client.Database(dbRef.Database).Name()
		fmt.Printf("Created database: %s\n", name)
	}
	return err
}

func (conn *mongoConnection) Drop(dbRef database.DataRef) error {
	if !conn.Valid || conn.Client == nil {
		return errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::Drop %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		name := conn.Client.Database(dbRef.Database).Name()
		collName := conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).Name()
		err = conn.Client.Database(dbRef.Database).Collection(dbRef.Namespace).Drop(*conn.Context)
		if err == nil {
			fmt.Printf("Dropped database: %s collection: %s\n", name, collName)
		}
	}
	return err
}

func (conn *mongoConnection) DropDb(dbRef database.DataRef) error {
	if !conn.Valid || conn.Client == nil {
		return errors.New("Connection is closed or invalid")
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Connection::DropDb %v", r))
			conn.err = err
		}
	}()
	if conn.Context == nil {
		err = errors.New("Mongo Context unavailable")
	} else {
		name := conn.Client.Database(dbRef.Database).Name()
		err = conn.Client.Database(dbRef.Database).Drop(*conn.Context)
		if err == nil {
			fmt.Printf("Dropped database: %s\n", name)
		}
	}
	return err
}

func (conn *mongoConnection) GetLastError() error {
	return conn.err
}

func (conn *mongoConnection) Close() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Mongo-Driver::Disconnect %v", r))
		}
	}()
	if conn.Valid {
		if conn.Cancel != nil {
			conn.Cancel()
		}
		conn.Context = nil
		conn.Client = nil
		conn.Cancel = nil
	} else {
		err = errors.New("Connection is already closed")
	}
	return err
}
func (conn *mongoConnection) IsConnected() bool {
	return conn.Valid
}
