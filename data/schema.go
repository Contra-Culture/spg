package data

type (
	Schema struct {
		name       string
		id         []string
		attributes []string
		arrows     map[string]*Arrow
	}
	SchemaCfgr struct {
		graphCfgr *GraphCfgr
		schema    *Schema
	}
	Arrow struct {
		name              string
		hostSchema        string
		remoteSchema      string
		attributesMapping map[string]string
	}
)

func (c *SchemaCfgr) ID(id []string) {
	c.schema.id = id
}
func (c *SchemaCfgr) Attribute(n string) {
	for _, a := range c.schema.attributes {
		if a == n {
			return
		}
	}
	c.schema.attributes = append(c.schema.attributes, n)
}
func (c *SchemaCfgr) Arrow(
	n, rn string,
	l int,
	mapping map[string]string,
) {
	c.schema.arrows[n] = &Arrow{
		name:              n,
		hostSchema:        c.schema.name,
		remoteSchema:      rn,
		attributesMapping: mapping,
	}
}
