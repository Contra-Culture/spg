package report

import "strings"

type (
	recordKind int
	Context    struct {
		depth    int
		children []interface{} // interface{} is *Context or *Record
		title    string
	}
	Record struct {
		kind    recordKind
		message string
	}
)

const (
	_ recordKind = iota
	Error
	Info
	Warn
	Deprecation
)

var kindToPrefixMapping = map[recordKind]string{
	Error:       "\t[ error ] ",
	Info:        "\t[ info ] ",
	Warn:        "\t[ warn ] ",
	Deprecation: "\t[ deprecation ] ",
}

func New(t string) (c *Context) {
	return &Context{
		depth:    0,
		title:    t,
		children: []interface{}{},
	}
}

func (c *Context) Write(sb *strings.Builder) {
	sb.WriteString(c.title)
	sb.WriteRune('\n')
	for _, rawChild := range c.children {
		for i := 0; i <= c.depth; i++ {
			sb.WriteRune('\t')
		}
		switch child := rawChild.(type) {
		case *Context:
			child.Write(sb)
		case *Record:
			sb.WriteString(kindToPrefixMapping[child.kind])
			sb.WriteString(child.message)
			sb.WriteRune('\n')
		default:
			panic("wrong children type")
		}
	}
}
func (c *Context) Error(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Error,
			message: m,
		})
}
func (c *Context) Warn(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Warn,
			message: m,
		})
}
func (c *Context) Deprecation(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Deprecation,
			message: m,
		})
}
func (c *Context) Info(m string) {
	c.children = append(
		c.children,
		&Record{
			kind:    Info,
			message: m,
		})
}
func (c *Context) Context(t string) *Context {
	newContext := &Context{
		depth: c.depth + 1,
		title: t,
	}
	c.children = append(c.children, newContext)
	return newContext
}
