package vm_test

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	MemorySize = flamego.MB
	CacheSize  = flamego.KB
	LineWidth  = 8
	BusSize    = 4
	OffsetBits = 3
)

func TestCache_Read(t *testing.T) {
	address := uint64(0)
	data := []byte{0, 1, 2, 3}

	memory := vm.NewMemory(MemorySize)
	memory.Set(address, data)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	// Cache shouldn't contain data
	assertCacheReadMiss(t, cache, address)

	// Memory should been read
	assertLowerRead(t, cache, memory, address)

	cache.Clock(0)
	assert.True(t, memory.IsFree()) // Cache should have freed memory

	// Cache should contain data after reading from memory
	assertCacheReadHit(t, cache, address, data)
}

func TestCache_Read_withEviction(t *testing.T) {
	address := uint64(0)
	data := make([]byte, 2*flamego.KB)

	for i := range data {
		data[i] = byte(i)
	}

	memory := vm.NewMemory(MemorySize)
	memory.Set(address, data)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	// Fill cache with dirty data
	for ; address < flamego.KB; address += BusSize {
		assertCacheWriteHit(t, cache, address, data[address:address+BusSize])
	}

	// When cache is full and a read request misses,
	assertCacheReadMiss(t, cache, address)

	// The data is fetched from the memory,
	assertLowerRead(t, cache, memory, address)

	cache.Clock(0)
	assert.False(t, memory.IsFree()) // Cache should not have freed memory

	// A line is evicted and written back to memory
	assertLowerWrite(t, cache, memory, cache.CreateAddress(0, 0, 0))

	cache.Clock(0)
	assert.True(t, memory.IsFree()) // Cache should have freed memory

	// Another cache read request will also miss,
	assertCacheReadMiss(t, cache, address)

	// and the data is fetched from the memory,
	assertLowerRead(t, cache, memory, address)

	cache.Clock(0)
	assert.True(t, memory.IsFree()) // Cache should have freed memory

	// Cache should contain data after reading from memory
	assertCacheReadHit(t, cache, address, data[address:address+BusSize])
}

func TestCache_Write(t *testing.T) {
	address := uint64(0)
	data := []byte{0, 1, 2, 3}

	memory := vm.NewMemory(MemorySize)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	assertCacheWriteHit(t, cache, address, data)
}

func TestCache_Write_withEviction_Writeback(t *testing.T) {
	address := uint64(0)
	data := make([]byte, 2*flamego.KB)

	for i := range data {
		data[i] = byte(i)
	}

	memory := vm.NewMemory(MemorySize)
	memory.Set(address, data)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	// Fill cache with dirty data
	for ; address < flamego.KB; address += BusSize {
		assertCacheWriteHit(t, cache, address, data[address:address+BusSize])
	}

	// When cache is full and a write request misses,
	assertCacheWriteMiss(t, cache, address, data[address:address+BusSize])

	// A line is evicted and written back to memory
	assertLowerWrite(t, cache, memory, cache.CreateAddress(0, 0, 0))

	cache.Clock(0)
	assert.True(t, memory.IsFree()) // Cache should have freed memory

	// Cache should have a free line now to write data
	assertCacheWriteHit(t, cache, address, data[address:address+BusSize])
}

func TestCache_Write_withEviction_Repurpose(t *testing.T) {
	address := uint64(0)
	data := make([]byte, 2*flamego.KB)

	for i := range data {
		data[i] = byte(i)
	}

	memory := vm.NewMemory(MemorySize)
	memory.Set(address, data)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	// Fill cache with clean data
	for ; address < flamego.KB; address += BusSize {
		if address%LineWidth == 0 {
			// Cache shouldn't contain data
			assertCacheReadMiss(t, cache, address)

			// Memory should been read
			assertLowerRead(t, cache, memory, address)

			cache.Clock(0)
			assert.True(t, memory.IsFree()) // Cache should have freed memory

			// Cache should contain data after reading from memory
			assertCacheReadHit(t, cache, address, data[address:address+BusSize])
		} else {
			// Cache should contain data after a previous read from memory
			assertCacheReadHit(t, cache, address, data[address:address+BusSize])
		}
	}

	// When cache is full and a write request misses,
	// a line must get evicted, if it is not dirty it doesn't need to be written back to memory
	// the line can then be used for the write so the request hits
	assertCacheWriteHit(t, cache, address, data[address:address+BusSize])
}

