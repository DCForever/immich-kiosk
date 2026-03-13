# Specify Prompt: Dynamic Collage Engine

**Purpose**: Use this text as the feature description when invoking `/speckit.specify`.

**How to use**: Run `/speckit.specify` and paste the content from the "Copy-paste block" section below as your feature description.

---

## Copy-paste block (for `/speckit.specify`)

---

I want to replace the current memories collage implementation (fixed CSS grid layouts) with a **dynamic collage engine** that produces more varied, visually pleasing arrangements. The engine should:

1. **Adapt layouts algorithmically** — Generate arrangements that adapt to the number of photos (3–16 per spec 008) and to the actual aspect ratios of the images, rather than forcing photos into fixed grid cells. Layouts should feel organic and varied, not uniform grids.
2. **Preserve aspect ratios** — Avoid cropping or distorting photos. Images should be scaled to fit while maintaining their original proportions. Gaps or spacing between images are acceptable.
3. **Support multiple layout styles** — The system should be able to produce different visual styles (e.g., justified rows like Google Photos, treemap-style packing, hero-plus-cluster arrangements, or randomized organic layouts). At least 3 distinct layout variants per photo count should be possible, with one selected at random per collage render.
4. **Integrate with the existing memories flow** — The output must work with the current kiosk: memories collage in slideshow, on-demand button, 45-second timer, caption, and lightbox-on-tap. The collage is displayed in the browser; the engine can generate either a pre-rendered image (e.g., PNG) or HTML/CSS that the frontend consumes.
5. **Include a research and evaluation phase** — Before implementation, the specification should drive a decision on *how* to achieve dynamic collages. Options to evaluate include:
  - **photocollage** (Python library) — automatic arrangement, aspect-ratio preservation, poster layout
  - **Pillow (PIL)** with custom packing algorithms — treemap, justified rows, recursive packing
  - **ImageMagick** montage — flexible tiling via command-line
  - **Other open-source libraries or tools** — any that support dynamic layout, aspect-ratio preservation, and integration with a Go-based kiosk (e.g., via subprocess, HTTP service, or native bindings)
   The chosen approach must be feasible for this project: Go backend, Immich/PhotoPrism photo sources, and a web frontend. The evaluation should consider maturity, maintainability, output format (image vs. HTML), and performance.
6. **Deliver implementation** — After the evaluation, implement the chosen solution and integrate it into the memories feature so that collages are generated dynamically when memories are requested (slideshow or on-demand).

---

## Notes for the specify skill

- The spec should define **what** the dynamic collage engine must achieve (user value, layout characteristics, integration points) without prescribing a specific technology.
- The research/evaluation of tools (photocollage, Pillow, ImageMagick, etc.) can be captured as a user story or requirement that the *plan* phase will expand into concrete tasks (e.g., "Evaluate and select collage generation approach").
- Success criteria should be measurable and technology-agnostic (e.g., "Collages preserve aspect ratios", "Layouts vary per render", "Integration with memories flow unchanged").

