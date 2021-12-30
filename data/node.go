package data

import (
	"strconv"
	"strings"

	"github.com/Contra-Culture/report"
)

type (
	nodePath        []string
	namesOrderedSet struct {
		names map[string]bool
		order []string
	}
	nodesOrderedSet struct {
		nodes map[string]*node
		order []string
	}
	linksOrderedSet struct {
		links map[string]*link
		order []string
	}
	namesOrderedMapping struct {
		mapping map[string]string
		order   []string
	}
	nodeKind int
	// Schema represents a type of data objects, like blog posts, arguments, rubrics, categories, etc.
	node struct {
		kind      nodeKind
		path      *nodePath
		pk        *namesOrderedSet
		nodes     *nodesOrderedSet
		links     *linksOrderedSet
		storePath string
	}

	NodeCfgr struct {
		graphCfgr *GraphCfgr
		node      *node
		report    report.Node
	}
	// Arrow sets relations between schemas and objects that belong to that schemas.
	link struct {
		limit          int
		name           string
		hostNodePath   *nodePath
		remoteNodePath *nodePath
		remoteLinkName string
		counterCache   string
		mapping        *namesOrderedMapping
	}
	LinkCfgr struct {
		link     *link
		nodeCfgr *NodeCfgr
		report   report.Node
	}
)

const (
	_ nodeKind = iota
	stringNode
	orderedMapNode
	mapNode
	arrayNode
)

