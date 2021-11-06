package spg

import (
	"github.com/Contra-Culture/go2html"
	"github.com/Contra-Culture/spg/data"
)

type (
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
)

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
