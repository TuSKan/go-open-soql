package builder

import (
	"strings"

	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
)

func binaryConditionOp(input string, op types.SoqlConditionOpcode, val any) (binCond []types.SoqlCondition) {
	field, err := parser.ParseWhereField(input)
	if err != nil {
		panic(err)
	}
	value, err := soqlValueInfoParser(val)
	if err != nil {
		panic(err)
	}
	binCond = append(binCond,
		types.SoqlCondition{
			Opcode: types.SoqlConditionOpcode_FieldInfo,
			Value:  *field,
		},
		types.SoqlCondition{
			Opcode: types.SoqlConditionOpcode_FieldInfo,
			Value:  value,
		},
		types.SoqlCondition{Opcode: op})
	return binCond
}

func (w SoqlWhere) And(left, right SoqlWhere) SoqlWhere {
	w = append(w, left...)
	w = append(w, right...)
	w = append(w, types.SoqlCondition{Opcode: types.SoqlConditionOpcode_And})
	return w
}

func (w SoqlWhere) Or(left, right SoqlWhere) SoqlWhere {
	w = append(w, left...)
	w = append(w, right...)
	w = append(w, types.SoqlCondition{Opcode: types.SoqlConditionOpcode_Or})
	return w
}

func (w SoqlWhere) Not(right SoqlWhere) SoqlWhere {
	w = append(w, right...)
	w = append(w, types.SoqlCondition{Opcode: types.SoqlConditionOpcode_Not})
	return w
}

func (w SoqlWhere) Equal(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Eq, value)...)
}

func (w SoqlWhere) NotEqual(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_NotEq, value)...)
}

func (w SoqlWhere) Greater(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Ge, value)...)
}

func (w SoqlWhere) GreaterThan(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Gt, value)...)
}

func (w SoqlWhere) Less(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Le, value)...)
}

func (w SoqlWhere) LessThan(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Lt, value)...)
}

func (w SoqlWhere) Like(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Like, value)...)
}

func (w SoqlWhere) NotLike(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_NotLike, value)...)
}

func (w SoqlWhere) In(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_In, value)...)
}

func (w SoqlWhere) NotIn(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_NotIn, value)...)
}

func (w SoqlWhere) Includes(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Includes, value)...)
}

func (w SoqlWhere) Excludes(field string, value string) SoqlWhere {
	return append(w, binaryConditionOp(field, types.SoqlConditionOpcode_Excludes, value)...)
}

func (w SoqlWhere) SOQL(b *strings.Builder) {
	if len(w) >= 2 {
		switch w[len(w)-1].Opcode {
		case types.SoqlConditionOpcode_And:
			SoqlWhere(w[:len(w)/2]).SOQL(b)
			b.WriteString(" AND ")
			SoqlWhere(w[len(w)/2 : len(w)-1]).SOQL(b)
		case types.SoqlConditionOpcode_Or:
			SoqlWhere(w[:len(w)/2]).SOQL(b)
			b.WriteString(" OR ")
			SoqlWhere(w[len(w)/2 : len(w)-1]).SOQL(b)
		case types.SoqlConditionOpcode_Not:
			b.WriteString(" NOT ")
			SoqlWhere(w[:len(w)-1]).SOQL(b)
		default:
			b.WriteString(" " + soqlFieldInfoBuilder(w[len(w)-3].Value))
			b.WriteString(" " + soqlConditionOpcodeBuilder(w[len(w)-1].Opcode))
			b.WriteString(" " + soqlFieldInfoBuilder(w[len(w)-2].Value))
		}
	}
}
