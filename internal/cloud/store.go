package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/schollz/jsonstore"
)

// DefaultStateFile is the gzip-JSON file the cloud persists to (cwd-relative).
const DefaultStateFile = "./tfsim.json.gz"

// storeKey is the jsonstore key under which the cloud snapshot lives.
const storeKey = "cloud"

// ChangeEvent describes a single mutation of the cloud, handed to subscribers.
type ChangeEvent struct {
	Op         string // "create" | "update" | "delete" | "reset"
	Type       string // Terraform type, e.g. "nimbus_subnet" ("" for reset)
	Collection string // REST collection, e.g. "subnets" ("" for reset)
	ID         string // affected id ("" for reset)
	Serial     int    // resource serial after the mutation
}

// ChangeHook is notified after every committed mutation.
type ChangeHook func(ChangeEvent)

// ErrNotFound is returned when an id does not exist in a collection.
var ErrNotFound = errors.New("resource not found")

// Store is the in-memory Nimbus cloud with gzip-JSON persistence. It is safe for
// concurrent use; the simulator reads it in-process while the provider mutates it
// over HTTP.
type Store struct {
	mu      sync.RWMutex
	objects map[string]map[string]Resource // collection -> id -> resource
	counter int
	path    string
	hooks   []ChangeHook
}

// NewStore returns an empty store backed by the given file path.
func NewStore(path string) *Store {
	s := &Store{path: path, objects: map[string]map[string]Resource{}}
	for _, k := range kinds {
		s.objects[k.Collection] = map[string]Resource{}
	}
	return s
}

// Subscribe registers a hook fired after each committed mutation. Hooks run
// synchronously under no lock (the store lock is released before firing).
func (s *Store) Subscribe(h ChangeHook) {
	s.mu.Lock()
	s.hooks = append(s.hooks, h)
	s.mu.Unlock()
}

func (s *Store) fire(ev ChangeEvent) {
	s.mu.RLock()
	hooks := append([]ChangeHook(nil), s.hooks...)
	s.mu.RUnlock()
	for _, h := range hooks {
		h(ev)
	}
}

// --- persistence ---------------------------------------------------------------

type persisted struct {
	Counter     int                                   `json:"counter"`
	Collections map[string]map[string]json.RawMessage `json:"collections"`
}

