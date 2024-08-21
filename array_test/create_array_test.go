package array_test

import (
	"github.com/dop251/goja"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"testing"
)

// 测试创建数组
func BenchmarkCreateArray(b *testing.B) {

	b.Run("GojaCreateArray", func(b *testing.B) {
		script := `
		const arraySize = 100000;
		// 创建数组
		function createArray() {
			var arr = new Array(arraySize);
			for (let i = 0; i < arraySize; i++) {
				arr[i] = i;
			}
	   		return arr;
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
			_, err := vm.RunString(`createArray();`)
			if err != nil {
				b.Fatal(err)
			}
		}

	})

	b.Run("YaegiCreateArray", func(b *testing.B) {
		script := `
		package main
		const arraySize = 100000;
		func createArray() []int {
			arr := make([]int, arraySize)
			for i := 0; i < arraySize; i++ {
				arr[i] = i
			}
			return arr
		}`
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		_, err := i.Eval(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err := i.Eval(`createArray()`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// 测试遍历数组
func BenchmarkTraverseArray(b *testing.B) {

	b.Run("GojaTraverseArray", func(b *testing.B) {
		script := `
		const arraySize = 10000;
		// 创建数组
		function createArray() {
			return Array.from({ length: arraySize }, (_, i) => i);
		}
		const arr = createArray();
	`
		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := vm.RunString(`
			for (let j = 0; j < arr.length; j++) {
				arr[j];
			}
		`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("YaegiTraverseArray", func(b *testing.B) {
		script := `
		package main
		const arraySize = 10000;
		func createArray() []int {
			arr := make([]int, arraySize)
			for i := 0; i < arraySize; i++ {
				arr[i] = i
			}
			return arr
		}
		const arr = createArray()
		`

		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		_, err := i.Eval(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err := i.Eval(` for i := 0; i < len(arr); i++ { arr[i]; }`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// 测试过滤数组
func BenchmarkFilterArray(b *testing.B) {

	b.Run("GojaFilterArray", func(b *testing.B) {
		script := `
		const arraySize = 10000;
		// 创建数组
		function createArray() {
			return Array.from({ length: arraySize }, (_, i) => i);
		}
		const arr = createArray();
	`
		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := vm.RunString(`
			arr.filter((item) => item % 2 === 0);
		`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("YaegiFilterArray", func(b *testing.B) {
		script := `
		package main
		const arraySize = 10000;
		func createArray() []int {
			arr := make([]int, arraySize)
			for i := 0; i < arraySize; i++ {
				arr[i] = i
			}
			return arr
		}
		const arr = createArray()
		`

		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		_, err := i.Eval(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err := i.Eval(` for i := 0; i < len(arr); i++ { if arr[i] % 2 == 0 { arr[i] } }`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
