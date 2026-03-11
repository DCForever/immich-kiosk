# Quickstart: Click name to add, change, or remove birthdate

## 1. Configure birthdate files

1. Ensure the kiosk config points to the JSON birthdate and mapping files, for example:
   - `birthdate_file_path: ./birthdates.json`
   - `birthdate_mapping_path: ./birthdate-mapping.json`
2. Restart or reload the kiosk so the existing age-display feature can load these files.

## 2. Build and run the kiosk

```bash
go test ./...
go run ./cmd/kiosk
```

Open the kiosk UI in a browser pointing at the configured server address.

## 3. Add a birthdate by clicking a name

1. Navigate to a photo or view where people’s names are shown with the new interactive labels.
2. Click a person’s name who currently has no age shown.
3. In the birthdate dialog or panel:
   - Enter their date of birth in `YYYY-MM-DD` format.
   - Submit/save the form.
4. On the next refresh, confirm that the person’s age appears on photos where ages are supported.

## 4. Change an existing birthdate

1. Find a person whose age is already shown.
2. Click their name to open the edit flow.
3. Update the date of birth and save.
4. Verify that their age updates on all relevant photos.

## 5. Remove a birthdate

1. Click the name of a person with a stored birthdate.
2. Choose the remove/delete option in the UI and confirm.
3. Confirm that their age is no longer displayed on any photo.

## 6. Error handling

- If the date is invalid or outside allowed bounds, the UI shows a clear validation message and does not save.
- If the birthdate source is missing or read-only, the UI explains that the change cannot be saved and leaves existing data unchanged.

