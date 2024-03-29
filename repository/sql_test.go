package repository

import (
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"testing"
)

func TestWhere(t *testing.T) {
	c := &config.Conf{}
	s := NewSql(c)

	valList := []*entity.Value{
		entity.NewValue("id", "1"),
	}

	rowList := []*entity.Row{}
	rowList = append(rowList, entity.NewRow(valList))

	list := s.Where(rowList, false)
	if len(list) != 1 {
		t.Error("Fail count result")
	}

	exp := "(`id` = '1')"
	if list[0] != exp {
		t.Errorf("Fail result from where func Exp: %s Got: %s", exp, list[0])
	}

	valList = append(valList, entity.NewValue("category_id", "4"))
	valList = append(valList, entity.NewValue("filter_id", "7"))

	rowList = append(rowList, entity.NewRow(valList))
	list = s.Where(rowList, false)
	if len(list) != 2 {
		t.Error("Fail count result")
	}

	exp = "(`id` = '1' AND `category_id` = '4' AND `filter_id` = '7')"
	if list[1] != exp {
		t.Errorf("Fail result from where func Exp: %s Got: %s", exp, list[1])
	}

	list = s.Where(rowList, true)

	exp = "(\\`id\\` = '1')"
	if list[0] != exp {
		t.Errorf("Fail result from where func Exp: %s Got: %s", exp, list[0])
	}
}

func TestWhereSlice(t *testing.T) {
	c := &config.Conf{
		LimitCli: 2,
	}
	s := NewSql(c)

	rowList := []*entity.Row{}
	for _, i := range []string{"3", "5", "6", "7", "8", "2", "9"} {
		valList := []*entity.Value{
			entity.NewValue("uuid", "xxxx-rrrr-"+i),
			entity.NewValue("cat_id", i),
		}

		rowList = append(rowList, entity.NewRow(valList))
	}

	list := s.WhereSlice(rowList, false)
	if len(list) != 4 {
		t.Error("Fail count result")
	}

	exp := "(`uuid` = 'xxxx-rrrr-3' AND `cat_id` = '3') OR (`uuid` = 'xxxx-rrrr-5' AND `cat_id` = '5')"
	if list[0] != exp {
		t.Errorf("Fail result from where func Exp: %s Got: %s", exp, list[0])
	}

}

func TestQueryLoadIds(t *testing.T) {
	c := &config.Conf{}
	s := NewSql(c)

	specs := new(config.Specs)
	prKeyList := []string{}

	_, err := s.QueryLoadIds(
		"user",
		specs,
		prKeyList,
	)

	if err == nil {
		t.Error("Should be set prKeyList field for Sort")
	}

	prKeyList = []string{"id"}
	exp := "SELECT `id` FROM `user` ORDER BY `id` DESC LIMIT 0"
	got, er := s.QueryLoadIds(
		"user",
		specs,
		prKeyList,
	)
	if er != nil {
		t.Error(er)
	}

	if exp != got {
		t.Errorf("Fail for QueryLoadIds Exp: %s Got: %s", exp, got)
	}

	specs = &config.Specs{
		Condition: "`updated_at` > NOW()",
		Sort:      []string{"updated_at", "created_at"},
		Limit:     3,
	}

	exp = "SELECT `id` FROM `user` WHERE `updated_at` > NOW() ORDER BY `updated_at`, `created_at` DESC LIMIT 3"
	got, er = s.QueryLoadIds(
		"user",
		specs,
		prKeyList,
	)
	if er != nil {
		t.Error(er)
	}

	if exp != got {
		t.Errorf("Fail for QueryLoadIds Exp: %s Got: %s", exp, got)
	}
}

func TestWrapKeys(t *testing.T) {
	c := &config.Conf{}
	s := NewSql(c)

	keys := []string{"test", "word"}
	list := s.wrapKeys(keys, "#")
	for i, item := range list {
		exp := "#" + keys[i] + "#"
		if item != exp {
			t.Errorf("Fail for wrapKeys Exp: %s Got: %s", exp, item)
		}
	}
}
