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

func New(cfg func(*GraphCfgr)) (g *Graph, err error) {
	g = &Graph{
		schemas: map[string]*Schema{},
		objects: map[string]map[string]interface{}{},
	}
	cfg(
		&GraphCfgr{
			graph: g,
		})
	return
}
func (c *GraphCfgr) Schema(n string, cfg func(*SchemaCfgr)) {
	schema := &Schema{
		name:       n,
		id:         []string{},
		attributes: []string{},
		arrows:     map[string]*Arrow{},
	}
	c.graph.schemas[n] = schema
	c.graph.objects[n] = map[string]interface{}{}
	cfg(
		&SchemaCfgr{
			graphCfgr: c,
			schema:    schema,
		})
}
func (g *Graph) Update(s string, attrs map[string]interface{}) (err error) {
	//schema := r.schemas[s]
	//objects := r.objects[s]
	return
}
