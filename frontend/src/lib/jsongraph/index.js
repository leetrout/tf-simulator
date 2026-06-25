// jsongraph — zero-dependency interactive node/edge graph renderer.
//
// Graph any JSON `{ nodes, edges }` on a <canvas>: pan / zoom / drag,
// hover-to-highlight relationships, click-to-pin focus, PNG/SVG export.
//
//   import { createGraph } from "jsongraph";
//   const { diagram } = createGraph(canvas, { nodes, edges }, { theme: "dark" });
import { Diagram } from "./diagram.js";
import { graphJSON } from "./graph-json.js";
import { layout } from "./layout.js";
/**
 * One-call setup. `source` is either a `{ nodes, edges }` JSON graph (mapped via
 * {@link graphJSON}) or a prebuilt {@link GraphModel}. Runs layout, constructs a
 * {@link Diagram}, wires callbacks, starts the render loop and fits the view.
 */
export function createGraph(canvas, source, opts = {}) {
    const model = source && "nodes" in source && !("tables" in source)
        ? graphJSON(source, opts.graph || {})
        : source;
    layout(model, { dir: opts.dir, spacing: opts.spacing });
    const diagram = new Diagram(canvas);
    if (opts.theme)
        diagram.setTheme(opts.theme);
    // read-only by default for JSON/model input; flip via opts.editable
    diagram.editable = opts.editable ?? false;
    if (opts.onEdit)
        diagram.onEdit = opts.onEdit;
    if (opts.onAddColumn)
        diagram.onAddColumn = opts.onAddColumn;
    if (opts.onLayoutChange)
        diagram.onLayoutChange = opts.onLayoutChange;
    if (opts.onZoom)
        diagram.onZoom = opts.onZoom;
    if (opts.onSelectionChange)
        diagram.onSelectionChange = opts.onSelectionChange;
    if (opts.onHiddenChange)
        diagram.onHiddenChange = opts.onHiddenChange;
    diagram.setModel(model);
    diagram.start();
    if (opts.fit !== false)
        diagram.fit();
    return { diagram, model };
}
// ---- public API ----
export { Diagram } from "./diagram.js";
export { graphJSON } from "./graph-json.js";
export { layout, removeOverlaps } from "./layout.js";
export { exportSVG } from "./svg-export.js";
export { THEMES, measureTable, rasterizeTable, columnY, ROW_H, HEADER_H } from "./renderer.js";
export { makeAnnotation, sanitizeAnnotations, NOTE_COLORS, GROUP_COLORS } from "./annotations.js";
//# sourceMappingURL=index.js.map