package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hellgate75/go-services/database"
	"reflect"
	"strings"
	"time"
)

type mySqlConnection struct {
	Configuration database.DbConfig
	DB            *sql.DB
	Context       *context.Context
	Valid         bool
	Cancel        context.CancelFunc
	err           error
}

func prepareCondition(cond database.Condition) string {
	fld := cond.Field
	operationSymbol := byte(cond.Operation)
	not := false
	if operationSymbol > byte(database.Not) {
		operationSymbol = operationSymbol - byte(database.Not)
		not = true
	}
	switch operationSymbol {
	case byte(database.LessThan):
		if not {
			return fld + " > ?"
		} else {
			return fld + " < ?"
		}
	case byte(database.LessThanEquals):
		if not {
			return fld + " >= ?"
		} else {
			return fld + " <= ?"
		}
	case byte(database.GraterThan):
		if not {
			return fld + " < ?"
		} else {
			return fld + " > ?"
		}
	case byte(database.GraterThanEquals):
		if not {
			return fld + " <= ?"
		} else {
			return fld + " >= ?"
		}
	case byte(database.Like):
		if not {
			return fld + " NOT LIKE ?"
		} else {
			return fld + " LIKE ?"
		}
	case byte(database.In):
		if not {
			return fld + " NOT IN (?)"
		} else {
			return fld + " IN (?)"
		}
	case byte(database.Null):
		if not {
			return fld + " IS NOT NULL"
		} else {
			return fld + " IS NULL"
		}
	default:
		if not {
			return fld + " <> ?"
		} else {
			return fld + " = ?"
		}
	}

}

func toMySqlTypeInstance(typeName string) (reflect.Type, interface{}) {
	switch strings.ToLower(typeName) {
	case "tinyint":
		v := byte(0)
		return reflect.TypeOf(v), v
	case "integer", "int", "smallint":
		v := int(0)
		return reflect.TypeOf(v), v
	case "mediumint":
		v := int32(0)
		return reflect.TypeOf(v), v
	case "bigint":
		v := int64(0)
		return reflect.TypeOf(v), v
	case "float":
		v := float64(0)
		return reflect.TypeOf(v), v
	case "double":
		v := float64(0)
		return reflect.TypeOf(v), v
	case "real":
		v := uint64(0)
		return reflect.TypeOf(v), v
	case "char", "varchar", "text":
		v := ""
		return reflect.TypeOf(v), v
	case "date", "time", "datetime", "timestamp":
		v := time.Now()
		return reflect.TypeOf(v), v
	case "binary", "varbinary", "blob", "enum", "set", "year",
		"longblob", "longtext":
		return reflect.TypeOf(""), sql.RawBytes{}
	default:
		return reflect.TypeOf(""), sql.RawBytes{}
	}
}

