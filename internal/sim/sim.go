// Package sim is the observer: it joins the in-process Nimbus cloud store to the
// on-disk terraform.tfstate, classifies every resource (in_sync/drift/phantom/
// untracked/moved), builds graphs, serves /api/sim/*, injects scenarios, and
// pushes structured websocket events. It never executes terraform.
package sim

import (
	"sort"
	"sync"

	"github.com/leetrout/terraform-sim/internal/cloud"
	"github.com/leetrout/terraform-sim/internal/ws"
)

// Status taxonomy — the diff vocabulary. Values match the frontend legend.
const (
	StatusInSync    = "in_sync"
	StatusDrift     = "drift"
	StatusPhantom   = "phantom"
	StatusUntracked = "untracked"
	StatusMoved     = "moved"
)

// CloudResource is the uniform JSON shape for a Nimbus resource in API output.
type CloudResource struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Serial     int               `json:"serial"`
	Attributes map[string]string `json:"attributes"`
}

// AttrDiff is a single attribute that differs between state and cloud.
type AttrDiff struct {
	Field string `json:"field"`
	State string `json:"state"`
	Cloud string `json:"cloud"`
}

// ResourceStatus is the per-resource classification.
type ResourceStatus struct {
	ID      string     `json:"id"`
	Type    string     `json:"type"`
	Name    string     `json:"name"`
	Address string     `json:"address,omitempty"`
	Status  string     `json:"status"`
	Detail  string     `json:"detail"`
	Diffs   []AttrDiff `json:"diffs,omitempty"`
}

// Counts aggregates statuses for the status bar and tab badges.
type Counts struct {
	Total    int            `json:"total"`
	Tracked  int            `json:"tracked"`
	Problems int            `json:"problems"`
	ByStatus map[string]int `json:"byStatus"`
}

// Meta carries environment + simulator context.
type Meta struct {
	WorkDir          string `json:"workDir"`
	StateFile        string `json:"stateFile"`
	StateExists      bool   `json:"stateExists"`
	TerraformVersion string `json:"terraformVersion"`
	Provider         string `json:"provider"`
	Region           string `json:"region"`
	Backend          string `json:"backend"`
	LoadedScenario   string `json:"loadedScenario,omitempty"`
	LastEvent        string `json:"lastEvent,omitempty"`
}

// Snapshot is the full /api/sim/snapshot payload.
type Snapshot struct {
	Cloud  []CloudResource  `json:"cloud"`
	State  []StateResource  `json:"state"`
	Status []ResourceStatus `json:"status"`
	Counts Counts           `json:"counts"`
	Meta   Meta             `json:"meta"`
}

// entry is the internal joined view of one resource (by cloud id), reused by the
// status and graph builders.
type entry struct {
	ID      string
	Type    string
	Name    string
	Address string
	Status  string
	InCloud bool
	InState bool
	Attrs   map[string]string // effective attributes (cloud wins; else state)
	Diffs   []AttrDiff
}

// Simulator ties the cloud store to the on-disk state and serves the sim API.
type Simulator struct {
	store     *cloud.Store
	workDir   string
	stateFile string
	region    string

	mu             sync.RWMutex
	loadedScenario string
	lastEvent      string

	watcher *watcher
}

// New constructs a Simulator over the given store and working directory.
func New(store *cloud.Store, workDir string) *Simulator {
	return &Simulator{
		store:     store,
		workDir:   workDir,
		stateFile: stateFilePath(workDir),
		region:    "us-west-1",
	}
}

// Start subscribes to cloud mutations and begins watching the state file. Both
// sources publish structured websocket events.
func (s *Simulator) Start() {
	s.store.Subscribe(func(ev cloud.ChangeEvent) {
		s.setLastEvent("cloud.changed: " + ev.Op + " " + ev.Type)
		ws.Broadcast(ws.Event{
			Type:         "cloud.changed",
			Scope:        "cloud",
			ResourceType: ev.Type,
			ID:           ev.ID,
			Serial:       ev.Serial,
			Detail:       ev.Op,
		})
	})
	s.watcher = newWatcher(s.stateFile, func() {
		s.setLastEvent("state.changed")
		ws.Broadcast(ws.Event{Type: "state.changed", Scope: "state"})
	})
	s.watcher.start()
}

// Stop halts the state watcher.
func (s *Simulator) Stop() {
	if s.watcher != nil {
		s.watcher.stop()
	}
}

func (s *Simulator) setLastEvent(msg string) {
	s.mu.Lock()
	s.lastEvent = msg
	s.mu.Unlock()
}