func (np *nodePath) slice() []string {
	return []string(*np)
}
func (np *nodePath) append(n string) *nodePath {
	path := nodePath(append(np.slice(), n))
	return &path
}
func (np *nodePath) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('[')
	lastIdx := len(*np) - 1
	for i, s := range *np {
		sb.WriteRune('"')
		sb.WriteString(s)
		sb.WriteRune('"')
		if i < lastIdx {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}

func newNames(attrs []string, handleExistance func(string)) (set *namesOrderedSet) {
	set = &namesOrderedSet{
		names: map[string]bool{},
		order: attrs,
	}
	for _, attr := range attrs {
		set.add(attr, handleExistance)
	}
	return
}
func (set *namesOrderedSet) add(newAttr string, handleExistance func(string)) {
	if set.names[newAttr] {
		handleExistance(newAttr)
		return
	}
	set.order = append(set.order, newAttr)
	set.names[newAttr] = true
}
func (set *namesOrderedSet) JSONString() string {
	if set == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteRune('[')
	for i, n := range set.order {
		sb.WriteRune('"')
		sb.WriteString(n)
		sb.WriteRune('"')
		if i < len(set.order)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}
func newLinks() *linksOrderedSet {
	return &linksOrderedSet{
		links: map[string]*link{},
		order: []string{},
	}
}
func (set *linksOrderedSet) add(a *link, handleExistance func(string)) {
	_, exists := set.links[a.name]
	if exists {
		handleExistance(a.name)
		return
	}
	set.links[a.name] = a
	set.order = append(set.order, a.name)
}
func (set *linksOrderedSet) JSONString() string {
	if set == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteRune('[')
	for i, n := range set.order {
		a := set.links[n]
		sb.WriteString("{\"limit\":\"")
		sb.WriteString(strconv.Itoa(a.limit))
		sb.WriteString("\",\"name\":\"")
		sb.WriteString(a.name)
		sb.WriteString("\",\"hostNode\":")
		sb.WriteString(a.hostNodePath.JSONString())
		sb.WriteString(",\"remoteNode\":")
		sb.WriteString(a.remoteNodePath.JSONString())
		sb.WriteString(",\"remoteLink\":\"")
		sb.WriteString(a.remoteLinkName)
		sb.WriteString("\",\"counterCache\":\"")
		sb.WriteString(a.counterCache)
		sb.WriteString("\",\"mapping\":")
		sb.WriteString(a.mapping.JSONString())
		sb.WriteString("}")
		if i < len(set.order)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}
func newNamesOrderedMapping() *namesOrderedMapping {
	return &namesOrderedMapping{
		mapping: map[string]string{},
		order:   []string{},
	}
}
func (set *namesOrderedMapping) add(ha, ra string, handleExistance func(string)) {
	if _, exists := set.mapping[ha]; exists {
		handleExistance(ha)
		return
	}
	set.order = append(set.order, ha)
	set.mapping[ha] = ra
}
func (set *namesOrderedMapping) JSONString() string {
	if set == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteRune('[')
	lastIdx := len(set.order) - 1
	for i, n := range set.order {
		sb.WriteString("[\"")
		sb.WriteString(n)
		sb.WriteString("\",\"")
		sb.WriteString(set.mapping[n])
		sb.WriteString("\"]")
		if i < lastIdx {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}

func newNodesOrderedSet() *nodesOrderedSet {
	return &nodesOrderedSet{
		nodes: map[string]*node{},
		order: []string{},
	}
}
func (set *nodesOrderedSet) add(s *node, handleExistence func([]string)) {
	path := s.path.slice()
	name := path[len(path)-1]
	if _, exists := set.nodes[name]; exists {
		handleExistence(path)
		return
	}
	set.order = append(set.order, name)
	set.nodes[name] = s
}
func (set *nodesOrderedSet) JSONString() string {
	if set == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("{")
	lastIdx := len(set.nodes) - 1
	for i, nodeName := range set.order {
		schema := set.nodes[nodeName]
		sb.WriteRune('"')
		sb.WriteString(nodeName)
		sb.WriteString("\":")
		sb.WriteString(schema.JSONString())
		if i < lastIdx {
			sb.WriteRune(',')
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (set *nodesOrderedSet) get(path *nodePath, handleNotFound func([]string)) (n *node, ok bool) {
	nodes := set.nodes
	lastIdx := len(path.slice()) - 1
	for i, chunk := range path.slice() {
		n, ok = nodes[chunk]
		if !ok {
			handleNotFound(path.slice()[:i+1])
			return
		}
		if i == lastIdx {
			return
		}
		nodes = n.nodes.nodes
	}
	return
}

// .PK() - specifies property names and their order for unique (primary) key calculation for the schema (data) object.
func (c *NodeCfgr) PK(id []string) {
	c.node.pk = newNames(id, func(n string) {
		c.report.Error("wrong ID: %s", n)
	})
}

// .Attribute() - specifies a schema (data) object's property.
func (c *NodeCfgr) String(n string) {
	c.node.nodes.add(
		&node{
			kind: stringNode,
			path: c.node.path.append(n),
		},
		func(np []string) {
			c.report.Error("attribute %s.%s already specified", c.node.path, n)
		})
}

// .Arrow() - specifies a relation between node objects, when one object has a reference to another one/other ones.
func (c *NodeCfgr) Link(name string, rnPath []string, cfg func(*LinkCfgr)) {
	nodePath := nodePath(rnPath)
	link := &link{
		name:           name,
		hostNodePath:   c.node.path,
		remoteNodePath: &nodePath,
		mapping:        newNamesOrderedMapping(),
	}
	cfg(
		&LinkCfgr{
			link:     link,
			nodeCfgr: c,
			report:   c.report.Structure("link: %s", name),
		})
	c.node.links.add(link, func(n string) {
		c.report.Error("link %s>-%s->... already specified", c.node.path, n)
	})
}
func (s *node) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('{')
	sb.WriteString("\"pk\":")
	sb.WriteString(s.pk.JSONString())
	sb.WriteRune(',')
	sb.WriteString("\"nodes\":")
	sb.WriteString(s.nodes.JSONString())
	sb.WriteRune(',')
	sb.WriteString("\"links\":")
	sb.WriteString(s.links.JSONString())
	sb.WriteString("}")
	return sb.String()
}
func (c *NodeCfgr) check() {

}

// .Map() - specifies a properties mapping for objects linking.
// For example, blog posts may be linked to a rubric, depending on their rubric_slug and rubric's slug properties
func (c *LinkCfgr) Map(ha, ra string) {
	c.link.mapping.add(ha, ra, func(n string) {
		c.nodeCfgr.report.Error(
			"%s>-%s->%s attributes mapping already specified",
			c.link.hostNodePath,
			c.link.name,
			c.link.remoteNodePath)
	})
}

// .RemoteLink() - specifies an arrow on remote schema that allows to link objects through another arrow.
func (c *LinkCfgr) RemoteLink(n string) {
	if len(c.link.remoteLinkName) > 0 {
		c.nodeCfgr.report.Error(
			"%s>-%s->%s remote link already specified",
			c.link.hostNodePath,
			c.link.name,
			c.link.remoteNodePath)
		return
	}
	c.link.remoteLinkName = n
}

// .Limit() - allows to limit a number of linked objects.
func (c *LinkCfgr) Limit(l int) {
	if c.link.limit > 0 {
		c.nodeCfgr.report.Error(
			"%s>-%s->%s limit already specified",
			c.link.hostNodePath,
			c.link.name,
			c.link.remoteNodePath)
		return
	}
	c.link.limit = l
}

// .CounterCache() - allows to specify a field with the count of associated objects.
// The value of object counter is using for limiting a number of linked objects (.Limit()).
func (c *LinkCfgr) CounterCache(n string) {
	if len(c.link.counterCache) > 0 {
		c.nodeCfgr.report.Error(
			"%s>-%s->%s counter-cache attribute already specified",
			c.link.hostNodePath,
			c.link.name,
			c.link.remoteNodePath)
		return
	}
	c.link.counterCache = n
}
