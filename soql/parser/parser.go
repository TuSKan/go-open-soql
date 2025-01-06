// Open source implementation of the SOQL parser.
package parser

import (
	"errors"
	"strconv"
	"time"

	"github.com/shellyln/go-open-soql-parser/soql/parser/core"
	"github.com/shellyln/go-open-soql-parser/soql/parser/postprocess"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

var (
	queryParser       ParserFn
	selectFieldParser ParserFn
	fromParser        ParserFn
	whereFieldParser  ParserFn
	whereValueParser  ParserFn
	groupByParser     ParserFn
	orderByParser     ParserFn
	havingFieldParser ParserFn
	havingValueParser ParserFn
)

func init() {
	queryParser = core.Query()
	selectFieldParser = core.SelectField()
	fromParser = core.From()
	whereFieldParser = core.WhereField()
	whereValueParser = core.WhereValue()
	groupByParser = core.GroupBy()
	orderByParser = core.OrderBy()
	havingFieldParser = core.HavingField()
	havingValueParser = core.HavingValue()
}

func errorParser(err error, s string, source SourcePosition) error {
	pos := GetLineAndColPosition(s, source, 4)
	return errors.New(
		err.Error() +
			"\n --> Line " + strconv.Itoa(pos.Line) +
			", Col " + strconv.Itoa(pos.Col) + "\n" +
			pos.ErrSource)
}

func parser(s string, fn ParserFn) (any, error) {
	out, err := fn(*NewStringParserContext(s))
	if err == nil && out.MatchStatus != MatchStatus_Matched {
		err = errors.New("Parse failed")
	}
	if err != nil {
		return nil, errorParser(err, s, out.SourcePosition)
	}
	return out.AstStack[0].Value, nil
}

func Parse(s string) (*types.SoqlQuery, error) {
	meta := &types.SoqlQueryMeta{
		Version: "0.9",
		Date:    time.Now().UTC(),
		Source:  s,
	}

	out, err := parser(s, queryParser)
	if err != nil {
		return nil, err
	}

	q := out.(types.SoqlQuery)

	q.Meta = meta

	if err := postprocess.Normalize(&q); err != nil {
		return nil, err
	}

	endDate := time.Now()
	q.Meta.ElapsedTime = endDate.Sub(q.Meta.Date)

	return &q, nil
}

var ParseQuery = Parse

func ParseSelectField(s string) (*types.SoqlFieldInfo, error) {
	if s == "*" {
		s = "FIELDS(ALL)"
	}
	out, err := parser(s, selectFieldParser)
	if err != nil {
		return nil, err
	}
	q := out.(types.SoqlFieldInfo)

	return &q, nil
}

func ParseFrom(s string) (*types.SoqlObjectInfo, error) {
	out, err := parser(s, fromParser)
	if err != nil {
		return nil, err
	}
	q := out.([]string)

	return &types.SoqlObjectInfo{Name: q}, nil
}

func ParseWhereField(s string) (*types.SoqlFieldInfo, error) {
	out, err := parser(s, whereFieldParser)
	if err != nil {
		return nil, err
	}
	q := out.(types.SoqlFieldInfo)

	return &q, nil
}

func ParseWhereValue(s string) (*types.SoqlFieldInfo, error) {
	out, err := parser(s, whereValueParser)
	if err != nil {
		return nil, err
	}
	q := out.(types.SoqlFieldInfo)

	return &q, nil
}

func ParseGroupBy(s string) (*types.SoqlFieldInfo, error) {
	out, err := parser(s, groupByParser)
	if err != nil {
		return nil, err
	}
	q := out.([]string)

	return &types.SoqlFieldInfo{Type: types.SoqlFieldInfo_Field, Name: q}, nil
}

func ParseOrderBy(s string) (*types.SoqlOrderByInfo, error) {
	out, err := parser(s, orderByParser)
	if err != nil {
		return nil, err
	}
	q := out.([]string)

	return &types.SoqlOrderByInfo{Field: types.SoqlFieldInfo{Type: types.SoqlFieldInfo_Field, Name: q}}, nil
}

func ParseHavingField(s string) (*types.SoqlFieldInfo, error) {
	out, err := parser(s, havingFieldParser)
	if err != nil {
		return nil, err
	}
	q := out.(types.SoqlFieldInfo)

	return &q, nil
}

func ParseHavingValue(s string) (*types.SoqlFieldInfo, error) {
	out, err := parser(s, havingValueParser)
	if err != nil {
		return nil, err
	}
	q := out.(types.SoqlFieldInfo)

	return &q, nil
}
