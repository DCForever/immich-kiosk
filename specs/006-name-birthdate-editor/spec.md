# Feature Specification: Click name to add, change, or remove birthdate

**Feature Branch**: `006-name-birthdate-editor`  
**Created**: 2026-03-10  
**Status**: Draft  
**Input**: User description: "I want to support clicking on a name in the UI to add/change/remove their birthdate"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add a birthdate by clicking a person's name (Priority: P1)

As a user of the kiosk, I want to click on a person's name when they have no birthdate set so that I can enter their date of birth and have ages appear for that person on photos.

**Why this priority**: Adding birthdates is the primary way to populate or extend birthdate data; without it, the feature cannot improve coverage.

**Independent Test**: Open the kiosk, navigate to a photo that shows a person who does not yet have a birthdate stored; click that person's name, enter a valid date of birth, and confirm. Verify that the birthdate is saved and that on this photo (and others containing that person) the person's age is now shown where supported.

**Acceptance Scenarios**:

1. **Given** a photo is displayed with at least one person whose name is shown and who has no stored birthdate, **When** the user clicks that person's name, **Then** the user can enter a date of birth for that person and save it.
2. **Given** the user has entered a valid date of birth for a person and confirmed, **When** the save completes, **Then** the birthdate is stored and subsequent display of that person (on the same or other photos) shows their age when the birthdate source is in use.

---

### User Story 2 - Change an existing birthdate by clicking a person's name (Priority: P2)

As a user of the kiosk, I want to click on a person's name when they already have a birthdate so that I can correct or update it.

**Why this priority**: Corrections and updates are necessary for data quality; the feature must support full lifecycle of birthdate data.

**Independent Test**: With at least one person who already has a stored birthdate, click that person's name, change the date to a different valid value, and save. Verify that the new date is stored and that ages shown for that person reflect the updated birthdate.

**Acceptance Scenarios**:

1. **Given** a photo is displayed with a person who already has a stored birthdate, **When** the user clicks that person's name, **Then** the user can see the current birthdate and change it to a new value and save.
2. **Given** the user has changed a person's birthdate and saved, **When** the save completes, **Then** the updated birthdate is stored and ages for that person on all relevant photos reflect the new date.

---

### User Story 3 - Remove a birthdate by clicking a person's name (Priority: P3)

As a user of the kiosk, I want to click on a person's name and remove their stored birthdate so that ages are no longer shown for that person when I no longer want that information stored.

**Why this priority**: Removal supports privacy and data correction; completes the add/change/remove set.

**Independent Test**: With a person who has a stored birthdate, click their name, choose to remove the birthdate, and confirm. Verify that the birthdate is no longer stored and that the person's age is no longer displayed on photos.

**Acceptance Scenarios**:

1. **Given** a photo is displayed with a person who has a stored birthdate, **When** the user clicks that person's name and chooses to remove the birthdate, **Then** the user can confirm removal and the birthdate is deleted from storage.
2. **Given** the user has removed a person's birthdate and confirmed, **When** the removal completes, **Then** that person no longer has a birthdate stored and no age is shown for them on any photo.

---

### Edge Cases

- **Birthdate source not configured or read-only**: When the system cannot persist birthdate data (e.g. no configured birthdate source or the storage is read-only), the user is informed that add/change/remove is not available or could not be saved; no partial or silent failure.
- **Invalid date entered**: When the user enters a date that is invalid (e.g. not a real date, or outside the accepted range such as future date or age over 120 years), the system does not save it and informs the user what is wrong so they can correct it.
- **User cancels without saving**: When the user opens the add/change flow by clicking a name but cancels or closes without saving, no change is made to stored data.
- **Same person identified under different names**: If the same person can appear under more than one identifier in the system, adding or changing a birthdate for one representation is scoped as defined by the existing mapping (e.g. one entry per mapped identity); behavior is consistent with how ages are resolved for display.
- **Concurrent or external edit**: If the birthdate source is modified by another process or user while the kiosk is open, the system either reflects the external change after the next refresh or clearly indicates that the stored data may have changed and allows the user to retry or reload.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow the user to trigger an add, change, or remove birthdate flow by clicking on a person's name where that name is displayed in the UI (e.g. on a photo).
- **FR-002**: The system MUST allow the user to enter or edit a date of birth for that person when adding or changing, and MUST validate the date (e.g. valid calendar date, within an accepted range such as not in the future and age at photo time not over 120 years) before saving.
- **FR-003**: The system MUST allow the user to remove a stored birthdate for that person when one exists, with explicit confirmation so that removal is intentional.
- **FR-004**: The system MUST persist added or changed birthdates to the same configured birthdate source used for displaying ages, and MUST remove the birthdate from that source when the user removes it, so that display of ages stays consistent with stored data.
- **FR-005**: The system MUST associate the added or changed birthdate with the correct person (the one whose name was clicked) using the same mapping mechanism used for age display, so that the right record is updated.
- **FR-006**: When the birthdate source or mapping is not configured, not writable, or persistence fails, the system MUST inform the user that the action could not be completed and MUST NOT silently fail or leave data in an inconsistent state.
- **FR-007**: When the user cancels or closes the add/change/remove flow without confirming save or remove, the system MUST make no change to stored birthdate data.

### Key Entities

- **Person**: A human individual shown on a photo, identified by a stable identifier and a display name; the same entity as in the existing age-display feature.
- **Birthdate source**: The configured store for dates of birth (e.g. a file or service) that is used both for displaying ages and for saving add/change/remove actions from this feature.
- **Person–birthdate mapping**: The configured mapping that links persons shown in the UI to keys in the birthdate source; used to resolve which record to update when the user clicks a name.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can add a birthdate for a person by clicking their name and completing the flow in under 30 seconds when the birthdate source is configured and writable.
- **SC-002**: After adding, changing, or removing a birthdate, the stored data matches the user's action and ages shown for that person (or absence of age after removal) are correct on the next display refresh.
- **SC-003**: When the user enters an invalid date or attempts to save when persistence is unavailable, 100% of such attempts result in a clear message to the user and no incorrect or partial data written.
- **SC-004**: In a usability check, at least 85% of users successfully complete add, change, and remove flows on first attempt without assistance when the birthdate source is properly configured.

## Dependencies and Assumptions

- **Dependency**: This feature assumes the existence of a configured birthdate source and mapping used to show ages for people on photos (as in the show-person-age feature). Click-to-edit reads from and writes to that same source.
- **Assumption**: Users who can view the kiosk are permitted to add, change, and remove birthdates; no separate permission or role is required unless otherwise specified.
- **Assumption**: The birthdate source supports the operations required (read, add/update, delete) for the keys and mapping used by the kiosk; if the backing store is append-only or read-only, the feature surfaces that limitation to the user instead of silently failing.
