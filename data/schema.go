package data

import (
	"github.com/Contra-Culture/report"
)

type (
	// Schema represents a type of data objects, like blog posts, arguments, rubrics, categories, etc.
	Schema struct {
		name       string
		id         []string
		attributes map[string]bool
		arrows     map[string]*Arrow
	}
	SchemaCfgr struct {
		graphCfgr *GraphCfgr
		schema    *Schema
		report    report.Node
	}
	// Arrow sets relations between schemas and objects that belong to that schemas.
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
		report     report.Node
	}
)

// .ID() - specifies property names and their order for unique (primary) key calculation for the schema (data) object.
func (c *SchemaCfgr) ID(id []string) {
	c.schema.id = id
}

// .Attribute() - specifies a schema (data) object's property.
func (c *SchemaCfgr) Attribute(n string) {
	_, exists := c.schema.attributes[n]
	if exists {
		c.report.Error("attribute %s.%s already specified", c.schema.name, n)
		return
	}
	c.schema.attributes[n] = true
}

// .Arrow() - specifies a relation between schema objects, when one object has a reference to another one/other ones.
func (c *SchemaCfgr) Arrow(n, rn string, cfg func(*ArrowCfgr)) {
	_, exists := c.schema.arrows[n]
	if exists {
		c.report.Error("arrow %s>-%s->... already specified", c.schema.name, n)
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
			report:     c.report.Structure("arrow: %s", n),
		})
	c.schema.arrows[n] = arrow
}

// .MapWith() - specifies a properties mapping for objects linking.
// For example, blog posts may be linked to a rubric, depending on their rubric_slug and rubric's slug properties
func (c *ArrowCfgr) MapWith(m map[string]string) {
	if c.arrow.mapping != nil {
		c.schemaCfgr.report.Error(
			"%s>-%s->%s attributes mapping already specified",
			c.arrow.hostSchema,
			c.arrow.name,
			c.arrow.remoteSchema)
		return
	}
	c.arrow.mapping = m
}

// .RemoteArrow() - specifies an arrow on remote schema that allows to link objects through another arrow.
func (c *ArrowCfgr) RemoteArrow(n string) {
	if len(c.arrow.remoteArrow) > 0 {
		c.schemaCfgr.report.Error(
			"%s>-%s->%s remote arrow already specified",
			c.arrow.hostSchema,
			c.arrow.name,
			c.arrow.remoteSchema)
		return
	}
	c.arrow.remoteArrow = n
}

// .Limit() - allows to limit a number of linked objects.
func (c *ArrowCfgr) Limit(l int) {
	if c.arrow.limit > 0 {
		c.schemaCfgr.report.Error(
			"%s>-%s->%s limit already specified",
			c.arrow.hostSchema,
			c.arrow.name,
			c.arrow.remoteSchema)
		return
	}
	c.arrow.limit = l
}

// .CounterCache() - allows to specify a field with the count of associated objects.
// The value of object counter is using for limiting a number of linked objects (.Limit()).
func (c *ArrowCfgr) CounterCache(n string) {
	if len(c.arrow.counterCache) > 0 {
		c.schemaCfgr.report.Error(
			"%s>-%s->%s counter-cache attribute already specified",
			c.arrow.hostSchema,
			c.arrow.name,
			c.arrow.remoteSchema)
		return
	}
	c.arrow.counterCache = n
}
