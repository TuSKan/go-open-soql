package builder

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
)

type SoqlFields []types.SoqlFieldInfo
type SoqlFrom []types.SoqlObjectInfo
type SoqlWhere []types.SoqlCondition
type SoqlGroupBy []types.SoqlFieldInfo
type SoqlHaving []types.SoqlCondition
type SoqlOrderBy []types.SoqlOrderByInfo

var conditionOpcodeBuilder = map[types.SoqlConditionOpcode]string{
	types.SoqlConditionOpcode_Not:      "NOT",
	types.SoqlConditionOpcode_And:      "AND",
	types.SoqlConditionOpcode_Or:       "OR",
	types.SoqlConditionOpcode_Eq:       "=",
	types.SoqlConditionOpcode_NotEq:    "!=",
	types.SoqlConditionOpcode_Lt:       "<",
	types.SoqlConditionOpcode_Le:       "<=",
	types.SoqlConditionOpcode_Gt:       ">",
	types.SoqlConditionOpcode_Ge:       ">=",
	types.SoqlConditionOpcode_Like:     "LIKE",
	types.SoqlConditionOpcode_NotLike:  "NOTLIKE",
	types.SoqlConditionOpcode_In:       "IN",
	types.SoqlConditionOpcode_NotIn:    "NOTIN",
	types.SoqlConditionOpcode_Includes: "INCLUDES",
	types.SoqlConditionOpcode_Excludes: "EXCLUDES",
}

func soqlObjectInfoBuilder(object types.SoqlObjectInfo) string {
	return strings.Join(object.Name, ".")
}

func soqlConditionOpcodeBuilder(op types.SoqlConditionOpcode) string {
	return conditionOpcodeBuilder[op]
}

func soqlFieldInfoBuilder(field types.SoqlFieldInfo) string {
	switch field.Type {
	case types.SoqlFieldInfo_Field: // field name
		return strings.Join(field.Name, ".")
	case types.SoqlFieldInfo_FieldSet: // fieldset name
		return strings.Join(field.Name, ".")
	case types.SoqlFieldInfo_Function: // function name and parameters
		var params []string
		for i := range field.Parameters {
			params = append(params, soqlFieldInfoBuilder(field.Parameters[i]))
		}
		return fmt.Sprintf("%s(%s)) %s", field.Name[0], strings.Join(params, ", "), field.AliasName)
	case types.SoqlFieldInfo_SubQuery: // SoqlQuery
		return strings.Join(field.Name, ".")
	case types.SoqlFieldInfo_Literal_Null: // nil
		return fmt.Sprintf("'%s'", "nil")
	case types.SoqlFieldInfo_Literal_Int: // int64
		return strconv.FormatInt(field.Value.(int64), 10)
	case types.SoqlFieldInfo_Literal_Float: // float64
		return strconv.FormatFloat(field.Value.(float64), 'f', -1, 64)
	case types.SoqlFieldInfo_Literal_Bool: // bool
		return strconv.FormatBool(field.Value.(bool))
	case types.SoqlFieldInfo_Literal_String: // string
		return fmt.Sprintf("'%s'", field.Value.(string))
	case types.SoqlFieldInfo_Literal_Blob: // []byte
		return fmt.Sprintf("'%s'", string(field.Value.([]byte)))
	case types.SoqlFieldInfo_Literal_Date: // timer.Time
		return field.Value.(time.Time).UTC().Format(time.DateOnly)
	case types.SoqlFieldInfo_Literal_DateTime: // timer.Time
		return field.Value.(time.Time).UTC().Format(time.RFC3339)
	case types.SoqlFieldInfo_Literal_Time: // timer.Time
		return field.Value.(time.Time).UTC().Format(time.TimeOnly)
	case types.SoqlFieldInfo_Literal_DateTimeRange: // SoqlTimeRange
		return strings.Join(field.Name, ".")
	case types.SoqlFieldInfo_Literal_List: // []SoqlListItem
	case types.SoqlFieldInfo_ParameterizedValue: // string
		return strings.Join(field.Name, ".")
	case types.SoqlFieldInfo_DateTimeLiteralName:
		t := field.Value.(types.SoqlDateTimeLiteralName)
		return fmt.Sprintf("%s(%d)", t.Name, t.N)
	}
	return ""
}

func soqlValueInfoParser(value any) (field types.SoqlFieldInfo, err error) {
	switch v := value.(type) {
	case nil:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_Null,
			Value: nil,
		}
	case float32:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_Float,
			Value: v,
		}
	case float64:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_Float,
			Value: v,
		}
	case int32:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_Int,
			Value: v,
		}
	case int64:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_Int,
			Value: v,
		}
	case bool:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_Bool,
			Value: v,
		}
	case time.Time:
		field = types.SoqlFieldInfo{
			Type:  types.SoqlFieldInfo_Literal_DateTime,
			Value: v.UTC(),
		}
	default:
		pv, err := parser.ParseWhereValue(fmt.Sprint(v))
		if err != nil {
			return field, err
		}
		field = *pv
	}
	return field, nil
}
