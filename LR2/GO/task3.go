package main

import (
	"bufio"
	"fmt"
	"io"
)

type Position struct{ a, b int }

func task3(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		fmt.Fprintln(out, "ОШИБКА: ожидалось n")
		return
	}
	if n <= 0 {
		fmt.Fprintln(out, "ОШИБКА: n > 0")
		return
	}
	var positions Array[Position]
	for i := 0; i < n; i++ {
		var p Position
		if _, err := fmt.Fscan(reader, &p.a, &p.b); err != nil {
			fmt.Fprintf(out, "ОШИБКА: ожидались a b для черепахи %d\n", i+1)
			return
		}
		if p.a < 0 || p.b < 0 || p.a > n-1 || p.b > n-1 || p.a+p.b != n-1 {
			continue
		}
		found := false
		for j := 0; j < positions.Size(); j++ {
			if positions.Get(j) == p {
				found = true
				break
			}
		}
		if !found {
			positions.Push(p)
		}
	}
	fmt.Fprintln(out, positions.Size())
}
