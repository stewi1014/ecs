package ecs

import "sync/atomic"

var idSrc = uint64(1)

// Entity is an interface to an entity.
// A struct embedding ID will implement it.
type Entity interface {
	GetID() *ID
	Equal(e Entity) bool
	IsZero() bool
}

// ID is the base component all entities have.
type ID struct {
	ID uint64
}

// NewID creates a new ID component.
func NewID() *ID {
	return &ID{
		ID: atomic.AddUint64(&idSrc, 1),
	}
}

// These functions are written in an unecceceraly verbose way to keep with the ecs ideology

// GetID implements Entity
func (id *ID) GetID() *ID {
	return id
}

// Equal returns true if the IDs are equal
func (id *ID) Equal(e Entity) bool {
	return id.ID == e.GetID().ID
}

// IsZero returns true if the ID is uninitalised
func (id *ID) IsZero() bool {
	return id.ID == 0
}

// Namer is an interface to an entity with a Name component
type Namer interface {
	GetName() *Name
}

// Name is an entity component describing its name
type Name struct {
	Name string
}

// NewName returns a new Name component
func NewName(name string) *Name {
	return &Name{
		Name: name,
	}
}

// GetName implements Namer
func (n *Name) GetName() *Name {
	return n
}
