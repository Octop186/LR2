package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func applySetCommand(s *ChainHash, op, value string, out io.Writer) bool {
	switch op {
	case "SETADD":
		s.Add(value)
	case "SETDEL":
		s.Del(value)
	case "SET_AT":
		if s.Contains(value) {
			fmt.Fprintln(out, "TRUE")
		} else {
			fmt.Fprintln(out, "FALSE")
		}
	default:
		fmt.Fprintf(out, "ОШИБКА: неизвестная операция '%s'\n", op)
		return false
	}
	return true
}

func task2(args []string, in io.Reader, out io.Writer) {
	var file, query string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--file":
			if i+1 >= len(args) {
				fmt.Fprintln(out, "ОШИБКА: после --file нужно указать имя файла")
				return
			}
			i++
			file = args[i]
		case "--query":
			if i+1 >= len(args) {
				fmt.Fprintln(out, "ОШИБКА: после --query нужна команда")
				return
			}
			i++
			query = args[i]
		default:
			fmt.Fprintf(out, "ОШИБКА: неизвестный аргумент '%s'\n", args[i])
			return
		}
	}

	var set ChainHash
	if file != "" {
		inputFile, err := os.Open(file)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(out, "ОШИБКА: не удалось открыть файл %s\n", file)
			return
		}
		if err == nil {
			scanner := bufio.NewScanner(inputFile)
			for scanner.Scan() {
				fields := strings.Fields(scanner.Text())
				if len(fields) == 0 {
					continue
				}
				if len(fields) != 2 {
					fmt.Fprintf(out, "ОШИБКА: неправильная строка '%s'\n", scanner.Text())
					continue
				}
				// В старом C++ коде SET_AT ошибочно записывалась в историю.
				// При загрузке игнорируем эту старую команду чтения.
				if fields[0] == "SET_AT" {
					continue
				}
				applySetCommand(&set, fields[0], fields[1], out)
			}
			if scanner.Err() != nil {
				fmt.Fprintln(out, "ОШИБКА: чтение файла")
				inputFile.Close()
				return
			}
			inputFile.Close()
		}
	}

	if query == "" {
		reader := bufio.NewReader(in)
		for {
			var op, value string
			if _, err := fmt.Fscan(reader, &op, &value); err != nil {
				break
			}
			applySetCommand(&set, op, value, out)
		}
		return
	}

	fields := strings.Fields(query)
	if len(fields) != 2 {
		fmt.Fprintln(out, "ОШИБКА: запрос должен быть <op> <value>")
		return
	}
	op, value := fields[0], fields[1]
	if !applySetCommand(&set, op, value, out) {
		return
	}

	// Исправление ошибки C++: сохраняем только изменения, но не SET_AT.
	if file != "" && (op == "SETADD" || op == "SETDEL") {
		outputFile, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintln(out, "ОШИБКА: не удалось открыть файл для записи")
			return
		}
		fmt.Fprintf(outputFile, "%s %s\n", op, value)
		if err := outputFile.Close(); err != nil {
			fmt.Fprintln(out, "ОШИБКА: запись файла")
		}
	}
}
