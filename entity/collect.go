package entity

type CollectInterface interface {
	PushTable(string)
	Tables() TableList
	PushRelation(RelationInterface)
	RelList(string) []RelationInterface
	PushTab(string)
	PushKeyList(string, string, []string)
	Tab(string) TabInterface
	PushPkList(string, []string)
	PkList(string) []string
	IsPk(string, string) bool
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

func (c *Collect) RelList(tabName string) []RelationInterface {
	return c.relList[tabName]
}

func (c *Collect) PushTab(tabName string) {
	c.tabList[tabName] = NewTab(tabName)
}

func (c *Collect) PushKeyList(tab, col string, list []string) {
	for _, key := range list {
		//fmt.Println(c.tabList[tab], tab, col, key)
		c.tabList[tab].Push(col, key)
	}
}

func (c *Collect) Tab(tab string) TabInterface {
	return c.tabList[tab]
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
