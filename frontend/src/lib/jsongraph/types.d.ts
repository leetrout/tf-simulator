/** A single row inside a node card. */
export interface Column {
    name: string;
    /** Display string shown on the right of the row (e.g. a value or data type). */
    type?: string;
    /** Show a "PK" badge. */
    pk?: boolean;
    /** Show an "FK" badge. */
    fk?: boolean;
    /** Not-null marker (carried through; not rendered by default). */
    nn?: boolean;
    /** Unique marker (carried through; not rendered by default). */
    unique?: boolean;
}
/**
 * A node, rendered as a card. `x/y/w/h` are assigned by {@link layout}; they
 * are `NaN`/`0` until a node has been laid out (checked via `Number.isFinite`).
 */
export interface Table {
    /** Header label. */
    name: string;
    /** Unique id, **lowercased**. Relations reference nodes by this key. */
    key: string;
    columns: Column[];
    /** Top-left x (world coords). NaN until laid out. */
    x: number;
    /** Top-left y (world coords). NaN until laid out. */
    y: number;
    /** Width (set by measureTable/layout). */
    w: number;
    /** Height (set by measureTable/layout). */
    h: number;
    rowH?: number;
    headerH?: number;
    /** The original input object this node was built from (adapter passthrough). */
    data?: unknown;
}
/** A connector between two nodes. Endpoints reference {@link Table.key}. */
export interface Relation {
    fromTable: string;
    fromCols: string[];
    toTable: string;
    toCols: string[];
    /** Optional edge label (carried through). */
    label?: string;
    /** True when the target node could not be resolved. */
    toMissing?: boolean;
}
/** The renderer model produced by adapters and consumed by everything else. */
export interface GraphModel {
    tables: Table[];
    relations: Relation[];
    /** Non-fatal problems encountered while building the model. */
    errors: string[];
}
export type ThemeName = "dark" | "light";
export interface Theme {
    bg: string;
    grid: string;
    tableBg: string;
    tableBorder: string;
    header: string;
    headerText: string;
    rowText: string;
    typeText: string;
    rowAlt: string;
    pk: string;
    fk: string;
    edge: string;
    edgeHi: string;
    shadow: string;
    divider: string;
}
export type LayoutDirection = "LR" | "TB";
export type LayoutSpacing = "compact" | "comfortable" | "spacious";
export interface LayoutOptions {
    /** Left-to-right (default) or top-to-bottom. */
    dir?: LayoutDirection;
    spacing?: LayoutSpacing;
}
export interface Camera {
    x: number;
    y: number;
    scale: number;
}
/** A user-drawn / inferred connection. */
export interface ManualLink {
    from: {
        table: string;
        col: string;
    };
    to: {
        table: string;
        col: string;
    };
}
/** Reported by the `onEdit` callback when a label/type is edited inline. */
export interface EditEvent {
    kind: "table" | "column-name" | "column-type";
    tableKey: string;
    colName?: string;
    value: string;
}
/** Sticky note / group-box annotation. */
export interface Annotation {
    id: string;
    type: "group" | "note";
    x: number;
    y: number;
    w: number;
    h: number;
    text: string;
    color: string;
}
//# sourceMappingURL=types.d.ts.map