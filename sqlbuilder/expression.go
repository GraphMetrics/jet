package sqlbuilder

import (
	"github.com/dropbox/godropbox/errors"
)

// An expression
type expression interface {
	clause
	projection
	groupByClause
	orderByClause

	IN(subQuery selectStatement) boolExpression
	NOT_IN(subQuery selectStatement) boolExpression

	AS(alias string) projection
	IS_DISTINCT_FROM(expression expression) boolExpression
	IS_NULL() boolExpression
	ASC() orderByClause
	DESC() orderByClause
}

type expressionInterfaceImpl struct {
	parent expression
}

func (e *expressionInterfaceImpl) IN(subQuery selectStatement) boolExpression {
	return newBinaryBoolExpression(e.parent, subQuery, "IN")
}

func (e *expressionInterfaceImpl) NOT_IN(subQuery selectStatement) boolExpression {
	return newBinaryBoolExpression(e.parent, subQuery, "NOT IN")
}

func (e *expressionInterfaceImpl) AS(alias string) projection {
	return NewAlias(e.parent, alias)
}

func (e *expressionInterfaceImpl) IS_DISTINCT_FROM(expression expression) boolExpression {
	return newBinaryBoolExpression(e.parent, expression, "IS DISTINCT FROM")
}

func (e *expressionInterfaceImpl) IS_NULL() boolExpression {
	return nil
}

func (e *expressionInterfaceImpl) ASC() orderByClause {
	return &orderByClauseImpl{expression: e.parent, ascent: true}
}

func (e *expressionInterfaceImpl) DESC() orderByClause {
	return &orderByClauseImpl{expression: e.parent, ascent: false}
}

func (e *expressionInterfaceImpl) serializeForGroupBy(statement statementType, out *queryData) error {
	return e.parent.serialize(statement, out)
}

func (e *expressionInterfaceImpl) serializeForProjection(statement statementType, out *queryData) error {
	return e.parent.serialize(statement, out)
}

func (e *expressionInterfaceImpl) serializeAsOrderBy(statement statementType, out *queryData) error {
	return e.parent.serialize(statement, out)
}

// Representation of binary operations (e.g. comparisons, arithmetic)
type binaryExpression struct {
	lhs, rhs expression
	operator string
}

func newBinaryExpression(lhs, rhs expression, operator string, parent ...expression) binaryExpression {
	binaryExpression := binaryExpression{
		lhs:      lhs,
		rhs:      rhs,
		operator: operator,
	}

	return binaryExpression
}

func isSimpleOperand(expression expression) bool {
	if _, ok := expression.(*literalExpression); ok {
		return true
	}
	if _, ok := expression.(column); ok {
		return true
	}
	if _, ok := expression.(*numericFunc); ok {
		return true
	}

	return false
}

func (c *binaryExpression) serialize(statement statementType, out *queryData) error {
	if c.lhs == nil {
		return errors.Newf("nil lhs.")
	}
	if c.rhs == nil {
		return errors.Newf("nil rhs.")
	}

	wrap := !isSimpleOperand(c.lhs) && !isSimpleOperand(c.rhs)

	if wrap {
		out.writeString("(")
	}

	if err := c.lhs.serialize(statement, out); err != nil {
		return err
	}

	out.writeString(" " + c.operator + " ")

	if err := c.rhs.serialize(statement, out); err != nil {
		return err
	}

	if wrap {
		out.writeString(")")
	}

	return nil
}

// A not expression which negates a expression value
type prefixExpression struct {
	expression expression
	operator   string
}

func newPrefixExpression(expression expression, operator string) prefixExpression {
	prefixExpression := prefixExpression{
		expression: expression,
		operator:   operator,
	}

	return prefixExpression
}

func (p *prefixExpression) serialize(statement statementType, out *queryData) error {
	out.writeString(p.operator + " ")

	if p.expression == nil {
		return errors.Newf("nil prefix expression.")
	}
	if err := p.expression.serialize(statement, out); err != nil {
		return err
	}

	return nil
}
