package dsl

import (
	"errors"
	"strings"
)

// PayloadItem represents a structured narrative statement in an HCP dispatch.
type PayloadItem struct {
	Type    string `json:"type"`    // EVENT, INTEL, DIRECTIVE, QUERY, WARNING, COMPACT
	Content string `json:"content"` // Narrative content string
}

// Dispatch represents a parsed Hyper-Comm Protocol (HCP) message.
type Dispatch struct {
	ID             string            `json:"id"`
	Encryption     string            `json:"encryption"`
	Origin         string            `json:"origin"`
	Target         string            `json:"target"`
	Stardate       string            `json:"stardate"`
	Priority       string            `json:"priority"`
	Status         string            `json:"status,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Payload        []PayloadItem     `json:"payload"`
	Authentication string            `json:"authentication"`
}

// Validate checks that all required fields of an HCP dispatch are populated and valid.
func (d *Dispatch) Validate() error {
	if strings.TrimSpace(d.ID) == "" {
		return errors.New("dispatch missing required field: ID")
	}
	if strings.TrimSpace(d.Encryption) == "" {
		return errors.New("dispatch missing required field: ENCRYPTION")
	}
	if strings.TrimSpace(d.Origin) == "" {
		return errors.New("dispatch missing required field: ORIGIN")
	}
	if strings.TrimSpace(d.Target) == "" {
		return errors.New("dispatch missing required field: TARGET")
	}
	if strings.TrimSpace(d.Stardate) == "" {
		return errors.New("dispatch missing required field: STARDATE")
	}
	if strings.TrimSpace(d.Priority) == "" {
		return errors.New("dispatch missing required field: PRIORITY")
	}
	if len(d.Payload) == 0 {
		return errors.New("dispatch PAYLOAD cannot be empty")
	}
	if strings.TrimSpace(d.Authentication) == "" {
		return errors.New("dispatch missing required field: AUTHENTICATION")
	}
	return nil
}
