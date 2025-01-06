package builder_test

import (
	"testing"

	. "github.com/shellyln/go-open-soql-parser/soql/builder"
	"github.com/shellyln/go-open-soql-parser/soql/parser"
	"github.com/stretchr/testify/assert"
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
			want:    Select("CONCAT(1,2,3,4) integers").From("Contact"),
			wantErr: false,
		},
		{
			name:    "left join 1",
			args:    args{s: `SELECT Id FROM Contact WHERE LastName = 'foo' or Account.Name = 'bar'`},
			want:    Select("Id").From("Contact").Where(Or(Equal("LastName", "'foo'"), Equal("Account.Name", "'bar'"))),
			wantErr: false,
		}, {
			name:    "left join 2",
			args:    args{s: `SELECT Id, Account.Id, Account.Name FROM Contact ORDER BY Account.Name`},
			want:    Select("Id", "Account.Id", "Account.Name").From("Contact").OrderBy(Asc("Account.Name")),
			wantErr: false,
		}, {
			name:    "inner join 1",
			args:    args{s: `SELECT Id FROM Contact WHERE LastName = 'foo' and Account.Name = 'bar'`},
			want:    Select("Id").From("Contact").Where(And(Equal("LastName", "'foo'"), Equal("Account.Name", "'bar'"))),
			wantErr: false,
		}, {
			name:    "inner join 2",
			args:    args{s: `SELECT Id FROM Contact WHERE Account.Name = 'bar'`},
			want:    Select("Id").From("Contact").Where(Equal("Account.Name", "'bar'")),
			wantErr: false,
		}, {
			name:    "fieldset 1",
			args:    args{s: `SELECT FIELDS(all) FROM Contact`},
			want:    Select("*").From("Contact"),
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
				bGot := NewSoqlBuilder(got)
				if !assert.Equal(t, bGot.SoqlQuery, tt.want.SoqlQuery) {
					t.Errorf("\ngot:\t%s \nwant:\t%s", bGot.SOQL(), tt.want.SOQL())
				}
			}
		})
	}
}
