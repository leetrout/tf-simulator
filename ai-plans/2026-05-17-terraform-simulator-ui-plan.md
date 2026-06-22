# Terraform Simulator UI and Workflow Plan

Date: 2026-05-17

## Goal

Build a laptop-friendly Terraform behavior simulator that runs as a single binary, starts a local webserver, opens the browser, and presents a modern two-pane UI. The UI should show a graph editor on one side and a JSON editor on the other, with tabs for both Terraform state and fake API resources.

## Current Baseline

The project already has several foundations in place:

- A Go binary entrypoint that initializes storage and starts the webserver.
- An embedded static frontend served by the Go webserver.
- Existing fake API resources for entities and groups.
- Local persistence through `tfsim.json.gz`.
- A Svelte/Tailwind frontend with basic entity/group fetching.
- A websocket endpoint that can notify frontend clients after changes.

## Implementation Plan

### 1. Separate Terraform State from Fake API Resources

Separate the simulator's model into two clear concepts:

1. Fake provider API resources
   - Represent remote service objects Terraform manages.
   - Existing `Entity` and `Group` resources are a good starting point.
   - This side should support intentional drift by allowing users to edit or delete remote resources outside Terraform.

2. Terraform state model
   - Add a Terraform-state-shaped model with fields such as:
     - `version`
     - `terraform_version`
     - `serial`
     - `lineage`
     - `outputs`
     - `resources`
     - `instances`
     - `attributes`
     - `dependencies`
   - Store it separately from fake API resources.
   - Keep the existing persistence mechanism initially, but evolve the persisted document into a richer simulator state.

Recommended initial shape:

```go
type SimulatorState struct {
    APIResources   APIResources
    TerraformState TerraformState
    Events         []Event
}
```

### 2. Add Simulator API Endpoints

Keep the existing provider-like entity and group endpoints, but add simulator-focused endpoints for the new UI.

Read endpoints:

- `GET /api/sim/state` returns the complete simulator state.
- `GET /api/sim/terraform-state` returns a Terraform-state-shaped JSON document.
- `GET /api/sim/api-resources` returns fake remote API resources.
- `GET /api/sim/graph?view=terraform|api` returns graph-ready nodes and edges.

Write/action endpoints:

- `POST /api/sim/reset` resets to empty or seeded examples.
- `POST /api/sim/apply` simulates Terraform apply.
- `POST /api/sim/refresh` simulates Terraform refresh.
- `POST /api/sim/plan` returns proposed changes without applying them.
- `PATCH /api/sim/terraform-state` accepts validated JSON editor changes.
- `PATCH /api/sim/api-resources` accepts validated fake API JSON changes.
- `POST /api/sim/drift` applies common drift scenarios.

### 3. Make the Binary Open the Web UI

The entrypoint already computes host and port and has a commented browser-opening hook. Restore this behavior in a robust way.

Add or confirm these CLI flags:

- `--addr :9321`
- `--no-open`
- `--state-file ./tfsim.json.gz`
- `--seed demo|empty`

Implementation notes:

- Open `http://localhost:<port>` after the server is actually listening.
- Keep interactive mode available for debugging.
- Use a small cross-platform browser-opening helper or dependency.

### 4. Restructure the Frontend Before the UI Rewrite

Before adding graph/editor complexity, add frontend structure:

- `src/lib/api/client.ts` for typed API calls.
- `src/lib/stores/simulator.ts` for Svelte stores.
- `src/lib/types/simulator.ts` for shared UI types.
- `src/lib/components/layout/` for shell, tabs, toolbar, and status bar.
- `src/lib/components/graph/` for graph canvas and inspector.
- `src/lib/components/json/` for JSON editor wrapper and validation messages.

Potential dependencies:

- Graph editing/viewing: Svelte Flow / XYFlow or Cytoscape.js.
- JSON editing: Monaco Editor or a lighter JSON editor.
- JSON schema validation: `ajv`.
- Styling: continue using Tailwind.

### 5. Build the Two-Pane UI Shell

The requested UI should have:

- A left graph pane.
- A right JSON editor pane.
- Top-level tabs for:
  - Terraform State
  - Fake API Resources
- A toolbar for simulator actions.
- A status bar for websocket connection, dirty state, last event, and state file path.

Suggested layout:

```text
┌────────────────────────────────────────────────────────────┐
│ TF Simulator        [Terraform State] [Fake API Resources] │
├────────────────────────────────────────────────────────────┤
│ Toolbar: Seed | Reset | Plan | Apply | Refresh | Drift     │
├─────────────────────────────┬──────────────────────────────┤
│ Graph pane                  │ JSON pane                    │
│                             │                              │
│ Nodes/edges                 │ Monaco JSON editor           │
│ Selection interactions      │ Validation + apply button    │
│                             │                              │
├─────────────────────────────┴──────────────────────────────┤
│ Status: websocket, dirty state, last event, state file path │
└────────────────────────────────────────────────────────────┘
```