func TestCache_Clear(t *testing.T) {
	address := uint64(0)
	data := []byte{0, 1, 2, 3}

	stale := []byte{3, 2, 1, 0}

	memory := vm.NewMemory(MemorySize)
	memory.Set(address, data)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	// Clear should always hit, even if it the data wasn't in the cache
	assertCacheClearHit(t, cache, address)

	// Write stale data to cache
	assertCacheWriteHit(t, cache, address, stale)

	// Clear stale data from cache
	assertCacheClearHit(t, cache, address)

	// Cache shouldn't contain data
	assertCacheReadMiss(t, cache, address)

	// Memory should been read
	assertLowerRead(t, cache, memory, address)

	cache.Clock(0)
	assert.True(t, memory.IsFree()) // Cache should have freed memory

	// Cache should contain data after reading from memory
	assertCacheReadHit(t, cache, address, data)
}

func TestCache_Flush(t *testing.T) {
	address := uint64(0)
	data := []byte{0, 1, 2, 3}

	memory := vm.NewMemory(MemorySize)
	memory.Set(address, data)
	assert.False(t, memory.IsBusy())
	assert.True(t, memory.IsFree())

	cache := vm.NewCache(CacheSize, LineWidth, BusSize, OffsetBits, memory)
	assert.False(t, cache.IsBusy())
	assert.True(t, cache.IsFree())

	// Flush should hit if the data wasn't in the cache
	assertCacheFlushHit(t, cache, address)

	// Write data to cache
	assertCacheWriteHit(t, cache, address, data)

	// Flush data from cache should miss until lower write was successful
	assertCacheFlushMiss(t, cache, address)

	// Ensure data was written to memory
	memory.Clock(0)
	d := memory.Data()
	for i, b := range data {
		assert.Equal(t, b, d[i])
	}

	// Ensure cache updates dirty flags
	cache.Clock(0)
	tag, index, offset := cache.ParseAddress(address)
	line := cache.Lines()[index]
	assert.Equal(t, tag, line.Tag())

	for i, b := range data {
		j := i + int(offset)
		assert.True(t, line.IsValid(j))
		assert.False(t, line.IsDirty(j))
		assert.Equal(t, b, line.Read(j))
	}

	// Flush should hit if the data in the cache is clean
	assertCacheFlushHit(t, cache, address)
}

func TestCache_Address(t *testing.T) {
	address := uint64(0xab54a98ceb1f0ad2)

	memory := vm.NewMemory(flamego.SizeMemory)
	l3Cache := vm.NewL3Cache(flamego.SizeL3Cache, memory)
	l2Cache := vm.NewL2Cache(flamego.SizeL2Cache, l3Cache)
	l1Cache := vm.NewL1Cache(flamego.SizeL1Cache, l2Cache)

	assert.Equal(t, flamego.SizeL3Cache, l3Cache.Size())
	assert.Equal(t, flamego.LineWidthL3Cache, l3Cache.LineWidth())
	assert.Equal(t, 44, l3Cache.TagBits())
	assert.Equal(t, 8, l3Cache.IndexBits())
	assert.Equal(t, 12, l3Cache.OffsetBits())
	tag, index, offset := l3Cache.ParseAddress(address)
	assert.Equal(t, uint64(0xab54a98ceb1), tag)
	assert.Equal(t, uint64(0xf0), index)
	assert.Equal(t, uint64(0xad2), offset)

	assert.Equal(t, address, l3Cache.CreateAddress(tag, index, offset))

	assert.Equal(t, flamego.SizeL2Cache, l2Cache.Size())
	assert.Equal(t, flamego.LineWidthL2Cache, l2Cache.LineWidth())
	assert.Equal(t, 49, l2Cache.TagBits())
	assert.Equal(t, 6, l2Cache.IndexBits())
	assert.Equal(t, 9, l2Cache.OffsetBits())
	tag, index, offset = l2Cache.ParseAddress(address)
	assert.Equal(t, uint64(0x156a95319d63e), tag)
	assert.Equal(t, uint64(0x5), index)
	assert.Equal(t, uint64(0xd2), offset)

	assert.Equal(t, address, l2Cache.CreateAddress(tag, index, offset))

	assert.Equal(t, flamego.SizeL1Cache, l1Cache.Size())
	assert.Equal(t, flamego.LineWidthL1Cache, l1Cache.LineWidth())
	assert.Equal(t, 54, l1Cache.TagBits())
	assert.Equal(t, 4, l1Cache.IndexBits())
	assert.Equal(t, 6, l1Cache.OffsetBits())
	tag, index, offset = l1Cache.ParseAddress(address)
	assert.Equal(t, uint64(0x2ad52a633ac7c2), tag)
	assert.Equal(t, uint64(0xb), index)
	assert.Equal(t, uint64(0x12), offset)

	assert.Equal(t, address, l1Cache.CreateAddress(tag, index, offset))
}

