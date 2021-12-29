package data

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Contra-Culture/report"
)

type (
	// Graph - stores a meta-data (schemas) and data objects.
	Graph struct {
		schemas *schemasOrderedSet
		objects map[*schema]map[string]*Object
		dataDir string
	}
	GraphCfgr struct {
		graph  *Graph
		report report.Node
	}
	schemasOrderedSet struct {
		schemas map[string]*schema
		order   []string
	}
)

func newSchemasOrderedSet() *schemasOrderedSet {
	return &schemasOrderedSet{
		schemas: map[string]*schema{},
		order:   []string{},
	}
}
func (set *schemasOrderedSet) add(s *schema, handleExistence func(string)) {
	if _, exists := set.schemas[s.name]; exists {
		handleExistence(s.name)
		return
	}
	set.order = append(set.order, s.name)
	set.schemas[s.name] = s
}
func (set *schemasOrderedSet) JSONString() string {
	var sb strings.Builder
	sb.WriteString("{")
	lastIdx := len(set.schemas) - 1
	for i, schemaName := range set.order {
		schema := set.schemas[schemaName]
		sb.WriteRune('"')
		sb.WriteString(schemaName)
		sb.WriteString("\":")
		sb.WriteString(schema.JSONString())
		if i < lastIdx {
			sb.WriteRune(',')
		}
	}
	sb.WriteString("}")
	return sb.String()
}

// New() - creates a new data graph.
func New(rc report.Node, dataDir string, cfg func(*GraphCfgr)) (g *Graph) {
	g = &Graph{
		schemas: newSchemasOrderedSet(),
		objects: map[*schema]map[string]*Object{},
		dataDir: dataDir,
	}
	cfgr := &GraphCfgr{
		graph:  g,
		report: rc,
	}
	_, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, 0766)
		if err != nil {
			panic(err)
		}
		path := fmt.Sprintf("%s/schemas", dataDir)
		err = os.Mkdir(path, 0766)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("%s/meta.json", dataDir)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
	}
	if err != nil {
		panic(err)
	}
	defer file.Close()

	cfg(cfgr)
	cfgr.check()
	_, err = file.WriteString(g.JSONString())
	if err != nil {
		panic(err)
	}
	return
}

// .Schema() - specifies a data schema within a data graph.
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	s := &schema{
		name:       n,
		id:         newAttributes(nil, nil),
		attributes: newAttributes(nil, nil),
		arrows:     newArrows(),
		absPath:    strings.Join([]string{c.graph.dataDir, "schemas", n}, "/"),
	}
	cfgr := &SchemaCfgr{
		graphCfgr: c,
		schema:    s,
		report:    c.report.Structure("schema: %s", n),
	}
	cfg(cfgr)
	cfgr.check()
	c.graph.schemas.add(s, func(n string) {
		c.report.Error("schema \"%s\" already specified", n)
	})
	_, err := os.Stat(s.absPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(s.absPath, 0766)
	}
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/meta.json", s.absPath)
	metaFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		metaFile, err = os.Create(path)
	}
	if err != nil {
		panic(err)
	}
	defer metaFile.Close()
	_, err = metaFile.WriteString(s.JSONString())
	if err != nil {
		panic(err)
	}
	c.graph.objects[s] = map[string]*Object{}
}

// .Get() - returns a data object by its schema name an unique (primary) key.
func (g *Graph) Get(s, id string) (object *Object, err error) {
	schema, ok := g.schemas.schemas[s]
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
	schema, ok := g.schemas.schemas[s]
	if !ok {
		err = fmt.Errorf("schema \"%s\" does not exist", s)
		return
	}
	objects := g.objects[schema]
	o := &Object{
		schema:    schema,
		updatedAt: time.Now(),
		props:     props,
	}
	id = o.ID()
	objPath := strings.Join([]string{g.dataDir, "schemas", schema.name, id}, "/")
	objFile, err := os.OpenFile(objPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		objFile, err = os.Create(objPath)
	}
	if err != nil {
		return
	}
	defer objFile.Close()
	_, err = objFile.WriteString(o.JSONString())
	if err != nil {
		return
	}
	objects[id] = o
	return
}

// .check() - checks data graph (link) consistency at the very end of data graph configuration.
func (c *GraphCfgr) check() {
	var (
		exists       bool
		remoteSchema *schema
	)
	for _, schemaName := range c.graph.schemas.order {
		schema := c.graph.schemas.schemas[schemaName]
		for _, arrowName := range schema.arrows.order {
			arrow := schema.arrows.arrows[arrowName]
			remoteSchema, exists = c.graph.schemas.schemas[arrow.remoteSchema]
			if !exists {
				c.report.Error("schema \"%s\" is not specified", arrow.remoteSchema)
				continue
			}
			for _, remoteArrowName := range remoteSchema.arrows.order {
				if remoteArrowName != arrow.remoteArrow {
					c.report.Error(
						"%s>-%s->%s>-[ %s ]->... arrow is not specified",
						schemaName,
						arrowName,
						arrow.remoteSchema,
						arrow.remoteArrow)
				}
			}
			for _, hostAttr := range arrow.mapping.order {
				remoteAttr := arrow.mapping.mapping[hostAttr]
				for _, a := range schema.attributes.order {
					if a != hostAttr {
						c.report.Error("%s.%s attribute is not specified", schemaName, hostAttr)
					}
				}
				for _, a := range remoteSchema.attributes.order {
					if a != remoteAttr {
						c.report.Error("%s.%s attribute is not specified", remoteSchema.name, remoteAttr)
					}
				}
			}
		}
	}
}
func (g *Graph) JSONString() string {
	var sb strings.Builder
	sb.WriteString("{\"schemas\":")
	sb.WriteString(g.schemas.JSONString())
	sb.WriteString("}")
	return sb.String()
}
