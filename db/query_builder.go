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
    selects           string
    whereConditions   []map[string]interface{}
    join              string
    group             string
    having            string
    order             string
    limit             int
    offset            int
}

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

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
    qb.offset = offset
    return qb
}

func (qb *QueryBuilder) Order(order string) *QueryBuilder {
    qb.order = order
    return qb
}

func (qb *QueryBuilder) Join(joinOperator string, tableName string, condition string) *QueryBuilder {
    qb.join = " " + joinOperator + " JOIN " + tableName + " ON " + condition
    return qb
}

func (qb *QueryBuilder) Group(column string) *QueryBuilder {
    qb.group = " GROUP BY " + column
    return qb
}

// Order("id DESC")
func (qb *QueryBuilder) Having(condition string) *QueryBuilder {
    qb.having = " Having " + condition
    return qb
}



/* TODO 暂时默认 query 只会是string，以后再扩充功能 */
/* TODO 类型断言太恶心了，需重构 */
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

func (qb *QueryBuilder) buildWhereSQL() (string) {
    var whereConditions []string
    for _, cond := range qb.whereConditions {
        if sql := qb.buildWhereCondition(cond); sql != "" {
            whereConditions = append(whereConditions, sql)
        }
    }
    whereSQL := ""
    if len(qb.join) > 0 {
        whereSQL += " " + qb.join
    }
    if len(whereConditions) > 0 {
        whereSQL +=  " WHERE " + strings.Join(whereConditions, " AND ")
    }
    if len(qb.group) > 0 {
        whereSQL += " " + qb.group
    }
    if len(qb.having) > 0 {
        whereSQL += " " + qb.having
    }
    if qb.order != "" {
        whereSQL += fmt.Sprintf(" ORDER BY %v", qb.order)
    }
    if qb.limit > 0 {
        whereSQL += fmt.Sprintf(" LIMIT %d", qb.limit)
    }
    return whereSQL
}

func (qb *QueryBuilder) buildSelect() (sql string) {
    if "" == qb.selects {
        qb.selects = "*"
    }
    sql = fmt.Sprintf("SELECT %v FROM %v%v", qb.selects, qb.tableName, qb.buildWhereSQL())
    return sql
}

func (qb *QueryBuilder) ToSQL() string {
    return qb.buildSelect()
}
