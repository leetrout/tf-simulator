# Local patches to vendored jsongraph

This is the built `dist/` of jsongraph (from `scratch/jsongraph`), vendored so the
app can bundle it without an npm dependency.

## renderer.js — per-node status color

`rasterizeTable` was extended to read `t.data.statusColor` (the adapter passes the
original node object through as `table.data`). When present it draws a left status
stripe, tints the header, and colors the border with that hex. This is what lets
the resource graph color nodes by `in_sync`/`drift`/`phantom`/`untracked`/`moved`
— jsongraph's stock renderer only supports a single global theme color.

Search `renderer.js` for `statesim patch` to find the three edits. Re-apply them if
jsongraph is re-vendored from an upstream build.
