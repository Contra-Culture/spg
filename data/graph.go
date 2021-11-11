package data

import (
	"errors"
	"fmt"
	"strings"
)

type (
	Graph struct {
		schemas map[string]*Schema
		objects map[string]map[string]interface{} // interface{} is map[string]interface{} or *Object
	}
	GraphCfgr struct {
		graph  *Graph
		errors []error
	}
)

func New(cfg func(*GraphCfgr)) (g *Graph, err error) {
	g = &Graph{
		schemas: map[string]*Schema{},
		objects: map[string]map[string]interface{}{},
	}
	cfgr := &GraphCfgr{
		graph: g,
	}
	cfg(cfgr)
	err = cfgr.check()
	if err != nil {
		g = nil
	}
	return
}
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	_, exists := c.graph.schemas[n]
	if exists {
		c.errors = append(c.errors, fmt.Errorf("schema \"%s\" already specified", n))
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
		})
}
func (g *Graph) Update(s string, attrs map[string]interface{}) (err error) {
	// schema := r.schemas[s]
	// objects := r.objects[s]
	return
}
func (c *GraphCfgr) check() (err error) {
	var (
		exists       bool
		remoteSchema *Schema
	)
	for schemaName, schema := range c.graph.schemas {
		for arrowName, arrow := range schema.arrows {
			remoteSchema, exists = c.graph.schemas[arrow.remoteSchema]
			if exists {
				c.errors = append(c.errors, fmt.Errorf("schema \"%s\" is not specified", arrow.remoteSchema))
			}
			if len(arrow.remoteArrow) > 0 {
				_, exists = remoteSchema.arrows[arrow.remoteArrow]
				if exists {
					c.errors = append(
						c.errors,
						fmt.Errorf(
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
					c.errors = append(
						c.errors,
						fmt.Errorf("%s.%s attribute is not specified", schemaName, hostAttr))
				}
				_, exists = remoteSchema.attributes[remoteAttr]
				if exists {
					c.errors = append(
						c.errors,
						fmt.Errorf("%s.%s attribute is not specified", remoteSchema.name, remoteAttr))
				}
			}
		}
	}
	if len(c.errors) > 0 {
		var sb strings.Builder
		sb.WriteString("graph specification error\n")
		for _, err := range c.errors {
			sb.WriteRune('\n')
			sb.WriteString(err.Error())
			sb.WriteRune('\n')
		}
		err = errors.New(sb.String())
	}
	return
}
