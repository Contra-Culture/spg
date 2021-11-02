package spg

import (
	"fmt"

	"github.com/Contra-Culture/go2html"
)

type (
	PageGenerator struct {
		name     string
		layout   *go2html.Template
		screen   *go2html.Template
		schemas  []string
		children map[string]*PageGenerator
	}
	PageGeneratorCfgr struct {
		hostCfgr *HostCfgr
		pg       *PageGenerator
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
			hostCfgr: c.hostCfgr,
		})
	c.pg.children[name] = pg
}
func (c *PageGeneratorCfgr) Layout(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.pg.layout = go2html.NewTemplate(c.pg.name, cfg)
}
func (c *PageGeneratorCfgr) Screen(cfg func(*go2html.TemplateConfiguringProxy)) {
	c.pg.screen = go2html.NewTemplate(c.pg.name, cfg)
}
func (c *PageGeneratorCfgr) Schema(n string, arrows []string, cfg func(*go2html.TemplateConfiguringProxy)) {
	c.pg.schemas = append(c.pg.schemas, n)
	c.hostCfgr.checkers = append(
		c.hostCfgr.checkers,
		func() error {
			_, ok := c.hostCfgr.host.repo.schemas[n]
			if ok {
				return nil
			}
			return fmt.Errorf("wrong schema `%s`", n)
		})
}
