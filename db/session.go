package db

import (
    "reflect"
    // "log"
    "strings"
    "fmt"
)

type Session struct {
    db     *DB
    qb     *QueryBuilder    // SQL组装器
    // table  *Table        // 目标对象结构信息
    // value  interface{}   // 需要操作的数据
    fields  []*Field
    SQL     string
    SQLVars []interface{}
}

func (s *Session) New(value interface{}) *Session {
    r  := reflect.ValueOf(value)
    rt := r.Type().Elem()
    for i := 0; i < rt.NumField(); i++ {
        f := rt.Field(i)
        field := &Field{ Name: strings.ToLower(f.Name), Tag: f.Tag, Value: reflect.Indirect(r).FieldByName(f.Name) }
        s.fields = append(s.fields, field)
    }
    s.qb = &QueryBuilder{ db: s.db, tableName: strings.ToLower(rt.Name()) }
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
