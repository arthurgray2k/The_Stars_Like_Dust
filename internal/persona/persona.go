package persona

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Metaopinion represents an exploratory ideological or philosophical reflection.
type Metaopinion struct {
	Topic                  string  `json:"topic"`
	CanonicalStance        string  `json:"canonical_stance"`
	ExploratoryMetaopinion string  `json:"exploratory_metaopinion"`
	CertaintyScore         float64 `json:"certainty_score"`
	NarrativeWeight        string  `json:"narrative_weight"`
}

// PsychologicalDrivers outlines the inner motives and character growth arc.
type PsychologicalDrivers struct {
	PrimaryGoal   string `json:"primary_goal"`
	Vulnerability string `json:"vulnerability"`
	GrowthArc     string `json:"growth_arc"`
}

// Persona represents a loaded character configuration from metadata/.
type Persona struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	Title                string               `json:"title"`
	Allegiance           string               `json:"allegiance"`
	HomeWorld            string               `json:"home_world"`
	Education            string               `json:"education"`
	CanonicalTraits      []string             `json:"canonical_traits"`
	PsychologicalDrivers PsychologicalDrivers `json:"psychological_drivers"`
	Metaopinions         []Metaopinion        `json:"metaopinions"`
	Relationships        map[string]string    `json:"relationships"`
}

// Manager loads and manages character personas.
type Manager struct {
	MetadataDir string
	personas    map[string]*Persona
}

// NewManager creates a persona manager targeting metadataDir.
func NewManager(metadataDir string) *Manager {
	return &Manager{
		MetadataDir: metadataDir,
		personas:    make(map[string]*Persona),
	}
}

// LoadAll loads all persona JSON files from the metadata directory.
func (m *Manager) LoadAll() (map[string]*Persona, error) {
	entries, err := os.ReadDir(m.MetadataDir)
	if err != nil {
		return nil, fmt.Errorf("failed reading metadata directory %s: %w", m.MetadataDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") || entry.Name() == "sector_map.json" {
			continue
		}

		path := filepath.Join(m.MetadataDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed reading persona file %s: %w", path, err)
		}

		var p Persona
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("failed unmarshaling persona %s: %w", path, err)
		}
		if p.ID == "" {
			p.ID = strings.TrimSuffix(entry.Name(), ".json")
		}
		m.personas[p.ID] = &p
	}

	return m.personas, nil
}

// Get retrieves a persona by ID.
func (m *Manager) Get(id string) (*Persona, bool) {
	p, ok := m.personas[id]
	return p, ok
}

// List returns all loaded personas.
func (m *Manager) List() []*Persona {
	list := make([]*Persona, 0, len(m.personas))
	for _, p := range m.personas {
		list = append(list, p)
	}
	return list
}
