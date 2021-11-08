package gennode

import (
	"github.com/Contra-Culture/spg/data"
)

type (
	Node struct {
		layout     []string
		screen     []string
		mainSchema string
		schemas    map[string]interface{} // interface{} is nil, []string or map[string]interface{}
		path       interface{}            // []string or func(*data.Object) []string
		children   map[string]*Node
	}
	NodeCfgr struct {
		key  string
		node *Node
		err  string
	}
)

func New(cfgs ...func(*NodeCfgr)) *Node {
	var (
		node = &Node{
			children: map[string]*Node{},
			schemas:  map[string]interface{}{},
			path: func(_ *data.Object) []string {
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
		schemas:  map[string]interface{}{},
	}
	cfg(
		&NodeCfgr{
			key:  key,
			node: node,
		})
	c.node.children[key] = node
}
func (c *NodeCfgr) Path(rawPath interface{}) {
	if c.err != "" {
		return
	}
	if c.node.path != nil {
		c.err = "path is already specified"
		return
	}
	switch path := rawPath.(type) {
	case []string:
		c.node.path = path
	case func(*data.Object) []string:
		c.node.path = path
	default:
		c.err = "wrong path, should be of type []string or func(data.Object)[]string"
	}
}
func (c *NodeCfgr) Layout(path []string) {
	if c.err != "" {
		return
	}
	if c.node.layout != nil {
		c.err = "layout is already specified"
		return
	}
	c.node.layout = path
}
func (c *NodeCfgr) Screen(path []string) {
	if c.err != "" {
		return
	}
	if c.node.screen != nil {
		c.err = "screen is already specified"
		return
	}
	c.node.screen = path
}
func (c *NodeCfgr) MainSchema(n string) {
	if c.err != "" {
		return
	}
	if c.node.mainSchema != "" {
		c.err = "main schema is already specified"
		return
	}
	c.node.mainSchema = n
}
func (c *NodeCfgr) Schema(n string, arrows map[string]interface{}) {
	if c.err != "" {
		return
	}
	if c.node.mainSchema != "" {
		c.err = "main schema is already specified"
		return
	}
	c.node.schemas[n] = arrows
}
