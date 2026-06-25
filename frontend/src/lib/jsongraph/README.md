# jsongraph

Zero-dependency, interactive **node/edge graph renderer** for the browser, written in TypeScript.
Feed it any JSON `{ nodes, edges }` and get a pannable / zoomable canvas diagram — each node a card of
its fields, each edge a connector — with hover-to-highlight, click-to-pin focus, drag-to-reposition, a
self-contained auto-layout, and PNG/SVG export. No runtime dependencies. Ships ES modules + `.d.ts`.

> Originally extracted from the canvas engine of
> [royalbhati/sqltoerdiagram](https://github.com/royalbhati/sqltoerdiagram) (MIT); the SQL parser was
> dropped and the dagre layout replaced with a dependency-free layered layout.

## Install & build

This is a source package — drop the `jsongraph/` folder into your project (or a monorepo workspace), then:

```bash
npm install      # installs the TypeScript dev dependency
npm run build    # tsc -> dist/ (ESM .js + .d.ts)
```

Then import it. With a bundler (Vite, esbuild, webpack, Next, etc.) or Node ESM:

```ts
import { createGraph } from "jsongraph"; // resolves to dist/index.js + dist/index.d.ts
```

If you're not publishing to a registry, point your `tsconfig`/bundler at the folder, e.g. a path alias
`"jsongraph": ["./jsongraph/dist/index.js"]`, or import the built file directly
(`import { createGraph } from "./jsongraph/dist/index.js"`).

### No build step? Use the browser bundle

`npm run build` also emits `dist/jsongraph.global.js` — a single classic-script (UMD) file with no
imports, exposing a `window.jsongraph` global. Load it with a plain `<script>` tag (works from `file://`,
CDNs, CodePen, etc.):

```html
<canvas id="g" style="width:100%;height:600px"></canvas>
<script src="jsongraph.global.js"></script>
<script>
  const { createGraph } = window.jsongraph;
  createGraph(document.getElementById("g"), { nodes, edges }, { theme: "dark" });
</script>
```

`examples/demo.html` uses exactly this, so you can just double-click it to open — no server needed.

## Quick start

```ts
import { createGraph, type GraphInput } from "jsongraph";

const data: GraphInput = {
  nodes: [
    { id: "api",   label: "API Gateway",  lang: "Go",   replicas: 3 },
    { id: "users", label: "Users DB",     engine: "PostgreSQL", size_gb: 120 },
  ],
  edges: [
    { source: "api", target: "users", label: "reads/writes" },
  ],
};

const canvas = document.querySelector<HTMLCanvasElement>("#graph")!;
const { diagram, model } = createGraph(canvas, data, { theme: "dark" });
```

The canvas should sit inside a **positioned** container (`position: relative`) — inline editors are
appended to the canvas's parent when `editable` is on. See `examples/demo.html` for a full page.

## `createGraph(canvas, source, opts?)`

One-call setup: maps the input (if it's `{ nodes, edges }`) via `graphJSON`, runs `layout`, constructs a
`Diagram`, applies the theme, wires callbacks, starts the render loop and fits the view. Returns
`{ diagram, model }`.

`source` is a `GraphInput` (`{ nodes, edges }`) **or** a prebuilt `GraphModel`.

```ts
interface CreateGraphOptions {
  theme?: "dark" | "light";        // default "dark"
  dir?: "LR" | "TB";               // layout direction, default "LR"
  spacing?: "compact" | "comfortable" | "spacious";
  fit?: boolean;                   // fit to canvas after build, default true
  editable?: boolean;              // see "Editable" below — default false for JSON
  graph?: GraphJSONOptions;        // forwarded to graphJSON (id/label/source/target/…)
  onEdit?: (e: EditEvent) => void;
  onAddColumn?: (tableKey: string) => void;
  onLayoutChange?: () => void;
  onZoom?: (scale: number) => void;
  onSelectionChange?: () => void;
  onHiddenChange?: () => void;
}
```

## Mapping your JSON — `graphJSON(data, opts?)`

Every node becomes a card (its fields are the rows); every edge becomes a connector. Accepts
`{ nodes, edges }`, d3-style `{ nodes, links }`, `{ vertices, edges }`, or a bare node array.

```ts
const model = graphJSON(data, {
  id: "id",          // node id field or (node) => id        (default "id")
  label: "label",    // header field or (node) => string     (default label/name/title/id)
  fields: ["count"], // rows to show: names, column specs, or (node)=>specs (default: all but id+label)
  source: "source",  // edge "from" field   ("from" also tried)
  target: "target",  // edge "to" field     ("to" also tried)
  sourceField: "sourceField", // optional: anchor edge to a column on the source node
  targetField: "targetField", // optional: anchor edge to a column on the target node
  maxValueLength: 60,
});
```

Notes: every node needs a unique id (matched **case-insensitively**); nodes/edges that are missing ids or
reference unknown nodes are skipped and recorded in `model.errors`. The original object is preserved at
`table.data`. Object/array values render as compact JSON. An edge endpoint may be an id or the node object
itself (d3-style). To fully control a card's rows, give a node a `columns` array of `Column`s.

## The `editable` switch

`diagram.editable` is the master read/write toggle: when `true` the user can rename/edit fields inline,
use "+ add column", **and draw new connections** by dragging a column's connector dots; when `false` the
diagram is structurally read-only (dots don't appear, link-drags can't start) while pan/zoom/drag/hover/
pin still work. `createGraph` defaults it to `false`; pass `{ editable: true }` or set it at runtime
(`diagram.editable = true`). The programmatic link APIs (`addManualLink`, `setManualLinks`, `inferLinks`)
work regardless.

## Re-running layout (e.g. an "Arrange" button)

`setModel` preserves existing node positions (so edits don't jump). To force a fresh arrangement, run
`layout` on the live model in place, then redraw:

```ts
import { layout } from "jsongraph";
layout(diagram.model, { dir: "TB" });
diagram.markDirty();
diagram.fit();
```

## Exports

Vector and raster:

```ts
import { exportSVG } from "jsongraph";
const svg: string | null = exportSVG(diagram.model, "dark");   // standalone SVG string
const png: string | null = diagram.exportPNG(2);               // PNG data URL @2x
```

## API surface

- `createGraph(canvas, source, opts?)` → `{ diagram, model }`
- `graphJSON(data, opts?)` → `GraphModel`
- `layout(model, opts?, hidden?)`, `removeOverlaps(model, hidden?)`
- `exportSVG(model, themeName, annotations?, hidden?)`
- `class Diagram` — the controller (see below)
- `THEMES`, `measureTable`, `rasterizeTable`, `columnY`, `ROW_H`, `HEADER_H`
- `makeAnnotation`, `sanitizeAnnotations`, `NOTE_COLORS`, `GROUP_COLORS`
- **Types**: `GraphInput`, `GraphJSONOptions`, `FieldSpec`, `JSONObject`, `JSONValue`, `Column`, `Table`,
  `Relation`, `GraphModel`, `Theme`, `ThemeName`, `LayoutOptions`, `LayoutDirection`, `LayoutSpacing`,
  `Camera`, `ManualLink`, `EditEvent`, `Annotation`, `TableDims`, `CreateGraphOptions`

### `Diagram` (most-used members)

`editable`, `model`, `setModel(model)`, `setTheme(name)`, `start()`, `resize()`, `markDirty()`,
`fit(padding?)`, `zoomBy(f)`, `resetZoom()`, `setCamera(cam)`, `centerOn(key)`, `exportPNG(scale?)`,
`hideTable(t?)` / `setTableHidden(key,bool)` / `showAllHidden()`, `selectByKey(key,on)` /
`clearSelection()`, `inferLinks()`, `addManualLink(fk,fc,tk,tc)` / `setManualLinks(arr)`,
`setAnnotations(arr)` / `addAnnotation("note"|"group")`, plus the `on*` callbacks.

## The model

`GraphModel` is what every function consumes — build it yourself to drive the renderer from any source:

```ts
interface Column { name: string; type?: string; pk?: boolean; fk?: boolean; nn?: boolean; unique?: boolean }
interface Table  { name: string; key: string; columns: Column[]; x: number; y: number; w: number; h: number; data?: unknown }
interface Relation { fromTable: string; fromCols: string[]; toTable: string; toCols: string[]; label?: string }
interface GraphModel { tables: Table[]; relations: Relation[]; errors: string[] }
```

`key` is the lowercased node id that relations reference; `x/y/w/h` are filled in by `layout` (NaN/0 until
laid out).

## License

MIT. Canvas engine derived from sqltoerdiagram by Royal Bhati.
