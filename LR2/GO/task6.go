package main

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

type CuckooHash struct {
	t1, t2                                                 Array[string]
	seed1, seed2                                           uint64
	attempts, displacements, rebuildAttempts, lookupChecks uint64
}

func newCuckooHash(capacity int) *CuckooHash {
	c := &CuckooHash{seed2: 1}
	c.initTables(capacity)
	return c
}

func (c *CuckooHash) initTables(capacity int) {
	c.t1.Clear()
	c.t2.Clear()
	for i := 0; i < capacity; i++ {
		c.t1.Push("")
		c.t2.Push("")
	}
}

func (c *CuckooHash) h1(key string) int {
	return int((hashString(key) + c.seed1) % uint64(c.t1.Size()))
}

func (c *CuckooHash) h2(key string) int {
	return int((hashString(key)*33 + c.seed2) % uint64(c.t2.Size()))
}

func (c *CuckooHash) Contains(key string) bool {
	c.lookupChecks++
	if c.t1.Get(c.h1(key)) == key {
		return true
	}
	c.lookupChecks++
	return c.t2.Get(c.h2(key)) == key
}

func (c *CuckooHash) place(key string) bool {
	cur := key
	first := true
	maxLoop := 2 * c.t1.Size()
	for i := 0; i < maxLoop; i++ {
		c.attempts++
		if first {
			p := c.h1(cur)
			displaced := c.t1.Get(p)
			if displaced != "" {
				c.displacements++
			}
			c.t1.Set(p, cur)
			cur = displaced
		} else {
			p := c.h2(cur)
			displaced := c.t2.Get(p)
			if displaced != "" {
				c.displacements++
			}
			c.t2.Set(p, cur)
			cur = displaced
		}
		if cur == "" {
			return true
		}
		first = !first
	}
	return false
}

func (c *CuckooHash) rebuild(keys *Array[string]) {
	capacity := c.t1.Size()
	failedAttempts := 0
	for {
		c.initTables(capacity)
		ok := true
		for i := 0; i < keys.Size(); i++ {
			if !c.place(keys.Get(i)) {
				ok = false
				break
			}
		}
		if ok {
			return
		}
		c.rebuildAttempts++
		c.seed1++
		c.seed2 += 2
		failedAttempts++
		if failedAttempts%5 == 0 {
			capacity *= 2
		}
	}
}

func (c *CuckooHash) Add(key string) {
	if key == "" || c.Contains(key) {
		return
	}
	var keys Array[string]
	for i := 0; i < c.t1.Size(); i++ {
		if val := c.t1.Get(i); val != "" {
			keys.Push(val)
		}
	}
	for i := 0; i < c.t2.Size(); i++ {
		if val := c.t2.Get(i); val != "" {
			keys.Push(val)
		}
	}
	keys.Push(key)
	c.rebuild(&keys)
}

func isomorphic(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var ab, ba [256]int
	for i := 0; i < 256; i++ {
		ab[i], ba[i] = -1, -1
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if ab[x] != -1 && ab[x] != int(y) {
			return false
		}
		if ba[y] != -1 && ba[y] != int(x) {
			return false
		}
		ab[x], ba[y] = int(y), int(x)
	}
	return true
}

func task6(in io.Reader, out io.Writer) {
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

	var chain ChainHash
	cuckoo := newCuckooHash(16)
	insertStart := time.Now()
	for i := 0; i < n; i++ {
		var key string
		if _, err := fmt.Fscan(reader, &key); err != nil {
			fmt.Fprintf(out, "ОШИБКА: ожидался ключ #%d\n", i+1)
			return
		}
		chain.Add(key)
		cuckoo.Add(key)
	}
	insertTime := time.Since(insertStart).Nanoseconds()

	var searchKey string
	if _, err := fmt.Fscan(reader, &searchKey); err != nil {
		fmt.Fprintln(out, "ОШИБКА: ожидался ключ для поиска")
		return
	}
	findStart := time.Now()
	chainFound := chain.Contains(searchKey)
	chainFindTime := time.Since(findStart).Nanoseconds()
	findStart = time.Now()
	cuckooFound := cuckoo.Contains(searchKey)
	cuckooFindTime := time.Since(findStart).Nanoseconds()

	fmt.Fprintln(out)
	if chainFound {
		fmt.Fprintln(out, "Chain: FOUND")
	} else {
		fmt.Fprintln(out, "Chain: NOT FOUND")
	}
	if cuckooFound {
		fmt.Fprintln(out, "Cuckoo: FOUND")
	} else {
		fmt.Fprintln(out, "Cuckoo: NOT FOUND")
	}

	var a, b string
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		fmt.Fprintln(out, "ОШИБКА: ожидались две строки для проверки изоморфизма")
		return
	}
	isoStart := time.Now()
	isoResult := isomorphic(a, b)
	isoTime := time.Since(isoStart).Nanoseconds()
	if isoResult {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}

	fmt.Fprintln(out, "\nСтатистика выполнения:")
	fmt.Fprintf(out, "Добавление %d ключей в обе таблицы: %d нс\n", n, insertTime)
	fmt.Fprintf(out, "Поиск в Chain Hash: %d нс; сравнений: %d\n", chainFindTime, chain.attempts)
	fmt.Fprintf(out, "Поиск в Cuckoo Hash: %d нс; проверок ячеек: %d\n", cuckooFindTime, cuckoo.lookupChecks)
	fmt.Fprintf(out, "Cuckoo Hash: попыток размещения: %d; вытеснений: %d; перестроек: %d\n", cuckoo.attempts, cuckoo.displacements, cuckoo.rebuildAttempts)
	fmt.Fprintf(out, "Проверка изоморфизма: %d нс\n", isoTime)

	fmt.Fprintln(out, "\nЭмпирический анализ (время в микросекундах)")
	fmt.Fprintln(out, "N\tChain insert\tChain find\tChain comparisons\tCuckoo insert\tCuckoo find\tCuckoo placements\tDisplacements\tRebuilds")
	for _, size := range []int{100, 1000, 5000} {
		var ch ChainHash
		cap := size
		if cap < 16 {
			cap = 16
		}
		ck := newCuckooHash(cap)
		var keys Array[string]
		for i := 0; i < size; i++ {
			keys.Push(fmt.Sprintf("key_%d", i))
		}

		start := time.Now()
		for i := 0; i < size; i++ {
			ch.Add(keys.Get(i))
		}
		chainInsert := time.Since(start).Microseconds()
		start = time.Now()
		for i := 0; i < size; i++ {
			ch.Contains(keys.Get(i))
		}
		chainFind := time.Since(start).Microseconds()
		chainComparisons := ch.attempts

		start = time.Now()
		for i := 0; i < size; i++ {
			ck.Add(keys.Get(i))
		}
		cuckooInsert := time.Since(start).Microseconds()
		start = time.Now()
		found := 0
		for i := 0; i < size; i++ {
			if ck.Contains(keys.Get(i)) {
				found++
			}
		}
		cuckooFind := time.Since(start).Microseconds()
		// found фиксирует, что поиск во всех ключах был реально произведён.
		_ = found

		fmt.Fprintf(out, "%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", size,
			chainInsert, chainFind, chainComparisons, cuckooInsert, cuckooFind,
			ck.attempts, ck.displacements, ck.rebuildAttempts)
	}
}
