// Package birthdate provides loading and lookup of person birthdates from
// external JSON files (birthdate file + mapping file). Used by the PhotoPrism
// provider to enrich DisplayAsset.People with BirthDate when ShowPersonAge is
// enabled and paths are configured.
package birthdate

import "time"

// Loader holds in-memory birthdate and mapping data loaded from JSON files.
// Safe for concurrent read (Lookup); load/reload should be done at startup or
// config reload.
type Loader struct {
	birthdates map[string]string // birthdate-file key -> "YYYY-MM-DD"
	mapping    map[string]string // person ID or name (normalized) -> birthdate-file key
}

// Load reads the birthdate and mapping JSON files and populates the loader.
// If either path is empty, no files are loaded. If a file is missing or
// unreadable, a warning is logged and the loader is left empty (or partially
// filled). Invalid DOB entries are skipped. Returns an error only for
// logging; callers should treat a non-nil error as "no birthdates available".
func Load(birthdatePath, mappingPath string) (*Loader, error) {
	return nil, nil
}

// Lookup returns the date of birth (YYYY-MM-DD) for the person identified by
// personID or personName. It looks up by personID first, then by personName.
// Returns empty string if not found or if the loader is nil.
func (l *Loader) Lookup(personID, personName string) string {
	if l == nil {
		return ""
	}
	return ""
}

// Reload re-reads the files at the given paths and replaces the loader's
// in-memory data. Same semantics as Load.
func (l *Loader) Reload(birthdatePath, mappingPath string) error {
	return nil
}

// ParseDOB validates and parses a YYYY-MM-DD string into time.Time (date-only).
func ParseDOB(s string) (time.Time, error) {
	return time.Time{}, nil
}
