package main

import (
	"bufio"
	"fmt"
	"io"
)

type CacheItem struct{ key, value, freq, last int }

type LFUCache struct {
	data     Array[CacheItem]
	capacity int
	timer    int
}

func newLFUCache(capacity int) *LFUCache { return &LFUCache{capacity: capacity} }

func (c *LFUCache) find(key int) int {
	for i := 0; i < c.data.Size(); i++ {
		if c.data.Get(i).key == key {
			return i
		}
	}
	return -1
}

func (c *LFUCache) get(key int) int {
	index := c.find(key)
	if index == -1 {
		return -1
	}
	item := c.data.Get(index)
	c.timer++
	item.freq++
	item.last = c.timer
	c.data.Set(index, item)
	return item.value
}

func (c *LFUCache) set(key, value int) {
	if c.capacity <= 0 {
		return
	}
	index := c.find(key)
	if index != -1 {
		item := c.data.Get(index)
		c.timer++
		item.value = value
		item.freq++
		item.last = c.timer
		c.data.Set(index, item)
		return
	}
	if c.data.Size() == c.capacity {
		del := 0
		for i := 1; i < c.data.Size(); i++ {
			current, target := c.data.Get(i), c.data.Get(del)
			if current.freq < target.freq || (current.freq == target.freq && current.last < target.last) {
				del = i
			}
		}
		c.data.Remove(del)
	}
	c.timer++
	c.data.Push(CacheItem{key: key, value: value, freq: 1, last: c.timer})
}

func task7(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	var cap, q int
	if _, err := fmt.Fscan(reader, &cap, &q); err != nil {
		fmt.Fprintln(out, "ОШИБКА: ожидалось cap q")
		return
	}
	if cap < 0 || q < 0 {
		fmt.Fprintln(out, "ОШИБКА: cap >= 0, q >= 0")
		return
	}
	cache := newLFUCache(cap)
	for i := 0; i < q; i++ {
		var op string
		var key int
		if _, err := fmt.Fscan(reader, &op, &key); err != nil {
			fmt.Fprintf(out, "ОШИБКА: ожидалось <op> <key> в запросе #%d\n", i+1)
			return
		}
		if op == "SET" {
			var value int
			if _, err := fmt.Fscan(reader, &value); err != nil {
				fmt.Fprintln(out, "ОШИБКА: SET требует ключ и значение")
				return
			}
			cache.set(key, value)
		} else if op == "GET" {
			fmt.Fprintln(out, cache.get(key))
		} else {
			fmt.Fprintf(out, "ОШИБКА: неизвестная операция '%s'\n", op)
			return
		}
	}
}
