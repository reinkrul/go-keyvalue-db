package bolt

import (
	"github.com/reinkrul/go-keyvalue-db/server/spi"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test(t *testing.T) {
	test(t, func(store spi.DataStore) {
		err := store.Set("foo", "hello", "world")
		assert.NoError(t, err)
		value, err := store.Get("foo", "hello")
		assert.NoError(t, err)
		assert.Equal(t, "world", value)
	})
}

func Test_Bucket_Doenst_Exist(t *testing.T) {
	test(t, func(store spi.DataStore) {
		value, err := store.Get("foo", "test")
		assert.NoError(t, err)
		assert.Empty(t, value)
	})
}

func Test_Value_Doenst_Exist(t *testing.T) {
	test(t, func(store spi.DataStore) {
		err := store.Set("foo", "hello", "world")
		assert.NoError(t, err)
		value, err := store.Get("foo", "test")
		assert.NoError(t, err)
		assert.Empty(t, value)
	})
}

func test(t *testing.T, f func(dataStore spi.DataStore)) {
	os.Remove("unittest.db")
	store, err := Connect("unittest.db")
	defer store.Close()
	assert.NoError(t, err)
	f(store)
}