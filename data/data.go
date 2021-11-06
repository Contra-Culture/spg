package data

import (
	"time"
)

type (
	arrowKind int
	Schema    struct {
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
	Graph struct {
		schemas map[string]*Schema
		objects map[string]map[string]interface{} // interface{} is map[string]interface{} or *Object
	}
	GraphCfgr struct {
		graph *Graph
	}
	Object struct {
		schema     *Schema
		sha        []byte
		updatedAt  time.Time
		attrs      map[string]string
		meta       map[string]interface{}
		embeddings map[string]interface{} // interface{} is *Object or []*Object
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

func New(cfg func(*GraphCfgr)) *Graph {
	graph := &Graph{
		schemas: map[string]*Schema{},
		objects: map[string]map[string]interface{}{},
	}
	cfg(
		&GraphCfgr{
			graph: graph,
		})
	return graph
}
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	schema := &Schema{
		name:       n,
		pk:         []string{},
		attributes: []string{},
		arrows:     map[string]*Arrow{},
	}
	c.graph.schemas[n] = schema
	c.graph.objects[n] = map[string]interface{}{}
	cfg(
		&SchemaCfgr{
			graph:  c.graph,
			schema: schema,
		})
}
func (g *Graph) Update(s string, attrs map[string]interface{}) {
	//schema := r.schemas[s]
	//objects := r.objects[s]
}
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
func (c *ArrowCfgr) Schema(n string, mapper func(*Object, map[string]interface{}) bool) {
	c.arrow.mapper = mapper
	c.arrow.toSchema = n
}
func (o *Object) Attr(n string) string {
	return o.attrs[n]
}
