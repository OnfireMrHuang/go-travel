package sqlparser

import (
	"testing"
	//"vitess.io/vitess/go/vt/sqlparser"
)

func TestSqlParse(t *testing.T) {
	sql := "select name, age, sex, class from user where name='bob' and age=10 order by class offset 0 limit 10"
	t.Log(sql)
}
