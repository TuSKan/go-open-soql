// Open source implementation of the SOQL parser (Lexical analysis and parsing phase).
package core

import (
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

func Query() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			selectStatement(),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func SelectField() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			complexSelectFieldName(),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func From() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			complexSymbolName(),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func WhereField() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			Trans(
				FlatGroup(
					First(
						selectFieldFunctionCall(),
						complexSymbolName(),
						Error("Unexpected token aheads near by the 'where' clause (unknown operand1)"),
					),
					// Dummy alias name
					Zero(Ast{Value: ""}),
				),
				transComplexSelectFieldName,
			),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func WhereValue() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			Trans(
				FlatGroup(
					// TODO: operators (* / + - ||) and parentheses
					First(
						literalValue(),
						subQuery(),
						selectFieldFunctionCall(),
						complexSymbolName(),
						listValue(),
						Error("Unexpected token aheads near by the 'where' clause (unknown operand2)"),
					),
					// Dummy alias name
					Zero(Ast{Value: ""}),
				),
				transComplexSelectFieldName,
			),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func GroupBy() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			complexSymbolName(),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func OrderBy() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			complexSymbolName(),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}
func HavingField() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			selectFieldFunctionCall(),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}

func HavingValue() ParserFn {
	return FlatGroup(
		Start(),
		sp0(),
		First(
			Trans(
				FlatGroup(
					First(
						literalValue(),
						subQuery(),
						selectFieldFunctionCall(),
						listValue(),
						Error("Unexpected token aheads near by the 'having' clause (unknown operand2)"),
					),
					// Dummy alias name
					Zero(Ast{Value: ""}),
				),
				transComplexSelectFieldName,
			),
			Error("Unexpected token aheads"),
		),
		sp0(),
		End(),
	)
}
