package data

import (
	"fmt"
	"strings"
	"time"
)

type (
	// Object represents a data object, like a blog post or rubric.
	Object struct {
		schema    *schema
		updatedAt time.Time
		props     map[string]string
		absPath   string
	}
)

// .Prop() - is a read accessor to project's property by its name.
func (o *Object) Prop(n string) (p string, err error) {
	p, ok := o.props[n]
	if !ok {
		err = fmt.Errorf("property %s.%s is not specified", o.schema.name, n)
	}
	return
}

// .ID() - returns unique (primary) key for the object.
func (o *Object) ID() string {
	var sb strings.Builder
	for _, n := range o.schema.id.order {
		sb.WriteString(o.props[n])
	}
	return sb.String()
}

// .JSONString() - returns string representation of JSON.
func (o Object) JSONString() string {
	var sb strings.Builder
	sb.WriteRune('"')
	idx := 0
	for prop, val := range o.props {
		sb.WriteRune('"')
		sb.WriteString(prop)
		sb.WriteString("\":\"")
		sb.WriteString(val)
		lastIdx := len(o.props) - 1
		if idx < lastIdx {
			sb.WriteString("\",")
			idx++
		}
	}
	sb.WriteRune('}')
	return sb.String()
}
