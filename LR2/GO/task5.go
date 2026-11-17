package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type BSTNode struct {
	value int
	left  *BSTNode
	right *BSTNode
}

func insertBST(root *BSTNode, value int) *BSTNode {
	if root == nil {
		return &BSTNode{value: value}
	}
	if value < root.value {
		root.left = insertBST(root.left, value)
	}
	if value > root.value {
		root.right = insertBST(root.right, value)
	}
	return root
}

func printLeaves(root *BSTNode, out io.Writer) {
	if root == nil {
		return
	}
	printLeaves(root.left, out)
	if root.left == nil && root.right == nil {
		fmt.Fprintf(out, "%d ", root.value)
	}
	printLeaves(root.right, out)
}

func task5(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	var root *BSTNode
	for {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			break
		}
		value, err := strconv.Atoi(token)
		if err != nil {
			fmt.Fprintf(out, "ОШИБКА: '%s' — не целое число\n", token)
			return
		}
		if value == 0 {
			break
		}
		root = insertBST(root, value)
	}
	printLeaves(root, out)
	fmt.Fprintln(out)
}
