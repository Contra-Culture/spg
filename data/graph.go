package data

type (
	Graph struct {
		schemas map[string]*Schema
		objects map[string]map[string]interface{} // interface{} is map[string]interface{} or *Object
	}
	GraphCfgr struct {
		graph *Graph
	}
)

func New(cfg func(*GraphCfgr)) *Graph {
	graph := &Graph{
		schemas: map[string]*Schema{},
		objects: map[string]map[string]interface{}{},
	}
	cfg(
		&GraphCfgr{
			graph: graph,
		})
	return graph
}
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	schema := &Schema{
		name:         n,
		pk:           []string{},
		attributes:   []string{},
		associations: map[string]*Association{},
	}
	c.graph.schemas[n] = schema
	c.graph.objects[n] = map[string]interface{}{}
	cfg(
		&SchemaCfgr{
			graphCfgr: c,
			schema:    schema,
		})
}
func (g *Graph) Update(s string, attrs map[string]interface{}) {
	//schema := r.schemas[s]
	//objects := r.objects[s]
}
