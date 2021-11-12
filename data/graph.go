package data

import (
	"fmt"

	"github.com/Contra-Culture/report"
)

type (
	Graph struct {
		schemas map[string]*Schema
		objects map[string]map[string]interface{} // interface{} is map[string]interface{} or *Object
	}
	GraphCfgr struct {
		graph  *Graph
		report *report.RContext
	}
)

func New(rc *report.RContext, cfg func(*GraphCfgr)) (g *Graph) {
	g = &Graph{
		schemas: map[string]*Schema{},
		objects: map[string]map[string]interface{}{},
	}
	cfgr := &GraphCfgr{
		graph:  g,
		report: rc,
	}
	cfg(cfgr)
	cfgr.check()
	return
}
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	_, exists := c.graph.schemas[n]
	if exists {
		c.report.Error(fmt.Sprintf("schema \"%s\" already specified", n))
		return
	}
	schema := &Schema{
		name:       n,
		id:         []string{},
		attributes: map[string]bool{},
		arrows:     map[string]*Arrow{},
	}
	c.graph.schemas[n] = schema
	c.graph.objects[n] = map[string]interface{}{}
	cfg(
		&SchemaCfgr{
			graphCfgr: c,
			schema:    schema,
			report:    c.report.Context(fmt.Sprintf("schema: %s", n)),
		})
}
func (g *Graph) Update(s string, attrs map[string]interface{}) (err error) {
	// schema := r.schemas[s]
	// objects := r.objects[s]
	return
}
func (c *GraphCfgr) check() {
	var (
		exists       bool
		remoteSchema *Schema
	)
	for schemaName, schema := range c.graph.schemas {
		for arrowName, arrow := range schema.arrows {
			remoteSchema, exists = c.graph.schemas[arrow.remoteSchema]
			if exists {
				c.report.Error(fmt.Sprintf("schema \"%s\" is not specified", arrow.remoteSchema))
			}
			if len(arrow.remoteArrow) > 0 {
				_, exists = remoteSchema.arrows[arrow.remoteArrow]
				if exists {
					c.report.Error(fmt.Sprintf(
						"%s>-%s->%s>-[ %s ]->... arrow is not specified",
						schemaName,
						arrowName,
						arrow.remoteSchema,
						arrow.remoteArrow))
				}
			}
			for hostAttr, remoteAttr := range arrow.mapping {
				_, exists = schema.attributes[hostAttr]
				if exists {
					c.report.Error(fmt.Sprintf("%s.%s attribute is not specified", schemaName, hostAttr))
				}
				_, exists = remoteSchema.attributes[remoteAttr]
				if exists {
					c.report.Error(fmt.Sprintf("%s.%s attribute is not specified", remoteSchema.name, remoteAttr))
				}
			}
		}
	}
}
