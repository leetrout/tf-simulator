import type { GraphModel, LayoutOptions } from "./types.js";
type HiddenSet = Set<string> | null;
export declare function layout(model: GraphModel, opts?: LayoutOptions, hidden?: HiddenSet): void;
export declare function removeOverlaps(model: GraphModel, hidden?: HiddenSet): void;
export {};
//# sourceMappingURL=layout.d.ts.map