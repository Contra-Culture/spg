package data

import (
	"fmt"
	"time"

	"github.com/Contra-Culture/report"
)

type (
	// Graph - stores a meta-data (schemas) and data objects.
	Graph struct {
		schemas map[string]*Schema
		objects map[*Schema]map[string]*Object
	}
	GraphCfgr struct {
		graph  *Graph
		report report.Node
	}
)

// New() - creates a new data graph.
func New(rc report.Node, cfg func(*GraphCfgr)) (g *Graph) {
	g = &Graph{
		schemas: map[string]*Schema{},
		objects: map[*Schema]map[string]*Object{},
	}
	cfgr := &GraphCfgr{
		graph:  g,
		report: rc,
	}
	cfg(cfgr)
	cfgr.check()
	return
}

// .Schema() - specifies a data schema within a data graph.
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	_, exists := c.graph.schemas[n]
	if exists {
		c.report.Error("schema \"%s\" already specified", n)
		return
	}
	schema := &Schema{
		name:       n,
		id:         []string{},
		attributes: map[string]bool{},
		arrows:     map[string]*Arrow{},
	}
	c.graph.schemas[n] = schema
	c.graph.objects[schema] = map[string]*Object{}
	cfg(
		&SchemaCfgr{
			graphCfgr: c,
			schema:    schema,
			report:    c.report.Structure("schema: %s", n),
		})
}

// .Get() - returns a data object by its schema name an unique (primary) key.
func (g *Graph) Get(s, id string) (object *Object, err error) {
	schema, ok := g.schemas[s]
	if !ok {
		err = fmt.Errorf("schema \"%s\" does not exist", s)
		return
	}
	objects := g.objects[schema]
	object, ok = objects[id]
	if !ok {
		err = fmt.Errorf("object %s[%s] does not exist", s, id)
		return
	}
	return
}

// .Update() - updates data graph's objects repository with the new or updated object.
func (g *Graph) Update(s string, props map[string]string) (id string, err error) {
	schema, ok := g.schemas[s]
	if !ok {
		err = fmt.Errorf("schema \"%s\" does not exist", s)
		return
	}
	object := &Object{
		schema:    schema,
		updatedAt: time.Now(),
		props:     props,
	}
	objects := g.objects[schema]
	id = object.ID()
	objects[id] = object
	return
}

// .check() - checks data graph (link) consistency at the very end of data graph configuration.
func (c *GraphCfgr) check() {
	var (
		exists       bool
		remoteSchema *Schema
	)
	for schemaName, schema := range c.graph.schemas {
		for arrowName, arrow := range schema.arrows {
			remoteSchema, exists = c.graph.schemas[arrow.remoteSchema]
			if exists {
				c.report.Error("schema \"%s\" is not specified", arrow.remoteSchema)
			}
			if len(arrow.remoteArrow) > 0 {
				_, exists = remoteSchema.arrows[arrow.remoteArrow]
				if exists {
					c.report.Error(
						"%s>-%s->%s>-[ %s ]->... arrow is not specified",
						schemaName,
						arrowName,
						arrow.remoteSchema,
						arrow.remoteArrow)
				}
			}
			for hostAttr, remoteAttr := range arrow.mapping {
				_, exists = schema.attributes[hostAttr]
				if exists {
					c.report.Error("%s.%s attribute is not specified", schemaName, hostAttr)
				}
				_, exists = remoteSchema.attributes[remoteAttr]
				if exists {
					c.report.Error("%s.%s attribute is not specified", remoteSchema.name, remoteAttr)
				}
			}
		}
	}
}
