package builder

import (
	"strings"

	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
)

func (f SoqlFields) SelectCols(cols ...string) (fields SoqlFields) {
	for i := range cols {
		field, err := parser.ParseSelectField(cols[i])
		if err != nil {
			panic(err)
		}
		field.ColumnId = i
		fields = append(fields, *field)
	}
	return
}

func (f SoqlFields) Min(field, alias string) SoqlFields {
	fieldP, err := soqlValueInfoParser(field)
	if err != nil {
		panic(err)
	}
	return append(f, types.SoqlFieldInfo{
		Type:       types.SoqlFieldInfo_Function,
		Name:       []string{"MIN"},
		AliasName:  alias,
		Parameters: []types.SoqlFieldInfo{fieldP},
	})
}

func (f SoqlFields) Max(field, alias string) SoqlFields {
	fieldP, err := soqlValueInfoParser(field)
	if err != nil {
		panic(err)
	}
	return append(f, types.SoqlFieldInfo{
		Type:       types.SoqlFieldInfo_Function,
		Name:       []string{"MAX"},
		AliasName:  alias,
		Parameters: []types.SoqlFieldInfo{fieldP},
	})
}

func (f SoqlFields) Count(field, alias string) SoqlFields {
	fieldP, err := soqlValueInfoParser(field)
	if err != nil {
		panic(err)
	}
	return append(f, types.SoqlFieldInfo{
		Type:       types.SoqlFieldInfo_Function,
		Name:       []string{"COUNT"},
		AliasName:  alias,
		Parameters: []types.SoqlFieldInfo{fieldP},
	})
}

func (f SoqlFields) SOQL(b *strings.Builder) {
	if len(f) > 0 {
		b.WriteString("SELECT ")
		var fields []string
		for i := range f {
			fields = append(fields, soqlFieldInfoBuilder(f[i]))
		}
		b.WriteString(strings.Join(fields, ", "))
	} else {
		b.WriteString("SELECT FIELDS(ALL) ")
	}
}
