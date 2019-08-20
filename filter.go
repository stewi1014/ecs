package ecs

import (
	"reflect"
)

// Filter is an interface to a type that filters entities.
type Filter interface {
	// Satisfies returns true if the entity passes the filters.
	Satisfies(Entity) bool
}

// NewMultiFilter retrurns a new MultiFilter with the given list of filters.
func NewMultiFilter(filters ...Filter) *MultiFilter {
	return &MultiFilter{
		Filters: filters,
	}
}

// MultiFilter is a Filter that has multiple Filters.
// All filters must pass for Satisfies to return true.
type MultiFilter struct {
	Filters []Filter
}

// Satisfies implements filter.
func (m *MultiFilter) Satisfies(e Entity) bool {
	for _, f := range m.Filters {
		if !f.Satisfies(e) {
			return false
		}
	}
	return true
}

// NewIDFilter returns a Filter for an Entity with a specifc ID.
func NewIDFilter(b *ID) Filter {
	return &IDFilter{
		id: b,
	}
}

// IDFilter is a filter for an entity with a specifc ID
type IDFilter struct {
	id *ID
}

// Satisfies implements Filter.
func (d *IDFilter) Satisfies(e Entity) bool {
	return d.id.Equal(e)
}

// NewTypeFilter returns a filter for a specific type.
// It takes a reflect.Type, struct, pointer to struct, or pointer to interface.
func NewTypeFilter(t interface{}) *TypeFilter {
	var ty reflect.Type
	if i, ok := t.(reflect.Type); ok {
		ty = i
	} else {
		ty = reflect.TypeOf(t)
	}

	if ty.Kind() == reflect.Ptr && ty.Elem().Kind() == reflect.Interface {
		ty = ty.Elem()
	}

	return &TypeFilter{
		t: ty,
	}
}

// TypeFilter is a dependancy on a type (struct or interface)
type TypeFilter struct {
	t reflect.Type
}

// Satisfies implements Filter
func (t *TypeFilter) Satisfies(e Entity) bool {
	et := reflect.TypeOf(e)
	if t.t.Kind() == reflect.Interface {
		return et.Implements(t.t)
	}

	return t.t == et
}

// NewNameFilter returns a new NameFilter for entities witht he given name.
func NewNameFilter(n *Name) *NameFilter {
	return &NameFilter{
		Name: n,
	}
}

// NameFilter is a filter for entities with a given name.
type NameFilter struct {
	Name *Name
}

// Satisfies implements Filter
func (n *NameFilter) Satisfies(e Entity) bool {
	if namer, ok := e.(Namer); ok {
		return namer.GetName().Name == n.Name.Name
	}
	return false
}
