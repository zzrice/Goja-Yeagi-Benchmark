package algorithm_test

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"go/build"
	"math/rand"
	"strconv"
	"testing"
)

//测试算法性能

func BenchmarkFib(b *testing.B) {

	// 斐波那契数列（动态规划）
	b.Run("GojaAlgorithm", func(b *testing.B) {
		script := `
		function fibonacci(n) {
			const fib = [0, 1];
			for (let i = 2; i <= n; i++) {
				fib[i] = fib[i - 1] + fib[i - 2];
			}
			return fib[n];
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
			x := rand.Intn(1000)
			_, err := vm.RunString(`fibonacci(` + strconv.Itoa(x) + `);`)
			if err != nil {
				b.Fatal(err)
			}
		}

	})

	b.Run("YaegiAlgorithm", func(b *testing.B) {
		fmt.Println("1")
		script := `
		package main
		import "fmt"
		func fibonacci(n int) int {
			if n <= 1 {
				return n
			}
			fib := make([]int, n+1)
			fib[0], fib[1] = 0, 1
			for i := 2; i <= n; i++ {
				fib[i] = fib[i-1] + fib[i-2]
			}
			fmt.Println("1")
			return fib[n]
		}`
		i := interp.New(interp.Options{GoPath: build.Default.GOPATH})
		_, err := i.Eval(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			x := rand.Intn(1000)
			_, err := i.Eval(`fibonacci(` + strconv.Itoa(x) + `)`)
			if err != nil {
				b.Fatal(err)
			}
		}

	})
}

func BenchmarkConvexHull(b *testing.B) {

	b.Run("GojaAlgorithm", func(b *testing.B) {
		script := `
		function convexHull(points) {
			// 按照 y 坐标排序，若 y 坐标相同则按 x 坐标排序
			points.sort((a, b) => a.y - b.y || a.x - b.x);
		
			// 计算凸包的下半部分
			const lower = [];
			for (const point of points) {
				while (lower.length >= 2 && cross(lower[lower.length - 2], lower[lower.length - 1], point) <= 0) {
					lower.pop();
				}
				lower.push(point);
			}
		
			// 计算凸包的上半部分
			const upper = [];
			for (let i = points.length - 1; i >= 0; i--) {
				const point = points[i];
				while (upper.length >= 2 && cross(upper[upper.length - 2], upper[upper.length - 1], point) <= 0) {
					upper.pop();
				}
				upper.push(point);
			}
		
			// 去掉重复的点
			upper.pop();
			lower.pop();
		
			return lower.concat(upper);
		}
		
		// 计算叉积
		function cross(o, a, b) {
			return (a.x - o.x) * (b.y - o.y) - (a.y - o.y) * (b.x - o.x);
		}
		
		// 示例使用
		const points = [
			{ x: 0, y: 0 },
			{ x: 1, y: 1 },
			{ x: 2, y: 2 },
			{ x: 1, y: 0 },
			{ x: 0, y: 1 },
			{ x: 2, y: 0 },
			{ x: 0, y: 2 },
			{ x: 2, y: 1 },
		];
	`
		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := vm.RunString(`convexHull(points);`)
			if err != nil {
				b.Fatal(err)
			}
		}

	})

	b.Run("YaegiAlgorithm", func(b *testing.B) {
		script := `
		package main
		// 定义一个点结构体
		type Point struct {
			x, y int
		}
		
		// 计算叉积
		func cross(o, a, b Point) int {
			return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
		}
		
		// 手动排序点
		func manualSort(points []Point) []Point {
			n := len(points)
			for i := 0; i < n; i++ {
				for j := 0; j < n-i-1; j++ {
					if points[j].x > points[j+1].x || (points[j].x == points[j+1].x && points[j].y > points[j+1].y) {
						points[j], points[j+1] = points[j+1], points[j]
					}
				}
			}
			return points
		}
		
		// 凸包计算
		func convexHull(points []Point) []Point {
			points = manualSort(points)
		
			// 创建下半部分
			lower := []Point{}
			for _, p := range points {
				for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) < 0 {
					lower = lower[:len(lower)-1]
				}
				lower = append(lower, p)
			}
		
			// 创建上半部分
			upper := []Point{}
			for i := len(points) - 1; i >= 0; i-- {
				p := points[i]
				for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) < 0 {
					upper = upper[:len(upper)-1]
				}
				upper = append(upper, p)
			}
		
			// 去掉上半部分的最后一个点，因为它与下半部分的第一个点重复
			return append(lower[:len(lower)-1], upper...)
		}
		
		const points = []Point{
			{0, 0},
			{1, 1},
			{2, 2},
			{1, 0},
			{0, 1},
			{2, 0},
			{0, 2},
			{2, 1},
		}
`
		i := interp.New(interp.Options{})
		_, err := i.Eval(script)
		i.Use(stdlib.Symbols)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.ReportAllocs()

		for j := 0; j < b.N; j++ {
			_, err := i.Eval(`convexHull(points)`)
			if err != nil {
				b.Fatal(err)
			}
		}

	})
}
