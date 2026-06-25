import { Diagram } from "./diagram.js";
import type { EditEvent, GraphModel, LayoutOptions, ThemeName } from "./types.js";
import type { GraphInput, GraphJSONOptions } from "./graph-json.js";
export interface CreateGraphOptions extends LayoutOptions {
    /** "dark" (default) or "light". */
    theme?: ThemeName;
    /** Fit the graph to the canvas after building (default `true`). */
    fit?: boolean;
    /**
     * Allow inline editing, "+ add column", and **drawing new connections**.
     * Defaults to `false` for JSON/model input (read-only). Set `true` to enable.
     */
    editable?: boolean;
    /** Options forwarded to {@link graphJSON} when `source` is a `{ nodes, edges }` object. */
    graph?: GraphJSONOptions;
    onEdit?: (e: EditEvent) => void;
    onAddColumn?: (tableKey: string) => void;
    onLayoutChange?: () => void;
    onZoom?: (scale: number) => void;
    onSelectionChange?: () => void;
    onHiddenChange?: () => void;
}
/**
 * One-call setup. `source` is either a `{ nodes, edges }` JSON graph (mapped via
 * {@link graphJSON}) or a prebuilt {@link GraphModel}. Runs layout, constructs a
 * {@link Diagram}, wires callbacks, starts the render loop and fits the view.
 */
export declare function createGraph(canvas: HTMLCanvasElement, source: GraphInput | GraphModel, opts?: CreateGraphOptions): {
    diagram: Diagram;
    model: GraphModel;
};
export { Diagram } from "./diagram.js";
export { graphJSON } from "./graph-json.js";
export { layout, removeOverlaps } from "./layout.js";
export { exportSVG } from "./svg-export.js";
export { THEMES, measureTable, rasterizeTable, columnY, ROW_H, HEADER_H } from "./renderer.js";
export { makeAnnotation, sanitizeAnnotations, NOTE_COLORS, GROUP_COLORS } from "./annotations.js";
export type { GraphInput, GraphJSONOptions, FieldSpec, JSONObject, JSONValue } from "./graph-json.js";
export type { Column, Table, Relation, GraphModel, Theme, ThemeName, LayoutOptions, LayoutDirection, LayoutSpacing, Camera, ManualLink, EditEvent, Annotation, } from "./types.js";
export type { TableDims } from "./renderer.js";
//# sourceMappingURL=index.d.ts.map