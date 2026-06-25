package sim

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/leetrout/terraform-sim/internal/cloud"
	"github.com/leetrout/terraform-sim/internal/ws"
)

// ScenarioInfo is the catalog entry shown in the Scenarios tab.
type ScenarioInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Badge       string `json:"badge"`
	Status      string `json:"status"` // the status this scenario manufactures
	Description string `json:"description"`
	Affects     string `json:"affects"`
	Command     string `json:"command"` // teaching / fix command
	Loaded      bool   `json:"loaded"`
}

// LoadResult is returned from POST /api/sim/scenarios/{id}/load.
type LoadResult struct {
	Loaded  string `json:"loaded"`
	Title   string `json:"title"`
	Affects string `json:"affects"`
	Command string `json:"command"`
	Note    string `json:"note,omitempty"`
}

// scenarioCatalog is the static metadata for the five teaching scenarios.
var scenarioCatalog = []ScenarioInfo{
	{
		ID: "in_sync", Title: "In Sync", Badge: "BASELINE", Status: StatusInSync,
		Description: "All resources are tracked and current. terraform plan reports no changes.",
		Affects:     "— none —",
		Command:     "terraform plan  # 0 to add, 0 to change, 0 to destroy",
	},
	{
		ID: "untracked", Title: "Missing in State", Badge: "UNTRACKED", Status: StatusUntracked,
		Description: "A storage bucket was created outside Terraform — unmanaged until imported.",
		Affects:     "nimbus_storage_bucket.uploads",
		Command:     "terraform import nimbus_storage_bucket.uploads <id>",
	},
	{
		ID: "phantom", Title: "Phantom in State", Badge: "DEPLOYED", Status: StatusPhantom,
		Description: "nimbus_storage_bucket.assets is in state but was deleted in the cloud. A plan will try to recreate it.",
		Affects:     "nimbus_storage_bucket.assets",
		Command:     "terraform state rm nimbus_storage_bucket.assets",
	},
	{
		ID: "drift", Title: "Config Drift", Badge: "DRIFT", Status: StatusDrift,
		Description: "nimbus_subnet.private.cidr_block changed in the cloud. Applying your config would overwrite it back.",
		Affects:     "nimbus_subnet.private",
		Command:     "terraform apply  # reconcile drift",
	},
	{
		ID: "moved", Title: "Renamed / Moved", Badge: "MOVED", Status: StatusMoved,
		Description: "A network interface was renamed in the cloud. Without a moved block it is destroyed and recreated.",
		Affects:     "nimbus_network_interface.web0",
		Command:     "moved { from = ...web0  to = ...web_primary }",
	},
}

// Scenarios returns the catalog with the loaded flag set for the active scenario.
func (s *Simulator) Scenarios() []ScenarioInfo {
	s.mu.RLock()
	loaded := s.loadedScenario
	s.mu.RUnlock()
	out := make([]ScenarioInfo, len(scenarioCatalog))
	copy(out, scenarioCatalog)
	for i := range out {
		out[i].Loaded = out[i].ID == loaded
	}
	return out
}

// LoadScenario performs the out-of-band cloud mutation that manufactures a
// scenario's condition and returns the teaching command(s).
func (s *Simulator) LoadScenario(id string) (LoadResult, error) {
	info, ok := scenarioByID(id)
	if !ok {
		return LoadResult{}, fmt.Errorf("unknown scenario %q", id)
	}
	res := LoadResult{Loaded: id, Title: info.Title, Affects: info.Affects, Command: info.Command}

	switch id {
	case "in_sync":
		if err := s.seedBaseline(true); err != nil {
			return res, err
		}
	case "untracked":
		region := s.region
		if existing, ok := s.store.List("storage-buckets"); ok && len(existing) > 0 {
			region = cloud.Attributes(existing[0])["region"]
		}
		created, err := s.store.Create(&cloud.StorageBucket{Name: "uploads", Region: region})
		if err != nil {
			return res, err
		}
		res.Command = "terraform import nimbus_storage_bucket.uploads " + created.GetID()
		res.Note = "Created bucket " + created.GetID() + " in the cloud, outside Terraform."
	case "phantom":
		target := s.findByName("storage-buckets", "assets")
		if target == nil {
			return res, fmt.Errorf("no storage bucket to delete; apply the example or seed demo first")
		}
		if err := s.store.Delete("storage-buckets", target.GetID()); err != nil {
			return res, err
		}
		res.Note = "Deleted bucket " + target.GetID() + " from the cloud; state still references it."
	case "drift":
		target := s.findByName("subnets", "private")
		if target == nil {
			return res, fmt.Errorf("no subnet to mutate; apply the example or seed demo first")
		}
		if _, err := s.store.Patch("subnets", target.GetID(), map[string]any{"cidr_block": "10.0.99.0/24"}); err != nil {
			return res, err
		}
		res.Note = "Changed cidr_block on " + target.GetID() + " in the cloud."
	case "moved":
		target := s.findByName("network-interfaces", "web0")
		if target == nil {
			return res, fmt.Errorf("no network interface to rename; apply the example or seed demo first")
		}
		if _, err := s.store.Patch("network-interfaces", target.GetID(), map[string]any{"name": "web_primary"}); err != nil {
			return res, err
		}
		res.Note = "Renamed " + target.GetID() + " to web_primary in the cloud."
	}

	s.mu.Lock()
	s.loadedScenario = id
	s.lastEvent = "scenario.loaded: " + id
	s.mu.Unlock()
	ws.Broadcast(ws.Event{Type: "scenario.loaded", Scope: "sim", Detail: id})
	return res, nil
}

