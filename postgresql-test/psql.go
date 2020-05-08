package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stoewer/go-strcase"
)

// https://gobyexample.com/maps
var columnMap = map[string]string{
	"bool":              "bool",
	"char":              "byte",
	"smallint":          "uint16",
	"smallserial":       "uint16",
	"int":               "uint32",
	"serial":            "uint32",
	"oid":               "uint32",
	"bigint":            "uint64",
	"bigserial":         "uint64",
	"real":              "float32",
	"double precision":  "float64",
	"character varying": "string",
	"integer":           "int32",
	// "json": "serde_json::Value",
	// "jsonb": "serde_json::Value",
	"text":                        "string",
	"timestamp with time zone":    "time.Time",
	"timestamp without time zone": "time.Time",
	"uuid":                        "string",
}

// https://stackoverflow.com/questions/36688008/couldt-convert-nil-into-type
// https://golang.org/pkg/database/sql/#NullString
var nullableColumnMap = map[string]string{
	"bool":              "sql.NullBool",
	"char":              "sql.NullString",
	"smallint":          "sql.NullInt32",
	"smallserial":       "sql.NullInt32",
	"int":               "sql.NullInt32",
	"serial":            "sql.NullInt32",
	"oid":               "sql.NullInt32",
	"bigint":            "sql.NullInt64",
	"bigserial":         "sql.NullInt64",
	"real":              "sql.NullFloat64",
	"double precision":  "sql.NullFloat64",
	"character varying": "sql.NullString",
	"integer":           "sql.NullInt32",
	// "json": "serde_json::Value",
	// "jsonb": "serde_json::Value",
	"text":                        "sql.NullString",
	"timestamp with time zone":    "sql.NullTime",
	"timestamp without time zone": "sql.NullTime",
	"uuid":                        "sql.NullString",
}

func main() {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	tablenames, err := db.Query(
		"SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY tablename ASC;",
	)
	if err != nil {
		panic(err)
	}

	output := ""

	for tablenames.Next() {
		var tablename string
		err := tablenames.Scan(&tablename)
		if err != nil {
			panic(err)
		}

		output += "type " + strcase.UpperCamelCase(tablename) + "Row struct {\n"
		// fmt.Print(output)
		output += getTable(db, tablename)
		output += "}\n\n"
	}

	fmt.Print(output)
}

func getTable(db *sql.DB, tablename string) string {
	columns, err := db.Query(
		"SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = $1",
		tablename,
	)

	if err != nil {
		panic(err)
	}

	output := ""
	for columns.Next() {
		var column_name string
		var pg_type string
		var is_nullable string
		err := columns.Scan(&column_name, &pg_type, &is_nullable)
		if err != nil {
			panic(err)
		}

		go_type := columnMap[pg_type]
		null_go_type, exists := nullableColumnMap[pg_type]
		if is_nullable == "YES" && exists {
			go_type = null_go_type
		}
		output += fmt.Sprintf(
			"  %s %s `db:\"%s\"` // %s\n",
			strcase.UpperCamelCase(column_name),
			go_type,
			column_name,
			pg_type,
		)
	}

	return output
}
