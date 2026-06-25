export const NOTE_COLORS = {
    yellow: { fill: "#ffe49c", text: "#3a3320" },
    blue: { fill: "#bcd9ff", text: "#1f2f45" },
    green: { fill: "#c4f0d0", text: "#1d3a28" },
    pink: { fill: "#ffc9d8", text: "#46202e" },
};
export const GROUP_COLORS = {
    blue: "#5aa7ff",
    green: "#4ec9b0",
    amber: "#f5b14a",
    purple: "#b488ff",
};
export const NOTE_ORDER = ["yellow", "blue", "green", "pink"];
export const GROUP_ORDER = ["blue", "green", "amber", "purple"];
let seq = 0;
function newId() {
    seq += 1;
    return "a" + Date.now().toString(36) + "-" + seq.toString(36);
}
export function makeAnnotation(type, x, y) {
    if (type === "group") {
        return { id: newId(), type: "group", x: x - 160, y: y - 120, w: 320, h: 240, text: "Group", color: "blue" };
    }
    return { id: newId(), type: "note", x: x - 90, y: y - 60, w: 180, h: 120, text: "", color: "yellow" };
}
// tolerant normaliser for loaded/shared annotations
export function sanitizeAnnotations(arr) {
    if (!Array.isArray(arr))
        return [];
    const out = [];
    for (const a of arr) {
        if (!a || (a.type !== "group" && a.type !== "note"))
            continue;
        if (![a.x, a.y, a.w, a.h].every((n) => Number.isFinite(n)))
            continue;
        const colors = a.type === "group" ? GROUP_COLORS : NOTE_COLORS;
        out.push({
            id: typeof a.id === "string" ? a.id : newId(),
            type: a.type,
            x: a.x,
            y: a.y,
            w: Math.max(80, a.w),
            h: Math.max(50, a.h),
            text: typeof a.text === "string" ? a.text : "",
            color: (a.color && colors[a.color]) ? a.color : a.type === "group" ? "blue" : "yellow",
        });
    }
    return out;
}
//# sourceMappingURL=annotations.js.map