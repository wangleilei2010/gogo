package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/wangleilei2010/gogo/collection"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ExcludedFieldName = "Table"
	TableName         = "TableName"
	ModelTagKeyName   = "db"
)

var tableScanPaths = make([]string, 0)
var modelTableNameMapping = make(map[string]string)

type ConnPool struct {
	db      *sql.DB
	connStr string
}

func OpenPool(connStr string, tbScanPaths ...string) (pool *ConnPool, err error) {
	if len(tbScanPaths) > 0 {
		tableScanPaths = tbScanPaths
		scanTableNames()
	}
	pool = &ConnPool{db: &sql.DB{}, connStr: connStr}
	err = pool.open()
	return
}

func (pool *ConnPool) open() error {
	var err error
	if pool.connStr == "" {
		return errors.New("DB connStr not set")
	}
	if pool.db, err = sql.Open("mysql", pool.connStr); err != nil {
		return err
	} else {
		pool.db.SetConnMaxLifetime(2000)
		pool.db.SetMaxIdleConns(10)
		err = pool.db.Ping()
		return err
	}
}

func (pool *ConnPool) Close() error {
	err := pool.db.Close()
	return err
}

func (pool *ConnPool) Count(sqlStmt string) int {
	rowNum := 0
	if rows, err := pool.db.Query(sqlStmt); err == nil {
		for rows.Next() {
			rowNum++
		}
		return rowNum
	}
	return -1
}

// Exec 支持数据库增/删/改
func (pool *ConnPool) Exec(query string, args ...any) (n int64, err error) {
	// example: pool.db.Exec("INSERT test SET name=?,age =?", "xiaowei", 18)
	var result sql.Result
	if result, err = pool.db.Exec(query, args...); err != nil {
		return
	} else {
		n, err = result.RowsAffected()
		return
	}
}

// FetchAll 支持查询多条数据
func FetchAll[M iTable](pool *ConnPool, whereOrQueryStmt string, args ...any) (collection.Slice[M], error) {
	var err error
	var rows *sql.Rows
	s := make(collection.Slice[M], 0)
	var q string
	t := reflect.TypeOf(new(M))
	if stmt := strings.ToUpper(whereOrQueryStmt); strings.HasPrefix(stmt, "SELECT") {
		q = whereOrQueryStmt
	} else {
		tn, modelFields := getFields(t)
		var dbFields = make([]string, 0)
		for _, mf := range modelFields {
			dbFields = append(dbFields, mf.Tag)
		}
		q = fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(dbFields, ","),
			tn, whereOrQueryStmt)
	}
	fExp := regexp.MustCompile(`(?i)(select|SELECT|from|FROM)`)
	matches := fExp.Split(q, 3)
	if len(matches) > 2 {
		if rows, err = pool.db.Query(q, args...); err == nil {
			defer rows.Close()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			oriFields, _ := rows.Columns()
			columns := len(oriFields)
			modelFields := make([]modelField, 0, columns)
			startPos := 0
			//LOOP:
			for j := 0; j < columns; j++ {
				for i := startPos; i < t.NumField(); i++ {
					if t.Field(i).Name != ExcludedFieldName && t.Field(i).Name != TableName &&
						strings.Contains(oriFields[j], t.Field(i).Tag.Get(ModelTagKeyName)) {
						modelFields = append(modelFields, modelField{
							Name: t.Field(i).Name,
							Type: t.Field(i).Type.Name(),
							Tag:  t.Field(i).Tag.Get(ModelTagKeyName),
							Pos:  j})
						//startPos = i + 1
						//continue LOOP
					}
				}
			}

			for rows.Next() {
				results := make([]interface{}, columns)
				data := make([]sql.NullString, columns)
				for i, _ := range results {
					results[i] = &data[i]
				}
				if err = rows.Scan(results...); err != nil {
					log.Println(err)
				}

				model := reflect.New(t)
				for _, v := range modelFields {
					switch v.Type {
					case "string":
						model.Elem().FieldByName(v.Name).SetString(data[v.Pos].String)
					case "int64", "int":
						intVal, _ := strconv.ParseInt(data[v.Pos].String, 10, 64)
						model.Elem().FieldByName(v.Name).SetInt(intVal)
					case "float64":
						floatVal, _ := strconv.ParseFloat(data[v.Pos].String, 64)
						model.Elem().FieldByName(v.Name).SetFloat(floatVal)
					default:
						model.Elem().FieldByName(v.Name).SetString(data[v.Pos].String)
					}
				}
				model.Elem().FieldByName(ExcludedFieldName).Set(reflect.ValueOf(Table{
					querySetIsNotNull: true,
				}))
				m := model.Interface()
				s.Push(*m.(*M))
			}
		}
	}
	return s, err
}

