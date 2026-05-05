package model

// Project represents a roadmap project rooted at a .liste/ directory.
type Project struct {
	Name   string  // display name from config
	Path   string  // absolute path to the .liste/ directory
	Config Config  // project configuration
	State  State   // internal state (next IDs, etc.)
	Items  []*Item // loaded items
}

// Config represents the .liste/config.yaml file.
type Config struct {
	Project    string   `yaml:"project"`
	Statuses   []string `yaml:"statuses"`
	Blocked    bool     `yaml:"blocked"`
	Types      []string `yaml:"types"`
	Priorities []string `yaml:"priorities"`
	Defaults   Defaults `yaml:"defaults"`
}

// Defaults represents default values for new items.
type Defaults struct {
	Status   string `yaml:"status"`
	Priority string `yaml:"priority"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig(name string) Config {
	return Config{
		Project:    name,
		Statuses:   []string{"idea", "planned", "active", "done", "cancelled"},
		Blocked:    true,
		Types:      []string{"feature", "bug", "idea", "task", "epic"},
		Priorities: []string{"critical", "high", "medium", "low"},
		Defaults: Defaults{
			Status:   "idea",
			Priority: "medium",
		},
	}
}

// State represents the .liste/.state.yaml file for internal bookkeeping.
type State struct {
	NextIDs map[string]int `yaml:"next_ids"` // type prefix -> next number
}

// DefaultState returns an empty state.
func DefaultState() State {
	return State{
		NextIDs: map[string]int{
			"FEAT": 1,
			"BUG":  1,
			"IDEA": 1,
			"TASK": 1,
			"EPIC": 1,
		},
	}
}

// IsValidStatus checks if a status is valid for this config.
func (c *Config) IsValidStatus(s string) bool {
	for _, valid := range c.Statuses {
		if s == valid {
			return true
		}
	}
	return false
}

// IsValidPriority checks if a priority is valid for this config.
func (c *Config) IsValidPriority(p string) bool {
	for _, valid := range c.Priorities {
		if p == valid {
			return true
		}
	}
	return false
}
