package main

import (
	"fmt"
	"io"
	"os"
)

func run(args []string, in io.Reader, out io.Writer) {
	if len(args) == 0 {
		fmt.Fprintln(out, "Укажите номер задания: 1-7")
		return
	}
	switch args[0] {
	case "1":
		task1(in, out)
	case "2":
		task2(args[1:], in, out)
	case "3":
		task3(in, out)
	case "4":
		task4(in, out)
	case "5":
		task5(in, out)
	case "6":
		task6(in, out)
	case "7":
		task7(in, out)
	default:
		fmt.Fprintln(out, "Нет такого задания")
	}
}

func main() {
	run(os.Args[1:], os.Stdin, os.Stdout)
}
