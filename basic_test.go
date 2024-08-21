package main

import (
	"github.com/dop251/goja"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	b.Run("GojaAdd", func(b *testing.B) {
		script := `function add(x, y) { return x + y; };`
		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			x, y := rand.Intn(10000), rand.Intn(10000)
			_, err := vm.RunString(`add(` + strconv.Itoa(x) + `, ` + strconv.Itoa(y) + `);`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("YaegiAdd", func(b *testing.B) {
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		_, err := i.Eval(`package main
		func add(x, y int) int { return x + y}`)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			x, y := rand.Intn(10000), rand.Intn(10000)
			_, err := i.Eval(`add(` + strconv.Itoa(x) + `, ` + strconv.Itoa(y) + `)`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
