package builder

import (
	"strings"

	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
)

func orderby(input string, desc, nullsLast bool) types.SoqlOrderByInfo {
	field, err := parser.ParseOrderBy(input)
	if err != nil {
		panic(err)
	}
	field.Desc = desc
	field.NullsLast = nullsLast
	return *field
}

func (o SoqlOrderBy) Asc(input string) SoqlOrderBy {
	return append(o, orderby(input, false, false))
}

func (o SoqlOrderBy) Desc(input string) SoqlOrderBy {
	return append(o, orderby(input, true, false))
}
func (o SoqlOrderBy) AscNullsLast(input string) SoqlOrderBy {
	return append(o, orderby(input, false, true))
}

func (o SoqlOrderBy) DescNullsLast(input string) SoqlOrderBy {
	return append(o, orderby(input, true, true))
}

func (o SoqlOrderBy) SOQL(b *strings.Builder) {
	if len(o) > 0 {
		b.WriteString(" ORDER BY ")
		var orderBy []string
		for i := range o {
			order := soqlFieldInfoBuilder(o[i].Field)
			if o[i].Desc {
				order = " DESC"
			}
			if o[i].NullsLast {
				order = " NULLS LAST"
			}
			orderBy = append(orderBy, order)
		}
		b.WriteString(strings.Join(orderBy, ", "))
	}
}
