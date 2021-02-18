package sqler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
)

const TagName = "sqler"

var (
	ctx context.Context
)

func init() {
	scanner.SetTagName(TagName)
	ctx = context.Background()
}

type (
	DB struct {
		*sql.DB
		ticker *time.Ticker
	}
	Options struct {
		Driver  string `yaml:"driver" mapstructure:"driver"`
		ConnStr string `yaml:"dsn" mapstructure:"dsn"`
		// 定时保活
		KeepAlive int `yaml:"keep_alive" mapstructure:"keep_alive"`
		// 最大可空闲连接数量
		MaxIdles int `yaml:"max_idles" mapstructure:"max_idles"`
		// 最大连接数量
		MaxOpens    int `yaml:"max_opens" mapstructure:"max_opens"`
		MaxLifeTime int `yaml:"max_life_time" mapstructure:"max_life_time"`
	}

	Where  map[string]interface{}
	Fields []string
)

// NewDB -
func NewDB(ops Options) (*DB, error) {
	db, err := sql.Open(ops.Driver, ops.ConnStr)
	if err != nil {
		return nil, err
	}
	if ops.MaxOpens > 0 {
		db.SetMaxOpenConns(ops.MaxOpens)
	}
	if ops.MaxIdles > 0 {
		db.SetMaxIdleConns(ops.MaxIdles)
	}
	if ops.MaxLifeTime > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(ops.MaxLifeTime))
	}

	rdb := &DB{DB: db}
	if ops.KeepAlive > 0 {
		rdb.keepAlive(time.Duration(ops.KeepAlive) * time.Second)
	}

	return rdb, nil
}

// 定时保活
func (db *DB) keepAlive(d time.Duration) {
	db.ticker = time.NewTicker(d)
	go func() {
		for range db.ticker.C {
			if err := db.Ping(); err != nil {
				fmt.Printf("数据库断开连接，%v", err)
			}
		}
	}()
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	if db.ticker != nil {
		db.ticker.Stop()
	}
	return db.DB.Close()
}

// Select 查询数据
func (db *DB) Select(name string, where Where, fields Fields) *Fruit {
	cond, values, err := builder.BuildSelect(name, where, fields)
	if err != nil {
		return &Fruit{
			err: err,
		}
	}

	rs, err := db.Query(cond, values...)
	return &Fruit{
		rows: rs,
		err:  err,
	}
}

// Modify 变更数据
func (db *DB) Modify(name string, where Where, update map[string]interface{}) (int64, error) {
	cond, values, err := builder.BuildUpdate(name, where, update)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(cond, values)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (db *DB) Delete(name string, where Where) (int64, error) {
	cond, values, err := builder.BuildDelete(name, where)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(cond, values)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Insert 单行插入
func (db *DB) Insert(name string, data map[string]interface{}) (int64, error) {
	return db.MultiInsert(name, []map[string]interface{}{data})
}

// MultiInsert 多行插入
func (db *DB) MultiInsert(name string, data []map[string]interface{}) (int64, error) {
	cond, values, err := builder.BuildInsert(name, data)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(cond, values)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Count 获取数据条数
func (db *DB) Count(name string, where Where, col string) (int64, error) {
	result, err := builder.AggregateQuery(ctx, db.DB, name, where, builder.AggregateCount(col))
	if err != nil {
		return 0, err
	}
	return result.Int64(), nil
}
