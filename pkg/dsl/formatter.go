package dsl

import (
	"fmt"
	"sort"
	"strings"
)

// Format serializes a Dispatch AST into standard HCP DSL syntax.
func Format(d *Dispatch) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("DISPATCH %s\n", d.ID))
	sb.WriteString(fmt.Sprintf("ENCRYPTION: %s\n", d.Encryption))
	sb.WriteString(fmt.Sprintf("ORIGIN: %s\n", d.Origin))
	sb.WriteString(fmt.Sprintf("TARGET: %s\n", d.Target))
	sb.WriteString(fmt.Sprintf("STARDATE: %s\n", d.Stardate))
	sb.WriteString(fmt.Sprintf("PRIORITY: %s\n", d.Priority))
	if d.Status != "" {
		sb.WriteString(fmt.Sprintf("STATUS: %s\n", d.Status))
	}

	// Sort extra headers for deterministic output
	keys := make([]string, 0, len(d.Headers))
	for k := range d.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, d.Headers[k]))
	}

	sb.WriteString("PAYLOAD:\n")
	for _, item := range d.Payload {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", item.Type, item.Content))
	}

	sb.WriteString(fmt.Sprintf("AUTHENTICATION: %s\n", d.Authentication))
	sb.WriteString("END_DISPATCH\n")

	return sb.String()
}
