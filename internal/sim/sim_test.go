package sim

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/leetrout/terraform-sim/internal/cloud"
)

// newDemoSim returns a simulator seeded with the baseline cloud and a matching
// demo terraform.tfstate written into a temp work dir.
func newDemoSim(t *testing.T) *Simulator {
	t.Helper()
	store := cloud.NewStore(filepath.Join(t.TempDir(), "cloud.json.gz"))
	s := New(store, t.TempDir())
	if err := s.SeedDemo(true); err != nil {
		t.Fatalf("seed demo: %v", err)
	}
	return s
}

func statusByType(statuses []ResourceStatus, tfType, name string) (ResourceStatus, bool) {
	for _, st := range statuses {
		if st.Type == tfType && st.Name == name {
			return st, true
		}
	}
	return ResourceStatus{}, false
}

func TestBaselineAllInSync(t *testing.T) {
	s := newDemoSim(t)
	statuses, counts := s.Status()
	if counts.Total != 7 {
		t.Fatalf("expected 7 resources, got %d", counts.Total)
	}
	if counts.Problems != 0 {
		t.Fatalf("expected 0 problems at baseline, got %d (%+v)", counts.Problems, counts.ByStatus)
	}
	if counts.Tracked != 7 {
		t.Fatalf("expected 7 tracked, got %d", counts.Tracked)
	}
	for _, st := range statuses {
		if st.Status != StatusInSync {
			t.Fatalf("%s expected in_sync, got %s", st.Address, st.Status)
		}
	}
}

func TestScenarioPhantom(t *testing.T) {
	s := newDemoSim(t)
	if _, err := s.LoadScenario("phantom"); err != nil {
		t.Fatalf("load phantom: %v", err)
	}
	statuses, counts := s.Status()
	if counts.Problems != 1 {
		t.Fatalf("expected 1 problem, got %d", counts.Problems)
	}
	st, ok := statusByType(statuses, "nimbus_storage_bucket", "assets")
	if !ok {
		t.Fatal("assets bucket missing from status")
	}
	if st.Status != StatusPhantom {
		t.Fatalf("expected phantom, got %s", st.Status)
	}
}

func TestScenarioDrift(t *testing.T) {
	s := newDemoSim(t)
	if _, err := s.LoadScenario("drift"); err != nil {
		t.Fatalf("load drift: %v", err)
	}
	statuses, _ := s.Status()
	st, _ := statusByType(statuses, "nimbus_subnet", "private")
	if st.Status != StatusDrift {
		t.Fatalf("expected drift, got %s", st.Status)
	}
	if len(st.Diffs) != 1 || st.Diffs[0].Field != "cidr_block" {
		t.Fatalf("expected a cidr_block diff, got %+v", st.Diffs)
	}
}

func TestScenarioMoved(t *testing.T) {
	s := newDemoSim(t)
	if _, err := s.LoadScenario("moved"); err != nil {
		t.Fatalf("load moved: %v", err)
	}
	statuses, _ := s.Status()
	// The interface is renamed in the cloud; only the name attribute differs.
	var found bool
	for _, st := range statuses {
		if st.Type == "nimbus_network_interface" && st.Status == StatusMoved {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a moved network interface, statuses=%+v", statuses)
	}
}

func TestScenarioUntracked(t *testing.T) {
	s := newDemoSim(t)
	if _, err := s.LoadScenario("untracked"); err != nil {
		t.Fatalf("load untracked: %v", err)
	}
	statuses, counts := s.Status()
	if counts.Total != 8 {
		t.Fatalf("expected 8 resources after untracked create, got %d", counts.Total)
	}
	st, ok := statusByType(statuses, "nimbus_storage_bucket", "uploads")
	if !ok {
		t.Fatal("uploads bucket missing")
	}
	if st.Status != StatusUntracked {
		t.Fatalf("expected untracked, got %s", st.Status)
	}
}

func TestResetReturnsToBaseline(t *testing.T) {
	s := newDemoSim(t)
	_, _ = s.LoadScenario("phantom")
	if err := s.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}
	_, counts := s.Status()
	if counts.Problems != 0 {
		t.Fatalf("expected 0 problems after reset, got %d", counts.Problems)
	}
}

func TestGraphHasNodesAndEdges(t *testing.T) {
	s := newDemoSim(t)
	g := s.Graph("")
	if len(g.Nodes) != 7 {
		t.Fatalf("expected 7 nodes, got %d", len(g.Nodes))
	}
	// subnet→network, nic→subnet, vm→nic = at least 4 edges (2 subnets).
	if len(g.Edges) < 4 {
		t.Fatalf("expected >=4 edges, got %d (%+v)", len(g.Edges), g.Edges)
	}
}

func TestStateWatcherDetectsWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "terraform.tfstate")
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed := make(chan struct{}, 1)
	w := newWatcher(path, func() { changed <- struct{}{} })
	w.interval = 20 * time.Millisecond
	w.start()
	defer w.stop()
	time.Sleep(40 * time.Millisecond) // let the watcher record the initial fingerprint
	if err := os.WriteFile(path, []byte(`{"version":4}`), 0o644); err != nil {
		t.Fatal(err)
	}
	select {
	case <-changed:
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not fire on write")
	}
}
