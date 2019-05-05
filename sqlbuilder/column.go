// Modeling of columns

package sqlbuilder

import (
	"strings"
)

type Column interface {
	Expression

	Name() string
	TableName() string

	DefaultAlias() Projection
	// Internal function for tracking tableName that a column belongs to
	// for the purpose of serialization
	setTableName(table string)
}

type NullableColumn bool

const (
	Nullable    NullableColumn = true
	NotNullable NullableColumn = false
)

type Collation string

const (
	UTF8CaseInsensitive Collation = "utf8_unicode_ci"
	UTF8CaseSensitive   Collation = "utf8_unicode"
	UTF8Binary          Collation = "utf8_bin"
)

// Representation of MySQL charsets
type Charset string

const (
	UTF8 Charset = "utf8"
)

// The base type for real materialized columns.
type baseColumn struct {
	expressionInterfaceImpl

	name      string
	nullable  NullableColumn
	tableName string
}

func newBaseColumn(name string, nullable NullableColumn, tableName string, parent Column) baseColumn {
	bc := baseColumn{
		name:      name,
		nullable:  nullable,
		tableName: tableName,
	}

	bc.expressionInterfaceImpl.parent = parent

	return bc
}

func (c *baseColumn) Name() string {
	return c.name
}

func (c *baseColumn) TableName() string {
	return c.tableName
}

func (c *baseColumn) setTableName(table string) {
	c.tableName = table
}

func (c *baseColumn) DefaultAlias() Projection {
	return c.As(c.tableName + "." + c.name)
}

func (c baseColumn) Serialize(out *queryData, options ...serializeOption) error {

	setOrderBy := out.statementType == set_statement && out.clauseType == order_by_clause

	if setOrderBy {
		out.WriteString(`"`)
	}

	if c.tableName != "" {
		out.WriteString(c.tableName)
		out.WriteString(".")
	}

	wrapColumnName := strings.Contains(c.name, ".") && !setOrderBy

	if wrapColumnName {
		out.WriteString(`"`)
	}

	out.WriteString(c.name)

	if wrapColumnName {
		out.WriteString(`"`)
	}

	//if contains(options, FOR_PROJECTION) && !contains(options, SKIP_DEFAULT_ALIASING) && c.tableName != "" {
	//	out.WriteString(" AS \"" + c.tableName + "." + c.name + `"`)
	//}

	if setOrderBy {
		out.WriteString(`"`)
	}

	return nil
}

//
//// This is a strict subset of the actual allowed identifiers
//var validIdentifierRegexp = regexp.MustCompile("^[a-zA-Z_]\\w*$")
//
//// Returns true if the given string is suitable as an identifier.
//func validIdentifierName(name string) bool {
//	return validIdentifierRegexp.MatchString(name)
//}

//
//// Pseudo Column type returned by tableName.C(name)
//type deferredLookupColumn struct {
//	isProjection
//	isExpression
//	tableName   *Table
//	colName string
//
//	cachedColumn NonAliasColumn
//}
//
//func (c *deferredLookupColumn) Name() string {
//	return c.colName
//}
//
//func (c *deferredLookupColumn) SerializeSqlForColumnList(
//	out *bytes.Buffer) error {
//
//	return c.Serialize(out)
//}
//
//func (c *deferredLookupColumn) Serialize(out *bytes.Buffer) error {
//	if c.cachedColumn != nil {
//		return c.cachedColumn.Serialize(out)
//	}
//
//	col, err := c.tableName.getColumn(c.colName)
//	if err != nil {
//		return err
//	}
//
//	c.cachedColumn = col
//	return col.Serialize(out)
//}
//
//func (c *deferredLookupColumn) setTableName(tableName string) error {
//	return errors.Newf(
//		"Lookup column '%s' should never have setTableName called on it",
//		c.colName)
//}
//
//func (c *deferredLookupColumn) Eq(rhs Expression) BoolExpression {
//	lit, ok := rhs.(*literalExpression)
//	if ok && sqltypes.Value(lit.value).IsNull() {
//		return newBoolExpression(c, rhs, []byte(" IS "))
//	}
//	return newBoolExpression(c, rhs, []byte(" = "))
//}
//
//func (c *deferredLookupColumn) Gte(rhs Expression) BoolExpression {
//	return Gte(c, rhs)
//}
//
//func (c *deferredLookupColumn) GteLiteral(rhs interface{}) BoolExpression {
//	return Gte(c, Literal(rhs))
//}
//
//func (c *deferredLookupColumn) Lte(rhs Expression) BoolExpression {
//	return Lte(c, rhs)
//}
//
//func (c *deferredLookupColumn) LteLiteral(literal interface{}) BoolExpression {
//	return Lte(c, Literal(literal))
//}
//
//func (c *deferredLookupColumn) Asc() OrderByClause {
//	return sqlbuilder.Asc(c)
//}
//
//func (c *deferredLookupColumn) Desc() OrderByClause {
//	return sqlbuilder.Desc(c)
//}
