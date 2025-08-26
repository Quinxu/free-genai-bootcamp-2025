package query

import (
	"database/sql"
	"fmt"
)

// QueryBuilder helps build SQL queries
type QueryBuilder struct {
	base    string
	where   []string
	args    []interface{}
	orderBy string
	limit   int
	offset  int
}

// New creates a new QueryBuilder
func New(base string) *QueryBuilder {
	return &QueryBuilder{
		base:  base,
		where: make([]string, 0),
		args:  make([]interface{}, 0),
	}
}

// Where adds a WHERE condition
func (q *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	q.where = append(q.where, condition)
	q.args = append(q.args, args...)
	return q
}

// OrderBy sets the ORDER BY clause
func (q *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	q.orderBy = orderBy
	return q
}

// Paginate adds LIMIT and OFFSET
func (q *QueryBuilder) Paginate(page, perPage int) *QueryBuilder {
	q.limit = perPage
	q.offset = (page - 1) * perPage
	return q
}

// Build returns the final query and arguments
func (q *QueryBuilder) Build() (string, []interface{}) {
	query := q.base

	if len(q.where) > 0 {
		query += " WHERE " + q.where[0]
		for _, w := range q.where[1:] {
			query += " AND " + w
		}
	}

	if q.orderBy != "" {
		query += " ORDER BY " + q.orderBy
	}

	if q.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", q.limit, q.offset)
	}

	return query, q.args
}

// buildWithoutOrderAndLimit builds the query without ORDER BY / LIMIT/OFFSET
func (q *QueryBuilder) buildWithoutOrderAndLimit() (string, []interface{}) {
	query := q.base

	if len(q.where) > 0 {
		query += " WHERE " + q.where[0]
		for _, w := range q.where[1:] {
			query += " AND " + w
		}
	}

	return query, q.args
}

// Count returns a count query based on the current conditions
func (q *QueryBuilder) Count() (string, []interface{}) {
	inner, args := q.buildWithoutOrderAndLimit()
	countQuery := "SELECT COUNT(*) FROM (" + inner + ") AS t"
	return countQuery, args
}

// Execute runs the query and returns rows
func (q *QueryBuilder) Execute(db *sql.DB) (*sql.Rows, error) {
	query, args := q.Build()
	return db.Query(query, args...)
}

// ExecuteCount runs the count query and returns total
func (q *QueryBuilder) ExecuteCount(db *sql.DB) (int, error) {
	query, args := q.Count()

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