// compute performs the cloud↔state join and returns everything the API needs.
func (s *Simulator) compute() (cloudOut []CloudResource, stateOut []StateResource, statuses []ResourceStatus, entries []entry, meta Meta) {
	cloudOut = []CloudResource{}
	stateOut = []StateResource{}
	statuses = []ResourceStatus{}
	cloudRes := s.store.All()
	state, _ := LoadState(s.stateFile)

	stateExists := state != nil
	tfVersion := ""
	stateByID := map[string]StateResource{}
	if state != nil {
		tfVersion = state.TerraformVersion
		stateOut = append(stateOut, state.Resources...)
		for _, sr := range state.Resources {
			if sr.ID != "" {
				if _, dup := stateByID[sr.ID]; !dup {
					stateByID[sr.ID] = sr
				}
			}
		}
	}

	consumed := map[string]bool{}

	for _, r := range cloudRes {
		k := cloud.KindOf(r)
		attrs := cloud.Attributes(r)
		cloudOut = append(cloudOut, CloudResource{
			Type:       k.Type,
			ID:         r.GetID(),
			Serial:     r.GetSerial(),
			Attributes: attrs,
		})
		e := entry{
			ID:      r.GetID(),
			Type:    k.Type,
			Name:    attrs["name"],
			InCloud: true,
			Attrs:   attrs,
		}
		if sr, ok := stateByID[r.GetID()]; ok {
			consumed[r.GetID()] = true
			e.InState = true
			e.Address = sr.Address
			e.Diffs = diffAttrs(sr.Attributes, attrs)
			e.Status = classifyMatched(e.Diffs)
		} else {
			e.Status = StatusUntracked
		}
		entries = append(entries, e)
	}

	// State resources with no matching cloud resource are phantoms.
	for _, sr := range stateOut {
		if sr.ID != "" && consumed[sr.ID] {
			continue
		}
		entries = append(entries, entry{
			ID:      sr.ID,
			Type:    sr.Type,
			Name:    sr.Name,
			Address: sr.Address,
			Status:  StatusPhantom,
			InState: true,
			Attrs:   sr.Attributes,
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type < entries[j].Type
		}
		return entries[i].ID < entries[j].ID
	})

	for _, e := range entries {
		statuses = append(statuses, ResourceStatus{
			ID:      e.ID,
			Type:    e.Type,
			Name:    e.Name,
			Address: e.Address,
			Status:  e.Status,
			Detail:  detailFor(e),
			Diffs:   e.Diffs,
		})
	}

	s.mu.RLock()
	meta = Meta{
		WorkDir:          s.workDir,
		StateFile:        s.stateFile,
		StateExists:      stateExists,
		TerraformVersion: tfVersion,
		Provider:         "Nimbus Cloud",
		Region:           s.region,
		Backend:          "local",
		LoadedScenario:   s.loadedScenario,
		LastEvent:        s.lastEvent,
	}
	s.mu.RUnlock()
	return cloudOut, stateOut, statuses, entries, meta
}

// Snapshot returns the full joined view.
func (s *Simulator) Snapshot() Snapshot {
	cloudOut, stateOut, statuses, _, meta := s.compute()
	return Snapshot{
		Cloud:  cloudOut,
		State:  stateOut,
		Status: statuses,
		Counts: countStatuses(statuses),
		Meta:   meta,
	}
}

// Status returns just the per-resource statuses and counts.
func (s *Simulator) Status() ([]ResourceStatus, Counts) {
	_, _, statuses, _, _ := s.compute()
	return statuses, countStatuses(statuses)
}

// diffAttrs compares cloud attributes against the state's attributes over the
// cloud attribute keys. A key absent from state is not treated as a difference
// (avoids false positives from computed/unknown values).
func diffAttrs(stateAttrs, cloudAttrs map[string]string) []AttrDiff {
	var diffs []AttrDiff
	keys := make([]string, 0, len(cloudAttrs))
	for k := range cloudAttrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sv, ok := stateAttrs[k]
		if !ok {
			continue
		}
		if sv != cloudAttrs[k] {
			diffs = append(diffs, AttrDiff{Field: k, State: sv, Cloud: cloudAttrs[k]})
		}
	}
	return diffs
}

// classifyMatched decides between in_sync, moved, and drift for a resource present
// in both state and cloud. A rename (only the name attribute differs) is "moved".
func classifyMatched(diffs []AttrDiff) string {
	if len(diffs) == 0 {
		return StatusInSync
	}
	if len(diffs) == 1 && diffs[0].Field == "name" {
		return StatusMoved
	}
	return StatusDrift
}

func detailFor(e entry) string {
	switch e.Status {
	case StatusInSync:
		return "Tracked and current; terraform plan reports no changes."
	case StatusDrift:
		fields := ""
		for i, d := range e.Diffs {
			if i > 0 {
				fields += ", "
			}
			fields += d.Field
		}
		return "Cloud differs from state (" + fields + "); apply would overwrite it back."
	case StatusPhantom:
		return "In state but absent from the cloud; a plan will try to recreate it."
	case StatusUntracked:
		return "In the cloud but absent from state; needs terraform import."
	case StatusMoved:
		return "Renamed in the cloud; without a moved block it is destroyed and recreated."
	}
	return ""
}

func countStatuses(statuses []ResourceStatus) Counts {
	c := Counts{ByStatus: map[string]int{
		StatusInSync: 0, StatusDrift: 0, StatusPhantom: 0, StatusUntracked: 0, StatusMoved: 0,
	}}
	for _, s := range statuses {
		c.Total++
		c.ByStatus[s.Status]++
		if s.Address != "" {
			c.Tracked++
		}
		if s.Status != StatusInSync {
			c.Problems++
		}
	}
	return c
}
