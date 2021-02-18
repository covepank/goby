package sqler

import (
	"database/sql"
	"errors"
	"reflect"
	"time"

	"github.com/didi/gendry/scanner"

	"github.com/sanbsy/goby/internal/strs"
)

// Mapping 将结构体转换成 map
func Mapping(data interface{}) map[string]interface{} {
	dataValue := reflect.ValueOf(data)
	if dataValue.Type().Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	dataType := dataValue.Type()
	result := make(map[string]interface{}, dataValue.NumField())

	for i := 0; i < dataValue.NumField(); i++ {
		result[getName(dataType.Field(i))] = dataValue.Field(i).Interface()
	}

	return result

}

// ScanToMap 将 SQL 查询结构 序列化到一个 map 中
func ScanToMap(rs *sql.Rows) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for rs.Next() {
		columns, err := rs.ColumnTypes()
		if err != nil {
			return nil, err
		}

		payload := make([]interface{}, 0, len(columns))
		for _, col := range columns {
			t, err := reflectType(col)
			if err != nil {
				return nil, err
			}
			payload = append(payload, reflect.New(t).Interface())
		}

		err = rs.Scan(payload...)
		if err != nil {
			return nil, err
		}

		items := make(map[string]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			items[columns[i].Name()] = reflect.ValueOf(payload[i]).Elem().Interface()
		}

		result = append(result, items)
	}

	return result, nil
}

func Scan(rs *sql.Rows, target interface{}) error {
	if err := scanner.Scan(rs, target); err != nil {
		return err
	}
	return nil
}

// 映射 数据库类型和 golang 类型
// 暂时不支持 复杂类型 如 enum 等
func reflectType(s *sql.ColumnType) (reflect.Type, error) {
	var vty reflect.Type

	switch s.DatabaseTypeName() {
	case "BIT", "TINYINT", "BOOL":
		vty = reflect.TypeOf(false)
	case "DATE", "DATETIME", "TIME", "TIMESTAMP":
		vty = reflect.TypeOf(time.Now())
	case "TEXT", "BLOB", "LONGBLOB", "LONGTEXT",
		"MEDIUMBLOB", "TINYBLOB", "TINYTEXT",
		"MEDIUMTEXT", "BINARY", "CHAR", "VARBINARY",
		"VARCHAR", "NVARCHAR":
		vty = reflect.TypeOf("")
	case "INT", "MEDIUMINT", "SMALLINT":
		vty = reflect.TypeOf(0)
	case "BIGINT":
		vty = reflect.TypeOf(int64(0))
	case "DOUBLE", "FLOAT", "DECIMAL":
		vty = reflect.TypeOf(float64(0))
	default:
		return nil, errors.New("can't resolve db type")
	}

	if r, ok := s.Nullable(); r && ok {
		vty = reflect.TypeOf(reflect.New(vty).Interface())
	}

	return vty, nil
}

func getName(ty reflect.StructField) string {
	if name, exist := ty.Tag.Lookup(TagName); exist {
		return name
	}

	return strs.LowerName(ty.Name)
}
