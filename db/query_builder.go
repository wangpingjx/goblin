package db

import (
    "fmt"
    "strings"
    "reflect"
    "strconv"
    // "database/sql"
    // "log"
)

/* SQL组装工具 */
type QueryBuilder struct {
    db                *DB
    tableName         string

    operation         string
    whereConditions   []map[string]interface{}
    selects           string
    limit             int
    order             string
}

// func (qb *QueryBuilder) init() *QueryBuilder {
//     qb.selects = "*"
//     return qb
// }

func (qb *QueryBuilder) Table(name string) *QueryBuilder {
    qb.tableName = name
    return qb
}

func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
    qb.whereConditions = append(qb.whereConditions, map[string]interface{}{"query": query, "args": args})
    return qb
}

func (qb *QueryBuilder) Select(selects string) *QueryBuilder {
    qb.operation = "SELECT"
    qb.selects   = selects
    return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.limit = limit
    return qb
}

// Order("id DESC")
func (qb *QueryBuilder) Order(order string) *QueryBuilder {
    qb.order = order
    return qb
}



/* TODO 暂时默认 query 只会是string，以后再扩充功能 */
/* TODO 类型断言太恶心了，要改 */
func (qb *QueryBuilder) buildWhereCondition(cond map[string]interface{}) (sql string) {
    sql = fmt.Sprintf("(%v)", cond["query"])

    args := cond["args"].([]interface{})
    for _, arg := range args {
        switch reflect.ValueOf(arg).Kind() {
        // 数组 Eg: where("id in (?)", []int(1,2))
        case reflect.Slice:
            values := reflect.ValueOf(arg)
            var tmps []string
            for i := 0; i < values.Len(); i++ {
                tmp := values.Index(i).Interface()
                if tmp, ok := tmp.(string); ok {
                    tmps = append(tmps, "'" + tmp + "'")
                }
                if tmp, ok := tmp.(int); ok {
                    tmps = append(tmps, strconv.Itoa(tmp))
                }
            }
            sql = strings.Replace(sql, "?", strings.Join(tmps, ","), 1)
        default:
            arg = reflect.ValueOf(arg).Interface()
            if arg, ok := arg.(string); ok {
                sql = strings.Replace(sql, "?", "'" + arg + "'", 1)
            }
            if arg, ok := arg.(int); ok {
                sql = strings.Replace(sql, "?", strconv.Itoa(arg), 1)
            }
        }
    }
    return sql
}

func (qb *QueryBuilder) buildOrderSQL() (string) {
    if qb.order != "" {
        return fmt.Sprintf(" ORDER BY %v", qb.order)
    } else {
        return ""
    }
}

func (qb *QueryBuilder) buildLimitSQL() (string) {
    if qb.limit > 0 {
        return fmt.Sprintf(" LIMIT %d", qb.limit)
    } else {
        return ""
    }
}

func (qb *QueryBuilder) buildWhereSQL() (string) {
    var whereConditions []string
    for _, cond := range qb.whereConditions {
        if sql := qb.buildWhereCondition(cond); sql != "" {
            whereConditions = append(whereConditions, sql)
        }
    }
    if len(whereConditions) > 0 {
        return "WHERE " + strings.Join(whereConditions, " AND ")
    } else {
        return ""
    }
}

func (qb *QueryBuilder) buildSelect() (sql string) {
    if "" == qb.selects {
        qb.selects = "*"
    }
    sql = fmt.Sprintf("SELECT %v FROM %v%v%v%v", qb.selects, qb.tableName, qb.buildWhereSQL(), qb.buildOrderSQL(), qb.buildLimitSQL())
    return sql
}

func (qb *QueryBuilder) ToSQL() string {
    return qb.buildSelect()
}
