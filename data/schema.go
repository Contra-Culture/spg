package data

import (
	"strconv"
	"strings"

	"github.com/Contra-Culture/report"
)

type (
	attributesOrderedSet struct {
		attributes map[string]bool
		order      []string
	}
	arrowsOrderedSet struct {
		arrows map[string]*arrow
		order  []string
	}
	attributesOrderedMapping struct {
		mapping map[string]string
		order   []string
	}
	// Schema represents a type of data objects, like blog posts, arguments, rubrics, categories, etc.
	schema struct {
		name       string
		id         *attributesOrderedSet
		attributes *attributesOrderedSet
		arrows     *arrowsOrderedSet
		absPath    string
	}
	SchemaCfgr struct {
		graphCfgr *GraphCfgr
		schema    *schema
		report    report.Node
	}
	// Arrow sets relations between schemas and objects that belong to that schemas.
	arrow struct {
		limit        int
		name         string
		hostSchema   string
		remoteSchema string
		remoteArrow  string
		counterCache string
		mapping      *attributesOrderedMapping
	}
	ArrowCfgr struct {
		arrow      *arrow
		schemaCfgr *SchemaCfgr
		report     report.Node
	}
)

func newAttributes(attrs []string, handleExistance func(string)) (set *attributesOrderedSet) {
	set = &attributesOrderedSet{
		attributes: map[string]bool{},
		order:      attrs,
	}
	for _, attr := range attrs {
		set.add(attr, handleExistance)
	}
	return
}
func (set *attributesOrderedSet) add(newAttr string, handleExistance func(string)) {
	if set.attributes[newAttr] {
		handleExistance(newAttr)
		return
	}
	set.order = append(set.order, newAttr)
	set.attributes[newAttr] = true
}
func (set *attributesOrderedSet) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, n := range set.order {
		sb.WriteRune('"')
		sb.WriteString(n)
		sb.WriteRune('"')
		if i < len(set.order)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}
func newArrows() *arrowsOrderedSet {
	return &arrowsOrderedSet{
		arrows: map[string]*arrow{},
		order:  []string{},
	}
}
func (set *arrowsOrderedSet) add(a *arrow, handleExistance func(string)) {
	_, exists := set.arrows[a.name]
	if exists {
		handleExistance(a.name)
		return
	}
	set.arrows[a.name] = a
	set.order = append(set.order, a.name)
}
func (set *arrowsOrderedSet) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, n := range set.order {
		a := set.arrows[n]
		sb.WriteString("{\"limit\":\"")
		sb.WriteString(strconv.Itoa(a.limit))
		sb.WriteString("\",\"name\":\"")
		sb.WriteString(a.name)
		sb.WriteString("\",\"hostSchema\":\"")
		sb.WriteString(a.hostSchema)
		sb.WriteString("\",\"remoteSchema\":\"")
		sb.WriteString(a.remoteSchema)
		sb.WriteString("\",\"remoteArrow\":\"")
		sb.WriteString(a.remoteArrow)
		sb.WriteString("\",\"counterCache\":\"")
		sb.WriteString(a.counterCache)
		sb.WriteString("\",\"mapping\":")
		sb.WriteString(a.mapping.JSONString())
		sb.WriteString("\"}")
		if i < len(set.order)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}
func newAttributesOrderedMapping() *attributesOrderedMapping {
	return &attributesOrderedMapping{
		mapping: map[string]string{},
		order:   []string{},
	}
}
func (set *attributesOrderedMapping) add(ha, ra string, handleExistance func(string)) {
	if _, exists := set.mapping[ha]; exists {
		handleExistance(ha)
		return
	}
	set.order = append(set.order, ha)
	set.mapping[ha] = ra
}
func (set *attributesOrderedMapping) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('[')
	lastIdx := len(set.order) - 1
	for i, n := range set.order {
		sb.WriteRune('[')
		sb.WriteString(n)
		sb.WriteRune(',')
		sb.WriteString(set.mapping[n])
		sb.WriteRune(']')
		if i < lastIdx {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}

// .ID() - specifies property names and their order for unique (primary) key calculation for the schema (data) object.
func (c *SchemaCfgr) ID(id []string) {
	c.schema.id = newAttributes(id, func(n string) {
		c.report.Error("wrong ID: %s", n)
	})
}

// .Attribute() - specifies a schema (data) object's property.
func (c *SchemaCfgr) Attribute(n string) {
	c.schema.attributes.add(n, func(n string) {
		c.report.Error("attribute %s.%s already specified", c.schema.name, n)
	})
}

// .Arrow() - specifies a relation between schema objects, when one object has a reference to another one/other ones.
func (c *SchemaCfgr) Arrow(n, rn string, cfg func(*ArrowCfgr)) {
	arrow := &arrow{
		name:         n,
		hostSchema:   c.schema.name,
		remoteSchema: rn,
		mapping:      newAttributesOrderedMapping(),
	}
	cfg(
		&ArrowCfgr{
			arrow:      arrow,
			schemaCfgr: c,
			report:     c.report.Structure("arrow: %s", n),
		})
	c.schema.arrows.add(arrow, func(n string) {
		c.report.Error("arrow %s>-%s->... already specified", c.schema.name, n)
	})
}
func (s *schema) JSONString() string {
	var sb strings.Builder
	sb.WriteString("{\"name\":\"")
	sb.WriteString(s.name)
	sb.WriteString("\",\"id\":")
	sb.WriteString(s.id.JSONString())
	sb.WriteString(",\"attributes\":")
	sb.WriteString(s.attributes.JSONString())
	sb.WriteString(",\"arrows\":")
	sb.WriteString(s.arrows.JSONString())
	sb.WriteString("}")
	return sb.String()
}
func (c *SchemaCfgr) check() {

}

// .Map() - specifies a properties mapping for objects linking.
// For example, blog posts may be linked to a rubric, depending on their rubric_slug and rubric's slug properties
func (c *ArrowCfgr) Map(ha, ra string) {
	c.arrow.mapping.add(ha, ra, func(n string) {
		c.schemaCfgr.report.Error(
			"%s>-%s->%s attributes mapping already specified",
			c.arrow.hostSchema,
			c.arrow.name,
			c.arrow.remoteSchema)
	})
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
