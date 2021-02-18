package sqler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func initDB() (*DB, error) {
	return NewDB(Options{
		Driver:    "mysql",
		ConnStr:   "root:123456!@tcp(127.0.0.1:3306)/test_db?charset=utf8&parseTime=True&loc=Local",
		KeepAlive: 10,
		MaxIdles:  10,
		MaxOpens:  10,
	})
}
func TestResolve(t *testing.T) {
	if testing.Short() {
		return
	}
	db, err := initDB()
	assert.Nil(t, err)

	data, err := db.Select("users", Where{
		"_limit": []uint{1},
	}, Fields{"*"}).ConvertToMap()
	assert.Len(t, data, 1)
}
