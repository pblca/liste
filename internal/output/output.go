package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pblca/liste/internal/model"
)

// Format represents the output format.
type Format int

const (
	FormatTable Format = iota
	FormatJSON
	FormatQuiet
)

// Formatter handles output rendering.
type Formatter struct {
	Writer io.Writer
	Format Format
}

// New creates a formatter for the given writer and format.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{Writer: w, Format: format}
}

// ItemCreated outputs the result of creating an item.
func (f *Formatter) ItemCreated(item *model.Item) {
	switch f.Format {
	case FormatJSON:
		f.json(map[string]any{
			"id":       item.ID,
			"type":     string(item.Type),
			"title":    item.Title,
			"status":   item.Status,
			"priority": item.Priority,
			"created":  item.Created.Format("2006-01-02"),
		})
	case FormatQuiet:
		fmt.Fprintln(f.Writer, item.ID)
	default:
		fmt.Fprintf(f.Writer, "Created %s: %s\n", item.ID, item.Title)
		fmt.Fprintf(f.Writer, "  Type: %s | Status: %s | Priority: %s\n", item.Type, item.Status, item.Priority)
	}
}

// ItemDetail outputs the full detail of an item.
func (f *Formatter) ItemDetail(item *model.Item, inverseLinks []InverseLinkDisplay) {
	switch f.Format {
	case FormatJSON:
		data := map[string]any{
			"id":       item.ID,
			"type":     string(item.Type),
			"title":    item.Title,
			"status":   item.Status,
			"priority": item.Priority,
			"created":  item.Created.Format("2006-01-02"),
			"updated":  item.Updated.Format("2006-01-02"),
			"tags":     item.Tags,
			"links":    item.Links,
			"body":     item.Body,
		}
		if item.Blocked != nil {
			data["blocked"] = item.Blocked
		}
		if len(inverseLinks) > 0 {
			data["referenced_by"] = inverseLinks
		}
		f.json(data)
	case FormatQuiet:
		fmt.Fprintln(f.Writer, item.ID)
	default:
		fmt.Fprintf(f.Writer, "%-8s %s\n", item.ID, item.Title)
		fmt.Fprintf(f.Writer, "Type: %s | Status: %s | Priority: %s\n", item.Type, item.Status, item.Priority)
		fmt.Fprintf(f.Writer, "Created: %s | Updated: %s\n", item.Created.Format("2006-01-02"), item.Updated.Format("2006-01-02"))

		if len(item.Tags) > 0 {
			fmt.Fprintf(f.Writer, "Tags: %s\n", strings.Join(item.Tags, ", "))
		}
		if item.Blocked != nil {
			reason := item.Blocked.Reason
			if reason == "" {
				reason = "(no reason)"
			}
			fmt.Fprintf(f.Writer, "BLOCKED: %s\n", reason)
		}
		if len(item.Links) > 0 {
			fmt.Fprintln(f.Writer, "\nLinks:")
			for _, l := range item.Links {
				proj := ""
				if l.Project != "" {
					proj = " [" + l.Project + "]"
				}
				fmt.Fprintf(f.Writer, "  %s %s%s\n", l.Type, l.Target, proj)
			}
		}
		if len(inverseLinks) > 0 {
			fmt.Fprintln(f.Writer, "\nReferenced by:")
			for _, l := range inverseLinks {
				fmt.Fprintf(f.Writer, "  %s %s\n", l.Type, l.SourceID)
			}
		}
		if item.Body != "" {
			fmt.Fprintf(f.Writer, "\n%s\n", item.Body)
		}
	}
}

// InverseLinkDisplay is a simplified inverse link for display.
type InverseLinkDisplay struct {
	Type     string `json:"type"`
	SourceID string `json:"source_id"`
}