// FetchOne 查询单条数据
func FetchOne[M iTable](pool *ConnPool, whereOrQueryStmt string, args ...any) (M, error) {
	if records, err := FetchAll[M](pool, whereOrQueryStmt, args...); err != nil {
		return *new(M), err
	} else {
		if getLen[M](records) > 0 {
			return records[0], nil
		} else {
			return *new(M), nil
		}
	}
}

func getLen[M iTable](rs []M) int {
	return len(rs)
}

func getFields(t reflect.Type) (string, []modelField) {
	var tableName string
	modelName := getModelName(t)
	if tn, ok := modelTableNameMapping[modelName]; ok {
		tableName = tn
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", nil
	}
	fieldNum := t.NumField()
	objFields := make([]modelField, 0, fieldNum)

	for i := 0; i < fieldNum; i++ {
		if t.Field(i).Name != ExcludedFieldName && t.Field(i).Name != TableName {
			objFields = append(objFields, modelField{
				Name: t.Field(i).Name,
				Type: t.Field(i).Type.Name(),
				Tag:  t.Field(i).Tag.Get(ModelTagKeyName)})
		}
		if tableName == "" && t.Field(i).Name == TableName {
			tableName = t.Field(i).Tag.Get(ModelTagKeyName)
		}
	}
	return tableName, objFields
}

type modelField struct {
	Name string
	Type string
	Tag  string
	Pos  int
}

type iTable interface {
	IsDataNotNull() bool
}

type Table struct {
	querySetIsNotNull bool
}

func (t Table) IsDataNotNull() bool {
	return t.querySetIsNotNull
}

func getModelName(t reflect.Type) string {
	fullName := t.String()
	splitFuncName := strings.Split(fullName, ".")
	return splitFuncName[len(splitFuncName)-1]
}

func scanSingleGoFile(fileName string) {
	fSet := token.NewFileSet()
	parsedAst, _ := parser.ParseFile(fSet, fileName, nil, parser.ParseComments)

	pkg := &ast.Package{
		Name:  "Any",
		Files: make(map[string]*ast.File),
	}
	pkg.Files[fileName] = parsedAst

	importPath, _ := filepath.Abs("/")
	myDoc := doc.New(pkg, importPath, doc.AllDecls)
	for _, t := range myDoc.Types {
		if strings.Contains(t.Doc, "@Table:") {
			modelTableNameMapping[t.Name] = handleModelDoc(t.Doc)
		}
	}
}

func handleModelDoc(d string) string {
	s := strings.TrimSpace(d)
	return strings.ReplaceAll(s, "@Table:", "")
}

func scanTableNames() {
	var scanPath, _ = os.Getwd()
	var paths []string
	if len(tableScanPaths) != 0 {
		paths = make([]string, 0)
		for _, sp := range tableScanPaths {
			paths = append(paths, filepath.Join(scanPath, sp))
		}
	} else {
		paths = []string{scanPath}
	}

	fmt.Println("tableName scan paths:", paths)
	var files = make([]string, 0)

	for _, p := range paths {
		if err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			fmt.Printf("scan files err: %v", err)
		} else {
			for _, f := range files {
				scanSingleGoFile(f)
			}
		}
	}
}
