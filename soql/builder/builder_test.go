package builder_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/shellyln/go-open-soql-parser/soql/builder"
	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/shellyln/go-open-soql-parser/soql/parser/types"
)

func TestBuilder(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		want     *SoqlBuilder
		wantErr  bool
		dbgBreak bool
	}{
		{
			name:    "1",
			args:    args{s: `SELECT CONCAT(1,2,3,4) integers FROM Contact`},
			want:    Select("CONCAT(1,2,3,4) integers").From("Contact").Normalize(),
			wantErr: false,
		},
		{
			name:    "left join 1",
			args:    args{s: `SELECT Id FROM Contact WHERE LastName = 'foo' or Account.Name = 'bar'`},
			want:    Select("Id").From("Contact").Where(Or(Equal("LastName", "'foo'"), Equal("Account.Name", "'bar'"))).Normalize(),
			wantErr: false,
		},
		{
			name:    "left join 2",
			args:    args{s: `SELECT Id, Account.Id, Account.Name FROM Contact ORDER BY Account.Name`},
			want:    Select("Id", "Account.Id", "Account.Name").From("Contact").OrderBy(Asc("Account.Name")).Normalize(),
			wantErr: false,
		},
		{
			name:    "inner join 1",
			args:    args{s: `SELECT Id FROM Contact WHERE LastName = 'foo' and Account.Name = 'bar'`},
			want:    Select("Id").From("Contact").Where(And(Equal("LastName", "'foo'"), Equal("Account.Name", "'bar'"))).Normalize(),
			wantErr: false,
		},
		{
			name:    "inner join 2",
			args:    args{s: `SELECT Id FROM Contact WHERE Account.Name = 'bar'`},
			want:    Select("Id").From("Contact").Where(Equal("Account.Name", "'bar'")).Normalize(),
			wantErr: false,
		},
		{
			name:    "fieldset 1",
			args:    args{s: `SELECT FIELDS(ALL) FROM Contact`},
			want:    Select("*").From("Contact").Normalize(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbgBreak {
				t.Log("debug")
			}
			got, err := parser.Parse(tt.args.s)
			if !tt.wantErr && err != nil {
				t.Errorf("%v", err)
			} else {
				if !cmp.Equal(got, &tt.want.SoqlQuery, cmpopts.IgnoreFields(types.SoqlQuery{}, "Meta")) {
					t.Error(cmp.Diff(got, &tt.want.SoqlQuery, cmpopts.IgnoreFields(types.SoqlQuery{}, "Meta")))
				}
			}
		})
	}
}
