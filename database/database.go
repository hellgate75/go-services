package database

import (
	"reflect"
	"strings"
)

// Operation enumeration type
type Operation byte

// DataType enumeration type
type DataType string

// DriverType enumeration type
type DriverType byte

const (
	// Equality comparator Operation enumeration type
	Equals Operation = iota + 1
	// Less comparator Operation enumeration type
	LessThan
	// Less or Equals comparator Operation enumeration type
	LessThanEquals
	// Grater comparator Operation enumeration type
	GraterThan
	// Grater or Equals comparator Operation enumeration type
	GraterThanEquals
	// Similitude comparator Operation enumeration type
	Like
	// List items comparator Operation enumeration type
	In
	// En item is null
	Null
	// Negation comparator Operation enumeration type
	Not Operation = 50
	// MongoDb DriverType enumeration type
	MongoDbDriver DriverType = iota + 1
	// MongoDb DriverType enumeration type
	MySQLDriver
)

func DriverToType(driver string) DriverType {
	switch strings.ToLower(driver) {
	case "mongo", "mongo-db", "mongodb":
		return MongoDbDriver
	case "mysql":
		return MySQLDriver
	default:
		return DriverType(0)
	}
}

//Database Configuration structure
type DbConfig struct {
	// Database name
	Driver string `json:"driver,omitempty" yaml:"driver,omitempty" xml:"driver,omitempty"`
	// Database name
	Database DataRef `json:"dbRef,omitempty" yaml:"dbRef,omitempty" xml:"db-ref,omitempty"`
	// Database url
	Url string `json:"connectUrl,omitempty" yaml:"connectUrl,omitempty" xml:"connect-url,omitempty"`
	// Database User name field
	Name string `json:"userName,omitempty" yaml:"userName,omitempty" xml:"user-name,omitempty"`
	// Database Password field
	Password string `json:"userPassword,omitempty" yaml:"userPassword,omitempty" xml:"user-password,omitempty"`
	// Database Password field
	DbPassword string `json:"dbPassword,omitempty" yaml:"dbPassword,omitempty" xml:"db-password,omitempty"`
	// MongoDb Certificate file path field
	Certificate string `json:"certificateFile,omitempty" yaml:"certificateFile,omitempty" xml:"certificate-file,omitempty"`
	// MongoDb Private Key file path field
	PrivateKey string `json:"privateKey,omitempty" yaml:"privateKey,omitempty" xml:"private-key,omitempty"`
	// MongoDb Host name field
	Host string `json:"hostname,omitempty" yaml:"hostname,omitempty" xml:"hostname,omitempty"`
	// MongoDb Port field
	Port int `json:"port,omitempty" yaml:"port,omitempty" xml:"port,omitempty"`
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

// Reference to a single ResultSet / Data Entity column
type Column struct {
	// Column name
	Name string
	// Column type
	Type DataType
	// Column type
	GoType reflect.Type
	// Columns size
	Length int64
	// Column numeric precision
	Precision int64
	// Column numeric precision
	Scale int64
}

// Reference to ResultSet / Data Entity structure
type MetaData struct {
	// Name of table or data entity
	EntityRef DataRef
	// ResultSet / Data Entity Columns
	Columns []Column
}

// Command Single Row Result descriptor structure
type Result struct {
	// Number of Columns
	Columns int64
	// Column Values
	Values []interface{}
	// Record Raw Value Value
	Document interface{}
}

// Result Set descriptor structure
type ResultSet struct {
	// Number of lines
	Lines int64
	// ResultSet / Data Entity metadata
	MetaData MetaData
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
	// Schema name
	Schema string
	// Data Query / Table SQL Reference
	SQL string
}

// Connection interface
type Connection interface {
	// Execute Query on the database instance
	Query(dbRef DataRef, fields []string, conditions []Condition, withAnd bool) (ResultSet, error)
	// Insert record on the database instance
	Insert(dbRef DataRef, fields []Field, values []Value) error
	// Update one or more records on the database instance
	Update(dbRef DataRef, conditions []Condition, fields []Field, values []Value, withAnd bool) (int64, error)
	// Delete one or more records on the database instance
	Delete(dbRef DataRef, conditions []Condition, withAnd bool) (int64, error)
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

