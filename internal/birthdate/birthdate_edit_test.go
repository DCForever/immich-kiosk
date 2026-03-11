package birthdate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateDOBForSave(t *testing.T) {
	ref := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		dob     string
		ref     time.Time
		wantErr bool
	}{
		{"valid", "1990-05-15", ref, false},
		{"valid old", "1905-01-01", ref, false},
		{"exactly 120", "1905-06-15", ref, false},
		{"over 120", "1900-01-01", ref, true},
		{"future", "2030-01-01", ref, true},
		{"empty", "", ref, true},
		{"invalid format", "15/05/1990", ref, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ref
			if tt.ref.Year() != 0 {
				r = tt.ref
			}
			err := ValidateDOBForSave(tt.dob, r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDOBForSave(%q) err = %v, wantErr %v", tt.dob, err, tt.wantErr)
			}
		})
	}
}

func TestSetBirthdate_newPerson(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	// Create empty files so SetBirthdate can read and then update them
	_ = os.WriteFile(birthPath, []byte(`{}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{}`), 0644)

	err := SetBirthdate(birthPath, mapPath, "pr-alice", "Alice Smith", "1990-05-15")
	if err != nil {
		t.Fatalf("SetBirthdate() err = %v", err)
	}

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if got := loader.Lookup("pr-alice", ""); got != "1990-05-15" {
		t.Errorf("Lookup(pr-alice) = %q, want 1990-05-15", got)
	}
	if got := loader.Lookup("", "Alice Smith"); got != "1990-05-15" {
		t.Errorf("Lookup(Alice Smith) = %q, want 1990-05-15", got)
	}
}

func TestSetBirthdate_updateExisting(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"pr-alice":"alice"}`), 0644)

	err := SetBirthdate(birthPath, mapPath, "pr-alice", "Alice", "1988-03-20")
	if err != nil {
		t.Fatalf("SetBirthdate() err = %v", err)
	}

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if got := loader.Lookup("pr-alice", ""); got != "1988-03-20" {
		t.Errorf("Lookup(pr-alice) = %q, want 1988-03-20", got)
	}
}

func TestSetBirthdate_invalidDOB(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")

	err := SetBirthdate(birthPath, mapPath, "pr-alice", "", "2030-01-01")
	if err == nil {
		t.Fatal("SetBirthdate(future date) expected error")
	}
	err = SetBirthdate(birthPath, mapPath, "pr-alice", "", "1800-01-01")
	if err == nil {
		t.Fatal("SetBirthdate(>120 years) expected error")
	}
}

func TestSetBirthdate_emptyPerson(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")

	err := SetBirthdate(birthPath, mapPath, "", "", "1990-05-15")
	if err == nil {
		t.Fatal("SetBirthdate(empty person) expected error")
	}
}

func TestRemoveBirthdate_existing(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"pr-alice":"alice","Alice":"alice"}`), 0644)

	err := RemoveBirthdate(birthPath, mapPath, "pr-alice", "Alice")
	if err != nil {
		t.Fatalf("RemoveBirthdate() err = %v", err)
	}

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if got := loader.Lookup("pr-alice", ""); got != "" {
		t.Errorf("Lookup(pr-alice) after remove = %q, want empty", got)
	}
}

func TestRemoveBirthdate_noOpWhenMissing(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"alice":"alice"}`), 0644)

	err := RemoveBirthdate(birthPath, mapPath, "unknown-person", "")
	if err != nil {
		t.Fatalf("RemoveBirthdate(unknown) err = %v", err)
	}
	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if got := loader.Lookup("alice", ""); got != "1990-05-15" {
		t.Errorf("RemoveBirthdate(unknown) should not change existing; Lookup(alice) = %q", got)
	}
}

func TestRemoveBirthdate_emptyPerson(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")

	err := RemoveBirthdate(birthPath, mapPath, "", "")
	if err == nil {
		t.Fatal("RemoveBirthdate(empty person) expected error")
	}
}

func TestWriteMapToJSONFile_atomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	m := map[string]string{"a": "1", "b": "2"}
	err := writeMapToJSONFile(path, m)
	if err != nil {
		t.Fatalf("writeMapToJSONFile() err = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() err = %v", err)
	}
	if len(data) == 0 {
		t.Error("file empty after write")
	}
	// Ensure valid JSON
	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decoding written file: %v", err)
	}
	if decoded["a"] != "1" || decoded["b"] != "2" {
		t.Errorf("decoded = %v", decoded)
	}
}
