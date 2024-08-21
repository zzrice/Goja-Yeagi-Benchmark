package interact_test

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"testing"
)

func createLargeData() map[string]interface{} {
	data := make(map[string]interface{})
	for i := 0; i < 10000; i++ {
		data[fmt.Sprintf("key%d", i)] = i
	}
	return data
}

func BenchmarkInteract(b *testing.B) {
	date := createLargeData()
	jsonData, _ := json.Marshal(date)

	b.Run("GojaInteract", func(b *testing.B) {
		script := `
		function parseJSON(data) {
			var data = JSON.parse(data);
		}`

		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			b.Fatal(err)
		}
		parseJSON, ok := goja.AssertFunction(vm.Get("parseJSON"))
		if !ok {
			b.Fatal(ok)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := parseJSON(goja.Undefined(), vm.ToValue(string(jsonData)))
			if err != nil {
				b.Fatal(err)
			}
		}

	})

	b.Run("YaegiInteract", func(b *testing.B) {
		script := `
		package main
		import "encoding/json"
		func parseJSON(data string) error {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(data), &result)
			if err != nil {
				return err
			}
			return nil
		}`

		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		_, err := i.Eval(script)
		v, err := i.Eval("main.parseJSON")
		if err != nil {
			b.Fatal(err)
		}

		parseJSON := v.Interface().(func(string) error)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err = parseJSON(string(jsonData))
			if err != nil {
				b.Fatal(err)
			}
		}

	})

}
