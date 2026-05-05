package model

import "time"

// Item represents a single trackable item (feature, bug, idea, task, epic).
type Item struct {
	ID       string    `yaml:"id"`
	Type     ItemType  `yaml:"type"`
	Title    string    `yaml:"title"`
	Status   string    `yaml:"status"`
	Priority string    `yaml:"priority"`
	Phase    *int      `yaml:"phase,omitempty"`
	Created  time.Time `yaml:"created"`
	Updated  time.Time `yaml:"updated"`
	Tags     []string  `yaml:"tags,omitempty"`
	Links    []Link    `yaml:"links,omitempty"`
	Blocked  *Blocked  `yaml:"blocked,omitempty"`
	Body     string    `yaml:"-"`
}

// ItemType is the kind of item.
type ItemType string

const (
	TypeFeature ItemType = "feature"
	TypeBug     ItemType = "bug"
	TypeIdea    ItemType = "idea"
	TypeTask    ItemType = "task"
	TypeEpic    ItemType = "epic"
)

// ValidTypes returns all valid item types.
func ValidTypes() []ItemType {
	return []ItemType{TypeFeature, TypeBug, TypeIdea, TypeTask, TypeEpic}
}

// ParseItemType converts a string to an ItemType, returns empty if invalid.
func ParseItemType(s string) (ItemType, bool) {
	t := ItemType(s)
	for _, valid := range ValidTypes() {
		if t == valid {
			return t, true
		}
	}
	return "", false
}

// Prefix returns the ID prefix for this item type.
func (t ItemType) Prefix() string {
	switch t {
	case TypeFeature:
		return "FEAT"
	case TypeBug:
		return "BUG"
	case TypeIdea:
		return "IDEA"
	case TypeTask:
		return "TASK"
	case TypeEpic:
		return "EPIC"
	default:
		return "ITEM"
	}
}

// Blocked represents the blocked state of an item.
type Blocked struct {
	Reason string `yaml:"reason,omitempty"`
}
