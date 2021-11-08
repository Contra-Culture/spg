package data

type (
	Schema struct {
		name         string
		pk           []string
		attributes   []string
		associations map[string]*Association
	}
	SchemaCfgr struct {
		graphCfgr *GraphCfgr
		schema    *Schema
	}
	Association struct {
		name             string
		hostSchema       string
		remoteSchema     string
		proxyAssociation string
		limit            int
		orderer          func([]*Object) []*Object
		mapper           func(*Object, map[string]interface{}) bool
	}
)

func (c *SchemaCfgr) PK(pk []string) {
	c.schema.pk = pk
}
func (c *SchemaCfgr) Attribute(n string) {
	for _, a := range c.schema.attributes {
		if a == n {
			return
		}
	}
	c.schema.attributes = append(c.schema.attributes, n)
}
func (c *SchemaCfgr) Association(
	n, rn, pa string,
	l int,
	orderer func([]*Object) []*Object,
	mapper func(*Object, map[string]interface{}) bool,
) {
	c.schema.associations[n] = &Association{
		name:             n,
		hostSchema:       c.schema.name,
		remoteSchema:     rn,
		proxyAssociation: pa,
		limit:            l,
		orderer:          orderer,
		mapper:           mapper,
	}
}
