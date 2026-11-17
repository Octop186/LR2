package main

import (
	"fmt"
	"io"
)

func task4(in io.Reader, out io.Writer) {
	input, err := io.ReadAll(in)
	if err != nil {
		fmt.Fprintln(out, "ОШИБКА: чтение ввода")
		return
	}
	var count [256]int
	for _, c := range input {
		if c != ' ' && c != '\n' {
			count[c]++
		}
	}
	maxCount := 0
	for _, c := range count {
		if c > maxCount {
			maxCount = c
		}
	}
	for row := maxCount; row > 0; row-- {
		for _, c := range count {
			if c > 0 {
				if c >= row {
					fmt.Fprint(out, "#")
				} else {
					fmt.Fprint(out, " ")
				}
			}
		}
		fmt.Fprintln(out)
	}
	for i, c := range count {
		if c > 0 {
			fmt.Fprintf(out, "%c", byte(i))
		}
	}
	fmt.Fprintln(out)
}
