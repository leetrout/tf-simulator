// Auto-layout — zero-dependency layered (Sugiyama-style) layout.
//
//   - sizes every node from its text metrics;
//   - ranks nodes with a Kahn topological pass (greedy cycle breaking handles
//     loops);
//   - orders nodes within each rank with a barycenter/median heuristic to
//     reduce edge crossings;
//   - assigns coordinates left-to-right (LR) or top-to-bottom (TB);
//   - packs disconnected nodes into a tidy grid off to the side;
//   - runs a separation pass so no two nodes overlap.
import { measureTable } from "./renderer.js";
const PRESETS = {
    comfortable: { nodesep: 36, ranksep: 130, edgesep: 24 },
    compact: { nodesep: 22, ranksep: 80, edgesep: 14 },
    spacious: { nodesep: 60, ranksep: 200, edgesep: 36 },
};
const MARGIN = 40;
export function layout(model, opts = {}, hidden = null) {
    const dir = opts.dir === "TB" ? "TB" : "LR";
    const preset = PRESETS[opts.spacing ?? "comfortable"] || PRESETS["comfortable"];
    const isHidden = (key) => !!(hidden && hidden.has(key));
    // size every node from its content
    for (const t of model.tables) {
        const dims = measureTable(t);
        t.w = dims.w;
        t.h = dims.h;
        t.rowH = dims.rowH;
        t.headerH = dims.headerH;
    }
    const nodes = model.tables.filter((t) => !isHidden(t.key));
    if (!nodes.length)
        return;
    const nodeKeys = new Set(nodes.map((t) => t.key));
    const byKey = new Map(nodes.map((t) => [t.key, t]));
    // de-duplicated directed edges between visible nodes (skip self-loops)
    const seen = new Set();
    const E = [];
    for (const r of model.relations) {
        const f = r.fromTable.toLowerCase();
        const t = r.toTable.toLowerCase();
        if (f === t || !nodeKeys.has(f) || !nodeKeys.has(t))
            continue;
        const k = f + " " + t;
        if (seen.has(k))
            continue;
        seen.add(k);
        E.push([f, t]);
    }
    // adjacency + in-degree
    const succ = new Map(nodes.map((t) => [t.key, []]));
    const pred = new Map(nodes.map((t) => [t.key, []]));
    const indeg = new Map(nodes.map((t) => [t.key, 0]));
    for (const [f, t] of E) {
        succ.get(f).push(t);
        pred.get(t).push(f);
        indeg.set(t, indeg.get(t) + 1);
    }
    // ---- ranking: Kahn topological longest-path, breaking cycles greedily ----
    const rank = new Map(nodes.map((t) => [t.key, 0]));
    const remaining = new Set(nodeKeys);
    const indegWork = new Map(indeg);
    let queue = nodes.filter((t) => indegWork.get(t.key) === 0).map((t) => t.key);
    while (remaining.size) {
        if (!queue.length) {
            // a cycle remains — force the least-constrained node to act as a source
            let best = null, bd = Infinity;
            for (const k of remaining) {
                const d = indegWork.get(k);
                if (d < bd) {
                    bd = d;
                    best = k;
                }
            }
            queue = [best];
            indegWork.set(best, 0);
        }
        const k = queue.shift();
        if (!remaining.has(k))
            continue;
        remaining.delete(k);
        for (const s of succ.get(k)) {
            if (!remaining.has(s))
                continue;
            if (rank.get(k) + 1 > rank.get(s))
                rank.set(s, rank.get(k) + 1);
            indegWork.set(s, indegWork.get(s) - 1);
            if (indegWork.get(s) === 0)
                queue.push(s);
        }
    }
    // ---- group into ranks ----
    let maxRank = 0;
    for (const t of nodes)
        maxRank = Math.max(maxRank, rank.get(t.key));
    const ranks = [];
    for (let r = 0; r <= maxRank; r++)
        ranks[r] = [];
    for (const t of nodes)
        ranks[rank.get(t.key)].push(t.key);
    // ---- crossing reduction: barycenter sweeps ----
    const posOf = () => {
        const pos = new Map();
        for (const arr of ranks)
            arr.forEach((k, i) => pos.set(k, i));
        return pos;
    };
    const bary = (k, list, pos) => {
        const ns = list.get(k);
        if (!ns.length)
            return pos.get(k);
        let s = 0;
        for (const n of ns)
            s += pos.get(n);
        return s / ns.length;
    };
    for (let sweep = 0; sweep < 8; sweep++) {
        const down = sweep % 2 === 0;
        const list = down ? pred : succ;
        const order = down ? [...ranks.keys()].slice(1) : [...ranks.keys()].slice(0, -1).reverse();
        for (const r of order) {
            const pos = posOf();
            ranks[r] = ranks[r]
                .map((k, i) => [k, bary(k, list, pos), i])
                .sort((a, b) => a[1] - b[1] || a[2] - b[2])
                .map((x) => x[0]);
        }
    }
    // ---- coordinate assignment ----
    const along = dir === "LR" ? "h" : "w"; // size used for within-rank stacking
    const cross = dir === "LR" ? "w" : "h"; // size used to advance between ranks
    const nodesep = preset.nodesep;
    const ranksep = preset.ranksep;
    const rankCross = ranks.map((arr) => arr.reduce((m, k) => Math.max(m, byKey.get(k)[cross]), 0));
    const rankStart = [];
    let acc = MARGIN;
    for (let r = 0; r < ranks.length; r++) {
        rankStart[r] = acc;
        acc += rankCross[r] + ranksep;
    }
    for (let r = 0; r < ranks.length; r++) {
        const arr = ranks[r];
        let total = 0;
        for (const k of arr)
            total += byKey.get(k)[along];
        total += Math.max(0, arr.length - 1) * nodesep;
        let cursor = -total / 2; // centre each rank about 0 along the stacking axis
        for (const k of arr) {
            const t = byKey.get(k);
            const crossPos = rankStart[r] + (rankCross[r] - t[cross]) / 2;
            if (dir === "LR") {
                t.x = crossPos;
                t.y = cursor;
            }
            else {
                t.y = crossPos;
                t.x = cursor;
            }
            cursor += t[along] + nodesep;
        }
    }
    placeOrphans(model, isHidden);
    removeOverlaps(model, hidden);
}
// Guarantee no two nodes overlap.
const OVERLAP_GAP = 18;
export function removeOverlaps(model, hidden = null) {
    const isHidden = (key) => !!(hidden && hidden.has(key));
    const ts = model.tables.filter((t) => Number.isFinite(t.x) && !isHidden(t.key));
    const n = ts.length;
    if (n < 2)
        return;
    const MAX_PASSES = 60;
    for (let pass = 0; pass < MAX_PASSES; pass++) {
        let moved = false;
        for (let i = 0; i < n; i++) {
            for (let j = i + 1; j < n; j++) {
                const a = ts[i], b = ts[j];
                const ox = Math.min(a.x + a.w, b.x + b.w) - Math.max(a.x, b.x) + OVERLAP_GAP;
                const oy = Math.min(a.y + a.h, b.y + b.h) - Math.max(a.y, b.y) + OVERLAP_GAP;
                if (ox <= 0 || oy <= 0)
                    continue; // not overlapping (with gap)
                if (ox < oy) {
                    const shift = ox / 2;
                    if (a.x < b.x) {
                        a.x -= shift;
                        b.x += shift;
                    }
                    else {
                        a.x += shift;
                        b.x -= shift;
                    }
                }
                else {
                    const shift = oy / 2;
                    if (a.y < b.y) {
                        a.y -= shift;
                        b.y += shift;
                    }
                    else {
                        a.y += shift;
                        b.y -= shift;
                    }
                }
                moved = true;
            }
        }
        if (!moved)
            break;
    }
}
// Nodes with no relations get packed into a compact grid beside the graph.
function placeOrphans(model, isHidden) {
    const connected = new Set();
    for (const r of model.relations) {
        connected.add(r.fromTable.toLowerCase());
        connected.add(r.toTable.toLowerCase());
    }
    const orphans = model.tables.filter((t) => !isHidden(t.key) && (!connected.has(t.key) || !Number.isFinite(t.x)));
    if (!orphans.length)
        return;
    const placed = model.tables.filter((t) => Number.isFinite(t.x) && connected.has(t.key) && !isHidden(t.key));
    let maxX = 0, minY = Infinity, maxY = -Infinity;
    for (const t of placed) {
        maxX = Math.max(maxX, t.x + t.w);
        minY = Math.min(minY, t.y);
        maxY = Math.max(maxY, t.y + t.h);
    }
    if (!Number.isFinite(minY)) {
        minY = 40;
        maxY = 40;
    }
    const startX = placed.length ? maxX + 100 : 40;
    const colW = Math.max(...orphans.map((t) => t.w), 160) + 36;
    const availH = Math.max(maxY - minY, 400);
    let x = startX, y = minY;
    for (const t of orphans) {
        if (y > minY && y + t.h > minY + availH) {
            y = minY;
            x += colW;
        }
        t.x = x;
        t.y = y;
        y += t.h + 36;
    }
}
//# sourceMappingURL=layout.js.map