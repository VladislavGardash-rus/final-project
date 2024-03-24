package lru_cach_service

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("a", 1) // [1]
		cache.Set("b", 2) // [1,2]
		cache.Set("c", 3) // [1,2,3]
		cache.Clear()

		_, ok := cache.Get("a")
		require.False(t, ok)

		_, ok = cache.Get("b")
		require.False(t, ok)

		_, ok = cache.Get("c")
		require.False(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("a", 1) // [1]
		cache.Set("b", 2) // [2,1]
		cache.Set("c", 3) // [3,2,1]
		cache.Set("d", 4) // [4,3,2]

		_, ok := cache.Get("a")
		require.False(t, ok)

		val, ok := cache.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = cache.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = cache.Get("d") // [4,2,3]
		require.True(t, ok)
		require.Equal(t, 4, val)

		cache.Set("e", 5) // [5,4,2]

		val, ok = cache.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = cache.Get("e")
		require.True(t, ok)
		require.Equal(t, 5, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
