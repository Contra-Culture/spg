package spg

import (
	"fmt"

	"github.com/Contra-Culture/go2html"
)

type (
	PageGenerator struct {
		name                  string
		layout                *go2html.Template
		screen                *go2html.Template
		schemas               []string
		relativePathGenerator func(*Object) []string
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
func (c *PageGeneratorCfgr) RelativePathGenerator(fn func(*Object) []string) {
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
	c.hostCfgr.checkers = append(
		c.hostCfgr.checkers,
		func() error {
			for _, s := range schemas {
				switch c.hostCfgr.host.repo.schemas[s] {
				case nil:
					return fmt.Errorf("wrong schema `%s`", s)
				default:
					return nil
				}
			}
			return nil
		})
}
