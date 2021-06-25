package xcache

import (
	"github.com/duchiporexia/goutils/xerr"
	"sync"
	"time"
)

type XCache struct {
	shards        []*XCacheShard
	shardIdxMask  int
	shardCount    int
	cleanMaxSize  int
	ringSlotCount int
	close         chan struct{}
}

type XCacheShard struct {
	sync.RWMutex
	itemMap map[string]*Item
	head    Item
	tail    Item
	currIdx uint64
}

func (s *XCacheShard) addItem(item *Item, ringSlotCount int) {
	s.Lock()
	oldItem, ok := s.itemMap[item.key]
	if ok {
		prevItem := oldItem.prev
		nextItem := oldItem.next
		prevItem.next = nextItem
		nextItem.prev = prevItem
	}

	tail := &s.tail
	lastItem := tail.prev

	item.prev = lastItem
	item.next = tail
	item.expiredAt = s.currIdx + uint64(ringSlotCount)

	lastItem.next = item
	tail.prev = item

	s.itemMap[item.key] = item
	s.Unlock()
}

type Item struct {
	key       string
	value     interface{}
	l         *sync.RWMutex
	isValid   bool
	prev      *Item
	next      *Item
	expiredAt uint64
}

type XCacheConfig struct {
	ShardCount   int
	EvictionTime time.Duration
	CleanPeriod  time.Duration
	CleanMaxSize int
}

func NewXCache(config *XCacheConfig) *XCache {
	shardCount := config.ShardCount
	if shardCount == 0 {
		shardCount = 128
	}
	evictionTime := config.EvictionTime
	if evictionTime == time.Duration(0) {
		evictionTime = 15 * time.Minute
	}
	cleanPeriod := config.CleanPeriod
	if cleanPeriod == time.Duration(0) {
		cleanPeriod = 10 * time.Second
	}
	cleanMaxSize := config.CleanMaxSize
	if cleanMaxSize == 0 {
		cleanMaxSize = 10000
	}

	ringSlotCount := int(evictionTime / cleanPeriod)
	shardIdxBit := 1

	for true {
		if 1<<shardIdxBit >= shardCount {
			break
		}
		shardIdxBit++
	}
	shardCount = 1 << shardIdxBit

	cache := XCache{
		shards:       make([]*XCacheShard, shardCount),
		shardCount:   shardCount,
		cleanMaxSize: cleanMaxSize,
		close:        make(chan struct{})}

	cache.shardIdxMask = 1<<shardIdxBit - 1
	cache.ringSlotCount = ringSlotCount

	for i := 0; i < shardCount; i++ {
		shard := &XCacheShard{
			itemMap: make(map[string]*Item)}
		shard.tail.prev = &shard.head
		shard.head.next = &shard.tail

		cache.shards[i] = shard
	}

	go func() {
		ticker := time.NewTicker(cleanPeriod)
		defer ticker.Stop()
		for {
			select {
			case _ = <-ticker.C:
				cache.cleanJob()
			case <-cache.close:
				return
			}
		}
	}()

	return &cache
}

func (m *XCache) cleanJob() {
	//fmt.Printf("run cleanJob ...\n")
	for i := 0; i < m.shardCount; i++ {
		shard := m.shards[i]
		shard.Lock()
		shard.currIdx++
		head := &shard.head
		tail := &shard.tail
		currItem := shard.head.next
		evictSize := 0
		for currItem != tail && currItem.expiredAt <= shard.currIdx && evictSize < m.cleanMaxSize {
			delete(shard.itemMap, currItem.key)
			currItem = currItem.next
			evictSize++
		}
		head.next = currItem
		currItem.prev = head
		shard.Unlock()
	}
}

func (m *XCache) Close() {
	close(m.close)
}

func (m *XCache) getShard(key string) *XCacheShard {
	hashKey := fnv32(key)
	return m.shards[hashKey&uint32(m.shardIdxMask)]
}

func (m *XCache) Set(key string, value interface{}) {
	shard := m.getShard(key)
	item := &Item{
		key:     key,
		value:   value,
		isValid: true,
	}
	shard.addItem(item, m.ringSlotCount)
}

func (m *XCache) Get(key string) (interface{}, bool) {
	shard := m.getShard(key)
	shard.RLock()
	item, ok := shard.itemMap[key]
	shard.RUnlock()
	if ok {
		return item.value, item.isValid
	}
	return nil, false
}

func (m *XCache) GetOrFetch(key string, cb func(key string) (interface{}, error)) (interface{}, bool) {
	shard := m.getShard(key)
	shard.RLock()
	item, ok := shard.itemMap[key]
	shard.RUnlock()
	if ok {
		l := item.l
		if l == nil {
			return item.value, item.isValid
		}
		l.RLock()
		val := item.value
		isValid := item.isValid
		l.RUnlock()
		return val, isValid
	}
	if cb == nil {
		return nil, false
	}
	var l sync.RWMutex
	l.Lock()
	item = &Item{
		key:   key,
		value: nil,
		l:     &l,
	}
	shard.addItem(item, m.ringSlotCount)
	val, err := safeCall(key, cb)
	isValid := err == nil
	item.value = val
	item.isValid = isValid
	item.l = nil
	l.Unlock()
	return val, isValid
}

func safeCall(key string, cb func(key string) (interface{}, error)) (interface{}, error) {
	var errx error
	defer func() {
		if err := recover(); err != nil {
			err = xerr.ErrSafeCall
		}
	}()
	val, errx := cb(key)
	return val, errx
}

func (m *XCache) Delete(key string) {
	shard := m.getShard(key)
	shard.Lock()
	oldItem, ok := shard.itemMap[key]
	if ok {
		prevItem := oldItem.prev
		nextItem := oldItem.next
		prevItem.next = nextItem
		nextItem.prev = prevItem

		delete(shard.itemMap, key)
	}
	shard.Unlock()
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
