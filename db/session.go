package db

import (
    "reflect"
    "log"
    "strings"
    "fmt"
    "database/sql"
)

type Session struct {
    db      *DB
    qb      *QueryBuilder    // SQL组装器
    Value   interface{}
    fields  []*Field
    SQL     string
    SQLVars []interface{}
}

// 废了一天功夫终于让它可以接受 *model.struct & *[]*model.struct 两种格式，待重构 TODO
func (s *Session) New(value interface{}) *Session {
    s.Value = value

    reflectValue := reflect.ValueOf(value)
    reflectType  := reflectValue.Type()

    var isSlice bool
    for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
        if reflectType.Kind() == reflect.Slice {
            isSlice = true
        }
        reflectType = reflectType.Elem()
    }
    for i := 0; i < reflectType.NumField(); i++ {
        f := reflectType.Field(i)
        if isSlice {
            reflectValue = reflect.ValueOf(reflect.New(reflectType).Interface())
        }
        field := &Field{ Name: strings.ToLower(f.Name), Tag: f.Tag, Value: reflect.Indirect(reflectValue).FieldByName(f.Name) }
        s.fields = append(s.fields, field)
    }
    s.qb = &QueryBuilder{ db: s.db, tableName: strings.ToLower(reflectType.Name()) }
    return s
}

func (s *Session) Create() *Session{
    // TODO before action
    var (
        columns       []string
        placeholders  []string
    )
    for _, field := range s.fields {
        // 主键不赋值，默认主键id TODO
        if field.Name == "id" {
            continue
        }
        columns      = append(columns, Quote(field.Name))
        placeholders = append(placeholders, "?")
        s.SQLVars    = append(s.SQLVars, field.Value.Interface())
    }
    s.SQL = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", Quote(s.qb.tableName), strings.Join(columns, ","), strings.Join(placeholders, ","))

    if _, err := s.db.db.Exec(s.SQL, s.SQLVars...); err != nil {
        fmt.Printf("%v", err)
    }
    // TODO after action
    return s
}

// 将查询结果渲染到 *Author{} or *[]*Author{} 即 s.Value
// 写完这个方法，感觉身体被掏空，最终完全使用了Gorm的处理方式...
func (s *Session) scan(rows *sql.Rows) {
    var isSlice, isPtr bool
    var resultType reflect.Type

    results := reflect.ValueOf(s.Value)
    for results.Kind() == reflect.Ptr {
        results = results.Elem()
    }

    if kind := results.Kind(); kind == reflect.Slice {
        isSlice = true
        resultType = results.Type().Elem()
        results.Set(reflect.MakeSlice(results.Type(), 0, 0))

        if resultType.Kind() == reflect.Ptr {
            isPtr = true
            resultType = resultType.Elem()
        }
    }
    defer rows.Close()
    for rows.Next() {
        elem := results
        if isSlice {
            elem = reflect.New(resultType).Elem()
        }
        session2 := s.New(elem.Addr().Interface())

        columns, _ := rows.Columns()
        onerow     := make([]interface{}, len(columns))
        for index, column := range columns {
            for _, field := range session2.fields {
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
