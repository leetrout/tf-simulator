import type { Annotation, Camera, EditEvent, GraphModel, ManualLink, Table, Theme, ThemeName } from "./types.js";
interface Rect {
    x: number;
    y: number;
    w: number;
    h: number;
}
interface EditTarget {
    table: Table;
    kind: EditEvent["kind"];
    colName?: string;
    value?: string;
    rect: Rect;
    align: CanvasTextAlign;
    weight: number;
}
interface EditingState {
    input: HTMLInputElement | HTMLTextAreaElement;
    target?: EditTarget;
    datalist?: HTMLDataListElement | null;
    anno?: Annotation;
}
type Drag = {
    t: Table;
    dx: number;
    dy: number;
    moved: boolean;
};
type DragGroup = {
    items: Array<{
        t: Table;
        sx0: number;
        sy0: number;
    }>;
    ax: number;
    ay: number;
    moved: boolean;
};
type Marquee = {
    ax: number;
    ay: number;
    x: number;
    y: number;
};
type Pan = {
    sx: number;
    sy: number;
    camx: number;
    camy: number;
    moved: boolean;
};
type Linking = {
    fromKey: string;
    fromCol: string;
    side: "left" | "right";
    wx: number;
    wy: number;
    cx: number;
    cy: number;
};
type HoverConn = {
    t: Table;
    colIndex: number;
};
type AnnoDrag = {
    a: Annotation;
    dx: number;
    dy: number;
    moved: boolean;
};
type AnnoResize = {
    a: Annotation;
    ox: number;
    oy: number;
    moved: boolean;
};
export declare class Diagram {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;
    dpr: number;
    cam: Camera;
    model: GraphModel;
    themeName: ThemeName;
    theme: Theme;
    bitmaps: Map<string, HTMLCanvasElement>;
    dirty: boolean;
    frameQueued: boolean;
    viewW: number;
    viewH: number;
    drag: Drag | null;
    dragGroup: DragGroup | null;
    marquee: Marquee | null;
    selected: Set<Table>;
    hidden: Set<string>;
    manualLinks: ManualLink[];
    hoverConn: HoverConn | null;
    linking: Linking | null;
    pan: Pan | null;
    hover: Table | null;
    pinned: Table | null;
    pinnedKeys: Set<string> | null;
    editing: EditingState | null;
    editable: boolean;
    typeSuggestions: string[];
    annotations: Annotation[];
    selectedAnno: Annotation | null;
    annoDrag: AnnoDrag | null;
    annoResize: AnnoResize | null;
    onHiddenChange: (() => void) | null;
    onSelectionChange: (() => void) | null;
    onZoom: ((scale: number) => void) | null;
    onLayoutChange: (() => void) | null;
    onEdit: ((e: EditEvent) => void) | null;
    onAddColumn: ((tableKey: string) => void) | null;
    private _tmap;
    private _tmapDirty;
    constructor(canvas: HTMLCanvasElement);
    setModel(model: GraphModel, { keepCamera }?: {
        keepCamera?: boolean;
    }): void;
    get focus(): Table | null;
    private _pin;
    private _setPin;
    pinByKey(key: string): void;
    editColumn(tableKey: string, colName: string): void;
    setTheme(name: ThemeName): void;
    resize(): void;
    start(): void;
    markDirty(): void;
    private _loop;
    screenToWorld(sx: number, sy: number): {
        x: number;
        y: number;
    };
    private _bitmap;
    private _render;
    private _connDots;
    private _annoVisible;
    private _drawGroup;
    private _drawNote;
    private _drawAnnoChrome;
    private _annoChromeRects;
    setAnnotations(arr: Annotation[]): void;
    addAnnotation(type: Annotation["type"]): Annotation;
    deleteSelectedAnnotation(): void;
    private _annoAt;
    private _noteAt;
    private _groupAt;
    private _grabAnno;
    private _annoChromeAt;
    private _beginEditAnnotation;
    private _addRect;
    private _drawAddButton;
    private _addButtonAt;
    private _drawGrid;
    private _edgeSeg;
    private _drawEdges;
    private _stroke;
    private _tableMap;
    tableAt(sx: number, sy: number): Table | null;
    private _bindInput;
    private _pointerDown;
    private _columnAtWorld;
    private _pointerMove;
    private _pointerUp;
    clearSelection(): void;
    hideTable(t?: Table | null): void;
    showAllHidden(): void;
    setHidden(keys: string[]): void;
    hiddenCount(): number;
    setTablesHidden(keys: string[], hidden: boolean): void;
    setTableHidden(key: string, hidden: boolean): void;
    selectByKey(key: string, on: boolean): void;
    isSelected(key: string): boolean;
    centerOn(key: string): void;
    private _connectorAt;
    private _columnAt;
    private _linkExists;
    addManualLink(fk: string, fc: string, tk: string, tc: string): boolean;
    setManualLinks(arr: ManualLink[]): void;
    linkAt(sx: number, sy: number): ManualLink | null;
    removeManualLink(link: ManualLink): void;
    clearManualLinks(): void;
    manualLinkCount(): number;
    inferLinks(): number;
    private _editAt;
    private _editTargetAt;
    private _beginEdit;
    private _commitEdit;
    private _cancelEdit;
    private _zoomAt;
    setCamera(cam: Camera | null | undefined): void;
    zoomBy(factor: number): void;
    resetZoom(): void;
    fit(padding?: number): void;
    bounds(padding?: number): {
        x0: number;
        y0: number;
        x1: number;
        y1: number;
        w: number;
        h: number;
    };
    exportPNG(scale?: number): string | null;
}
export {};
//# sourceMappingURL=diagram.d.ts.map