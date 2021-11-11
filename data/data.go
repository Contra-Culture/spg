package data

import (
	"time"
)

type (
	Object struct {
		schema     *Schema
		sha        []byte
		updatedAt  time.Time
		props      map[string]string
		meta       map[string]interface{}
		embeddings map[string]interface{} // interface{} is *Object or []*Object
	}
)

func (o *Object) Prop(n string) string {
	return o.props[n]
}
