package builder

import (
	"strings"
	"time"

	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/shellyln/go-open-soql-parser/soql/parser/postprocess"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
)

type SoqlBuilder struct {
	types.SoqlQuery
}

func NewSoqlBuilder(q types.SoqlQuery) *SoqlBuilder {
	if q.Meta == nil {
		q.Meta = &types.SoqlQueryMeta{
			Version: "0.9",
			Date:    time.Now().UTC(),
		}
	}
	return &SoqlBuilder{q}
}

func Select(cols ...string) *SoqlBuilder {
	return NewSoqlBuilder(types.SoqlQuery{}).Select(SelectCols(cols...))
}

func (b *SoqlBuilder) Select(fields SoqlFields) *SoqlBuilder {
	b.SoqlQuery.Fields = fields
	return b
}

func (b *SoqlBuilder) From(from ...string) *SoqlBuilder {
	var objects SoqlFrom
	for i := range from {
		obj, err := parser.ParseFrom(from[i])
		if err != nil {
			panic(err)
		}
		objects = append(objects, *obj)
	}
	b.SoqlQuery.From = objects
	return b
}

func (b *SoqlBuilder) FromAlias(from ...Alias) *SoqlBuilder {
	var objects SoqlFrom
	for i := range from {
		obj, err := parser.ParseFrom(from[i].Name)
		if err != nil {
			panic(err)
		}
		obj.AliasName = from[i].Alias
		objects = append(objects, *obj)
	}
	b.SoqlQuery.From = objects
	return b
}

func (b *SoqlBuilder) Where(cond SoqlWhere) *SoqlBuilder {
	b.SoqlQuery.Where = cond
	return b
}

func (b *SoqlBuilder) GroupBy(cols ...string) *SoqlBuilder {
	var fields SoqlFields
	for i := range cols {
		field, err := parser.ParseGroupBy(cols[i])
		if err != nil {
			panic(err)
		}
		fields = append(fields, *field)
	}
	b.SoqlQuery.GroupBy = fields
	return b
}

func (b *SoqlBuilder) OrderBy(orders SoqlOrderBy) *SoqlBuilder {
	b.SoqlQuery.OrderBy = orders
	return b
}

func (b *SoqlBuilder) Having(cond ...types.SoqlCondition) *SoqlBuilder {
	b.SoqlQuery.Having = cond
	return b
}

func (b *SoqlBuilder) Limit(limit int64) *SoqlBuilder {
	b.SoqlQuery.Parent.OffsetAndLimit.Limit = limit
	return b
}

func (b *SoqlBuilder) Offset(offset int64) *SoqlBuilder {
	b.SoqlQuery.OffsetAndLimit.Offset = offset
	return b
}

func (b *SoqlBuilder) Normalize(cond ...types.SoqlCondition) *SoqlBuilder {
	if err := postprocess.Normalize(&b.SoqlQuery); err != nil {
		panic(err)
	}
	b.SoqlQuery.Meta.ElapsedTime = time.Since(b.SoqlQuery.Meta.Date)
	return b
}

func (f SoqlFrom) SOQL(b *strings.Builder) {
	if len(f) > 0 {
		b.WriteString(" FROM ")
		var from []string
		for i := range f {
			from = append(from, soqlObjectInfoBuilder(f[i]))
		}
		b.WriteString(strings.Join(from, ", "))
	}
}

func (g SoqlGroupBy) SOQL(b *strings.Builder) {
	if len(g) > 0 {
		b.WriteString(" GROUP BY ")
		var groupBy []string
		for i := range g {
			groupBy = append(groupBy, soqlFieldInfoBuilder(g[i]))
		}
		b.WriteString(strings.Join(groupBy, ", "))
	}
}

func (b SoqlBuilder) SOQL() string {
	query := strings.Builder{}

	SoqlFields(b.SoqlQuery.Fields).SOQL(&query)

	SoqlFrom(b.SoqlQuery.From).SOQL(&query)

	if len(b.SoqlQuery.Where) > 0 {
		query.WriteString(" WHERE ")
		SoqlWhere(b.SoqlQuery.Where).SOQL(&query)
	}

	SoqlGroupBy(b.SoqlQuery.GroupBy).SOQL(&query)

	SoqlOrderBy(b.SoqlQuery.OrderBy).SOQL(&query)

	// TODO Having, Limit Offset

	return query.String()
}
