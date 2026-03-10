# Feature Specification: Show age for people on photos

**Feature Branch**: `005-show-person-age`  
**Created**: 2026-03-09  
**Status**: Draft  
**Input**: User description: "I want to add support for showing years old for people based on when the photos is taken. Photoprism doesn't support birthdates for people so we need another source for them."

## Clarifications

### Session 2026-03-09

- Q: What form should the external birthdate source take? → A: Single local file (e.g. JSON or CSV); path configured by administrator (config or env).
- Q: How should the kiosk match a Photoprism person to a record in the birthdate file? → A: By administrator-defined mapping (e.g. separate mapping file or table: Photoprism ID or name → birthdate file key).
- Q: If the birthdate file (or mapping file) is missing or unreadable, how should the kiosk behave? → A: Show photos and people as usual; do not show any ages until the file is present and valid. Log a warning but do not surface an error to the viewer.
- Q: What is the maximum age in years to treat as valid? → A: 120 years (fixed; any age above is treated as invalid).
- Q: For someone under 1 year old at photo time, how should the age be shown? → A: In months.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See age for each person on a photo (Priority: P1)

As a viewer of the photo kiosk, I want to see how old each identified person was when the photo was taken so that I can better appreciate the context and time in their life.

**Why this priority**: This is the core value of the feature and directly answers the user's request; without this, the feature delivers no visible benefit.

**Independent Test**: Display a photo with one or more identified people where birthdates are known from the external source; verify that each person has an age (in years, or in months when under 1 year) that matches the difference between the photo date and their birthdate.

**Acceptance Scenarios**:

1. **Given** a photo with a single identified person whose birthdate exists in the external source, **When** the photo is shown on the kiosk, **Then** that person’s age (in years, or in months when under 1) is shown alongside their name based on the photo’s taken date.  
2. **Given** a photo with multiple identified people whose birthdates all exist in the external source, **When** the photo is shown on the kiosk, **Then** each person’s age (years or months as applicable) is shown next to that person and all ages are calculated consistently from the same photo date.

---

### User Story 2 - Handle missing or partial birthdate data gracefully (Priority: P2)

As a viewer, I do not want to see confusing or obviously wrong ages when the system does not have full birthdate data, so that the display remains trustworthy.

**Why this priority**: Birthdate coverage may be incomplete; the feature must fail gracefully to avoid misleading information or cluttering the UI.

**Independent Test**: Show photos for people whose birthdate is missing or incomplete in the external source and confirm that no age is displayed and that the rest of the photo information still appears normally.

**Acceptance Scenarios**:

1. **Given** a photo with an identified person whose birthdate is not present in the external source, **When** the photo is displayed, **Then** the person’s name is still shown but no age value or placeholder age is displayed.  
2. **Given** a photo with an identified person whose birthdate is present but clearly invalid (e.g., in the future or outside a reasonable human range), **When** the photo is displayed, **Then** the system treats the birthdate as unusable and does not show an age for that person.

---

### User Story 3 - Respect configuration of external birthdate source (Priority: P3)

As the kiosk administrator, I want to configure and maintain a separate source of birthdate data for people so that ages can be shown even though Photoprism does not provide birthdates.

**Why this priority**: The feature depends on a secondary data source; administrators must be able to set it up and understand what is required for ages to appear.

**Independent Test**: Configure the external birthdate source with one new person record, restart or refresh the kiosk as required, and confirm that ages begin to appear for that person’s photos while other people without birthdates still show no age.

**Acceptance Scenarios**:

1. **Given** the external birthdate source is configured with a valid mapping between a person and their birthdate, **When** the kiosk is reloaded or refreshed according to documented behavior, **Then** photos containing that person show the correct age.  
2. **Given** the external birthdate source is edited to change a person’s birthdate to a different valid value, **When** the kiosk is reloaded or refreshed according to documented behavior, **Then** ages for that person on all photos update to reflect the new birthdate.

---

### Edge Cases

