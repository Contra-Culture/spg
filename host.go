package spg

import (
	"errors"

	"github.com/Contra-Culture/go2html"
	"github.com/Contra-Culture/spg/data"
	"github.com/Contra-Culture/spg/gennode"
)

type (
	Host struct {
		title     string
		host      string
		rootNode  *gennode.Node
		prepared  *Node
		dataGraph *data.Graph
		templates *go2html.TemplateRegistry
	}
	HostCfgr struct {
		host     *Host
		checkers []func() error
		err      error
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
		title:     t,
		host:      h,
		templates: reg,
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
func (c *HostCfgr) Root(cfg func(*gennode.NodeCfgr)) {
	if c.err != nil {
		return
	}
	if c.host.rootNode != nil {
		c.err = errors.New("root is already specified")
		return
	}
	c.host.rootNode = gennode.New(
		func(cfg *gennode.NodeCfgr) {
			cfg.Path(
				func(_ *data.Object) []string {
					return []string{"/"}
				})
		},
		cfg)
}
func (c *HostCfgr) DataGraph(cfg func(*data.GraphCfgr)) {
	dg, err := data.New(cfg)
	if err != nil {
		c.err = err

		return
	}
	c.host.dataGraph = dg
}
func (h *Host) Update(s string, attrs map[string]interface{}) {
	// schema := r.schemas[s]
	// objects := r.objects[s]
}
func (h *Host) Get(path []string) *Page {
	if len(path) == 1 && path[0] == "/" {
		return h.prepared.page
	}
	node := h.prepared
	var ok bool
	for _, chunk := range path[1:] {
		node, ok = node.children[chunk]
		if !ok {
			return nil
		}
	}
	return node.page
}
