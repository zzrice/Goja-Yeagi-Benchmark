package sql_test

import (
	"database/sql"
	"github.com/dop251/goja"
	_ "github.com/go-sql-driver/mysql"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
	"testing"
)

var (
	DB  *sql.DB
	err error
)

func initDB() {
	dsn := "root:root@tcp(127.0.0.1:3306)/fe2"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}

func SqlSelect(id int) string {
	initDB()
	defer DB.Close()

	sqlStr := "select status from fe_user where id=?;"
	rowObj := DB.QueryRow(sqlStr, id)

	var status string
	err := rowObj.Scan(&status)
	if err != nil {
		return "" // 返回空字符串表示错误
	}
	return status
}

func BenchmarkSqlSelect(b *testing.B) {
	b.Run("GojaSqlSelect", func(b *testing.B) {
		vm := goja.New()
		SqlSelect := func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(SqlSelect(int(call.Argument(0).ToInteger())))
		}
		vm.Set("SqlSelect", SqlSelect)
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err := vm.RunString(`SqlSelect(27)`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("YaegiSqlSelect", func(b *testing.B) {
		script := `
		package main
		import (
			"example.com/sql_test"
		)
		func test(id int) string {
			return sql_test.SqlSelect(id)
		}
		`
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		i.Use(map[string]map[string]reflect.Value{
			"example.com/sql_test/sql_test": {
				"SqlSelect": reflect.ValueOf(SqlSelect),
			},
		},
		)
		_, err := i.Eval(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err = i.Eval(`test(27)`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// go原生
	b.Run("GoSqlSelect", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			SqlSelect(27)
		}
	})
}
