package goblinDB

import (
    "log"
)

// const CommaSpace = ", "

type QueryBuilder struct {
    tableName   string

    whereConds  []map[string]interface{}
    selects     map[string]interface{}
    limit       interface{}
}

func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
    qb.whereConds = append(qb.whereConds, map[string]interface{}{"query": query, "args": args})
    return qb
}

func (qb *QueryBuilder) Select(query interface{}, args ...interface{}) *QueryBuilder {
    qb.selects = map[string]interface{}{"query": query, "args":args }
    return qb
}

func (qb *QueryBuilder) Limit(limit interface{}) *QueryBuilder {
    qb.limit = limit
    return qb
}

func (qb *QueryBuilder) Table(name string) *QueryBuilder {
    qb.tableName = name
    return qb
}

func (qb *QueryBuilder) ToString() string {
    log.Println("=> tableName: " + qb.tableName)
    log.Printf("=> WhereConds: %v", qb.whereConds)
    return "I'm a sql string"
}