func (c *mySqlConnection) Query(dbRef database.DataRef, fields []string, conditions []database.Condition, withAnd bool) (database.ResultSet, error) {
	resultSet := database.ResultSet{
		Records: make([]database.Result, 0),
		Lines:   0,
		MetaData: database.MetaData{
			EntityRef: dbRef,
			Columns:   make([]database.Column, 0),
		},
	}
	if c.DB == nil {
		return resultSet, errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	var err error
	var rows *sql.Rows
	if dbRef.SQL != "" {
		rows, err = c.DB.Query(dbRef.SQL)
	} else {
		var selCols = ""
		for _, f := range fields {
			if selCols != "" {
				selCols += ", "
			}
			selCols += f
		}
		if selCols != "" {
			selCols = "*"
		}
		if len(conditions) == 0 {
			rows, err = c.DB.Query(fmt.Sprintf("SELECT %s FROM %s", selCols, dbRef.Namespace))
		} else {
			var where string
			var values = make([]interface{}, 0)
			for _, c := range conditions {
				if where != "" {
					if withAnd {
						where += " AND "
					} else {
						where += " OR "
					}
				}
				cond := prepareCondition(c)
				if byte(c.Operation) != byte(database.Null) &&
					byte(c.Operation) != byte(database.Not)+byte(database.Null) {
					values = append(values, c.Value)
				}
				where += cond
			}
			if where != "" {
				where = " WHERE " + where
			}
			stmt, err := c.DB.Prepare(fmt.Sprintf("SELECT %s FROM %s%s", selCols, dbRef.Namespace, where))
			if err != nil {
				return resultSet, err
			}
			defer func() {
				_ = stmt.Close()
			}()
			rows, err = stmt.Query(values...)
		}
	}
	if err != nil {
		return resultSet, err
	}
	defer func() {
		_ = rows.Close()
	}()
	cols, _ := rows.Columns()
	colTypes, _ := rows.ColumnTypes()
	var values = make([]interface{}, len(cols))
	var queryArgs = make([]interface{}, len(cols))
	for i, colName := range cols {
		cType := colTypes[i]
		length, _ := cType.Length()
		precision, scale, _ := cType.DecimalSize()
		dataType := cType.Name()
		goType, defValue := toMySqlTypeInstance(cType.Name())
		values[i] = defValue
		queryArgs[i] = &values[i]
		resultSet.MetaData.Columns =
			append(resultSet.MetaData.Columns,
				database.Column{
					Name:      colName,
					Length:    length,
					Precision: precision,
					Scale:     scale,
					Type:      database.DataType(dataType),
					GoType:    goType,
				})
	}
	resultSet.Lines = 0
	for rows.Next() {
		err = rows.Scan(queryArgs)
		if err != nil {
			var resultValues = make([]interface{}, len(values))
			for i, _ := range values {
				resultValues[i] = values[i]
				if reflect.TypeOf(values[i]).Kind().String() == "sql.RawBytes" {
					resultValues[i] = string(values[i].(sql.RawBytes))
				}
			}
			resultSet.Lines++
			resultSet.Records = append(
				resultSet.Records, database.Result{
					Columns:  int64(len(resultSet.MetaData.Columns)),
					Document: nil,
					Values:   resultValues,
				})
		}
	}
	return resultSet, err
}

func (c *mySqlConnection) Insert(dbRef database.DataRef, fields []database.Field, values []database.Value) error {
	if c.DB == nil {
		return errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	if len(fields) != len(values) {
		return errors.New(fmt.Sprintf("Columns and values must have same length: %v <> %v", len(fields), len(values)))
	}
	if len(fields) == 0 || len(values) == 0 {
		return errors.New(fmt.Sprint("Insert statement needs list of Columns and Values of same length"))
	}
	cols := ""
	colValues := ""
	for _, f := range fields {
		if cols != "" {
			cols += ", "
		}
		cols += f.Name
		if colValues != "" {
			colValues += ", "
		}
		colValues += "?"
	}
	if cols != "" {
		cols = "(" + cols + ")"
	} else {
		return errors.New(fmt.Sprintf("Need column fields to create the insert"))
	}
	if colValues != "" {
		colValues = "VALUES(" + colValues + ")"
	} else {
		return errors.New(fmt.Sprintf("Need column fields to create the insert"))
	}
	sqlText := fmt.Sprintf("INSERT INTO %s%s%s"+dbRef.Namespace, cols, colValues)
	prep, err := c.DB.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer func() {
		_ = prep.Close()
	}()
	var sqlValues = make([]interface{}, 0)
	for _, val := range values {
		sqlValues = append(sqlValues, val.Value)
	}
	_, err = prep.Exec(sqlValues...)
	if err != nil {
		return err
	}
	return nil
}

func (c *mySqlConnection) Update(dbRef database.DataRef, conditions []database.Condition, fields []database.Field, values []database.Value, withAnd bool) (int64, error) {
	var records int64
	if c.DB == nil {
		return records, errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	if len(fields) != len(values) {
		return records, errors.New(fmt.Sprintf("Columns and values must have same length: %v <> %v", len(fields), len(values)))
	}
	if len(fields) == 0 || len(values) == 0 {
		return records, errors.New(fmt.Sprint("Insert statement needs list of Columns and Values of same length"))
	}
	cols := ""
	for _, f := range fields {
		if cols != "" {
			cols += ", "
		}
		cols += f.Name + " = ?"
	}
	if cols != "" {
		cols = " SET " + cols + ""
	} else {
		return records, errors.New(fmt.Sprintf("Need column fields to create the insert"))
	}
	var where string
	var whereValues = make([]interface{}, 0)
	for _, c := range conditions {
		if where != "" {
			if withAnd {
				where += " AND "
			} else {
				where += " OR "
			}
		}
		cond := prepareCondition(c)
		if byte(c.Operation) != byte(database.Null) &&
			byte(c.Operation) != byte(database.Not)+byte(database.Null) {
			whereValues = append(whereValues, c.Value)
		}
		where += cond
	}
	if where != "" {
		where = " WHERE " + where
	}
	sqlText := fmt.Sprintf("UPDATE %s%s%s"+dbRef.Namespace, cols, where)
	prep, err := c.DB.Prepare(sqlText)
	if err != nil {
		return records, err
	}
	defer func() {
		_ = prep.Close()
	}()
	var sqlValues = make([]interface{}, 0)
	for _, val := range values {
		sqlValues = append(sqlValues, val.Value)
	}
	sqlValues = append(sqlValues, whereValues...)
	r, err := prep.Exec(sqlValues...)
	if err != nil {
		return records, err
	}
	records, err = r.RowsAffected()
	return records, err
}

func (c *mySqlConnection) Delete(dbRef database.DataRef, conditions []database.Condition, withAnd bool) (int64, error) {
	var records int64
	if c.DB == nil {
		return records, errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	var where string
	var whereValues = make([]interface{}, 0)
	for _, c := range conditions {
		if where != "" {
			if withAnd {
				where += " AND "
			} else {
				where += " OR "
			}
		}
		cond := prepareCondition(c)
		if byte(c.Operation) != byte(database.Null) &&
			byte(c.Operation) != byte(database.Not)+byte(database.Null) {
			whereValues = append(whereValues, c.Value)
		}
		where += cond
	}
	if where != "" {
		where = " WHERE " + where
	}
	sqlText := fmt.Sprintf("DELETE FROM %s%s"+dbRef.Namespace, where)
	prep, err := c.DB.Prepare(sqlText)
	if err != nil {
		return records, err
	}
	defer func() {
		_ = prep.Close()
	}()
	var sqlValues = make([]interface{}, 0)
	sqlValues = append(sqlValues, whereValues...)
	r, err := prep.Exec(sqlValues...)
	if err != nil {
		return records, err
	}
	records, err = r.RowsAffected()
	return records, err
}

func (c *mySqlConnection) dropTable(name string) (int64, error) {
	if c.DB == nil {
		return 0, errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	_, err := c.DB.Query(fmt.Sprintf("DROP TABLE %s CASCADE", name))
	if err != nil {
		return 0, err
	}
	return 1, nil

}

func (c *mySqlConnection) truncateTable(name string) (int64, error) {
	if c.DB == nil {
		return 0, errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	_, err := c.DB.Query(fmt.Sprintf("TRUNCATE TABLE %s", name))
	if err != nil {
		return 0, err
	}
	return 1, nil

}

func (c *mySqlConnection) Purge(dbRef database.DataRef) (int64, error) {
	var err error
	var count int64
	if c.DB == nil {
		return count, errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	if dbRef.Namespace != "" {
		return c.truncateTable(dbRef.Namespace)
	} else if dbRef.Database != "" {
		rows, err := c.DB.Query(fmt.Sprint("show tables"))
		if err != nil {
			return count, err
		}
		var tableName string
		var scanArgs = make([]interface{}, 1)
		scanArgs[0] = &tableName
		for rows.Next() {
			err = rows.Scan(scanArgs)
			if err == nil {
				fmt.Println("Truncating table:", tableName)
				var n int64
				n, err = c.truncateTable(tableName)
				count += n
			}
		}
	} else {
		return count, errors.New(fmt.Sprint("Please choose truncate entity between Namespace for Table and Database for all Tables"))
	}
	if err != nil {
		return count, err
	}
	return 1, nil
}

func (c *mySqlConnection) Create(dbRef database.DataRef, fields []database.Field) error {
	if c.DB == nil {
		return errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	var err error
	if dbRef.Namespace != "" {
		//Create table
		//TODO: Implement MySql create table task
		return errors.New(fmt.Sprint("Create table not implemented yet"))
	} else if dbRef.FieldSetRef != "" {
		//Create table
		_, err = c.DB.Query(fmt.Sprintf("CREATE TABLESPACE %s", dbRef.Schema))
	} else if dbRef.Database != "" {
		err = c.CreateDb(dbRef)
	} else if dbRef.Schema != "" {
		err = c.CreateDb(dbRef)
	}
	return err
}

func (c *mySqlConnection) CreateDb(dbRef database.DataRef) error {
	var err error
	if c.DB == nil {
		return errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	if dbRef.Database != "" {
		_, err = c.DB.Query(fmt.Sprintf("CREATE DATABASE %s", dbRef.Database))
	} else if dbRef.Schema != "" {
		_, err = c.DB.Query(fmt.Sprintf("CREATE SCHEMA %s", dbRef.Schema))
	}
	return err
}

func (c *mySqlConnection) Drop(dbRef database.DataRef) error {
	var err error
	if c.DB == nil {
		return errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	if dbRef.Namespace != "" {
		_, err = c.dropTable(dbRef.Namespace)
	} else if dbRef.Database != "" {
		return c.DropDb(dbRef)
	} else if dbRef.FieldSetRef != "" {
		_, err = c.DB.Query(fmt.Sprintf("DROP TABLESPACE %s", dbRef.Namespace))
	} else if dbRef.Schema != "" {
		_, err = c.DB.Query(fmt.Sprintf("DROP SCHEMA %s", dbRef.Schema))
	} else {
		return errors.New(fmt.Sprint("Please choose drop entity between Namespace for Table, FieldSet for Tablespace and Database for all Tables"))
	}
	return err
}

func (c *mySqlConnection) DropDb(dbRef database.DataRef) error {
	if c.DB == nil {
		return errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	_, err := c.DB.Query(fmt.Sprintf("DROP DATABASE %s", dbRef.Database))
	return err
}

func (c *mySqlConnection) Close() error {
	if !c.IsConnected() {
		return errors.New(fmt.Sprint("Database connection is already closed"))
	}
	if c.DB == nil {
		return errors.New(fmt.Sprint("Database is closed,please reconnect before any operation"))
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		c.DB = nil
	}()
	err = c.DB.Close()
	return err
}

func (c *mySqlConnection) IsConnected() bool {
	return c.DB != nil
}

func (c *mySqlConnection) GetLastError() error {
	return nil
}
