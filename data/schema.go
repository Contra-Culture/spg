package data

import "fmt"

type (
	Schema struct {
		name       string
		id         []string
		attributes map[string]bool
		arrows     map[string]*Arrow
	}
	SchemaCfgr struct {
		graphCfgr *GraphCfgr
		schema    *Schema
	}
	Arrow struct {
		limit        int
		name         string
		hostSchema   string
		remoteSchema string
		remoteArrow  string
		counterCache string
		mapping      map[string]string
	}
	ArrowCfgr struct {
		arrow      *Arrow
		schemaCfgr *SchemaCfgr
	}
)

func (c *SchemaCfgr) ID(id []string) {
	c.schema.id = id
}
func (c *SchemaCfgr) Attribute(n string) {
	_, exists := c.schema.attributes[n]
	if exists {
		c.graphCfgr.errors = append(
			c.graphCfgr.errors,
			fmt.Errorf("attribute %s.%s already specified", c.schema.name, n))
		return
	}
	c.schema.attributes[n] = true
}
func (c *SchemaCfgr) Arrow(n, rn string, cfg func(*ArrowCfgr)) {
	_, exists := c.schema.arrows[n]
	if exists {
		c.graphCfgr.errors = append(
			c.graphCfgr.errors,
			fmt.Errorf("arrow %s>-%s->... already specified", c.schema.name, n))
		return
	}
	arrow := &Arrow{
		name:         n,
		hostSchema:   c.schema.name,
		remoteSchema: rn,
	}
	cfg(
		&ArrowCfgr{
			arrow:      arrow,
			schemaCfgr: c,
		})
	c.schema.arrows[n] = arrow
}

func (c *ArrowCfgr) MapWith(m map[string]string) {
	if c.arrow.mapping != nil {
		c.schemaCfgr.graphCfgr.errors = append(
			c.schemaCfgr.graphCfgr.errors,
			fmt.Errorf(
				"%s>-%s->%s attributes mapping already specified",
				c.arrow.hostSchema,
				c.arrow.name,
				c.arrow.remoteSchema))
		return
	}
	c.arrow.mapping = m
}
func (c *ArrowCfgr) RemoteArrow(n string) {
	if len(c.arrow.remoteArrow) > 0 {
		c.schemaCfgr.graphCfgr.errors = append(
			c.schemaCfgr.graphCfgr.errors,
			fmt.Errorf(
				"%s>-%s->%s remote arrow already specified",
				c.arrow.hostSchema,
				c.arrow.name,
				c.arrow.remoteSchema))
		return
	}
	c.arrow.remoteArrow = n
}
func (c *ArrowCfgr) Limit(l int) {
	if c.arrow.limit > 0 {
		c.schemaCfgr.graphCfgr.errors = append(
			c.schemaCfgr.graphCfgr.errors,
			fmt.Errorf(
				"%s>-%s->%s limit already specified",
				c.arrow.hostSchema,
				c.arrow.name,
				c.arrow.remoteSchema))
		return
	}
	c.arrow.limit = l
}
func (c *ArrowCfgr) CounterCache(n string) {
	if len(c.arrow.counterCache) > 0 {
		c.schemaCfgr.graphCfgr.errors = append(
			c.schemaCfgr.graphCfgr.errors,
			fmt.Errorf(
				"%s>-%s->%s counter-cache attribute already specified",
				c.arrow.hostSchema,
				c.arrow.name,
				c.arrow.remoteSchema))
		return
	}
	c.arrow.counterCache = n
}
