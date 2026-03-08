# Features by source

The kiosk supports two media backends: **Immich** and **PhotoPrism**. Most behaviour (slideshow, transitions, config, URL builder) is the same; some options apply only to one source.

| Feature | Immich | PhotoPrism |
|--------|--------|------------|
| **Photos** | ✅ | ✅ |
| **Videos** | ✅ | ✅ (if returned by API) |
| **Albums** | ✅ (album IDs) | ✅ (album UIDs) |
| **Tags / labels** | ✅ | ✅ (labels) |
| **Date range** | ✅ | ✅ |
| **Favourites** | ✅ | ✅ (favourite filter) |
| **People** | ✅ | ✅ (by subject name; subject list if API supports it) |
| **Memories** | ✅ | ❌ Not supported |
| **Star rating** | ✅ | ❌ Not supported |
| **Like / Hide / Tag actions** | ✅ (writes to Immich) | No-op (UI only) |
| **Archive** | ✅ | No-op |
| **URL Builder** | ✅ | ✅ (albums, labels) |
| **Caching** | ✅ | ✅ |
| **Prefetch** | ✅ | ✅ |
| **Offline mode** | ✅ | ✅ (same flow) |

When `source: photoprism`, config options for memories and rating are ignored. People filtering uses PhotoPrism subject names (config `people`); if GET /api/v1/subjects is available, the URL builder shows people in the dropdown.
