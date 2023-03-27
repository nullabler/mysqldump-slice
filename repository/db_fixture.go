package repository

import "mysqldump-slice/entity"

func FillPrimaryKeys(c ConfDbMock) {
	c["PrimaryKeys_user:1"] = NewDbMock()
	c["PrimaryKeys_user:1"].StrList([]string{"id"})

	c["PrimaryKeys_product:1"] = NewDbMock()
	c["PrimaryKeys_product:1"].StrList([]string{"uuid"})

	c["PrimaryKeys_category:1"] = NewDbMock()
	c["PrimaryKeys_category:1"].StrList([]string{"prod_id", "user_id"})
}

func FillValList(c ConfDbMock) {
	c["LoadIds_user:1"] = NewDbMock()
	c["LoadIds_user:1"].ValList([]*entity.Value{
		entity.NewValue("id", "1"),
		entity.NewValue("id", "2"),
		entity.NewValue("id", "3"),
		entity.NewValue("id", "4"),
	})

	c["LoadIds_product:1"] = NewDbMock()
	c["LoadIds_product:1"].ValList([]*entity.Value{
		entity.NewValue("uuid", "xxxx-kkkk-0001"),
		entity.NewValue("uuid", "xxxx-kkkk-0002"),
		entity.NewValue("uuid", "xxxx-kkkk-0003"),
	})

	c["LoadIds_category:1"] = NewDbMock()
	c["LoadIds_category:1"].ValList([]*entity.Value{
		entity.NewValue("user_id", "1"),
		entity.NewValue("prod_id", "7"),
		entity.NewValue("user_id", "2"),
	})
}
