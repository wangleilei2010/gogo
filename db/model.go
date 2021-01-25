package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/wangleilei2010/gogo/collection"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db            = &sql.DB{}
	globalConnStr string
)

const (
	ExcludedFieldName = "Table"
	ModelTagKeyName   = "db"
)

func SetGlobalConnStr(connStr string) {
	globalConnStr = connStr
}

func initDB() error {
	if globalConnStr != "" {
		var err error
		if db, err = sql.Open("mysql", globalConnStr); err != nil {
			log.Fatal("fail to connect database!")
			return err
		}
		db.SetConnMaxLifetime(500)
		db.SetMaxIdleConns(100)

		if err := db.Ping(); err != nil {
			log.Fatal("fail to ping database server!")
			return err
		}
		return nil
	} else {
		return errors.New("global connStr not set")
	}
}

func closeDB() error {
	err := db.Close()
	return err
}

type dbEngine struct {
}

func (e dbEngine) ExecSQL(sqlStmt string) bool {
	if err := initDB(); err != nil {
		return false
	}
	defer closeDB()

	if result, err := db.Exec(sqlStmt); err == nil {
		if _, err := result.RowsAffected(); err == nil {
			return true
		}
	}
	return false
}

func (e dbEngine) FetchAll(sqlStmt string, m iTable) collection.GoSlice {
	if err := initDB(); err != nil {
		return nil
	}
	defer closeDB()

	fExp := regexp.MustCompile(`(?i)(select|SELECT|from|FROM)`)
	matches := fExp.Split(sqlStmt, 3)
	if len(matches) > 2 {
		oriFields := strings.Split(matches[1], ",")
		if rows, err := db.Query(sqlStmt); err == nil {
			defer rows.Close()
			t := reflect.TypeOf(m)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			columns := len(oriFields)
			//dbFields := make([]string, 0, columns)
			modelFields := make([]modelField, 0, columns)

			for j := 0; j < columns; j++ {
				for i := 0; i < t.NumField(); i++ {
					if t.Field(i).Name != ExcludedFieldName &&
						strings.Contains(oriFields[j], t.Field(i).Tag.Get(ModelTagKeyName)) {
						//dbFields = append(dbFields, string(t.Field(i).Tag))
						modelFields = append(modelFields, modelField{
							Name: t.Field(i).Name,
							Type: t.Field(i).Type.Name(),
							Tag:  t.Field(i).Tag.Get(ModelTagKeyName)})
					}
				}
			}
			dbResults := make([]interface{}, 0)

			for rows.Next() {
				results := make([]interface{}, columns)
				data := make([]sql.NullString, columns)
				for i, _ := range results {
					results[i] = &data[i]
				}
				if err := rows.Scan(results...); err != nil {
					log.Fatal(err)
				}

				model := reflect.New(t)
				for i, v := range modelFields {
					switch v.Type {
					case "string":
						model.Elem().FieldByName(v.Name).SetString(data[i].String)
					case "int64", "int":
						intVal, _ := strconv.ParseInt(data[i].String, 10, 64)
						model.Elem().FieldByName(v.Name).SetInt(intVal)
					case "float64":
						floatVal, _ := strconv.ParseFloat(data[i].String, 64)
						model.Elem().FieldByName(v.Name).SetFloat(floatVal)
					default:
						model.Elem().FieldByName(v.Name).SetString(data[i].String)
					}
				}
				m := model.Interface()
				dbResults = append(dbResults, m)
			}
			return dbResults
		}
	}
	return nil
}

type Model struct {
	Table iTable
}

func (m Model) FetchAll(whereOrQueryStmt string) collection.GoSlice {
	e := dbEngine{}
	if strings.HasPrefix(whereOrQueryStmt, "SELECT") || strings.HasPrefix(whereOrQueryStmt, "select") {
		return e.FetchAll(whereOrQueryStmt, m.Table)
	} else {
		modelFields := getFields(m.Table)
		var dbFields = make([]string, 0)
		for _, m := range modelFields {
			dbFields = append(dbFields, m.Tag)
		}
		sqlStmt := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(dbFields, ","),
			m.Table.getTableName(), whereOrQueryStmt)
		return e.FetchAll(sqlStmt, m.Table)
	}
}

func (m Model) FetchOne(whereStmt string) interface{} {
	all := m.FetchAll(whereStmt)
	if len(all) > 0 {
		return all[0]
	} else {
		return nil
	}
}

func (m Model) ExecSQL(sql string) bool {
	e := dbEngine{}
	return e.ExecSQL(sql)
}

func NewDBModel(t iTable, tableName string) Model {
	t.setTableName(tableName)
	return Model{Table: t}
}

func getFields(m iTable) []modelField {
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	fieldNum := t.NumField()
	//dbFields := make([]string, 0, fieldNum)
	objFields := make([]modelField, 0, fieldNum)

	for i := 0; i < fieldNum; i++ {
		if t.Field(i).Name != ExcludedFieldName {
			//dbFields = append(dbFields, string(t.Field(i).Tag))
			objFields = append(objFields, modelField{
				Name: t.Field(i).Name,
				Type: t.Field(i).Type.Name(),
				Tag:  t.Field(i).Tag.Get(ModelTagKeyName)})
		}
	}
	return objFields
}

type modelField struct {
	Name string
	Type string
	Tag  string
}

type iTable interface {
	setTableName(tableName string)
	getTableName() string
}

type Table struct {
	tableName string
}

func (t *Table) setTableName(tableName string) {
	t.tableName = tableName
}

func (t Table) getTableName() string {
	return t.tableName
}
