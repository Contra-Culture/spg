package spg

import (
	"github.com/Contra-Culture/go2html/registry"
	"github.com/Contra-Culture/report"
	"github.com/Contra-Culture/spg/data"
	"github.com/Contra-Culture/spg/node"
)

type (
	// Host - is a container for a whole website.
	Host struct {
		title     string
		host      string
		rootNode  *node.Node
		prepared  *Node
		dataGraph *data.Graph
		templates registry.Registry
	}
	HostCfgr struct {
		host     *Host
		checkers []func() error
		report   report.Node
	}
)

// New() - creates new host (website project).
func New(t string, h string, cfg func(*HostCfgr)) *Host {
	reg := registry.New()
	reg.Mkdir([]string{"layouts"})
	reg.Mkdir([]string{"screens"})
	reg.Mkdirf(
		[]string{"schemas"},
		func(dir registry.Registry) {
			dir.Mkdir([]string{"fullViews"})
			dir.Mkdir([]string{"cardViews"})
			dir.Mkdir([]string{"listItemViews"})
			dir.Mkdir([]string{"linkViews"})
		})
	reg.Mkdirf(
		[]string{"associations"},
		func(dir registry.Registry) {
			dir.Mkdir([]string{"itemViews"})
			dir.Mkdir([]string{"collectionViews"})
		})
	host := &Host{
		title:     t,
		host:      h,
		templates: reg,
	}
	hostCfgr := &HostCfgr{
		host:     host,
		checkers: []func() error{},
		report:   report.New("host configuring"),
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

// .Root() - specifies a host's top-level (root) node (page generator).
func (c *HostCfgr) Root(cfg func(*node.NodeCfgr)) {
	if c.host.rootNode != nil {
		c.report.Error("root is already specified")
		return
	}
	c.host.rootNode = node.New(
		func(cfg *node.NodeCfgr) {
			cfg.Path(
				func(_ *data.Object) []string {
					return []string{"/"}
				})
		},
		cfg)
}

// .DataGraph() - specifies a data graph for the host.
func (c *HostCfgr) DataGraph(dataDir string, cfg func(*data.GraphCfgr)) {
	if c.host.dataGraph != nil {
		c.report.Error("root is already specified")
		return
	}
	c.host.dataGraph = data.New(c.report.Structure("data-graph"), dataDir, cfg)
}

// .Update() - updates host with a new or updated data object and re-renders related pages.
func (h *Host) Update(s string, props map[string]string) (string, error) {
	return h.dataGraph.Update(s, props)
}

// .Get() - provides a generated page by its path.
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