func assertCacheReadHit(t *testing.T, cache *vm.Cache, address uint64, data []byte) {
	t.Helper()

	cache.Read(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheRead, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())      // Cache has not yet been released
	assert.True(t, cache.IsSuccessful()) // Data was read

	bus := cache.Bus()
	for i, b := range data {
		assert.True(t, bus.IsValid(i))
		assert.False(t, bus.IsDirty(i))
		assert.Equal(t, b, bus.Read(i))
	}

	cache.Free()
	assert.True(t, cache.IsFree())
}

func assertCacheReadMiss(t *testing.T, cache *vm.Cache, address uint64) {
	t.Helper()

	cache.Read(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheRead, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())       // Cache has not yet been released
	assert.False(t, cache.IsSuccessful()) // Data was not in cache

	cache.Free()
	assert.True(t, cache.IsFree())
}

func assertLowerRead(t *testing.T, cache *vm.Cache, memory *vm.Memory, address uint64) {
	// Cache should have issued a request to read from memory
	assert.Equal(t, address, cache.LowerAddress())
	assert.Equal(t, flamego.CacheRead, cache.LowerOperation())
	assert.True(t, memory.IsBusy())
	assert.False(t, memory.IsFree())
	assert.Equal(t, address, memory.Address())
	assert.Equal(t, flamego.MemoryRead, memory.Operation())

	memory.Clock(0)
	assert.False(t, memory.IsBusy())
	assert.False(t, memory.IsFree())      // Memory has not yet been released
	assert.True(t, memory.IsSuccessful()) // Data was retrieved from memory
}

func assertLowerWrite(t *testing.T, cache *vm.Cache, memory *vm.Memory, address uint64) {
	// Cache should have issued a request to write to memory
	assert.Equal(t, address, cache.LowerAddress())
	assert.Equal(t, flamego.CacheWrite, cache.LowerOperation())
	assert.True(t, memory.IsBusy())
	assert.False(t, memory.IsFree())
	assert.Equal(t, address, memory.Address())
	assert.Equal(t, flamego.MemoryWrite, memory.Operation())

	memory.Clock(0)
	assert.False(t, memory.IsBusy())
	assert.False(t, memory.IsFree())      // Memory has not yet been released
	assert.True(t, memory.IsSuccessful()) // Data was stored in memory
}

func assertCacheWriteMiss(t *testing.T, cache *vm.Cache, address uint64, data []byte) {
	t.Helper()

	bus := cache.Bus()
	for i, b := range data {
		bus.Write(i, b)
	}

	cache.Write(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheWrite, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())       // Cache has not yet been released
	assert.False(t, cache.IsSuccessful()) // Data was not written

	cache.Free()
	assert.True(t, cache.IsFree())
}

func assertCacheWriteHit(t *testing.T, cache *vm.Cache, address uint64, data []byte) {
	t.Helper()

	bus := cache.Bus()
	for i, b := range data {
		bus.Write(i, b)
	}

	cache.Write(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheWrite, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())      // Cache has not yet been released
	assert.True(t, cache.IsSuccessful()) // Data was written

	cache.Free()
	assert.True(t, cache.IsFree())

	tag, index, offset := cache.ParseAddress(address)
	line := cache.Lines()[index]
	assert.Equal(t, tag, line.Tag())

	for i, b := range data {
		j := i + int(offset)
		assert.True(t, line.IsValid(j))
		assert.True(t, line.IsDirty(j))
		assert.Equal(t, b, line.Read(j))
	}
}

func assertCacheClearHit(t *testing.T, cache *vm.Cache, address uint64) {
	t.Helper()

	cache.Clear(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheClear, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())      // Cache has not yet been released
	assert.True(t, cache.IsSuccessful()) // Data was cleared
}

func assertCacheFlushHit(t *testing.T, cache *vm.Cache, address uint64) {
	t.Helper()

	cache.Flush(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheFlush, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())      // Cache has not yet been released
	assert.True(t, cache.IsSuccessful()) // Data was flushed
}

func assertCacheFlushMiss(t *testing.T, cache *vm.Cache, address uint64) {
	t.Helper()

	cache.Flush(address)
	assert.True(t, cache.IsBusy())
	assert.False(t, cache.IsFree())
	assert.Equal(t, address, cache.Address())
	assert.Equal(t, flamego.CacheFlush, cache.Operation())

	cache.Clock(0)
	assert.False(t, cache.IsBusy())
	assert.False(t, cache.IsFree())       // Cache has not yet been released
	assert.False(t, cache.IsSuccessful()) // Data is still being flushed
}