// Load reads the snapshot from disk if present; a missing file is not an error
// (the store stays empty and is written on the first mutation).
func (s *Store) Load() error {
	if _, err := os.Stat(s.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	ks, err := jsonstore.Open(s.path)
	if err != nil {
		return err
	}
	var p persisted
	if err := ks.Get(storeKey, &p); err != nil {
		// Key absent (e.g. legacy file) — treat as empty cloud.
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = p.Counter
	for _, k := range kinds {
		s.objects[k.Collection] = map[string]Resource{}
		for id, raw := range p.Collections[k.Collection] {
			r := newOf(k.Collection)
			if err := json.Unmarshal(raw, r); err != nil {
				return fmt.Errorf("decode %s/%s: %w", k.Collection, id, err)
			}
			s.objects[k.Collection][id] = r
		}
	}
	return nil
}

// save persists the current snapshot. Caller must hold at least a read lock.
func (s *Store) save() error {
	p := persisted{Counter: s.counter, Collections: map[string]map[string]json.RawMessage{}}
	for _, k := range kinds {
		col := map[string]json.RawMessage{}
		for id, r := range s.objects[k.Collection] {
			raw, err := json.Marshal(r)
			if err != nil {
				return err
			}
			col[id] = raw
		}
		p.Collections[k.Collection] = col
	}
	ks := new(jsonstore.JSONStore)
	if err := ks.Set(storeKey, p); err != nil {
		return err
	}
	return jsonstore.Save(ks, s.path)
}

// --- generic accessors ---------------------------------------------------------

// List returns a copy of every resource in a collection, sorted by id.
func (s *Store) List(collection string) ([]Resource, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.objects[collection]
	if !ok {
		return nil, false
	}
	out := make([]Resource, 0, len(m))
	for _, r := range m {
		out = append(out, r.clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetID() < out[j].GetID() })
	return out, true
}

// Get returns a copy of one resource by collection and id.
func (s *Store) Get(collection, id string) (Resource, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.objects[collection]
	if !ok {
		return nil, false
	}
	r, ok := m[id]
	if !ok {
		return nil, false
	}
	return r.clone(), true
}

// All returns a copy of every resource across all collections, in dependency
// order then id order. Used by the simulator for status and graph building.
func (s *Store) All() []Resource {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Resource
	for _, k := range kinds {
		m := s.objects[k.Collection]
		ids := make([]string, 0, len(m))
		for id := range m {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			out = append(out, m[id].clone())
		}
	}
	return out
}

// New returns a zero-valued resource for a collection (for handler decoding).
func New(collection string) (Resource, bool) {
	r := newOf(collection)
	return r, r != nil
}

// --- mutations -----------------------------------------------------------------

// Create assigns an id and serial, validates references, stores, persists, and
// fires the change hook. The passed resource is mutated in place with its id.
func (s *Store) Create(r Resource) (Resource, error) {
	k := r.kind()
	s.mu.Lock()
	if err := s.validateRefs(r); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	s.counter++
	id := fmt.Sprintf("%s%03x", k.Prefix, s.counter)
	r.SetID(id)
	r.SetSerial(1)
	s.objects[k.Collection][id] = r.clone()
	if err := s.save(); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	stored := s.objects[k.Collection][id].clone()
	s.mu.Unlock()
	s.fire(ChangeEvent{Op: "create", Type: k.Type, Collection: k.Collection, ID: id, Serial: stored.GetSerial()})
	return stored, nil
}

// Put stores a resource at an explicit id (used by replace/import-style writes and
// by scenario seeding). It bumps the serial and validates references.
func (s *Store) Put(collection, id string, r Resource) (Resource, error) {
	k, ok := KindByCollection(collection)
	if !ok {
		return nil, fmt.Errorf("unknown collection %q", collection)
	}
	s.mu.Lock()
	if err := s.validateRefs(r); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	r.SetID(id)
	prevSerial := 0
	if prev, ok := s.objects[collection][id]; ok {
		prevSerial = prev.GetSerial()
	}
	r.SetSerial(prevSerial + 1)
	s.objects[collection][id] = r.clone()
	if err := s.save(); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	stored := s.objects[collection][id].clone()
	s.mu.Unlock()
	s.fire(ChangeEvent{Op: "update", Type: k.Type, Collection: collection, ID: id, Serial: stored.GetSerial()})
	return stored, nil
}

// Patch applies a partial update (a subset of snake_case attribute fields) to an
// existing resource, bumps the serial, validates, persists, and fires the hook.
func (s *Store) Patch(collection, id string, fields map[string]any) (Resource, error) {
	k, ok := KindByCollection(collection)
	if !ok {
		return nil, fmt.Errorf("unknown collection %q", collection)
	}
	s.mu.Lock()
	existing, ok := s.objects[collection][id]
	if !ok {
		s.mu.Unlock()
		return nil, ErrNotFound
	}
	// Round-trip the existing object to JSON, overlay the patch fields (ignoring
	// id/serial), and decode back into a fresh typed resource.
	base, _ := json.Marshal(existing)
	var merged map[string]any
	_ = json.Unmarshal(base, &merged)
	for key, val := range fields {
		if key == "id" || key == "serial" {
			continue
		}
		merged[key] = val
	}
	raw, _ := json.Marshal(merged)
	updated := newOf(collection)
	if err := json.Unmarshal(raw, updated); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	if err := s.validateRefs(updated); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	updated.SetID(id)
	updated.SetSerial(existing.GetSerial() + 1)
	s.objects[collection][id] = updated.clone()
	if err := s.save(); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	stored := s.objects[collection][id].clone()
	s.mu.Unlock()
	s.fire(ChangeEvent{Op: "update", Type: k.Type, Collection: collection, ID: id, Serial: stored.GetSerial()})
	return stored, nil
}

// Delete removes a resource. Deletes are intentionally non-cascading so that
// "phantom" conditions are easy to manufacture.
func (s *Store) Delete(collection, id string) error {
	k, ok := KindByCollection(collection)
	if !ok {
		return fmt.Errorf("unknown collection %q", collection)
	}
	s.mu.Lock()
	if _, ok := s.objects[collection][id]; !ok {
		s.mu.Unlock()
		return ErrNotFound
	}
	delete(s.objects[collection], id)
	if err := s.save(); err != nil {
		s.mu.Unlock()
		return err
	}
	s.mu.Unlock()
	s.fire(ChangeEvent{Op: "delete", Type: k.Type, Collection: collection, ID: id})
	return nil
}

// Reset replaces the entire cloud with the provided resources (already typed,
// with ids set) and resets the id counter to cover them. Fires a single reset
// event.
func (s *Store) Reset(seed []Resource) error {
	s.mu.Lock()
	for _, k := range kinds {
		s.objects[k.Collection] = map[string]Resource{}
	}
	s.counter = 0
	for _, r := range seed {
		k := r.kind()
		s.objects[k.Collection][r.GetID()] = r.clone()
		if n := idCounter(r.GetID(), k.Prefix); n > s.counter {
			s.counter = n
		}
	}
	if err := s.save(); err != nil {
		s.mu.Unlock()
		return err
	}
	s.mu.Unlock()
	s.fire(ChangeEvent{Op: "reset"})
	return nil
}

// validateRefs checks that every parent reference resolves to an existing
// resource. Caller must hold the write lock.
func (s *Store) validateRefs(r Resource) error {
	for _, ref := range r.refs() {
		if ref.id == "" {
			return fmt.Errorf("%s requires %s", r.kind().Type, ref.field)
		}
		if _, ok := s.objects[ref.collection][ref.id]; !ok {
			return fmt.Errorf("%s references unknown %s %q", r.kind().Type, ref.field, ref.id)
		}
	}
	return nil
}

// idCounter parses the hex suffix of an id like "net-1a" into its integer value,
// so Reset can keep the counter ahead of seeded ids. Returns 0 if unparseable.
func idCounter(id, prefix string) int {
	if len(id) <= len(prefix) {
		return 0
	}
	var n int
	if _, err := fmt.Sscanf(id[len(prefix):], "%x", &n); err != nil {
		return 0
	}
	return n
}
