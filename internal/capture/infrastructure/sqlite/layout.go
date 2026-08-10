package sqlite

import "strings"

// column binds one table column to how its value is written from a source and read back into a row. A table's insert,
// select, and scan all derive from a single ordered layout of these, so their column order can never drift apart.
type column[Source any, Row any] struct {
	name  string
	read  func(*Row) any
	write func(Source) any
}

// layout is the ordered columns of one table: the single source for its column list,
// its placeholders, the values to insert, and the targets to scan into.
type layout[Source any, Row any] []column[Source, Row]

func (table layout[Source, Row]) columns() string {
	names := make([]string, len(table))
	for index, column := range table {
		names[index] = column.name
	}
	return strings.Join(names, ", ")
}

func (table layout[Source, Row]) placeholders() string {
	return strings.TrimSuffix(strings.Repeat("?, ", len(table)), ", ")
}

func (table layout[Source, Row]) values(source Source) []any {
	values := make([]any, len(table))
	for index, column := range table {
		values[index] = column.write(source)
	}
	return values
}

func (table layout[Source, Row]) targets(row *Row) []any {
	targets := make([]any, len(table))
	for index, column := range table {
		targets[index] = column.read(row)
	}
	return targets
}