// Reset returns the cloud (and demo state) to the known baseline.
func (s *Simulator) Reset() error {
	if err := s.seedBaseline(true); err != nil {
		return err
	}
	s.mu.Lock()
	s.loadedScenario = "in_sync"
	s.lastEvent = "reset"
	s.mu.Unlock()
	ws.Broadcast(ws.Event{Type: "reset", Scope: "sim"})
	return nil
}

// SeedDemo seeds the baseline cloud and, when writeState is true, writes a
// matching demo terraform.tfstate so the full loop is demoable without running
// terraform.
func (s *Simulator) SeedDemo(writeState bool) error {
	return s.seedBaseline(writeState)
}

// findByName returns the first resource in a collection whose name matches, or the
// first resource if none match (nil if the collection is empty).
func (s *Simulator) findByName(collection, name string) cloud.Resource {
	list, ok := s.store.List(collection)
	if !ok || len(list) == 0 {
		return nil
	}
	for _, r := range list {
		if cloud.NameOf(r) == name {
			return r
		}
	}
	return list[0]
}

func scenarioByID(id string) (ScenarioInfo, bool) {
	for _, sc := range scenarioCatalog {
		if sc.ID == id {
			return sc, true
		}
	}
	return ScenarioInfo{}, false
}

// --- baseline topology ---------------------------------------------------------

// baselineResources returns the seven-resource demo topology with fixed ids.
func baselineResources() []cloud.Resource {
	return []cloud.Resource{
		&cloud.Network{Meta: cloud.Meta{ID: "net-001", Serial: 1}, Name: "main", CIDRBlock: "10.0.0.0/16", Region: "us-west-1"},
		&cloud.Subnet{Meta: cloud.Meta{ID: "subnet-002", Serial: 1}, Name: "public", NetworkID: "net-001", CIDRBlock: "10.0.1.0/24"},
		&cloud.Subnet{Meta: cloud.Meta{ID: "subnet-003", Serial: 1}, Name: "private", NetworkID: "net-001", CIDRBlock: "10.0.2.0/24"},
		&cloud.NetworkInterface{Meta: cloud.Meta{ID: "nic-004", Serial: 1}, Name: "web0", SubnetID: "subnet-002", PrivateIP: "10.0.1.20"},
		&cloud.VMInstance{Meta: cloud.Meta{ID: "vm-005", Serial: 1}, Name: "web", MachineType: "n2-standard-2", NICID: "nic-004"},
		&cloud.StorageBucket{Meta: cloud.Meta{ID: "bkt-006", Serial: 1}, Name: "logs", Region: "us-west-1"},
		&cloud.StorageBucket{Meta: cloud.Meta{ID: "bkt-007", Serial: 1}, Name: "assets", Region: "us-west-1"},
	}
}

// seedBaseline resets the cloud to the baseline topology and optionally writes a
// matching demo state file.
func (s *Simulator) seedBaseline(writeState bool) error {
	base := baselineResources()
	if err := s.store.Reset(base); err != nil {
		return err
	}
	if writeState {
		return writeDemoState(s.stateFile, base)
	}
	return nil
}

// --- demo state writer ---------------------------------------------------------

type stateFileOut struct {
	Version          int           `json:"version"`
	TerraformVersion string        `json:"terraform_version"`
	Serial           int           `json:"serial"`
	Lineage          string        `json:"lineage"`
	Resources        []stateResOut `json:"resources"`
}

type stateResOut struct {
	Mode      string         `json:"mode"`
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Provider  string         `json:"provider"`
	Instances []stateInstOut `json:"instances"`
}

type stateInstOut struct {
	SchemaVersion int            `json:"schema_version"`
	Attributes    map[string]any `json:"attributes"`
}

const nimbusProvider = `provider["registry.terraform.io/leetrout/nimbus"]`

// writeDemoState renders a version-4 terraform.tfstate matching the given cloud
// resources, so the simulator classifies them all as in_sync.
func writeDemoState(path string, resources []cloud.Resource) error {
	out := stateFileOut{
		Version:          4,
		TerraformVersion: "1.9.2",
		Serial:           1,
		Lineage:          "statesim-demo",
	}
	for _, r := range resources {
		k := cloud.KindOf(r)
		attrs := map[string]any{"id": r.GetID(), "serial": r.GetSerial()}
		for key, val := range cloud.Attributes(r) {
			attrs[key] = val
		}
		out.Resources = append(out.Resources, stateResOut{
			Mode:     "managed",
			Type:     k.Type,
			Name:     cloud.NameOf(r),
			Provider: nimbusProvider,
			Instances: []stateInstOut{{
				SchemaVersion: 0,
				Attributes:    attrs,
			}},
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
