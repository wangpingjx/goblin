/*****************************
 *       处理SQL执行结果       *
 *****************************/

package db

import (
    "database/sql"
    "reflect"
    "log"
)

// 渲染查询结果, session.Value 支持的格式为 *Author{} or []*Author{} or *[]*Author{}
func (session *Session) scan(rows *sql.Rows) {
    var (
        isSlice, isPtr bool
        resultType     reflect.Type
    )
    results := reflect.ValueOf(session.Value)
    for results.Kind() == reflect.Ptr {
        results = results.Elem()
    }
    if kind := results.Kind(); kind == reflect.Slice {
        isSlice = true
        resultType = results.Type().Elem()
        results.Set(reflect.MakeSlice(results.Type(), 0, 0))

        if resultType.Kind() == reflect.Ptr {
            isPtr      = true
            resultType = resultType.Elem()
        }
    }
    defer rows.Close()
    for rows.Next() {
        elem := results
        if isSlice {
            elem = reflect.New(resultType).Elem()
        }
        modelStruct := session.GetModelStruct(elem.Addr().Interface())
        columns, _  := rows.Columns()
        onerow      := make([]interface{}, len(columns))
        for index, column := range columns {
            for _, field  := range modelStruct.Fields {
                if field.Name == column {
                    onerow[index] = field.Value.Addr().Interface()
                }
            }
        }
        err := rows.Scan(onerow...)
        if err != nil {
            log.Fatal(err)
        }
        if isSlice {
            if isPtr {
                results.Set(reflect.Append(results, elem.Addr()))
            } else {
                results.Set(reflect.Append(results, elem))
            }
        }
    }
}
