package ecs

import (
	"errors"
	"math/rand"
	"time"
)

// Depender is an interface to an entity that has dependancies.
type Depender interface {
	Entity
	GetDepends() *Depends
}

// NewDepends creates a new Depends component.
// Each element in dependancies is treated as a seperate dependancy and Entity.
func NewDepends(dependancies ...Filter) *Depends {
	return &Depends{
		Dependancies: dependancies,
	}
}

// Depends is an Entity component describing its dependancies.
type Depends struct {
	Dependancies []Filter
	Entities     []Entity
}

// GetDepends implements Depender.
func (d *Depends) GetDepends() *Depends {
	return d
}

// DependsSystem is a system for managing dependancies.
type DependsSystem struct {
	Priority int

	entities  []Entity
	dependers map[ID][]Depender
	shuffeler *rand.Rand
	state     *State
}

// NewDependsSystem returns a new DependsSystem
func NewDependsSystem() *DependsSystem {
	return &DependsSystem{
		shuffeler: rand.New(rand.NewSource(time.Now().Unix())),
		dependers: make(map[ID][]Depender),
	}
}

// Init implements Initaliser
func (d *DependsSystem) Init(s *State) {
	d.state = s
}

// Update implements System
func (d *DependsSystem) Update(u interface{}) {
	//Nothing to do
}

// ErrUnmetDependancy is returned by DependsSystem if a dependancy is not met.
var ErrUnmetDependancy = errors.New("dependancy for entity was not met")

// Add implements System
// If the entity is a Depender, it resolves dependancies
func (d *DependsSystem) Add(e Entity) error {
	d.entities = append(d.entities, e)

	if depender, ok := e.(Depender); ok {
		depends := depender.GetDepends()
		depends.Entities = make([]Entity, len(depends.Dependancies))
		var have int
		for _, e := range d.shuffeler.Perm(len(d.entities)) {
			for j := range depends.Dependancies {
				if depends.Entities[j] != nil {
					continue
				}
				if depends.Dependancies[j].Satisfies(d.entities[e]) {
					depends.Entities[j] = d.entities[e]
					d.dependers[*d.entities[e].GetID()] = append(d.dependers[*d.entities[e].GetID()], depender)
					have++
					if have == len(depends.Dependancies) {
						return nil
					}
				}
			}
		}

		// We didn't find all dependancies
		return ErrUnmetDependancy
	}

	return nil
}

// Remove implements System
func (d *DependsSystem) Remove(e Entity) {
	for i := range d.entities {
		if d.entities[i].Equal(e) {
			d.entities = append(d.entities[:i], d.entities[i+1:]...)
		}
	}

	if depender, ok := e.(Depender); ok {
		depends := depender.GetDepends()
		for i := range depends.Entities {
			dependers := d.dependers[*depends.Entities[i].GetID()]
			for j := range dependers {
				if dependers[j].Equal(depender) {
					if len(dependers) == 1 {
						delete(d.dependers, *depends.Entities[i].GetID())
						break
					}
					dependers = append(dependers[:j], dependers[j+1:]...)
				}
			}
			d.dependers[*depends.Entities[i].GetID()] = dependers
		}
	}

	if dependers, ok := d.dependers[*e.GetID()]; ok {
		for i := range dependers {
			d.state.Remove(dependers[i])
		}
		delete(d.dependers, *e.GetID())
	}
}

// GetPriority implements Prioritiser
func (d *DependsSystem) GetPriority() int {
	return d.Priority
}
