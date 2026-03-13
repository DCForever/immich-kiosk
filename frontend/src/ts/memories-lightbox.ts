import { getPausedElapsedMs, pausePolling, resumePolling } from "./polling";

const LIGHTBOX_ID = "memories-lightbox";
const LIGHTBOX_OVERLAY_CLASS = "memories-lightbox-overlay";
const LIGHTBOX_IMAGE_CLASS = "memories-lightbox-image";
const LIGHTBOX_CLOSE_CLASS = "memories-lightbox-close";

let clickHandler: ((e: Event) => void) | null = null;
let keyHandler: ((e: KeyboardEvent) => void) | null = null;

/**
 * Opens the lightbox with the given asset image.
 * Pauses polling (without showing menu) so the collage timer is frozen.
 */
export function openLightbox(assetId: string): void {
    const overlay = document.getElementById(LIGHTBOX_ID);
    if (overlay) return;

    const img = document.createElement("img");
    img.className = LIGHTBOX_IMAGE_CLASS;
    img.src = `/image/${assetId}`;
    img.alt = "Memory";

    const closeBtn = document.createElement("button");
    closeBtn.className = LIGHTBOX_CLOSE_CLASS;
    closeBtn.type = "button";
    closeBtn.setAttribute("aria-label", "Close");
    closeBtn.innerHTML = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512" width="32" height="32" aria-hidden="true">
            <path fill="currentColor" d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/>
        </svg>
    `;

    const container = document.createElement("div");
    container.className = LIGHTBOX_OVERLAY_CLASS;
    container.id = LIGHTBOX_ID;
    container.appendChild(img);
    container.appendChild(closeBtn);

    const onClose = () => {
        closeLightbox();
    };

    closeBtn.addEventListener("click", onClose);
    container.addEventListener("click", (e) => {
        if (e.target === container) onClose();
    });
    keyHandler = (e: KeyboardEvent) => {
        if (e.key === "Escape") {
            onClose();
        }
    };
    document.addEventListener("keydown", keyHandler);

    document.body.appendChild(container);

    pausePolling(false);
}

/**
 * Closes the lightbox and resumes polling with remaining time.
 */
export function closeLightbox(): void {
    const overlay = document.getElementById(LIGHTBOX_ID);
    if (!overlay) return;

    overlay.remove();
    document.removeEventListener("keydown", keyHandler!);
    keyHandler = null;

    const elapsed = getPausedElapsedMs();
    resumePolling(true, elapsed);
}

/**
 * Initializes lightbox: attaches click handlers to .collage-cell elements.
 * Call when a memories collage is rendered (data-memories="true" present).
 */
export function initMemoriesLightbox(container: HTMLElement): void {
    if (clickHandler) return;

    clickHandler = (e: Event) => {
        const cell = (e.target as HTMLElement).closest(".collage-cell");
        if (!cell) return;
        const assetId = cell.getAttribute("data-asset-id");
        if (assetId) {
            e.preventDefault();
            openLightbox(assetId);
        }
    };

    container.addEventListener("click", clickHandler);
}

/**
 * Cleans up lightbox handlers. Call when swapping away from memories view.
 */
export function cleanupMemoriesLightbox(container: HTMLElement): void {
    if (clickHandler) {
        container.removeEventListener("click", clickHandler);
        clickHandler = null;
    }
    const overlay = document.getElementById(LIGHTBOX_ID);
    if (overlay) {
        overlay.remove();
    }
}
