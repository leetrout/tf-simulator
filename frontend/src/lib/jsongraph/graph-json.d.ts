import type { Column, GraphModel } from "./types.js";
/** Any object. Nodes/edges may carry arbitrary fields. */
export type JSONValue = unknown;
export type JSONObject = Record<string, JSONValue>;
/** Accepted input shapes. Aliases (links/vertices/relationships) are supported. */
export interface GraphInput {
    nodes?: JSONObject[];
    edges?: JSONObject[];
    links?: JSONObject[];
    vertices?: JSONObject[];
    relationships?: JSONObject[];
}
/** A column spec a caller can return from `fields`. */
export type FieldSpec = string | (Partial<Column> & {
    name: string;
});
export interface GraphJSONOptions {
    /** Node id field name, or a function `(node) => id`. Default `"id"`. */
    id?: string | ((node: JSONObject) => unknown);
    /** Card header field name, or a function. Default: label/name/title/id. */
    label?: string | ((node: JSONObject) => unknown);
    /**
     * Which rows to show. Omit for "all fields except id and the label field".
     * Pass an array of field names / column specs, or a function returning them.
     */
    fields?: FieldSpec[] | ((node: JSONObject) => FieldSpec[]);
    /** Edge field naming the source endpoint. Default `"source"` (`"from"` also tried). */
    source?: string;
    /** Edge field naming the target endpoint. Default `"target"` (`"to"` also tried). */
    target?: string;
    /** Edge field naming a column on the source node to anchor to. Default `"sourceField"`. */
    sourceField?: string;
    /** Edge field naming a column on the target node to anchor to. Default `"targetField"`. */
    targetField?: string;
    /** Truncate field values longer than this. Default `60`. */
    maxValueLength?: number;
}
export declare function graphJSON(data: GraphInput | JSONObject[] | null | undefined, opts?: GraphJSONOptions): GraphModel;
//# sourceMappingURL=graph-json.d.ts.map