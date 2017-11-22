package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Table is a table identifier.
type Table struct {
	Statement
	Name  string
	Alias string
}

// NewTable returns a new Table instance.
func NewTable(name string) Table {
	return NewTableAlias(name, "")
}

// NewTableAlias returns a new Table instance with an alias.
func NewTableAlias(name, alias string) Table {
	return Table{
		Name:  name,
		Alias: alias,
	}
}

// As is used to give an alias name to the column.
func (table Table) As(alias string) Table {
	table.Alias = alias
	return table
}

// Write expose statement as a SQL query.
func (table Table) Write(ctx *types.Context) {
	ctx.Write(table.Name)
	if table.Alias != "" {
		ctx.Write(" AS ")
		ctx.Write(table.Alias)
	}
}

// IsEmpty return true if statement is undefined.
func (table Table) IsEmpty() bool {
	return table.Name == ""
}