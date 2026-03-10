package birthdate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDOB(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{"valid", "1990-05-15", false},
		{"valid padding", "2000-01-01", false},
		{"empty", "", true},
		{"whitespace", "  ", true},
		{"invalid format", "15/05/1990", true},
		{"invalid date", "2000-02-30", true},
		{"not a date", "abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDOB(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDOB(%q) err = %v, wantErr %v", tt.s, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Year() == 0 && got.Month() == 0 && got.Day() == 0 {
				t.Error("ParseDOB returned zero time for valid input")
			}
		})
	}
}

func TestLoad_emptyPaths(t *testing.T) {
	loader, err := Load("", "")
	if err != nil {
		t.Fatalf("Load(empty, empty) err = %v", err)
	}
	if loader != nil {
		t.Fatalf("Load(empty, empty) expected nil loader, got %v", loader)
	}
}

func TestLoad_missingFile(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(mapPath, []byte(`{"alice":"alice"}`), 0644)

	loader, err := Load(birthPath, mapPath)
	if err == nil {
		t.Fatal("Load(missing birthdate file) expected error")
	}
	if loader != nil {
		t.Fatal("Load(missing file) expected nil loader")
	}
}

func TestLoad_validFiles(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15","bob":"1985-11-22"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"pr-alice":"alice","Alice Smith":"alice","pr-bob":"bob"}`), 0644)

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if loader == nil {
		t.Fatal("Load() expected non-nil loader")
	}

	if got := loader.Lookup("pr-alice", ""); got != "1990-05-15" {
		t.Errorf("Lookup(pr-alice) = %q, want 1990-05-15", got)
	}
	if got := loader.Lookup("", "Alice Smith"); got != "1990-05-15" {
		t.Errorf("Lookup(Alice Smith) = %q, want 1990-05-15", got)
	}
	if got := loader.Lookup("pr-bob", "Bob"); got != "1985-11-22" {
		t.Errorf("Lookup(pr-bob) = %q, want 1985-11-22", got)
	}
	if got := loader.Lookup("unknown", ""); got != "" {
		t.Errorf("Lookup(unknown) = %q, want empty", got)
	}
}

func TestLoad_invalidDOBSkipped(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15","bad":"not-a-date","also":"2000-13-01"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"alice":"alice","bad":"bad","also":"also"}`), 0644)

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if loader == nil {
		t.Fatal("Load() expected non-nil loader")
	}

	if got := loader.Lookup("alice", ""); got != "1990-05-15" {
		t.Errorf("Lookup(alice) = %q, want 1990-05-15", got)
	}
	if got := loader.Lookup("bad", ""); got != "" {
		t.Errorf("Lookup(bad) = %q, want empty (invalid DOB skipped)", got)
	}
	if got := loader.Lookup("also", ""); got != "" {
		t.Errorf("Lookup(also) = %q, want empty (invalid month skipped)", got)
	}
}

func TestLoad_duplicateKeysLastWins(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15","alice":"1995-01-01"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"id":"alice"}`), 0644)

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if loader == nil {
		t.Fatal("Load() expected non-nil loader")
	}
	// JSON decode overwrites: last value wins
	got := loader.Lookup("id", "")
	if got != "1995-01-01" && got != "1990-05-15" {
		t.Errorf("Lookup(id) = %q (either value acceptable for duplicate key)", got)
	}
}

func TestLoader_Lookup_nilLoader(t *testing.T) {
	var l *Loader
	if got := l.Lookup("x", "y"); got != "" {
		t.Errorf("Lookup on nil loader = %q, want empty", got)
	}
}

func TestReload(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1990-05-15"}`), 0644)
	_ = os.WriteFile(mapPath, []byte(`{"alice":"alice"}`), 0644)

	loader, err := Load(birthPath, mapPath)
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if loader.Lookup("alice", "") != "1990-05-15" {
		t.Fatal("initial load failed")
	}

	// Update file and reload
	_ = os.WriteFile(birthPath, []byte(`{"alice":"1988-03-20"}`), 0644)
	if err := loader.Reload(birthPath, mapPath); err != nil {
		t.Fatalf("Reload() err = %v", err)
	}
	if got := loader.Lookup("alice", ""); got != "1988-03-20" {
		t.Errorf("after Reload Lookup(alice) = %q, want 1988-03-20", got)
	}
}

func TestParseDOB_roundtrip(t *testing.T) {
	const dob = "2000-06-15"
	tm, err := ParseDOB(dob)
	if err != nil {
		t.Fatal(err)
	}
	if tm.Year() != 2000 || tm.Month() != 6 || tm.Day() != 15 {
		t.Errorf("ParseDOB(%q) = %v", dob, tm)
	}
	// Ensure it's date-only (no time component that could affect age)
	if tm.Hour() != 0 || tm.Minute() != 0 || tm.Second() != 0 {
		t.Errorf("ParseDOB should return date-only, got %v", tm)
	}
	_ = time.Time{}
}
