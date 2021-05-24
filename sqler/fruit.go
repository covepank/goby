package sqler

import (
	"database/sql"

	"github.com/didi/gendry/scanner"

	"github.com/sanbsy/gopkg/errors"
)

type Fruit struct {
	rows *sql.Rows
	err  error
}

func (rs *Fruit) Error() error {
	if rs.err != nil {
		return rs.err
	}
	return nil
}

func (rs *Fruit) Rows() *sql.Rows {
	return rs.rows
}

func (rs *Fruit) ConvertTo(target interface{}) error {
	if rs.err != nil {
		return rs.err
	}

	defer rs.Close()

	if err := scanner.Scan(rs.rows, target); err != nil {
		return err
	}
	return nil
}

func (rs *Fruit) Close() {
	if rs.rows != nil {
		rs.rows.Close()
	}
}

func (rs *Fruit) ConvertToMap() ([]map[string]interface{}, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	defer rs.Close()

	return scanner.ScanMapDecode(rs.rows)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, scanner.ErrEmptyResult)
}
