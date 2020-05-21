package database

import "go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

// Operation enumeration type
type Operation byte

// DataType enumeration type
type DataType byte

// DriverType enumeration type
type DriverType byte

const (
	// Equality comparator Operation enumeration type
	Equals Operation = iota + 1
	// Less comparator Operation enumeration type
	LessThan
	// Grater comparator Operation enumeration type
	GraterThan
	// Similitude comparator Operation enumeration type
	Like
	// List items comparator Operation enumeration type
	In
	// Negation comparator Operation enumeration type
	Not Operation = 50
	// String data type DataType enumeration type
	StringType DataType = iota + 1
	// Numeric data type DataType enumeration type
	NumericType
	// Decimal data type DataType enumeration type
	DecimalType
	// List data type DataType enumeration type
	ListType
	// Structure data type DataType enumeration type
	StructType
	// MongoDb DriverType enumeration type
	MongoDbDriver DriverType = iota + 1
)

//Database Configuration structure
type DbConfig struct {
	// MongoDb User name field
	Name string
	// MongoDb Password field
	Password string
	// MongoDb Certificate file path field
	Certificate string
	// MongoDb Private Key file path field
	PrivateKey string
	// MongoDb Host name field
	Host string
	// MongoDb Port field
	Port int
}

// Field descriptor structure
type Field struct {
	// Field name
	Name string
	// Field type
	Type string
	// Field size
	Size int64
	// Field precision
	Precision int
}

// Value descriptor structure
type Value struct {
	// Value Type
	Type DataType
	// Value Content
	Value interface{}
}

// Condition descriptor structure
type Condition struct {
	// Condition Field name
	Field string
	// Condition operation type
	Operation Operation
	// Comparison Value or RegExp expression
	Value Value
}

// Command Single Row Result descriptor structure
type Result struct {
	// Number of Columns
	Columns int64
	// Column Values
	Values []interface{}
	// Record Value
	Document bsoncore.Document
}

// Result Set descriptor structure
type ResultSet struct {
	// Number of lines
	Lines int64
	// Results Records
	Records []Result
}

// Data Reference descriptor structure
type DataRef struct {
	// Database name
	Database string
	// Namespace value
	Namespace string
	// Field Set Reference
	FieldSetRef string
}

// Connection interface
type Connection interface {
	// Execute Query on the database instance
	Query(dbRef DataRef, fields []string, conditions []Condition) (ResultSet, error)
	// Insert record on the database instance
	Insert(dbRef DataRef, fields []Field, values []Value) error
	// Update one or more records on the database instance
	Update(dbRef DataRef, conditions []Condition, fields []Field, values []Value) (int64, error)
	// Delete one or more records on the database instance
	Delete(dbRef DataRef, conditions []Condition) (int64, error)
	// Purge one or more records on the database instance
	Purge(dbRef DataRef) (int64, error)
	// Create Namespace, Collection or Entity element on the database instance
	Create(dbRef DataRef, fields []Field) error
	// Create new database instance
	CreateDb(dbRef DataRef) error
	// Drop Namespace, Collection or Entity element on the database instance
	Drop(dbRef DataRef) error
	// Drop existing database instance
	DropDb(dbRef DataRef) error
	// Close connection
	Close() error
	// Is connection open
	IsConnected() bool
	// Get latest execution error
	GetLastError() error
}

// Driver interface
type Driver interface {
	// Connect to a database instance or database cluster instance
	Connect(config DbConfig) (Connection, error)
}
