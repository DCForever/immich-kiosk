# Quickstart: Show age for people on photos (PhotoPrism)

**Feature**: 005-show-person-age

## Prerequisites

- Kiosk configured with PhotoPrism as the source.
- `show_person_name` and/or `show_person_age` enabled (config or URL).

## 1. Create the birthdate file

Create a JSON file that maps **logical keys** to dates of birth (YYYY-MM-DD). Example:

**birthdates.json**:

```json
{
  "alice": "1990-05-15",
  "bob": "1985-11-22",
  "charlie": "2020-03-01"
}
```

Keys (`alice`, `bob`, `charlie`) are internal identifiers you will reference from the mapping file.

## 2. Create the mapping file

Create a JSON file that maps **Photoprism person ID or display name** to the key used in the birthdate file. Get person IDs from PhotoPrism (e.g. subject UID) or use the exact display name.

**birthdate-mapping.json**:

```json
{
  "pr-subj-abc123": "alice",
  "Alice Smith": "alice",
  "pr-subj-def456": "bob",
  "Bob Jones": "bob"
}
```

- Left-hand side: Photoprism subject UID or the person’s name as shown in PhotoPrism.
- Right-hand side: Key from the birthdate file.

If both ID and name are present, either can match; the loader looks up by person ID first, then by name.

## 3. Configure the kiosk

Set the paths to both files (config file or environment variables).

**config.yaml**:

```yaml
birthdate_file_path: /path/to/birthdates.json
birthdate_mapping_path: /path/to/birthdate-mapping.json
```

**Environment variables**:

- `KIOSK_BIRTHDATE_FILE_PATH=/path/to/birthdates.json`
- `KIOSK_BIRTHDATE_MAPPING_PATH=/path/to/birthdate-mapping.json`

Ensure `show_person_age` is enabled (e.g. `show_person_age: true` or URL `?show_person_age=true`).

## 4. Run and verify

Start or reload the kiosk. Photos that contain people with a mapping and valid birthdate will show ages (in years, or in months when under 1). If either file is missing or invalid, the kiosk logs a warning and shows no ages; photos and names still display.

- **Reload**: After editing the birthdate or mapping file, reload config (or restart) so the loader picks up changes.
- **120-year cap**: Ages over 120 years are not shown.
- **Under 1 year**: Shown in whole months (e.g. "6 months").

## Contract reference

- Birthdate file schema: [contracts/birthdate-file.json](./contracts/birthdate-file.json)
- Mapping file schema: [contracts/birthdate-mapping.json](./contracts/birthdate-mapping.json)