// ItemList outputs a list of items.
func (f *Formatter) ItemList(items []*model.Item) {
	switch f.Format {
	case FormatJSON:
		list := make([]map[string]any, 0, len(items))
		for _, item := range items {
			entry := map[string]any{
				"id":       item.ID,
				"type":     string(item.Type),
				"title":    item.Title,
				"status":   item.Status,
				"priority": item.Priority,
			}
			if item.Blocked != nil {
				entry["blocked"] = true
			}
			list = append(list, entry)
		}
		f.json(list)
	case FormatQuiet:
		for _, item := range items {
			fmt.Fprintln(f.Writer, item.ID)
		}
	default:
		if len(items) == 0 {
			fmt.Fprintln(f.Writer, "No items found.")
			return
		}
		// Table header
		fmt.Fprintf(f.Writer, "%-10s %-8s %-8s %-10s %s\n", "ID", "TYPE", "STATUS", "PRIORITY", "TITLE")
		fmt.Fprintf(f.Writer, "%-10s %-8s %-8s %-10s %s\n", "---", "----", "------", "--------", "-----")
		for _, item := range items {
			blocked := ""
			if item.Blocked != nil {
				blocked = " [BLOCKED]"
			}
			fmt.Fprintf(f.Writer, "%-10s %-8s %-8s %-10s %s%s\n",
				item.ID, item.Type, item.Status, item.Priority, item.Title, blocked)
		}
		fmt.Fprintf(f.Writer, "\n%d item(s)\n", len(items))
	}
}

// StatusSummary outputs a dashboard-style summary.
func (f *Formatter) StatusSummary(items []*model.Item, projectName string) {
	// Group by status
	groups := make(map[string][]*model.Item)
	for _, item := range items {
		status := item.Status
		if item.Blocked != nil {
			status = "blocked"
		}
		groups[status] = append(groups[status], item)
	}

	switch f.Format {
	case FormatJSON:
		summary := map[string]any{
			"project": projectName,
			"total":   len(items),
			"by_status": func() map[string]int {
				counts := make(map[string]int)
				for k, v := range groups {
					counts[k] = len(v)
				}
				return counts
			}(),
			"items": func() []map[string]any {
				list := make([]map[string]any, 0, len(items))
				for _, item := range items {
					list = append(list, map[string]any{
						"id":       item.ID,
						"type":     string(item.Type),
						"title":    item.Title,
						"status":   item.Status,
						"priority": item.Priority,
						"blocked":  item.Blocked != nil,
					})
				}
				return list
			}(),
		}
		f.json(summary)
	case FormatQuiet:
		fmt.Fprintf(f.Writer, "%d items\n", len(items))
	default:
		fmt.Fprintf(f.Writer, "Project: %s (%d items)\n\n", projectName, len(items))

		statusOrder := []string{"active", "planned", "blocked", "idea", "done", "cancelled"}
		for _, status := range statusOrder {
			group, ok := groups[status]
			if !ok || len(group) == 0 {
				continue
			}
			label := strings.ToUpper(status)
			fmt.Fprintf(f.Writer, "%s (%d)\n", label, len(group))
			for _, item := range group {
				fmt.Fprintf(f.Writer, "  %-10s [%s] %s\n", item.ID, item.Priority, item.Title)
			}
			fmt.Fprintln(f.Writer)
		}
	}
}

// ProjectList outputs the list of discovered projects.
func (f *Formatter) ProjectList(root string, subProjects []ProjectDisplay) {
	switch f.Format {
	case FormatJSON:
		f.json(map[string]any{
			"root":         root,
			"sub_projects": subProjects,
		})
	default:
		fmt.Fprintf(f.Writer, "Root: %s\n", root)
		if len(subProjects) > 0 {
			fmt.Fprintln(f.Writer, "\nSub-projects:")
			for _, p := range subProjects {
				fmt.Fprintf(f.Writer, "  %s (%d items)\n", p.Name, p.ItemCount)
			}
		}
	}
}

// ProjectDisplay is a simplified project for display.
type ProjectDisplay struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	ItemCount int    `json:"item_count"`
}

// Message outputs a simple message.
func (f *Formatter) Message(msg string) {
	switch f.Format {
	case FormatJSON:
		f.json(map[string]string{"message": msg})
	default:
		fmt.Fprintln(f.Writer, msg)
	}
}

// Error outputs an error message.
func (f *Formatter) Error(err error) {
	switch f.Format {
	case FormatJSON:
		f.json(map[string]string{"error": err.Error()})
	default:
		fmt.Fprintf(f.Writer, "Error: %s\n", err)
	}
}

func (f *Formatter) json(v any) {
	enc := json.NewEncoder(f.Writer)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
