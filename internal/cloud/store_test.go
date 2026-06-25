package cloud

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return NewStore(filepath.Join(t.TempDir(), "cloud.json.gz"))
}

func TestCreateAssignsIDAndSerial(t *testing.T) {
	s := newTestStore(t)
	created, err := s.Create(&Network{Name: "main", CIDRBlock: "10.0.0.0/16", Region: "us-west-1"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.GetID() == "" {
		t.Fatal("expected an assigned id")
	}
	if got := KindOf(created).Prefix; got != "net-" {
		t.Fatalf("expected net- prefix, got id %q", created.GetID())
	}
	if created.GetSerial() != 1 {
		t.Fatalf("expected serial 1, got %d", created.GetSerial())
	}
}

func TestReferentialValidation(t *testing.T) {
	s := newTestStore(t)
	// A subnet referencing an unknown network must be rejected.
	if _, err := s.Create(&Subnet{Name: "public", NetworkID: "net-999", CIDRBlock: "10.0.1.0/24"}); err == nil {
		t.Fatal("expected error for unknown network_id")
	}
	net, err := s.Create(&Network{Name: "main", CIDRBlock: "10.0.0.0/16", Region: "us-west-1"})
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	if _, err := s.Create(&Subnet{Name: "public", NetworkID: net.GetID(), CIDRBlock: "10.0.1.0/24"}); err != nil {
		t.Fatalf("expected valid subnet to be accepted, got %v", err)
	}
}

func TestPatchBumpsSerialAndPersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cloud.json.gz")
	s := NewStore(path)
	created, _ := s.Create(&StorageBucket{Name: "logs", Region: "us-west-1"})
	updated, err := s.Patch("storage-buckets", created.GetID(), map[string]any{"region": "us-east-1"})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if updated.GetSerial() != 2 {
		t.Fatalf("expected serial 2 after patch, got %d", updated.GetSerial())
	}
	if Attributes(updated)["region"] != "us-east-1" {
		t.Fatalf("patch did not apply: %v", Attributes(updated))
	}

	// Reload from disk and confirm persistence.
	s2 := NewStore(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	got, ok := s2.Get("storage-buckets", created.GetID())
	if !ok {
		t.Fatal("resource missing after reload")
	}
	if Attributes(got)["region"] != "us-east-1" {
		t.Fatalf("persisted region wrong: %v", Attributes(got))
	}
}

func TestChangeHookFires(t *testing.T) {
	s := newTestStore(t)
	var events []ChangeEvent
	s.Subscribe(func(ev ChangeEvent) { events = append(events, ev) })
	created, _ := s.Create(&StorageBucket{Name: "logs", Region: "us-west-1"})
	_ = s.Delete("storage-buckets", created.GetID())
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Op != "create" || events[1].Op != "delete" {
		t.Fatalf("unexpected event ops: %+v", events)
	}
}

func TestDeleteIsNonCascading(t *testing.T) {
	s := newTestStore(t)
	net, _ := s.Create(&Network{Name: "main", CIDRBlock: "10.0.0.0/16", Region: "us-west-1"})
	sub, _ := s.Create(&Subnet{Name: "public", NetworkID: net.GetID(), CIDRBlock: "10.0.1.0/24"})
	if err := s.Delete("networks", net.GetID()); err != nil {
		t.Fatalf("delete network: %v", err)
	}
	// The subnet must survive (delete does not cascade) so phantoms are easy to make.
	if _, ok := s.Get("subnets", sub.GetID()); !ok {
		t.Fatal("subnet should still exist after parent delete")
	}
}
