package sim

import "testing"

func TestParseStateV4(t *testing.T) {
	raw := []byte(`{
      "version": 4,
      "terraform_version": "1.9.2",
      "serial": 3,
      "resources": [
        {
          "mode": "managed",
          "type": "nimbus_subnet",
          "name": "public",
          "instances": [
            {"attributes": {"id": "subnet-7c1", "name": "public", "network_id": "net-1a20", "cidr_block": "10.0.1.0/24", "serial": 2, "tags": {"k":"v"}}}
          ]
        },
        {
          "mode": "data",
          "type": "nimbus_network",
          "name": "ignored",
          "instances": [{"attributes": {"id": "net-x"}}]
        },
        {
          "mode": "managed",
          "type": "random_pet",
          "name": "skip",
          "instances": [{"attributes": {"id": "fluffy"}}]
        }
      ]
    }`)
	st, err := ParseState(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if st.TerraformVersion != "1.9.2" {
		t.Fatalf("version: %q", st.TerraformVersion)
	}
	if len(st.Resources) != 1 {
		t.Fatalf("expected only the managed nimbus_subnet, got %d", len(st.Resources))
	}
	r := st.Resources[0]
	if r.Address != "nimbus_subnet.public" {
		t.Fatalf("address: %q", r.Address)
	}
	if r.ID != "subnet-7c1" {
		t.Fatalf("id: %q", r.ID)
	}
	if r.Attributes["cidr_block"] != "10.0.1.0/24" {
		t.Fatalf("cidr_block: %q", r.Attributes["cidr_block"])
	}
	if r.Attributes["serial"] != "2" {
		t.Fatalf("expected numeric serial flattened to \"2\", got %q", r.Attributes["serial"])
	}
	if _, ok := r.Attributes["tags"]; ok {
		t.Fatal("composite attribute should have been skipped")
	}
}
