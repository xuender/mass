package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuender/mass/app"
)

func TestSql_Where(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	tab := app.NewTable("select * from user where a=1")

	ass.Equal("a=1", tab.Where())

	tab = app.NewTable("select * from user ")
	ass.Equal("", tab.Where())
}

func TestGetTable(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	tab := app.NewTable("select * from user where a=1")

	ass.Equal("user", tab.Table())

	tab = app.NewTable("update user set b=2 where a=1")
	ass.Equal("user", tab.Table())
}

func TestIteratorWhere(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	pks := []string{"k1"}
	row := []any{3}

	ass.Equal(" AND (`k1` > 3 OR `k1` = 3)", app.IteratorWhere(pks, row))

	row = []any{"a"}
	ass.Equal(" AND (`k1` > 'a' OR `k1` = 'a')", app.IteratorWhere(pks, row))
}

func TestIteratorWhere_2(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	pks := []string{"k1", "k2"}
	row := []any{3, "a"}

	ass.Equal(" AND (`k1` > 3 OR (`k1` = 3 AND (`k2` > 'a' OR `k2` = 'a')))", app.IteratorWhere(pks, row))
}

func TestIteratorWhere_3(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	pks := []string{"k1", "k2", "k3"}
	row := []any{3, "a", 9}
	// nolint: lll
	ass.Equal(" AND (`k1` > 3 OR (`k1` = 3 AND (`k2` > 'a' OR (`k2` = 'a' AND (`k3` > 9 OR `k3` = 9)))))", app.IteratorWhere(pks, row))
}
