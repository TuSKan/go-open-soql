package builder

import (
	"fmt"
	"strings"

	"slices"

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
	if len(left) > 0 && len(right) > 0 {
		w = append(w, types.SoqlCondition{Opcode: types.SoqlConditionOpcode_And})
	}
	return w
}

func (w SoqlWhere) Or(left, right SoqlWhere) SoqlWhere {
	w = append(w, left...)
	w = append(w, right...)
	if len(left) > 0 && len(right) > 0 {
		w = append(w, types.SoqlCondition{Opcode: types.SoqlConditionOpcode_Or})
	}
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
	for i := 0; i < len(w); i++ {
		if w[i].Opcode != types.SoqlConditionOpcode_FieldInfo {
			if len(w[:i+1]) == 2 {
				w[i+0] = types.SoqlCondition{
					Opcode: types.SoqlConditionOpcode_FieldInfo,
					Value: types.SoqlFieldInfo{
						Type: types.SoqlFieldInfo_FieldSet,
						Name: []string{fmt.Sprintf(" (%s %s)",
							soqlConditionOpcodeBuilder(w[i].Opcode),
							soqlFieldInfoBuilder(w[i-1].Value))},
					},
				}
				w = slices.Delete(w, i-1, i)
				i--
			} else if len(w[:i+1]) > 2 {
				w[i+0] = types.SoqlCondition{
					Opcode: types.SoqlConditionOpcode_FieldInfo,
					Value: types.SoqlFieldInfo{
						Type: types.SoqlFieldInfo_FieldSet,
						Name: []string{fmt.Sprintf(" (%s %s %s)",
							soqlFieldInfoBuilder(w[i-1].Value),
							soqlConditionOpcodeBuilder(w[i+0].Opcode),
							soqlFieldInfoBuilder(w[i-2].Value))},
					},
				}
				w = slices.Delete(w, i-2, i)
				i = i - 2
			}
		} else if len(w[i:]) >= 1 && w[i+1].Opcode != types.SoqlConditionOpcode_FieldInfo {
			w[i+0] = types.SoqlCondition{
				Opcode: types.SoqlConditionOpcode_FieldInfo,
				Value: types.SoqlFieldInfo{
					Type: types.SoqlFieldInfo_FieldSet,
					Name: []string{fmt.Sprintf(" %s %s",
						soqlConditionOpcodeBuilder(w[i+1].Opcode),
						soqlFieldInfoBuilder(w[i+0].Value))},
				},
			}
			w = slices.Delete(w, i+1, i+2)
		} else if len(w[i:]) >= 2 && w[i+2].Opcode != types.SoqlConditionOpcode_FieldInfo {
			w[i+0] = types.SoqlCondition{
				Opcode: types.SoqlConditionOpcode_FieldInfo,
				Value: types.SoqlFieldInfo{
					Type: types.SoqlFieldInfo_FieldSet,
					Name: []string{fmt.Sprintf(" %s %s %s",
						soqlFieldInfoBuilder(w[i+0].Value),
						soqlConditionOpcodeBuilder(w[i+2].Opcode),
						soqlFieldInfoBuilder(w[i+1].Value))},
				},
			}
			w = slices.Delete(w, i+1, i+3)
		}
	}
	b.WriteString(soqlFieldInfoBuilder(w[0].Value))
}
