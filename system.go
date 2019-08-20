package ecs

// System is an interface to a type implementing a system.
type System interface {
	Update(interface{})
	Add(Entity) error
	Remove(Entity)
}

// Prioritizer is a System that is executed with a priority.
// Default Priority is 0, executed from smallest to largest.
type Prioritizer interface {
	System
	GetPriority() int
}

// Initaliser is a System that has an Init function.
// It is called on addition to the state.
type Initaliser interface {
	System
	Init(*State)
}
