package error_test

import (
	"github.com/dop251/goja"
	"github.com/traefik/yaegi/interp"
	"testing"
)

func BenchmarkError(b *testing.B) {

	b.Run("GojaError", func(b *testing.B) {
		script := `
		function test() {
			try {
				throw new Error("Test error");
			} catch (e) {
			}
		}
	`
		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := vm.RunString("test()")
			if err != nil {
				continue
			}
		}
	})

	b.Run("YaegiError", func(b *testing.B) {
		script := `
		package main
		func test() {
			defer func() {
				if r := recover(); r != nil {}
			}()
			panic("Test error")
		}
`
		i := interp.New(interp.Options{})
		_, err := i.Eval(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err := i.Eval("test()")
			if err != nil {
				continue
			}
		}
	})

	// go原生
	b.Run("GoError", func(b *testing.B) {
		e := func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			panic("Test error")
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			e()
		}
	})
}
