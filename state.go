package ecs

import (
	"sort"
)

// State is an Entity-Component-System.
// It is simply a slice of systems.
type State []System

// AddSystem adds a system to the State
func (s *State) AddSystem(system System) {
	*s = append(*s, system)
	sort.Sort(s)
}

// Update calls Update in all systems.
// u can be ignored, or used for run-levels, a time-delta, or whatever you like.
// it is passed to all systems
func (s State) Update(u interface{}) {
	for i := range s {
		s[i].Update(u)
	}
}

// Add adds an entity to all systems.
// If an error occurs, it removes the entity from already-added systems,
// and returns the error. If the Entity'd ID is 0, it will be generated.
func (s State) Add(e Entity) error {
	if e.GetID().ID == 0 {
		id := e.GetID()
		nid := NewID()
		*id = *nid
	}
	for i := range s {
		err := s[i].Add(e)
		if err != nil {
			for ; i >= 0; i-- {
				s[i].Remove(e)
			}
			return err
		}
	}
	return nil
}

// Remove removes an entity from all systems
func (s State) Remove(e Entity) {
	for i := range s {
		s[i].Remove(e)
	}
}

// These functions implement sort.Interface

func (s State) Len() int      { return len(s) }
func (s State) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s State) Less(i, j int) bool {
	var ip, jp int
	if p, ok := s[i].(Prioritizer); ok {
		ip = p.Priority()
	}
	if p, ok := s[j].(Prioritizer); ok {
		jp = p.Priority()
	}
	return ip < jp
}
