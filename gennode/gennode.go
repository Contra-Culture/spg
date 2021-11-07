package gennode

import (
	"github.com/Contra-Culture/go2html"
	"github.com/Contra-Culture/spg/data"
)

type (
	Node struct {
		layout                *go2html.Template
		screen                *go2html.Template
		schemas               []string
		relativePathGenerator func(*data.Object) []string
		children              map[string]*Node
	}
	NodeCfgr struct {
		key  string
		node *Node
	}
)

func New(cfgs ...func(*NodeCfgr)) *Node {
	var (
		node = &Node{
			children: map[string]*Node{},
			relativePathGenerator: func(_ *data.Object) []string {
				return []string{"/"}
			},
		}
		nodeCfgr = &NodeCfgr{
			node: node,
		}
	)
	for _, cfg := range cfgs {
		cfg(nodeCfgr)
	}
	return node
}

func (c *NodeCfgr) Node(key string, cfg func(*NodeCfgr)) {
	node := &Node{
		children: map[string]*Node{},
	}
	cfg(
		&NodeCfgr{
			key:  key,
			node: node,
		})
	c.node.children[key] = node
}
func (c *NodeCfgr) RelativePathGenerator(fn func(*data.Object) []string) {
	c.node.relativePathGenerator = fn
}
func (c *NodeCfgr) Layout(path []string) {
}
func (c *NodeCfgr) Screen(path []string) {

}
func (c *NodeCfgr) MainSchema(n string) {

}
func (c *NodeCfgr) Schema(n string, arrows []string) {
	c.node.schemas = append(c.node.schemas, n)
}
