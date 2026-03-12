/**
 * Birthdate edit: click on a person name to add/change/remove their birthdate.
 * Listens for clicks on .person-name-clickable and POSTs to /api/person-birthdate.
 */
import htmx from "htmx.org";

const API_SET = "/api/person-birthdate";
const API_DELETE = "/api/person-birthdate";

function getCsrfOrEmpty(): string {
    const meta = document.querySelector('meta[name="csrf-token"]');
    return meta?.getAttribute("content") ?? "";
}

function getKioskParams(): Record<string, string> {
    const params: Record<string, string> = {};
    document.querySelectorAll(".kiosk-param").forEach((el) => {
        if (el instanceof HTMLInputElement && el.name && el.value) {
            params[el.name] = el.value;
        }
    });
    return params;
}

async function setBirthdate(personId: string, birthDate: string): Promise<{ status: string; message?: string }> {
    const qs = new URLSearchParams(getKioskParams()).toString();
    const url = qs ? `${API_SET}?${qs}` : API_SET;
    const res = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
            ...(getCsrfOrEmpty() ? { "X-CSRF-Token": getCsrfOrEmpty() } : {}),
        },
        body: JSON.stringify({ personId, birthDate }),
    });
    const data = await res.json();
    if (!res.ok) {
        throw new Error(data.message ?? "Request failed");
    }
    return data;
}

async function deleteBirthdate(personId: string): Promise<{ status: string; message?: string }> {
    const qs = new URLSearchParams(getKioskParams()).toString();
    const url = qs ? `${API_DELETE}?${qs}` : API_DELETE;
    const res = await fetch(url, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
            ...(getCsrfOrEmpty() ? { "X-CSRF-Token": getCsrfOrEmpty() } : {}),
        },
        body: JSON.stringify({ personId }),
    });
    const data = await res.json();
    if (!res.ok) {
        throw new Error(data.message ?? "Request failed");
    }
    return data;
}

export function initBirthdateEdit(): void {
    document.body.addEventListener("click", (e: MouseEvent) => {
        const startEl =
            (e.target as Node).nodeType === Node.ELEMENT_NODE
                ? (e.target as HTMLElement)
                : (e.target as Node).parentElement;
        let target = startEl?.closest?.(".person-name-clickable") ?? null;
        if (!target && startEl) {
            const peopleList =
                startEl.closest(".asset--metadata--people-list") ??
                startEl.closest(".asset--metadata--has-icon")?.querySelector(".asset--metadata--people-list");
            if (peopleList) {
                const spans = peopleList.querySelectorAll<HTMLElement>(".person-name-clickable");
                for (const span of spans) {
                    const r = span.getBoundingClientRect();
                    if (
                        e.clientX >= r.left &&
                        e.clientX <= r.right &&
                        e.clientY >= r.top &&
                        e.clientY <= r.bottom
                    ) {
                        target = span;
                        break;
                    }
                }
            }
        }
        if (!target || !(target instanceof HTMLElement)) return;

        e.preventDefault();
        e.stopPropagation();

        const personId = target.getAttribute("data-person-id") ?? "";
        const personName = target.getAttribute("data-person-name") ?? "";
        const currentDob = target.getAttribute("data-birthdate") ?? "";

        if (!personId && !personName) return;

        const id = personId || personName;

        const hasExisting = currentDob.trim() !== "";

        if (hasExisting) {
            const action = window.confirm(
                "Change or remove birthdate?\n\nOK = Edit date\nCancel = Remove birthdate",
            );
            if (action) {
                const newDate = window.prompt("Enter new date of birth (YYYY-MM-DD):", currentDob);
                if (newDate === null) return;
                if (newDate.trim() === "") return;
                setBirthdate(id, newDate.trim())
                    .then(() => {
                        window.alert("Birthdate updated. The page will refresh.");
                        htmx.ajax("get", window.location.pathname + window.location.search, {
                            target: "body",
                            swap: "innerHTML",
                        });
                    })
                    .catch((err) => window.alert(err instanceof Error ? err.message : "Failed to update"));
            } else {
                if (!window.confirm("Remove this person's birthdate?")) return;
                deleteBirthdate(id)
                    .then(() => {
                        window.alert("Birthdate removed. The page will refresh.");
                        htmx.ajax("get", window.location.pathname + window.location.search, {
                            target: "body",
                            swap: "innerHTML",
                        });
                    })
                    .catch((err) => window.alert(err instanceof Error ? err.message : "Failed to remove"));
            }
        } else {
            const dob = window.prompt("Enter date of birth (YYYY-MM-DD):");
            if (dob === null || dob.trim() === "") return;
            setBirthdate(id, dob.trim())
                .then(() => {
                    window.alert("Birthdate saved. The page will refresh.");
                    htmx.ajax("get", window.location.pathname + window.location.search, {
                        target: "body",
                        swap: "innerHTML",
                    });
                })
                .catch((err) => window.alert(err instanceof Error ? err.message : "Failed to save"));
        }
    });
}