- A photo has no reliable taken date; the system cannot compute age and therefore does not display an age for any person on that photo.  
- A person appears in a photo taken before their birthdate (e.g., incorrect metadata); the system treats the age as invalid and does not show an age for that person on that photo.  
- The external birthdate source contains multiple entries for the same person (e.g., duplicate names or identifiers); the system uses a deterministic rule to choose a single birthdate and, if the entries conflict, prioritizes the one that appears valid and within a reasonable range.  
- The computed age exceeds 120 years; the system treats the birthdate or metadata as invalid and does not display an age.  
- Time zone differences between the photo taken date and birthdate representation do not cause the displayed age (years or months) to be off by one unit.  
- **Under 1 year old**: When the person’s age at photo time is less than 1 year, the system displays age in whole months (e.g. “6 months”) rather than years.  
- **Birthdate or mapping file missing or unreadable**: The system shows photos and people as usual and does not show any ages until the file(s) are present and valid; the system logs a warning and does not surface an error to the viewer.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST calculate a person’s age from the photo’s taken date and that person’s birthdate: in whole years when age ≥ 1 year, and in whole months when age &lt; 1 year.  
- **FR-002**: The system MUST display the calculated age together with the person’s existing label or name wherever people are shown on a photo (years when ≥ 1, months when &lt; 1).  
- **FR-003**: The system MUST obtain birthdates for people from a single local file (e.g. JSON or CSV) whose path is configured by the administrator (via config or environment variable), independent of Photoprism’s people model.  
- **FR-004**: The system MUST match people detected on photos to records in the birthdate file using an administrator-defined mapping (e.g. a separate mapping file or table that maps Photoprism person ID or name to the key used in the birthdate file).  
- **FR-005**: The system MUST avoid showing an age when either the photo’s taken date or the person’s birthdate is missing, clearly invalid, or when the computed age is greater than 120 years.  
- **FR-006**: The system MUST handle cases where a person is present on a photo but has no birthdate in the external source without causing errors or visual glitches in the photo display.  
- **FR-009**: If the birthdate file or mapping file is missing or unreadable, the system MUST show photos and people as usual without displaying ages, MUST log a warning, and MUST NOT show an error to the viewer.  
- **FR-007**: The system MUST allow the administrator to update the external birthdate source and have those changes reflected in subsequent age calculations without requiring a full recreation of the kiosk configuration.  
- **FR-008**: The system MUST ensure that any age display clearly corresponds to the time the photo was taken, not the current date.

### Key Entities *(include if feature involves data)*

- **Person**: A human individual recognized on photos, identified by a stable identifier and a display label or name.  
- **Birthdate Source**: A single local file (JSON or CSV) whose path is set via configuration or environment variable; contains records that associate a person identifier (birthdate-file key) with a date of birth.  
- **Person–birthdate mapping**: Administrator-defined mapping (e.g. separate mapping file or table) from Photoprism person ID or name to the key used in the birthdate file.  
- **Birthdate Source Record**: A record in the birthdate file associating a birthdate-file key with a date of birth and optional metadata.  
- **Photo**: An image with a taken date (or best available approximation) and zero or more associated people.  
- **Age Display**: The age derived from the photo taken date and the birthdate source record: whole years when the person is at least 1 year old, whole months when under 1 year old.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For a test set of photos where birthdates are known and valid for all people, at least 99% of displayed ages match the expected age (in years, or in months when under 1 year) based on the photo taken date.  
- **SC-002**: For photos where birthdates are missing or invalid for some people, 100% of those people show no age rather than an incorrect or placeholder value, while other people on the same photo with valid data still show ages.  
- **SC-003**: After configuring or updating the external birthdate source according to documentation, an administrator can see correct ages on relevant photos within one normal refresh or reload cycle of the kiosk.  
- **SC-004**: In a usability check with representative viewers, at least 90% report that the displayed ages are clear, understandable, and appear to reflect the person’s age at the time the photo was taken.

