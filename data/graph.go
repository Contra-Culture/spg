package data

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Contra-Culture/report"
)

type (
	// Graph - stores a meta-data (schemas) and data objects.
	Graph struct {
		nodes   *nodesOrderedSet
		objects map[*node]map[string]*Object
		dataDir string
	}
	GraphCfgr struct {
		graph  *Graph
		report report.Node
	}
)

// New() - creates a new data graph.
func New(rc report.Node, dataDir string, cfg func(*GraphCfgr)) (g *Graph) {
	g = &Graph{
		nodes:   newNodesOrderedSet(),
		objects: map[*node]map[string]*Object{},
		dataDir: dataDir,
	}
	cfgr := &GraphCfgr{
		graph:  g,
		report: rc,
	}
	_, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, 0766)
		if err != nil {
			panic(err)
		}
		path := fmt.Sprintf("%s/schemas", dataDir)
		err = os.Mkdir(path, 0766)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("%s/meta.json", dataDir)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
	}
	if err != nil {
		panic(err)
	}
	defer file.Close()

	cfg(cfgr)
	cfgr.check()
	_, err = file.WriteString(g.JSONString())
	if err != nil {
		panic(err)
	}
	return
}

// .Schema() - specifies a data schema within a data graph.
func (c *GraphCfgr) Node(n string, cfg func(*NodeCfgr)) {
	np := nodePath([]string{n})
	s := &node{
		path:      &np,
		pk:        newNames(nil, nil),
		nodes:     newNodesOrderedSet(),
		links:     newLinks(),
		storePath: strings.Join([]string{c.graph.dataDir, "dataNodes", n}, "/"),
	}
	cfgr := &NodeCfgr{
		graphCfgr: c,
		node:      s,
		report:    c.report.Structure("data node: %s", n),
	}
	cfg(cfgr)
	cfgr.check()
	c.graph.nodes.add(s, func(n []string) {
		c.report.Error("data node \"%s\" already specified", n)
	})
	_, err := os.Stat(s.storePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(s.storePath, 0766)
	}
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/meta.json", s.storePath)
	metaFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		metaFile, err = os.Create(path)
	}
	if err != nil {
		panic(err)
	}
	defer metaFile.Close()
	_, err = metaFile.WriteString(s.JSONString())
	if err != nil {
		panic(err)
	}
	c.graph.objects[s] = map[string]*Object{}
}

// .Get() - returns a data object by its schema name an unique (primary) key.
func (g *Graph) Get(s, id string) (object *Object, err error) {
	node, ok := g.nodes.nodes[s]
	if !ok {
		err = fmt.Errorf("data node \"%s\" does not exist", s)
		return
	}
	objects := g.objects[node]
	object, ok = objects[id]
	if !ok {
		err = fmt.Errorf("object %s[%s] does not exist", s, id)
		return
	}
	return
}

// .Update() - updates data graph's objects repository with the new or updated object.
func (g *Graph) Update(s string, props map[string]string) (pk string, err error) {
	node, ok := g.nodes.nodes[s]
	if !ok {
		err = fmt.Errorf("node \"%s\" does not exist", s)
		return
	}
	objects := g.objects[node]
	o := &Object{
		node:      node,
		updatedAt: time.Now(),
		props:     props,
	}
	pk = o.PK()
	np := node.path.slice()
	objPath := strings.Join([]string{g.dataDir, "nodes", strings.Join(np, "/"), pk}, "/")
	objFile, err := os.OpenFile(objPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		objFile, err = os.Create(objPath)
	}
	if err != nil {
		return
	}
	defer objFile.Close()
	_, err = objFile.WriteString(o.JSONString())
	if err != nil {
		return
	}
	objects[pk] = o
	return
}

// .check() - checks data graph (link) consistency at the very end of data graph configuration.
func (c *GraphCfgr) check() {
	var (
		exists     bool
		remoteNode *node
	)
	for _, nodeName := range c.graph.nodes.order {
		node := c.graph.nodes.nodes[nodeName]
		for _, linkName := range node.links.order {
			link := node.links.links[linkName]
			remoteNode, exists = c.graph.nodes.get(link.remoteNodePath, func(path []string) {
				c.report.Error("node \"%s\" is not specified", link.remoteNodePath)
			})
			if !exists {
				continue
			}
			for _, remoteLinkName := range remoteNode.links.order {
				if remoteLinkName != link.remoteLinkName {
					c.report.Error(
						"%s>-%s->%s>-[ %s ]->... link is not specified",
						nodeName,
						linkName,
						link.remoteNodePath,
						link.remoteLinkName)
				}
			}
			for _, hostAttr := range link.mapping.order {
				remoteAttr := link.mapping.mapping[hostAttr]
				for _, a := range node.nodes.order {
					if a != hostAttr {
						c.report.Error("%s.%s attribute is not specified", nodeName, hostAttr)
					}
				}
				for _, a := range remoteNode.nodes.order {
					if a != remoteAttr {
						c.report.Error("%s.%s attribute is not specified", remoteNode.path, remoteAttr)
					}
				}
			}
		}
	}
}
func (g *Graph) JSONString() string {
	var sb strings.Builder
	sb.WriteString("{\"dataNodes\":")
	sb.WriteString(g.nodes.JSONString())
	sb.WriteString("}")
	return sb.String()
}
