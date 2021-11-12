package data

import (
	"fmt"
	"strings"
	"time"
)

type (
	Object struct {
		schema    *Schema
		updatedAt time.Time
		props     map[string]string
	}
)

func (o *Object) Prop(n string) (p string, err error) {
	p, ok := o.props[n]
	if !ok {
		err = fmt.Errorf("property %s.%s is not specified", o.schema.name, n)
	}
	return
}
func (o *Object) ID() string {
	var sb strings.Builder
	for _, pName := range o.schema.id {
		sb.WriteString(o.props[pName])
	}
	return sb.String()
}
