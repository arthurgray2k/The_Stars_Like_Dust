package dsl

import (
	"bufio"
	"fmt"
	"strings"
)

// Parse parses a raw HCP DSL string into a Dispatch structure.
func Parse(content string) (*Dispatch, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	dispatch := &Dispatch{
		Headers: make(map[string]string),
		Payload: make([]PayloadItem, 0),
	}

	inPayload := false
	sawDispatch := false
	sawEndDispatch := false

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip comments and empty lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if strings.HasPrefix(trimmed, "DISPATCH") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				dispatch.ID = parts[1]
			} else {
				dispatch.ID = "UNKNOWN"
			}
			sawDispatch = true
			continue
		}

		if trimmed == "END_DISPATCH" {
			sawEndDispatch = true
			inPayload = false
			continue
		}

		if trimmed == "PAYLOAD:" {
			inPayload = true
			continue
		}

		if inPayload {
			// Check if we exited payload section via a top-level key like AUTHENTICATION:
			if strings.HasPrefix(trimmed, "AUTHENTICATION:") {
				inPayload = false
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					dispatch.Authentication = strings.TrimSpace(parts[1])
				}
				continue
			}

			// Payload lines are typically indented: TYPE: "content"
			idx := strings.Index(trimmed, ":")
			if idx > 0 {
				itemType := strings.ToUpper(strings.TrimSpace(trimmed[:idx]))
				rawContent := strings.TrimSpace(trimmed[idx+1:])
				cleanContent := strings.Trim(rawContent, "\"")
				dispatch.Payload = append(dispatch.Payload, PayloadItem{
					Type:    itemType,
					Content: cleanContent,
				})
			}
			continue
		}

		// Top-level key: value
		idx := strings.Index(trimmed, ":")
		if idx > 0 {
			key := strings.ToUpper(strings.TrimSpace(trimmed[:idx]))
			val := strings.TrimSpace(trimmed[idx+1:])
			val = strings.Trim(val, "\"")

			switch key {
			case "ENCRYPTION":
				dispatch.Encryption = val
			case "ORIGIN":
				dispatch.Origin = val
			case "TARGET":
				dispatch.Target = val
			case "STARDATE":
				dispatch.Stardate = val
			case "PRIORITY":
				dispatch.Priority = val
			case "STATUS":
				dispatch.Status = val
			case "AUTHENTICATION":
				dispatch.Authentication = val
			default:
				dispatch.Headers[key] = val
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading HCP stream: %w", err)
	}

	if !sawDispatch {
		return nil, fmt.Errorf("invalid HCP: missing starting DISPATCH header")
	}
	if !sawEndDispatch {
		return nil, fmt.Errorf("invalid HCP: missing END_DISPATCH marker")
	}

	if err := dispatch.Validate(); err != nil {
		return nil, err
	}

	return dispatch, nil
}
