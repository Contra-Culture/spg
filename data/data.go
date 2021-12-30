package data

import (
	"fmt"
	"strings"
	"time"
)

type (
	// Object represents a data object, like a blog post or rubric.
	Object struct {
		node      *node
		updatedAt time.Time
		props     map[string]string
		absPath   string
	}
)

// .Prop() - is a read accessor to project's property by its name.
func (o *Object) Prop(n string) (p string, err error) {
	p, ok := o.props[n]
	if !ok {
		err = fmt.Errorf("property %s.%s is not specified", o.node.path, n)
	}
	return
}

// .ID() - returns unique (primary) key for the object.
func (o *Object) PK() string {
	var sb strings.Builder
	for _, n := range o.node.pk.order {
		sb.WriteString(o.props[n])
	}
	return sb.String()
}

// .JSONString() - returns string representation of JSON.
func (o Object) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('{')
	idx := 0
	lastIdx := len(o.props) - 1
	for prop, val := range o.props {
		sb.WriteRune('"')
		sb.WriteString(prop)
		sb.WriteString("\":\"")
		sb.WriteString(val)
		sb.WriteRune('"')
		if idx < lastIdx {
			sb.WriteRune(',')
			idx++
		}
	}
	sb.WriteRune('}')
	return sb.String()
}
