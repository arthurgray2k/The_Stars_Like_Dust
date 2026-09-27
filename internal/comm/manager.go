package comm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arthurgray2k/The_Stars_Like_Dust/pkg/dsl"
)

// Manager manages reading, writing, and validating HCP dispatches in the comm/ directory.
type Manager struct {
	CommDir string
}

// NewManager creates a new communication manager targeting the given directory.
func NewManager(commDir string) *Manager {
	return &Manager{CommDir: commDir}
}

// LoadAll reads and parses all .hcp dispatch files in the comm directory.
func (m *Manager) LoadAll() ([]*dsl.Dispatch, error) {
	entries, err := os.ReadDir(m.CommDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*dsl.Dispatch{}, nil
		}
		return nil, fmt.Errorf("failed to read comm dir %s: %w", m.CommDir, err)
	}

	dispatches := make([]*dsl.Dispatch, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hcp") {
			continue
		}
		path := filepath.Join(m.CommDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}
		disp, err := dsl.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("failed to parse dispatch in %s: %w", path, err)
		}
		dispatches = append(dispatches, disp)
	}

	return dispatches, nil
}

// SaveDispatch writes a dispatch to the comm directory in HCP format.
func (m *Manager) SaveDispatch(filename string, d *dsl.Dispatch) error {
	if err := os.MkdirAll(m.CommDir, 0755); err != nil {
		return err
	}
	content := dsl.Format(d)
	path := filepath.Join(m.CommDir, filename)
	return os.WriteFile(path, []byte(content), 0644)
}

// FilterByTarget finds dispatches intended for a specific target actor or broadcast.
func FilterByTarget(dispatches []*dsl.Dispatch, target string) []*dsl.Dispatch {
	result := make([]*dsl.Dispatch, 0)
	targetUpper := strings.ToUpper(target)
	for _, d := range dispatches {
		dt := strings.ToUpper(d.Target)
		if strings.Contains(dt, targetUpper) || strings.Contains(dt, "ALL") || strings.Contains(dt, "BROADCAST") {
			result = append(result, d)
		}
	}
	return result
}
