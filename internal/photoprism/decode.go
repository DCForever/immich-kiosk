package photoprism

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// listWrapperKeys are object keys that may wrap a collection.
// Stock PhotoPrism returns a bare array; some proxies and forks wrap it.
var listWrapperKeys = []string{
	"photos", "albums", "labels", "subjects", "results", "items", "data",
}

func indexRaw(raw map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(raw))
	for key, value := range raw {
		out[strings.ToLower(key)] = value
	}
	return out
}

func objectFields(b []byte) (map[string]json.RawMessage, error) {
	if isNull(b) {
		return nil, nil
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return indexRaw(raw), nil
}

func decodeBody(body []byte, dest any) error {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return emptyBody(dest)
	}
	err := json.Unmarshal(trimmed, dest)
	if err == nil {
		return nil
	}
	if !isSlicePtr(dest) {
		return fmt.Errorf("photoprism api json: %w", err)
	}
	wrapped, ok := unwrapList(trimmed)
	if !ok {
		return fmt.Errorf("photoprism api json: %w", err)
	}
	if isNull(wrapped) {
		return setEmptySlice(dest)
	}
	if err := json.Unmarshal(wrapped, dest); err != nil {
		return fmt.Errorf("photoprism api json: %w", err)
	}
	return nil
}

func emptyBody(dest any) error {
	if isSlicePtr(dest) {
		return setEmptySlice(dest)
	}
	return fmt.Errorf("photoprism api json: empty body")
}

func isSlicePtr(dest any) bool {
	rv := reflect.ValueOf(dest)
	return rv.Kind() == reflect.Pointer && !rv.IsNil() && rv.Elem().Kind() == reflect.Slice
}

func setEmptySlice(dest any) error {
	rv := reflect.ValueOf(dest)
	if !isSlicePtr(dest) {
		return nil
	}
	elem := rv.Elem()
	elem.Set(reflect.MakeSlice(elem.Type(), 0, 0))
	return nil
}

func unwrapList(body []byte) ([]byte, bool) {
	fields, err := objectFields(body)
	if err != nil || fields == nil {
		return nil, false
	}
	for _, key := range listWrapperKeys {
		if raw, ok := fields[key]; ok {
			return raw, true
		}
	}
	return nil, false
}

func decodeSlice[T any](fields map[string]json.RawMessage, key string) []T {
	raw, ok := fields[key]
	if !ok || isNull(raw) {
		return nil
	}
	var items []json.RawMessage
	if json.Unmarshal(raw, &items) != nil {
		return nil
	}
	out := make([]T, 0, len(items))
	for _, item := range items {
		var value T
		if json.Unmarshal(item, &value) != nil {
			continue
		}
		out = append(out, value)
	}
	return out
}

func (p *Photo) UnmarshalJSON(b []byte) error {
	fields, err := objectFields(b)
	if err != nil || fields == nil {
		return err
	}
	assignPhotoIdentity(p, fields)
	assignPhotoCapture(p, fields)
	assignPhotoPlace(p, fields)
	p.Files = decodeSlice[File](fields, "files")
	return nil
}

func assignPhotoIdentity(p *Photo, fields map[string]json.RawMessage) {
	p.ID = flexString(decodeString(fields, "id"))
	p.UID = decodeString(fields, "uid")
	p.Type = decodeString(fields, "type")
	p.Hash = decodeString(fields, "hash")
	p.Title = decodeString(fields, "title")
	p.Description = decodeString(fields, "description")
	p.Favorite = decodeBool(fields, "favorite")
	p.FileUID = decodeString(fields, "fileuid")
	p.FileName = decodeString(fields, "filename")
}

func assignPhotoCapture(p *Photo, fields map[string]json.RawMessage) {
	p.TakenAt = decodeTime(fields, "takenat")
	p.TakenAtLocal = decodeTime(fields, "takenatlocal")
	p.TimeZone = decodeString(fields, "timezone")
	p.Year = decodeInt(fields, "year")
	p.Month = decodeInt(fields, "month")
	p.Day = decodeInt(fields, "day")
	p.Iso = decodeInt(fields, "iso")
	p.FocalLength = decodeInt(fields, "focallength")
	p.FNumber = decodeFloat(fields, "fnumber")
	p.Exposure = decodeString(fields, "exposure")
	p.Width = decodeInt(fields, "width")
	p.Height = decodeInt(fields, "height")
	p.Portrait = decodeBool(fields, "portrait")
}

