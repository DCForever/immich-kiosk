package photoprism

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// flexString unmarshals from either a JSON string or number.
// PhotoPrism list responses use a string ID; detail responses use a number.
type flexString string

func (s *flexString) UnmarshalJSON(b []byte) error {
	if isNull(b) {
		*s = ""
		return nil
	}
	*s = flexString(rawToString(b))
	return nil
}

func isNull(b []byte) bool {
	trimmed := bytes.TrimSpace(b)
	return len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null"))
}

func rawToString(b []byte) string {
	if isNull(b) {
		return ""
	}
	var s string
	if json.Unmarshal(b, &s) == nil {
		return s
	}
	var n json.Number
	if json.Unmarshal(b, &n) == nil {
		return n.String()
	}
	return ""
}

func decodeString(fields map[string]json.RawMessage, key string) string {
	raw, ok := fields[key]
	if !ok {
		return ""
	}
	return rawToString(raw)
}

func decodeBool(fields map[string]json.RawMessage, key string) bool {
	raw, ok := fields[key]
	if !ok || isNull(raw) {
		return false
	}
	var b bool
	if json.Unmarshal(raw, &b) == nil {
		return b
	}
	if n, ok := decodeNumber(raw); ok {
		f, err := n.Float64()
		return err == nil && f != 0
	}
	switch strings.ToLower(strings.TrimSpace(rawToString(raw))) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

func decodeInt(fields map[string]json.RawMessage, key string) int {
	return int(decodeInt64(fields, key))
}

func decodeInt64(fields map[string]json.RawMessage, key string) int64 {
	raw, ok := fields[key]
	if !ok || isNull(raw) {
		return 0
	}
	if n, ok := decodeNumber(raw); ok {
		return numberToInt64(n)
	}
	s := strings.TrimSpace(rawToString(raw))
	if s == "" {
		return 0
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(f)
}

func decodeFloat(fields map[string]json.RawMessage, key string) float64 {
	raw, ok := fields[key]
	if !ok || isNull(raw) {
		return 0
	}
	if n, ok := decodeNumber(raw); ok {
		f, err := n.Float64()
		if err != nil {
			return 0
		}
		return f
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(rawToString(raw)), 64)
	if err != nil {
		return 0
	}
	return f
}

func decodeNumber(raw json.RawMessage) (json.Number, bool) {
	var n json.Number
	if json.Unmarshal(raw, &n) != nil {
		return "", false
	}
	return n, true
}

func numberToInt64(n json.Number) int64 {
	if i, err := n.Int64(); err == nil {
		return i
	}
	f, err := n.Float64()
	if err != nil {
		return 0
	}
	return int64(f)
}

func decodeTime(fields map[string]json.RawMessage, key string) time.Time {
	raw, ok := fields[key]
	if !ok || isNull(raw) {
		return time.Time{}
	}
	if n, ok := decodeNumber(raw); ok {
		return unixFromNumber(n)
	}
	return parseTime(strings.TrimSpace(rawToString(raw)))
}

func unixFromNumber(n json.Number) time.Time {
	f, err := n.Float64()
	if err != nil || f <= 0 {
		return time.Time{}
	}
	if f > 1e12 {
		return time.UnixMilli(int64(f)).UTC()
	}
	return time.Unix(int64(f), 0).UTC()
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}
