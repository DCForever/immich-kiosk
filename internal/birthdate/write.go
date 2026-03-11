// Package birthdate: write.go provides atomic write and edit helpers for birthdate and mapping JSON files.

package birthdate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxAgeYears = 120

// loadBirthdateFileForEdit reads the birthdate file; if it does not exist returns (nil, nil).
func loadBirthdateFileForEdit(path string) (map[string]string, error) {
	m, err := loadBirthdateFile(path)
	if err != nil && (os.IsNotExist(err) || errors.Is(err, os.ErrNotExist)) {
		return nil, nil
	}
	return m, err
}

// loadMappingFileForEdit reads the mapping file; if it does not exist returns (nil, nil).
func loadMappingFileForEdit(path string) (map[string]string, error) {
	m, err := loadMappingFile(path)
	if err != nil && (os.IsNotExist(err) || errors.Is(err, os.ErrNotExist)) {
		return nil, nil
	}
	return m, err
}

// writeMapToJSONFile writes m to path atomically (temp file + rename).
func writeMapToJSONFile(path string, m map[string]string) error {
	if path == "" {
		return fmt.Errorf("birthdate: empty path")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	tmpPath := f.Name()
	defer os.Remove(tmpPath)
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	if err := enc.Encode(m); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// ValidateDOBForSave checks that dob is YYYY-MM-DD, not in the future, and not more than maxAgeYears ago from ref (e.g. time.Now()).
func ValidateDOBForSave(dob string, ref time.Time) error {
	t, err := ParseDOB(dob)
	if err != nil {
		return err
	}
	if t.After(ref) {
		return fmt.Errorf("birth date cannot be in the future")
	}
	refDate := time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, ref.Location())
	age := refDate.Year() - t.Year()
	if refDate.YearDay() < t.YearDay() {
		age--
	}
	if age > maxAgeYears {
		return fmt.Errorf("birth date implies age over %d years", maxAgeYears)
	}
	return nil
}

// SetBirthdate updates or adds the birthdate for the person identified by personID and/or personName.
// It reads the existing birthdate and mapping files, updates in-memory maps, and writes both back atomically.
// Use personID and/or personName; at least one must be non-empty. The birthdate-file key is taken from
// the existing mapping if present, otherwise derived as personID (or personName if personID is empty).
// The next Load() (e.g. on the next request) will see the updated files; no shared loader reload is needed.
func SetBirthdate(birthdatePath, mappingPath, personID, personName, dob string) error {
	personID = strings.TrimSpace(personID)
	personName = strings.TrimSpace(personName)
	if personID == "" && personName == "" {
		return fmt.Errorf("person identifier or name required")
	}
	if err := ValidateDOBForSave(dob, time.Now()); err != nil {
		return err
	}
	dob = strings.TrimSpace(dob)
	if _, err := ParseDOB(dob); err != nil {
		return err
	}

	birthdates, err := loadBirthdateFile(birthdatePath)
	if err != nil {
		return err
	}
	mapping, err := loadMappingFile(mappingPath)
	if err != nil {
		return err
	}
	if birthdates == nil {
		birthdates = make(map[string]string)
	}
	if mapping == nil {
		mapping = make(map[string]string)
	}

	var bkey string
	for _, k := range []string{personID, personName} {
		if k == "" {
			continue
		}
		if bk, ok := mapping[k]; ok {
			bkey = bk
			break
		}
	}
	if bkey == "" {
		if personID != "" {
			bkey = personID
		} else {
			bkey = personName
		}
	}

	birthdates[bkey] = dob
	if personID != "" {
		mapping[personID] = bkey
	}
	if personName != "" && personName != personID {
		mapping[personName] = bkey
	}

	if err := writeMapToJSONFile(birthdatePath, birthdates); err != nil {
		return err
	}
	return writeMapToJSONFile(mappingPath, mapping)
}

// RemoveBirthdate removes the birthdate for the person identified by personID and/or personName.
// It reads the existing files, removes the mapping entries and the corresponding birthdate key, and writes both back.
// If the person has no stored birthdate, this is a no-op and returns nil.
func RemoveBirthdate(birthdatePath, mappingPath, personID, personName string) error {
	personID = strings.TrimSpace(personID)
	personName = strings.TrimSpace(personName)
	if personID == "" && personName == "" {
		return fmt.Errorf("person identifier or name required")
	}

	birthdates, err := loadBirthdateFileForEdit(birthdatePath)
	if err != nil {
		return err
	}
	mapping, err := loadMappingFileForEdit(mappingPath)
	if err != nil {
		return err
	}
	if birthdates == nil {
		birthdates = make(map[string]string)
	}
	if mapping == nil {
		mapping = make(map[string]string)
	}

	var bkey string
	for _, k := range []string{personID, personName} {
		if k == "" {
			continue
		}
		if bk, ok := mapping[k]; ok {
			bkey = bk
			break
		}
	}
	if bkey == "" {
		return nil // no mapping; nothing to remove
	}

	delete(birthdates, bkey)
	if personID != "" {
		delete(mapping, personID)
	}
	if personName != "" {
		delete(mapping, personName)
	}

	if err := writeMapToJSONFile(birthdatePath, birthdates); err != nil {
		return err
	}
	return writeMapToJSONFile(mappingPath, mapping)
}
