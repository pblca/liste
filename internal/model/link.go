package model

// Link represents a typed relationship between items.
type Link struct {
	Type    LinkType `yaml:"type"`
	Target  string   `yaml:"target"`
	Project string   `yaml:"project,omitempty"` // cross-project reference
}

// LinkType defines the kind of relationship.
type LinkType string

const (
	LinkDependsOn LinkType = "depends-on"
	LinkBlocks    LinkType = "blocks"
	LinkParentOf  LinkType = "parent-of"
	LinkChildOf   LinkType = "child-of"
	LinkRelatesTo LinkType = "relates-to"
)

// ValidLinkTypes returns all valid link types.
func ValidLinkTypes() []LinkType {
	return []LinkType{LinkDependsOn, LinkBlocks, LinkParentOf, LinkChildOf, LinkRelatesTo}
}

// ParseLinkType converts a string to a LinkType.
func ParseLinkType(s string) (LinkType, bool) {
	t := LinkType(s)
	for _, valid := range ValidLinkTypes() {
		if t == valid {
			return t, true
		}
	}
	return "", false
}

// Inverse returns the inverse link type for bidirectional resolution.
func (t LinkType) Inverse() LinkType {
	switch t {
	case LinkDependsOn:
		return LinkBlocks
	case LinkBlocks:
		return LinkDependsOn
	case LinkParentOf:
		return LinkChildOf
	case LinkChildOf:
		return LinkParentOf
	case LinkRelatesTo:
		return LinkRelatesTo
	default:
		return ""
	}
}
