# Features by source

The kiosk supports two media backends: **Immich** and **PhotoPrism**. Most behaviour (slideshow, transitions, config, URL builder) is the same; some options apply only to one source.

| Feature | Immich | PhotoPrism |
|--------|--------|------------|
| **Photos** | ✅ | ✅ |
| **Videos** | ✅ | ✅ (if returned by API) |
| **Albums** | ✅ (album IDs) | ✅ (album UIDs, `type=album`; owned and shared use the same list) |
| **Tags / labels** | ✅ | ✅ (labels) |
| **Date range** | ✅ | ✅ |
| **Favourites** | ✅ | ✅ (`favorite=true` filter; like/unlike does not write) |
| **People** | ✅ | ✅ (`GET /api/v1/subjects?type=person`; `person:`, `people:`, `face:`) |
| **Memories** | ✅ | ✅ (search collage: exact day, then week, then month; no memories API) |
| **Star rating** | ✅ | ❌ Not supported (PhotoPrism quality is 1–7, not Immich stars) |
| **Like / Hide / Tag actions** | ✅ (writes to Immich) | No-op (does not write to PhotoPrism) |
| **Archive** | ✅ | No-op (does not write to PhotoPrism) |
| **URL Builder** | ✅ | ✅ (albums, labels) |
| **Caching** | ✅ | ✅ |
| **Prefetch** | ✅ | ✅ |
| **Offline mode** | ✅ | ✅ (same flow) |

When `source: photoprism`, the kiosk talks only to stock PhotoPrism `/api/v1`. Daniel’s PhotoPrism fork adds local vision (OpenAI / LM Studio) and no public API for albums, people, memories, or ratings, so the kiosk does not call those vision endpoints.

Star rating stays unsupported: PhotoPrism quality is a 1–7 scale, not Immich stars, and some app passwords hide photos with quality below 3. Like, hide, tag, and archive stay intentional no-ops. Favourites are selected with `favorite=true` and are not written back.

People come from `GET /api/v1/subjects?type=person`. Search uses `person:"Name"`, `people:"A & B"`, or `face:<uid>` for a subject UID. Subject records have no birth date; age uses the kiosk birthdate mapping only. Albums come from `GET /api/v1/albums?type=album`. Owned and shared album picks use that same list. Photos use `merged=true`, `primary=true`, and album UIDs in `s=`.

Memories are a collage, not an Immich memories API: `taken:"date"` plus `type:image`, then `after:`/`before:` for the ISO week, then the month. A single random memory asset is unsupported.

The client reads `GET /api/v1/config` once per server and logs `version` when debug logging is on. A failed probe does not block the slideshow. JSON decoding accepts `ID` as a string or number, `Files`/`files`, and `Markers`/`markers`, plus empty dates, numeric booleans, and wrapped lists. Preview and download tokens come from `X-Preview-Token` and `X-Download-Token`. Thumbnails use `/api/v1/t/{hash}/{token}/{size}` (`fit_720`, then `fit_1280`, then `fit_1920`). Video uses `/api/v1/videos/{hash}/{token}/avc`.
