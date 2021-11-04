package spg

import (
	"fmt"
	"time"

	"github.com/Contra-Culture/go2html"
)

type (
	orderedMap struct {
		data  map[string]string
		order []string
	}
	arrowKind int
	Schema    struct {
		name       string
		pk         []string
		attributes []string
		arrows     map[string]*Arrow
	}
	SchemaCfgr struct {
		hostCfgr *HostCfgr
		schema   *Schema
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
		view       *go2html.Template
	}
	Repo struct {
		schemas map[string]*Schema
		objects map[string]map[string]interface{} // interface{} is map[string]interface{} or *Object
	}
	RepoCfgr struct {
		hostCfgr *HostCfgr
	}
	Object struct {
		schema     *Schema
		sha        []byte
		updatedAt  time.Time
		attrs      *orderedMap
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

func (c *RepoCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	schema := &Schema{
		name:       n,
		pk:         []string{},
		attributes: []string{},
		arrows:     map[string]*Arrow{},
	}
	c.hostCfgr.host.repo.schemas[n] = schema
	c.hostCfgr.host.repo.objects[n] = map[string]interface{}{}
	cfg(
		&SchemaCfgr{
			hostCfgr: c.hostCfgr,
			schema:   schema,
		})
}
func (r *Repo) Update(s string, attrs map[string]interface{}) {
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
func (c *SchemaCfgr) FullView(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.hostCfgr.host.templates.Add(
		go2html.NewTemplate(c.schema.name, cfg),
		c.schema.FullView(),
	)
}
func (c *SchemaCfgr) ListItemView(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.hostCfgr.host.templates.Add(
		go2html.NewTemplate(c.schema.name, cfg),
		c.schema.ListItemView(),
	)
}
func (c *SchemaCfgr) CardView(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.hostCfgr.host.templates.Add(
		go2html.NewTemplate(c.schema.name, cfg),
		c.schema.CardView(),
	)
}
func (c *SchemaCfgr) LinkView(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.hostCfgr.host.templates.Add(
		go2html.NewTemplate(c.schema.name, cfg),
		c.schema.LinkView(),
	)
}
func (s *Schema) FullView() []string {
	return []string{"schemas", "fullViews", s.name}
}
func (s *Schema) ListItemView() []string {
	return []string{"schemas", "listItemViews", s.name}
}
func (s *Schema) CardView() []string {
	return []string{"schemas", "cardViews", s.name}
}
func (s *Schema) LinkView() []string {
	return []string{"schemas", "linkViews", s.name}
}
func (c *ArrowCfgr) Schema(n string, mapper func(*Object, map[string]interface{}) bool) {
	c.arrow.mapper = mapper
	c.arrow.toSchema = n
	hostCfgr := c.schemaCfgr.hostCfgr
	hostCfgr.checkers = append(
		hostCfgr.checkers,
		func() error {
			_, ok := hostCfgr.host.repo.schemas[n]
			if ok {
				return nil
			}
			return fmt.Errorf("wrong schema name `%s`", n)
		})
}
func (c *ArrowCfgr) ItemView(cfg func(*go2html.TemplateConfiguringProxy)) {
	path := c.arrow.CollectionView()
	template := go2html.NewTemplate(path[len(path)-1], cfg)
	c.schemaCfgr.hostCfgr.host.templates.Add(template, path)
}
func (c *ArrowCfgr) CollectionView(cfg func(*go2html.TemplateConfiguringProxy)) {
	path := c.arrow.CollectionView()
	template := go2html.NewTemplate(path[len(path)-1], cfg)
	c.schemaCfgr.hostCfgr.host.templates.Add(template, path)
}
func (a *Arrow) ItemView() []string {
	name := fmt.Sprintf(
		"%s>-(%s)->%s",
		a.fromSchema,
		a.name,
		a.toSchema,
	)
	return []string{"associations", "itemViews", name}
}
func (a *Arrow) CollectionView() []string {
	name := fmt.Sprintf(
		"%s>-(%s)->%s",
		a.fromSchema,
		a.name,
		a.toSchema,
	)
	return []string{"associations", "collectionViews", name}
}
func (o *Object) Attr(n string) string {
	return o.attrs.data[n]
}
