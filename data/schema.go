package data

type (
	Schema struct {
		name              string
		pk                []string
		attributes        []string
		arrows            map[string]*Arrow
		pageGeneratorPath []string
	}
	SchemaCfgr struct {
		graph  *Graph
		schema *Schema
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
func (c *SchemaCfgr) HasMany(n string, cfg func(*ArrowCfgr)) {
	arrow := &Arrow{
		name:       n,
		kind:       hasMANY,
		fromSchema: c.schema.name,
	}
	c.schema.arrows[n] = arrow
	cfg(
		&ArrowCfgr{
			schemaCfgr: c,
			arrow:      arrow,
		})
}
func (c *SchemaCfgr) HasOne(n string, cfg func(*ArrowCfgr)) {
	arrow := &Arrow{
		name:       n,
		kind:       hasONE,
		fromSchema: c.schema.name,
	}
	c.schema.arrows[n] = arrow
	cfg(
		&ArrowCfgr{
			schemaCfgr: c,
			arrow:      arrow,
		})
}
func (c *SchemaCfgr) HasManyThrough(n string, cfg func(*ArrowCfgr)) {
	arrow := &Arrow{
		name:       n,
		kind:       hasMANY_THROUGH,
		fromSchema: c.schema.name,
	}
	c.schema.arrows[n] = arrow
	cfg(
		&ArrowCfgr{
			schemaCfgr: c,
			arrow:      arrow,
		})
}
func (c *SchemaCfgr) HasOneThrough(n string, cfg func(*ArrowCfgr)) {
	arrow := &Arrow{
		name:       n,
		kind:       hasONE_THROUGH,
		fromSchema: c.schema.name,
	}
	c.schema.arrows[n] = arrow
	cfg(
		&ArrowCfgr{
			schemaCfgr: c,
			arrow:      arrow,
		})
}
func (c *SchemaCfgr) BelongsTo(n string, cfg func(*ArrowCfgr)) {
	arrow := &Arrow{
		name:       n,
		kind:       belongsTO,
		fromSchema: c.schema.name,
	}
	c.schema.arrows[n] = arrow
	cfg(
		&ArrowCfgr{
			schemaCfgr: c,
			arrow:      arrow,
		})
}