Interaction goals:

- Selecting a graph node highlights the related JSON fragment.
- Selecting a JSON resource highlights the graph node.
- Dragging nodes updates visual layout metadata initially.
- Later, adding or removing graph edges can mutate dependencies or group membership.

### 6. Generate Graph Data on the Backend

Do not make the frontend infer relationships from raw Terraform JSON. The backend should expose graph-ready data.

Terraform state graph nodes:

- Resource blocks, such as `tfsim_entity.example`.
- Resource instances, such as `tfsim_entity.example[0]`.
- Data sources if added later.
- Outputs if useful.

Terraform state graph edges:

- Terraform dependencies.
- Group-to-entity relationships.
- Synthetic state relationships.

Fake API graph nodes:

- Entities.
- Groups.
- Optional API root or namespace node.

Fake API graph edges:

- `Group -> Entity` edges based on `EntitySet`.

Recommended DTO:

```go
type GraphView struct {
    Nodes []GraphNode `json:"nodes"`
    Edges []GraphEdge `json:"edges"`
}
```

### 7. Add JSON Editing with Validation

The JSON pane should be useful but safe.

Capabilities:

- Pretty-print JSON.
- Validate syntax client-side before submission.
- Validate schema and invariants server-side.
- Show errors inline.
- Format, revert, apply, download, and copy JSON.

Behavior:

- Editing fake API JSON creates drift.
- Editing Terraform state JSON simulates direct state surgery.
- Every mutation records an event and notifies websocket clients.

### 8. Add Simulator Workflows

Implement common Terraform workflows so users can explore behavior.

MVP workflows:

1. Create desired config from the graph.
2. Apply changes to fake API resources and Terraform state.
3. Edit fake API resources manually to create drift.
4. Refresh Terraform state from fake API resources.
5. Plan changes without applying them.
6. Edit Terraform state JSON directly and visualize consequences.

Later workflows:

- Import simulation.
- Taint and replace simulation.
- Moved blocks and state moves.
- Provider read failures.
- Computed and unknown values.

### 9. Improve the Websocket Protocol

Replace string-based websocket messages with structured JSON events.

Example event:

```json
{
  "type": "resource.updated",
  "scope": "api",
  "resourceType": "entity",
  "id": "...",
  "serial": 12,
  "timestamp": "..."
}
```

Frontend behavior:

- Receive structured event.
- Refetch the affected view or full simulator snapshot.
- Show the latest event in the status bar.
- Derive websocket URL from `window.location` instead of hardcoding localhost.
- Reconnect gracefully if the socket closes.

### 10. Add Example Scenarios

Add seed scenarios that make the simulator easy to understand.

Suggested examples:

- Empty workspace.
- Simple entity.
- Group with two entities.
- Drifted entity attribute.
- Missing remote object.
- State has an object the API does not.
- API has unmanaged object.
- Dependency chain.

Suggested implementation:

- `internal/examples` package.
- `GET /api/sim/examples` endpoint.
- `POST /api/sim/seed/{example}` endpoint.

### 11. Testing Strategy

Backend tests:

- Store load/save round trips.
- Terraform state JSON generation.
- Graph generation from Terraform state.
- Graph generation from fake API resources.
- Plan/apply/refresh behavior.
- API handler validation.

Frontend checks/tests:

- `npm run check`.
- Component tests for tab switching.
- Component tests for JSON validation state.
- Component tests for graph selection and selected JSON path.
- Optional Playwright smoke test once the UI stabilizes.

Build checks:

- `go test ./...`.
- `npm run check` from `frontend`.
- `npm run build` from `frontend`.
- `script/build.sh` if that is the intended release path.

## Suggested Milestones

### Milestone 1: Backend Simulator Model

- Extend store state to hold fake API resources and Terraform-state-like data.
- Add graph DTOs.
- Add `/api/sim/*` read endpoints.
- Add seed/reset support.
- Keep existing entity/group endpoints working.

### Milestone 2: Single Binary Launch Polish

- Restore browser opening.
- Add `--no-open`.
- Make websocket URL dynamic in the frontend.
- Confirm embedded static build output still works.

### Milestone 3: Frontend App Shell

- Replace the simple entity/group page with the two-pane shell.
- Add top tabs for Terraform State and Fake API Resources.
- Add JSON viewer/editor in the right pane.
- Show a graph placeholder in the left pane.

### Milestone 4: Real Graph Viewer

- Implement backend graph generation.
- Add graph visualization. We already have a library for this in typescript. see scratch/graph.
- Add node selection and JSON synchronization.
- Add visual states for managed, unmanaged, drifted, missing, and pending changes.

### Milestone 5: Editing and Workflows

- Implement JSON apply/revert.
- Implement graph edits.
- Implement plan/apply/refresh/drift actions.
- Add event log and structured websocket messages.

### Milestone 6: Polish and Documentation

- Add example scenarios.
- Add README quickstart.
- Add screenshots or GIFs.
- Add tests and release build instructions.

