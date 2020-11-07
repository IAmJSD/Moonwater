package fastsql

import (
	"fastsql/driver"
	"github.com/auttaja/go-tlru"
	"time"
)

// QueryManager is the manager which is responsible for making/managing queries. Note this should only be made with NewQueryManager.
type QueryManager struct {
	d driver.DBDriver
	cache *tlru.Cache
}

// Handles the translation cache after compilation.
func (q *QueryManager) handleTranslationCache(stmt interface{}) (interface{}, error) {
	item, ok := q.cache.Get(stmt)
	if ok {
		return item, nil
	}
	translation, err := q.d.Translate(stmt)
	if err != nil {
		return nil, err
	}
	q.cache.Set(stmt, translation)
	return translation, nil
}

// NewQueryManager is used to create a new query manager based on a driver.
func NewQueryManager(DriverInit func(string) (driver.DBDriver, error), ConnectionString string) (*QueryManager, error) {
	d, err := DriverInit(ConnectionString)
	if err != nil {
		return nil, err
	}
	return &QueryManager{cache: tlru.NewCache(1000, 1000000, time.Minute*10), d: d}, nil
}
