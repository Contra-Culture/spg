package spg

import "github.com/Contra-Culture/go2html"

type (
	Host struct {
		title     string
		host      string
		root      *PageGenerator
		schemas   map[string]*Schema
		templates *go2html.TemplateRegistry
	}
	HostCfgr struct {
		host *Host
	}
)

func New(t string, h string, configure func(*HostCfgr)) *Host {
	host := &Host{
		title:     t,
		host:      h,
		schemas:   map[string]*Schema{},
		templates: go2html.Reg(t),
	}
	configure(&HostCfgr{host})
	return host
}
func (c *HostCfgr) Root(cfg func(*PageGeneratorCfgr)) {
	if c.host.root != nil {
		panic("root already defined")
	}
	pg := &PageGenerator{}
	cfg(&PageGeneratorCfgr{
		pageGenerator: pg,
	})
	c.host.root = pg
}
func (c *HostCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	schema := &Schema{
		name: n,
	}
	cfg(&SchemaCfgr{
		host:   c.host,
		schema: schema,
	})
	c.host.schemas[n] = schema
}
