import type { Table, Theme, ThemeName } from "./types.js";
declare const ROW_H = 26;
declare const HEADER_H = 34;
export declare const THEMES: Record<ThemeName, Theme>;
export interface TableDims {
    w: number;
    h: number;
    rowH: number;
    headerH: number;
}
export declare function measureTable(t: Table): TableDims;
export declare function rasterizeTable(t: Table, theme: Theme, dpr: number): HTMLCanvasElement;
export declare function columnY(t: Table, colName?: string): number;
export { ROW_H, HEADER_H };
//# sourceMappingURL=renderer.d.ts.map