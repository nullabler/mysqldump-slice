package entity

type CollectInterface interface {
	PushTable(tabName string)
	Tables() TableList
	PushRelation(rel RelationInterface)
	AllRelList() map[string][]RelationInterface
	RelList(tabName string) []RelationInterface
	PushTab(tabName string)
	PushValList(tabName string, list [][]*Value)
	Tab(tabName string) TabInterface
	PushPkList(tabName string, pkList []string)
	PkList(tabName string) []string
	IsPk(tabName, tabCol string) bool
}

type Collect struct {
	tables  TableList
	relList map[string][]RelationInterface
	tabList map[string]TabInterface
	pkList  map[string][]string
}

func NewCollect() *Collect {
	return &Collect{
		relList: make(map[string][]RelationInterface),
		tabList: make(map[string]TabInterface),
		pkList:  make(map[string][]string),
	}
}

func (c *Collect) PushTable(tabName string) {
	c.tables = append(c.tables, NewTable(tabName))
}

func (c *Collect) Tables() TableList {
	return c.tables
}

func (c *Collect) PushRelation(rel RelationInterface) {
	c.relList[rel.Tab()] = append(c.relList[rel.Tab()], rel)
}

func (c *Collect) AllRelList() map[string][]RelationInterface {
	return c.relList
}

func (c *Collect) RelList(tabName string) []RelationInterface {
	return c.relList[tabName]
}

func (c *Collect) PushTab(tabName string) {
	c.tabList[tabName] = NewTab(tabName)
}

func (c *Collect) PushValList(tabName string, list [][]*Value) {
	if len(list) == 0 {
		return
	}

	for _, valList := range list {
		c.Tab(tabName).Push(valList)
	}
}

func (c *Collect) Tab(tabName string) TabInterface {
	return c.tabList[tabName]
}

func (c *Collect) PushPkList(tabName string, pkList []string) {
	c.pkList[tabName] = pkList
}

func (c *Collect) PkList(tabName string) []string {
	return c.pkList[tabName]
}

func (c *Collect) IsPk(tabName, tabCol string) bool {
	for _, pk := range c.PkList(tabName) {
		if pk == tabCol {
			return true
		}
	}

	return false
}
