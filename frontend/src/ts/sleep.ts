import fullyKiosk from "./fullykiosk";
import immichFrame from "./immichframe";
import { triggerNewAsset } from "./polling";

let isSleeping = false;

/**
 * Enables or disables UI sleep mode and optionally controls the device screensaver state. Triggers a new asset when waking up from sleep.
 *
 * When sleep mode is enabled, the "sleep" CSS class is added to the document body; when disabled, it is removed.
 * If {@link screensaver} is true, the device screensaver state is set to match {@link turnOn} using the Fully Kiosk API.
 * If {@link runningInImmichFrame} is true, the Immich Frame screensaver state is also updated.
 *
 * @param turnOn - Whether to enable (true) or disable (false) sleep mode.
 * @param screensaver - Whether to also control the device screensaver state.
 * @param runningInImmichFrame - Whether the app is running inside Immich Frame and should update its screensaver state.
 */
function sleepMode(
    turnOn: boolean,
    screensaver: boolean,
    runningInImmichFrame: boolean,
): void {
    if (turnOn) {
        document.body.classList.add("sleep");
        isSleeping = true;
    } else {
        document.body.classList.remove("sleep");
        if (isSleeping) {
            triggerNewAsset();
        }
        isSleeping = false;
    }

    if (screensaver) {
        fullyKiosk.setScreensaverState(turnOn);
        if (runningInImmichFrame) immichFrame.setScreensaverState(turnOn);
    }
}

export { sleepMode };
