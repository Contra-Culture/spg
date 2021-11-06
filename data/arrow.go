package data

type (
	arrowKind int

	Arrow struct {
		name          string
		kind          arrowKind
		mapper        func(*Object, map[string]interface{}) bool
		fromSchema    string
		throughSchema string
		toSchema      string
	}
	ArrowCfgr struct {
		schemaCfgr *SchemaCfgr
		arrow      *Arrow
	}
)

const (
	_ arrowKind = iota
	hasONE
	hasMANY
	hasONE_THROUGH
	hasMANY_THROUGH
	belongsTO
)

func (c *ArrowCfgr) Schema(n string, mapper func(*Object, map[string]interface{}) bool) {
	c.arrow.mapper = mapper
	c.arrow.toSchema = n
}
