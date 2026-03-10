// Package birthdate provides loading and lookup of person birthdates from
// external JSON files (birthdate file + mapping file). Used by the PhotoPrism
// provider to enrich DisplayAsset.People with BirthDate when ShowPersonAge is
// enabled and paths are configured.
package birthdate

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"log/slog"
)

// Loader holds in-memory birthdate and mapping data loaded from JSON files.
// Safe for concurrent read (Lookup); load/reload should be done at startup or
// config reload.
type Loader struct {
	birthdates map[string]string // birthdate-file key -> "YYYY-MM-DD"
	mapping    map[string]string // person ID or name (normalized) -> birthdate-file key
}

// Load reads the birthdate and mapping JSON files and populates the loader.
// If either path is empty, no files are loaded and (nil, nil) is returned.
// If a file is missing or unreadable, a warning is logged and (nil, err) is
// returned so callers treat as "no birthdates available". Invalid DOB entries
// are skipped. Returns an error only for logging; callers should treat a
// non-nil error as "no birthdates available".
func Load(birthdatePath, mappingPath string) (*Loader, error) {
	if strings.TrimSpace(birthdatePath) == "" || strings.TrimSpace(mappingPath) == "" {
		return nil, nil
	}

	birthdates, err := loadBirthdateFile(birthdatePath)
	if err != nil {
		slog.Warn("birthdate: could not load birthdate file", "path", birthdatePath, "err", err)
		return nil, err
	}

	mapping, err := loadMappingFile(mappingPath)
	if err != nil {
		slog.Warn("birthdate: could not load mapping file", "path", mappingPath, "err", err)
		return nil, err
	}

	return &Loader{birthdates: birthdates, mapping: mapping}, nil
}

func loadBirthdateFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		v = strings.TrimSpace(v)
		if _, err := ParseDOB(v); err != nil {
			continue // skip invalid
		}
		out[strings.TrimSpace(k)] = v
	}
	return out, nil
}

func loadMappingFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		key := strings.TrimSpace(k)
		val := strings.TrimSpace(v)
		if key != "" && val != "" {
			out[key] = val
		}
	}
	return out, nil
}

// Lookup returns the date of birth (YYYY-MM-DD) for the person identified by
// personID or personName. It looks up by personID first, then by personName.
// Returns empty string if not found or if the loader is nil.
func (l *Loader) Lookup(personID, personName string) string {
	if l == nil || l.birthdates == nil || l.mapping == nil {
		return ""
	}
	personID = strings.TrimSpace(personID)
	personName = strings.TrimSpace(personName)
	for _, key := range []string{personID, personName} {
		if key == "" {
			continue
		}
		if bkey, ok := l.mapping[key]; ok {
			if dob, ok := l.birthdates[bkey]; ok {
				return dob
			}
		}
	}
	return ""
}

// Reload re-reads the files at the given paths and replaces the loader's
// in-memory data. Same semantics as Load. If Load would return nil, Reload
// clears the loader's maps.
func (l *Loader) Reload(birthdatePath, mappingPath string) error {
	if l == nil {
		return nil
	}
	if strings.TrimSpace(birthdatePath) == "" || strings.TrimSpace(mappingPath) == "" {
		l.birthdates = nil
		l.mapping = nil
		return nil
	}
	birthdates, err := loadBirthdateFile(birthdatePath)
	if err != nil {
		slog.Warn("birthdate: could not reload birthdate file", "path", birthdatePath, "err", err)
		return err
	}
	mapping, err := loadMappingFile(mappingPath)
	if err != nil {
		slog.Warn("birthdate: could not reload mapping file", "path", mappingPath, "err", err)
		return err
	}
	l.birthdates = birthdates
	l.mapping = mapping
	return nil
}

// ParseDOB validates and parses a YYYY-MM-DD string into time.Time (date-only).
func ParseDOB(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("empty date")
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
