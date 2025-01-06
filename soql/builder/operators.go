package builder

type Alias struct {
	Name  string
	Alias string
}

func As(name, alias string) Alias {
	return Alias{Name: name, Alias: alias}
}

func Asc(input string) (o SoqlOrderBy) {
	return o.Asc(input)
}

func Desc(input string) (o SoqlOrderBy) {
	return o.Desc(input)
}

func AscNullsLast(input string) (o SoqlOrderBy) {
	return o.AscNullsLast(input)
}

func DescNullsLast(input string) (o SoqlOrderBy) {
	return o.DescNullsLast(input)
}

func And(left, right SoqlWhere) (w SoqlWhere) {
	return w.And(left, right)
}

func Or(left, right SoqlWhere) (w SoqlWhere) {
	return w.Or(left, right)
}

func Not(right SoqlWhere) (w SoqlWhere) {
	return w.Not(right)
}

func Equal(field string, value string) (w SoqlWhere) {
	return w.Equal(field, value)
}

func NotEqual(field string, value string) (w SoqlWhere) {
	return w.NotEqual(field, value)
}

func Greater(field string, value string) (w SoqlWhere) {
	return w.Greater(field, value)
}

func GreaterThan(field string, value string) (w SoqlWhere) {
	return w.GreaterThan(field, value)
}

func Less(field string, value string) (w SoqlWhere) {
	return w.Less(field, value)
}

func LessThan(field string, value string) (w SoqlWhere) {
	return w.LessThan(field, value)
}

func Like(field string, value string) (w SoqlWhere) {
	return w.Like(field, value)
}

func NotLike(field string, value string) (w SoqlWhere) {
	return w.NotLike(field, value)
}

func In(field string, value string) (w SoqlWhere) {
	return w.In(field, value)
}

func NotIn(field string, value string) (w SoqlWhere) {
	return w.NotIn(field, value)
}

func Includes(field string, value string) (w SoqlWhere) {
	return w.Includes(field, value)
}

func Excludes(field string, value string) (w SoqlWhere) {
	return w.Excludes(field, value)
}
