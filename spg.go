package spg

import (
	"github.com/Contra-Culture/go2html"
	"github.com/Contra-Culture/spg/data"
)

type (
	Host struct {
		title                          string
		host                           string
		rootPageGenerator              *PageGenerator
		rootNode                       *Node
		graph                          *data.Graph
		templates                      *go2html.TemplateRegistry
		schemaPageGeneratorPathMapping map[string][][]string
	}
	HostCfgr struct {
		host     *Host
		checkers []func() error
	}
	PageGenerator struct {
		name                  string
		layout                *go2html.Template
		screen                *go2html.Template
		schemas               []string
		relativePathGenerator func(*data.Object) []string
		children              map[string]*PageGenerator
	}
	PageGeneratorCfgr struct {
		pg       *PageGenerator
		path     []string
		hostCfgr *HostCfgr
	}
	Page struct {
		Headers map[string]string
		Body    string
	}
	Node struct {
		page     *Page
		children map[string]*Node
	}
)

func New(t string, h string, cfg func(*HostCfgr)) *Host {
	reg := go2html.Reg(t)
	reg.Mkdir([]string{"layouts"})
	reg.Mkdir([]string{"screens"})
	reg.Mkdir([]string{"schemas"})
	reg.Mkdir([]string{"schemas", "fullViews"})
	reg.Mkdir([]string{"schemas", "cardViews"})
	reg.Mkdir([]string{"schemas", "listItemViews"})
	reg.Mkdir([]string{"schemas", "linkViews"})
	reg.Mkdir([]string{"associations"})
	reg.Mkdir([]string{"associations", "itemViews"})
	reg.Mkdir([]string{"associations", "collectionViews"})
	host := &Host{
		title:                          t,
		host:                           h,
		templates:                      reg,
		schemaPageGeneratorPathMapping: map[string][][]string{},
	}
	hostCfgr := &HostCfgr{
		host:     host,
		checkers: []func() error{},
	}
	cfg(hostCfgr)
	for _, check := range hostCfgr.checkers {
		err := check()
		if err != nil {
			panic(err)
		}
	}
	return host
}
func (c *HostCfgr) Root(cfg func(*PageGeneratorCfgr)) {
	if c.host.rootPageGenerator != nil {
		panic("root already defined")
	}
	pg := &PageGenerator{
		name:     "root",
		children: map[string]*PageGenerator{},
		relativePathGenerator: func(_ *data.Object) []string {
			return []string{"/"}
		},
	}
	c.host.rootPageGenerator = pg
	cfg(
		&PageGeneratorCfgr{
			pg:       pg,
			path:     []string{"/"},
			hostCfgr: c,
		})
}
func (c *HostCfgr) Repo(cfg func(*data.GraphCfgr)) {
	c.host.graph = data.New(cfg)
}
func (h *Host) Update(s string, attrs map[string]interface{}) {
	//schema := r.schemas[s]
	//objects := r.objects[s]
}
func (h *Host) Get(path []string) *Page {
	if len(path) == 1 && path[0] == "/" {
		return h.rootNode.page
	}
	node := h.rootNode
	var ok bool
	for _, chunk := range path[1:] {
		node, ok = node.children[chunk]
		if !ok {
			return nil
		}
	}
	return node.page
}
func (c *PageGeneratorCfgr) PageGenerator(name string, cfg func(*PageGeneratorCfgr)) {
	pg := &PageGenerator{
		name:     name,
		children: map[string]*PageGenerator{},
	}
	cfg(
		&PageGeneratorCfgr{
			pg:       pg,
			path:     append(c.path, name),
			hostCfgr: c.hostCfgr,
		})
	c.pg.children[name] = pg
}
func (c *PageGeneratorCfgr) RelativePathGenerator(fn func(*data.Object) []string) {
	c.pg.relativePathGenerator = fn
}
func (c *PageGeneratorCfgr) Layout(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.pg.layout = go2html.NewTemplate(c.pg.name, cfg)
}
func (c *PageGeneratorCfgr) Screen(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.pg.screen = go2html.NewTemplate(c.pg.name, cfg)
}
func (c *PageGeneratorCfgr) Schema(n string, arrows []string, cfg func(*go2html.TemplateConfiguringProxy)) {
	c.pg.schemas = append(c.pg.schemas, n)
	host := c.hostCfgr.host
	schemas := append(arrows, n)
	for _, s := range schemas {
		host.schemaPageGeneratorPathMapping[s] = append(
			host.schemaPageGeneratorPathMapping[s],
			c.path,
		)
	}
}
