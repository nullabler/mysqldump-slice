package repository

import "mysqldump-slice/entity"

func FillPrimaryKeys(c ConfDbMock) {
	c["PrimaryKeys_user:1"] = NewDbMock()
	c["PrimaryKeys_user:1"].StrList([]string{"id"})

	c["PrimaryKeys_product:1"] = NewDbMock()
	c["PrimaryKeys_product:1"].StrList([]string{"uuid"})

	c["PrimaryKeys_category:1"] = NewDbMock()
	c["PrimaryKeys_category:1"].StrList([]string{"id"})

	c["PrimaryKeys_order:1"] = NewDbMock()
	c["PrimaryKeys_order:1"].StrList([]string{"uuid"})
}

func FillValList(c ConfDbMock) {
	// Users
	c["LoadIds_user:1"] = NewDbMock()
	c["LoadIds_user:1"].ValList([]*entity.Value{
		entity.NewValue("id", "1"),
	})
	c["LoadIds_user:1"].ValList([]*entity.Value{
		entity.NewValue("id", "2"),
	})

	// Products
	c["LoadIds_product:1"] = NewDbMock()
	c["LoadIds_product:1"].ValList([]*entity.Value{
		entity.NewValue("uuid", "pxxx-kkkk-0001"),
	})
	c["LoadIds_product:1"].ValList([]*entity.Value{
		entity.NewValue("uuid", "pxxx-kkkk-0002"),
	})
	c["LoadIds_product:1"].ValList([]*entity.Value{
		entity.NewValue("uuid", "pxxx-kkkk-0003"),
	})
	c["LoadIds_product:1"].ValList([]*entity.Value{
		entity.NewValue("uuid", "pxxx-kkkk-0004"),
	})

	// Categories
	c["LoadIds_category:1"] = NewDbMock()
	c["LoadIds_category:1"].ValList([]*entity.Value{
		entity.NewValue("id", "1"),
	})

	// Orders
	c["LoadIds_order:1"] = NewDbMock()
	c["LoadIds_order:1"].ValList([]*entity.Value{
		entity.NewValue("uuid", "oxxx-kkkk-0001"),
	})
}