func assignPhotoPlace(p *Photo, fields map[string]json.RawMessage) {
	p.CameraMake = decodeString(fields, "cameramake")
	p.CameraModel = decodeString(fields, "cameramodel")
	p.LensModel = decodeString(fields, "lensmodel")
	p.Lat = decodeFloat(fields, "lat")
	p.Lng = decodeFloat(fields, "lng")
	p.PlaceCity = decodeString(fields, "placecity")
	p.PlaceState = decodeString(fields, "placestate")
	p.PlaceCountry = decodeString(fields, "placecountry")
	p.PlaceLabel = decodeString(fields, "placelabel")
}

func (f *File) UnmarshalJSON(b []byte) error {
	fields, err := objectFields(b)
	if err != nil || fields == nil {
		return err
	}
	f.UID = decodeString(fields, "uid")
	f.Hash = decodeString(fields, "hash")
	f.Width = decodeInt(fields, "width")
	f.Height = decodeInt(fields, "height")
	f.Primary = decodeBool(fields, "primary")
	f.FileType = decodeString(fields, "filetype")
	f.MediaType = decodeString(fields, "mediatype")
	f.Mime = decodeString(fields, "mime")
	f.Size = decodeInt64(fields, "size")
	f.Markers = decodeSlice[Marker](fields, "markers")
	return nil
}

func (m *Marker) UnmarshalJSON(b []byte) error {
	fields, err := objectFields(b)
	if err != nil || fields == nil {
		return err
	}
	m.UID = decodeString(fields, "uid")
	m.Name = decodeString(fields, "name")
	m.SubjUID = firstField(fields, "subjuid", "subjectuid")
	return nil
}

func (a *Album) UnmarshalJSON(b []byte) error {
	fields, err := objectFields(b)
	if err != nil || fields == nil {
		return err
	}
	a.UID = decodeString(fields, "uid")
	a.Title = decodeString(fields, "title")
	a.PhotoCount = decodeInt(fields, "photocount")
	return nil
}

func (l *Label) UnmarshalJSON(b []byte) error {
	fields, err := objectFields(b)
	if err != nil || fields == nil {
		return err
	}
	l.UID = decodeString(fields, "uid")
	l.Name = decodeString(fields, "name")
	l.Slug = decodeString(fields, "slug")
	return nil
}

func (s *Subject) UnmarshalJSON(b []byte) error {
	fields, err := objectFields(b)
	if err != nil || fields == nil {
		return err
	}
	s.UID = decodeString(fields, "uid")
	s.Name = decodeString(fields, "name")
	s.BirthDate = decodeString(fields, "birthdate")
	return nil
}

func firstField(fields map[string]json.RawMessage, keys ...string) string {
	for _, key := range keys {
		if value := decodeString(fields, key); value != "" {
			return value
		}
	}
	return ""
}

func compactPhotos(list []Photo) []Photo {
	out := make([]Photo, 0, len(list))
	for _, photo := range list {
		if photo.UID != "" || photo.Hash != "" {
			out = append(out, photo)
		}
	}
	return out
}

func headerMap(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for key, values := range h {
		if len(values) > 0 {
			out[strings.ToLower(key)] = values[0]
		}
	}
	return out
}

func absorbTokens(headers map[string]string, body []byte) {
	copyAlias(headers, "x-preview-token", "x-photoprism-preview-token", "preview-token")
	copyAlias(headers, "x-download-token", "x-photoprism-download-token", "download-token")
	if headers["x-preview-token"] != "" && headers["x-download-token"] != "" {
		return
	}
	fields, err := objectFields(body)
	if err != nil || fields == nil {
		return
	}
	if headers["x-preview-token"] == "" {
		headers["x-preview-token"] = firstField(fields, "previewtoken", "preview_token")
	}
	if headers["x-download-token"] == "" {
		headers["x-download-token"] = firstField(fields, "downloadtoken", "download_token")
	}
}

func copyAlias(headers map[string]string, canonical string, aliases ...string) {
	if headers[canonical] != "" {
		return
	}
	for _, alias := range aliases {
		if headers[alias] != "" {
			headers[canonical] = headers[alias]
			return
		}
	}
}
