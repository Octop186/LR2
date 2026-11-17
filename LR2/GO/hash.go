package main

// Та же полиномиальная хеш-функция, что и в hash.cpp (h = h*31 + c).
func hashString(s string) uint64 {
	var h uint64
	for i := 0; i < len(s); i++ {
		h = h*31 + uint64(s[i])
	}
	return h
}

const chainCapacity = 101

type ChainHash struct {
	table    [chainCapacity]ForwardList[string]
	attempts uint64
}

func (h *ChainHash) index(key string) int { return int(hashString(key) % chainCapacity) }

func (h *ChainHash) Add(key string) {
	index := h.index(key)
	for cur := h.table[index].head; cur != nil; cur = cur.next {
		h.attempts++
		if cur.value == key {
			return
		}
	}
	h.table[index].PushFront(key)
}

func (h *ChainHash) Del(key string) bool { return h.table[h.index(key)].DelValue(key) }

func (h *ChainHash) Contains(key string) bool {
	for cur := h.table[h.index(key)].head; cur != nil; cur = cur.next {
		h.attempts++
		if cur.value == key {
			return true
		}
	}
	return false
}
