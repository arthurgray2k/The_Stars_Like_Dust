package dsl

import (
	"strings"
	"testing"
)

func TestParseAndFormat(t *testing.T) {
	raw := `
# Interstellar Subspace Transmission
DISPATCH HCP-820-001
ENCRYPTION: TYRANN-CIPHER-7
ORIGIN: Rhodia-Orbit / Tyrann-Flagship
TARGET: Commissioner-Aratap
STARDATE: 820-G.E.329
PRIORITY: URGENT
STATUS: INTERCEPTED
CODENAME: NEBULA-SHADOW
PAYLOAD:
  EVENT: "Sub-ether tracer affixed to vessel Remembrance."
  INTEL: "Widemos heir Biron accompanied by Hinriad defectors."
  DIRECTIVE: "Maintain passive sensor trailing."
AUTHENTICATION: SIG-ARATAP-9941-OMEGA
END_DISPATCH
`

	d, err := Parse(raw)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if d.ID != "HCP-820-001" {
		t.Errorf("expected ID 'HCP-820-001', got %s", d.ID)
	}
	if d.Encryption != "TYRANN-CIPHER-7" {
		t.Errorf("expected Encryption 'TYRANN-CIPHER-7', got %s", d.Encryption)
	}
	if d.Origin != "Rhodia-Orbit / Tyrann-Flagship" {
		t.Errorf("expected Origin 'Rhodia-Orbit / Tyrann-Flagship', got %s", d.Origin)
	}
	if d.Target != "Commissioner-Aratap" {
		t.Errorf("expected Target 'Commissioner-Aratap', got %s", d.Target)
	}
	if d.Priority != "URGENT" {
		t.Errorf("expected Priority 'URGENT', got %s", d.Priority)
	}
	if d.Status != "INTERCEPTED" {
		t.Errorf("expected Status 'INTERCEPTED', got %s", d.Status)
	}
	if d.Headers["CODENAME"] != "NEBULA-SHADOW" {
		t.Errorf("expected Header CODENAME 'NEBULA-SHADOW', got %s", d.Headers["CODENAME"])
	}
	if len(d.Payload) != 3 {
		t.Fatalf("expected 3 payload items, got %d", len(d.Payload))
	}
	if d.Payload[0].Type != "EVENT" || d.Payload[0].Content != "Sub-ether tracer affixed to vessel Remembrance." {
		t.Errorf("unexpected payload item 0: %+v", d.Payload[0])
	}
	if d.Authentication != "SIG-ARATAP-9941-OMEGA" {
		t.Errorf("expected Authentication 'SIG-ARATAP-9941-OMEGA', got %s", d.Authentication)
	}

	// Format back
	formatted := Format(d)
	if !strings.Contains(formatted, "DISPATCH HCP-820-001") {
		t.Errorf("formatted text missing DISPATCH header")
	}
	if !strings.Contains(formatted, "END_DISPATCH") {
		t.Errorf("formatted text missing END_DISPATCH marker")
	}

	// Re-parse formatted text
	d2, err := Parse(formatted)
	if err != nil {
		t.Fatalf("failed to re-parse formatted dispatch: %v", err)
	}
	if d2.ID != d.ID || len(d2.Payload) != len(d.Payload) {
		t.Errorf("round trip mismatch: %+v vs %+v", d, d2)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "missing starting dispatch",
			input:   "ENCRYPTION: NONE\nEND_DISPATCH\n",
			wantErr: "missing starting DISPATCH",
		},
		{
			name:    "missing end dispatch",
			input:   "DISPATCH HCP-001\nENCRYPTION: NONE\n",
			wantErr: "missing END_DISPATCH marker",
		},
		{
			name: "missing required origin",
			input: `DISPATCH HCP-002
ENCRYPTION: NONE
TARGET: ALL
STARDATE: 820
PRIORITY: ROUTINE
PAYLOAD:
  EVENT: "test"
AUTHENTICATION: SIG-1
END_DISPATCH`,
			wantErr: "missing required field: ORIGIN",
		},
		{
			name: "empty payload",
			input: `DISPATCH HCP-003
ENCRYPTION: NONE
ORIGIN: Earth
TARGET: ALL
STARDATE: 820
PRIORITY: ROUTINE
PAYLOAD:
AUTHENTICATION: SIG-1
END_DISPATCH`,
			wantErr: "PAYLOAD cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error %q to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
