package database

import "go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

type Operation byte
type DataType byte
type DriverType byte

const (
	Equals			Operation = iota + 1
	LessThan
	GraterThan
	Like
	In
	Not				Operation = 50
	StringType		DataType = iota + 1
	NumericType
	DecimalType
	ListType
	StructType
	MongoDbDriver	DriverType = iota + 1
)


type DbConfig struct {
	Name		string
	Password	string
	Certificate	string
	PrivateKey	string
	Host		string
	Port		int
}

type Field struct {
	Name	string
	Type	string
	Size	int64
	Precision	int
}

type Value struct {
	Type	DataType
	Value	interface{}
}

type Condition struct {
	Field		string
	Operation	Operation
	Value		Value
}

type Result struct {
	Columns		int64
	Values		[]interface{}
	Document	bsoncore.Document
}

type ResultSet struct{
	Lines		int64
	Records		[]Result
}

type DataRef struct {
	Database	string
	Namespace	string
	FieldSetRef	string
}

type Connection interface {
	Query(dbRef DataRef, fields []string, conditions []Condition) (ResultSet, error)
	Insert(dbRef DataRef, fields  []Field, values[]Value) error
	Update(dbRef DataRef, conditions []Condition, fields  []Field, values[]Value) (int64, error)
	Delete(dbRef DataRef, conditions []Condition) (int64, error)
	Purge(dbRef DataRef) (int64, error)
	Create(dbRef DataRef, fields []Field) error
	CreateDb(dbRef DataRef) error
	Drop(dbRef DataRef) error
	DropDb(dbRef DataRef) error
	Close() error
	IsConnected() bool
	GetLastError()	error

}

type Driver interface {
	Connect(config DbConfig) (Connection, error)
}
