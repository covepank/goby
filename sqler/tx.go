package sqler

import (
	"context"
	"database/sql"

	"github.com/didi/gendry/builder"
	"github.com/sanbsy/gopkg/internal/idx"
)

type (
	// Tx 事务
	Tx struct {
		*sql.Tx
	}
	// TxOptions 事务配置项，类型别名
	TxOptions = sql.TxOptions
)

// Begin 开启一个事务
func (db *DB) Begin() (*Tx, error) {
	otx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{otx}, nil
}

// BeginCtx 开启一个事务
func (db *DB) BeginCtx(ctx context.Context, opts *TxOptions) (*Tx, error) {
	otx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{otx}, nil
}

// Select  execute select by transaction
func (tx *Tx) Select(name string, where Where, fields Fields) *Fruit {
	cond, values, err := builder.BuildSelect(name, where, fields)
	if err != nil {
		return &Fruit{
			err: err,
		}
	}

	rs, err := tx.Query(cond, values...)
	return &Fruit{
		rows: rs,
		err:  err,
	}
}

// Modify 修改数据
func (tx *Tx) Modify(name string, where Where, update map[string]interface{}) (int64, error) {
	cond, values, err := builder.BuildUpdate(name, where, update)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(cond, values)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (tx *Tx) Delete(name string, where Where) (int64, error) {
	cond, values, err := builder.BuildDelete(name, where)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(cond, values)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Insert 插入单行数据
func (tx *Tx) Insert(name string, data map[string]interface{}) (int64, error) {
	return tx.MultiInsert(name, []map[string]interface{}{data})
}

// MultiInsert 插入多行
func (tx *Tx) MultiInsert(name string, data []map[string]interface{}) (int64, error) {
	cond, values, err := builder.BuildInsert(name, data)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(cond, values)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// SetSavePoint Create SavePoint
func (tx *Tx) SetSavePoint() (string, error) {
	tag := idx.NewID().String()
	_, err := tx.Exec("RELEASE SAVEPOINT ?", tag)
	if err != nil {
		return "", err
	}
	return tag, nil
}

// RollbackSavePoint 回滚到指定 SavePoint
func (tx *Tx) RollbackSavePoint(tag string) error {
	_, err := tx.Exec("ROLLBACK TO SAVEPOINT ?", tag)
	if err != nil {
		return err
	}
	return nil
}

// Commit 事务提交
func (tx *Tx) Commit() error {
	return tx.Commit()
}

// Rollback 回滚事务
func (tx *Tx) Rollback() error {
	return tx.Rollback()
}
