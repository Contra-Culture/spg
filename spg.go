package spg

import "github.com/Contra-Culture/go2html"

type (
	Host struct {
		title                          string
		host                           string
		rootPageGenerator              *PageGenerator
		rootNode                       *Node
		repo                           *Repo
		templates                      *go2html.TemplateRegistry
		schemaPageGeneratorPathMapping map[string][][]string
	}
	HostCfgr struct {
		host     *Host
		checkers []func() error
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
		relativePathGenerator: func(_ *Object) []string {
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
func (c *HostCfgr) Repo(cfg func(*RepoCfgr)) {
	repo := &Repo{
		schemas: map[string]*Schema{},
		objects: map[string]map[string]interface{}{},
	}
	c.host.repo = repo
	cfg(
		&RepoCfgr{
			hostCfgr: c,
		})
}
func (h *Host) Update(s string, attrs map[string]interface{}) {
	//schema := r.schemas[s]
	//objects := r.objects[s]
}
