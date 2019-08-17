package postgres

import "github.com/go-jet/jet/internal/jet"

// SelectTable is interface for MySQL sub-queries
type SelectTable interface {
	ReadableTable
	jet.SelectTable
}

type selectTableImpl struct {
	jet.SelectTableImpl
	readableTableInterfaceImpl
}

func newSelectTable(selectStmt jet.StatementWithProjections, alias string) SelectTable {
	subQuery := &selectTableImpl{
		SelectTableImpl: jet.NewSelectTable(selectStmt, alias),
	}

	subQuery.readableTableInterfaceImpl.parent = subQuery

	return subQuery
}
