import type { Annotation } from "./types.js";
export declare const NOTE_COLORS: Record<string, {
    fill: string;
    text: string;
}>;
export declare const GROUP_COLORS: Record<string, string>;
export declare const NOTE_ORDER: string[];
export declare const GROUP_ORDER: string[];
export declare function makeAnnotation(type: Annotation["type"], x: number, y: number): Annotation;
export declare function sanitizeAnnotations(arr: unknown): Annotation[];
//# sourceMappingURL=annotations.d.ts.map