package main

import (
	"bufio"
	"fmt"
	"io"
)

func task1(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		fmt.Fprintln(out, "ОШИБКА: ожидалось n m")
		return
	}
	if n <= 0 || m < 0 {
		fmt.Fprintln(out, "ОШИБКА: n > 0, m >= 0")
		return
	}

	graph := make([]ForwardList[int], n)
	for i := 0; i < m; i++ {
		var a, b int
		if _, err := fmt.Fscan(reader, &a, &b); err != nil {
			fmt.Fprintln(out, "ОШИБКА: ожидалось ребро a b")
			return
		}
		if a < 1 || a > n || b < 1 || b > n {
			fmt.Fprintf(out, "ОШИБКА: вершина вне [1..%d]\n", n)
			return
		}
		a--
		b--
		graph[a].PushFront(b)
		graph[b].PushFront(a)
	}

	var start int
	if _, err := fmt.Fscan(reader, &start); err != nil {
		fmt.Fprintln(out, "ОШИБКА: ожидалась стартовая вершина")
		return
	}
	if start < 1 || start > n {
		fmt.Fprintf(out, "ОШИБКА: старт вне [1..%d]\n", n)
		return
	}

	used := make([]bool, n)
	var stack Stack[int]
	stack.Push(start - 1)
	for !stack.Empty() {
		v, _ := stack.Pop()
		if used[v] {
			continue
		}
		used[v] = true
		fmt.Fprintf(out, "%d ", v+1)
		for node := graph[v].head; node != nil; node = node.next {
			if !used[node.value] {
				stack.Push(node.value)
			}
		}
	}
	fmt.Fprintln(out)
}
