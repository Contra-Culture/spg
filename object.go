package spg

import (
	"fmt"

	"github.com/Contra-Culture/go2html"
)

type (
	arrowKind int
	Schema    struct {
		name       string
		view       *go2html.Template
		attributes []string
		arrows     map[string]*Arrow
	}
	SchemaCfgr struct {
		host   *Host
		schema *Schema
	}
	Arrow struct {
		name          string
		kind          arrowKind
		mapper        func(map[string]interface{}) bool
		fromSchema    *Schema
		throughSchema *Schema
		toSchema      *Schema
	}
	ArrowCfgr struct {
		host  *Host
		arrow *Arrow
		view  *go2html.Template
	}
)

const (
	_ arrowKind = iota
	hasONE
	hasMANY
)

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
		fromSchema: c.schema,
	}
	cfg(&ArrowCfgr{
		host:  c.host,
		arrow: arrow,
	})
}
func (c *SchemaCfgr) HasOne(n string, cfg func(*ArrowCfgr)) {
	arrow := &Arrow{
		name:       n,
		kind:       hasONE,
		fromSchema: c.schema,
	}
	cfg(&ArrowCfgr{
		arrow: arrow,
	})
}
func (c *SchemaCfgr) HasManyThrough(h string, cfg func(*ArrowCfgr)) {

}
func (c *SchemaCfgr) HasOneThrough(h string, cfg func(*ArrowCfgr)) {

}
func (c *SchemaCfgr) FullView(cfg func(*go2html.TemplateConfiguringProxy)) {

}
func (c *SchemaCfgr) ListItemView(cfg func(*go2html.TemplateConfiguringProxy)) {

}
func (c *SchemaCfgr) CardView(cfg func(*go2html.TemplateConfiguringProxy)) {

}
func (c SchemaCfgr) LinkView(cfg func(*go2html.TemplateConfiguringProxy)) {

}
func (c *ArrowCfgr) Schema(n string, mapper func(map[string]interface{}) bool) {
	if c.arrow.mapper != nil {
		panic("mapper already defined")
	}
	schema, ok := c.host.schemas[n]
	if !ok {
		panic("wrong schema name")
	}
	c.arrow.mapper = mapper
	c.arrow.toSchema = schema
}
func (c *ArrowCfgr) Component(cfg func(*go2html.TemplateConfiguringProxy)) {
	if c.arrow.view != nil {
		panic("template already defined")
	}
	template := go2html.NewTemplate(
		fmt.Sprintf(
			"%s>-(%s)->%s",
			c.arrow.fromSchema.name,
			c.arrow.name,
			c.arrow.toSchema.name,
		),
		cfg,
	)
	c.arrow.view = template
}
