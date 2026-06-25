// Package cloud implements "Nimbus", the fake cloud that backs statesim. It is
// the source of truth for "what is actually deployed". A real Terraform provider
// (see cmd/terraform-provider-nimbus) talks to it over HTTP via the handlers in
// this package, while the simulator (internal/sim) reads the store in-process.
package cloud

// Meta is embedded in every Nimbus resource. Because it is an anonymous embedded
// struct, its `id` and `serial` fields are promoted to the top level of the JSON
// object, keeping the wire format flat (e.g. {"id":...,"serial":...,"name":...}).
type Meta struct {
	// ID is the server-assigned identifier and the join key between the fake
	// cloud and terraform.tfstate. It is prefixed per type (net-, subnet-, ...).
	ID string `json:"id"`
	// Serial is bumped on every mutation so observers can detect changes.
	Serial int `json:"serial"`
}

func (m *Meta) GetID() string   { return m.ID }
func (m *Meta) SetID(id string) { m.ID = id }
func (m *Meta) GetSerial() int  { return m.Serial }
func (m *Meta) SetSerial(s int) { m.Serial = s }
func (m *Meta) bumpSerial()     { m.Serial++ }

// ref is a dependency reference: an attribute on a child resource pointing at the
// id of a parent in another collection. Used for referential validation and graph
// edges.
type ref struct {
	field      string // attribute name, e.g. "network_id"
	collection string // parent collection, e.g. "networks"
	id         string // the referenced id (empty means "no reference")
}

// Resource is the behaviour every Nimbus resource type implements. Concrete types
// embed Meta (for id/serial) and add their own attributes.
type Resource interface {
	GetID() string
	SetID(string)
	GetSerial() int
	SetSerial(int)
	bumpSerial()

	// kind returns the static descriptor for this resource type.
	kind() Kind
	// attributes returns the comparable attributes (everything the diff engine
	// joins on), excluding id and serial. Keys are snake_case wire names.
	attributes() map[string]string
	// refs returns this resource's parent references for validation/graph edges.
	refs() []ref
	// clone returns a deep copy so callers cannot mutate stored objects.
	clone() Resource
}

// Kind is the static descriptor tying a Terraform type to its REST collection and
// id prefix.
type Kind struct {
	Type       string // Terraform resource type, e.g. "nimbus_network"
	Collection string // REST collection path segment, e.g. "networks"
	Prefix     string // id prefix, e.g. "net-"
}

// kinds is the registry of all five resource types, in dependency order.
var kinds = []Kind{
	{Type: "nimbus_network", Collection: "networks", Prefix: "net-"},
	{Type: "nimbus_subnet", Collection: "subnets", Prefix: "subnet-"},
	{Type: "nimbus_network_interface", Collection: "network-interfaces", Prefix: "nic-"},
	{Type: "nimbus_vm_instance", Collection: "vm-instances", Prefix: "vm-"},
	{Type: "nimbus_storage_bucket", Collection: "storage-buckets", Prefix: "bkt-"},
}

// Kinds returns the registry (dependency order) for callers that need to iterate
// every resource type.
func Kinds() []Kind { return append([]Kind(nil), kinds...) }

// KindByCollection looks up a Kind by its REST collection segment.
func KindByCollection(collection string) (Kind, bool) {
	for _, k := range kinds {
		if k.Collection == collection {
			return k, true
		}
	}
	return Kind{}, false
}

// KindByType looks up a Kind by its Terraform type string.
func KindByType(t string) (Kind, bool) {
	for _, k := range kinds {
		if k.Type == t {
			return k, true
		}
	}
	return Kind{}, false
}

// ParentField describes a dependency reference of a resource type: the attribute
// that holds a parent id, and the parent's collection.
type ParentField struct {
	Field      string
	Collection string
}

// ParentFields returns the dependency reference fields for a Terraform type. It
// works without an instance, so the simulator can compute edges for state-only
// (phantom) resources too.
func ParentFields(tfType string) []ParentField {
	switch tfType {
	case "nimbus_subnet":
		return []ParentField{{"network_id", "networks"}}
	case "nimbus_network_interface":
		return []ParentField{{"subnet_id", "subnets"}}
	case "nimbus_vm_instance":
		return []ParentField{{"nic_id", "network-interfaces"}}
	}
	return nil
}

// Attributes returns a resource's comparable attributes (snake_case, excluding
// id/serial). Exported wrapper so other packages can diff resources.
func Attributes(r Resource) map[string]string { return r.attributes() }

// KindOf returns the static descriptor for a resource instance.
func KindOf(r Resource) Kind { return r.kind() }

// NameOf returns a resource's "name" attribute (empty if it has none).
func NameOf(r Resource) string { return r.attributes()["name"] }

// newOf returns a zero-valued, decodable resource for the given collection.
func newOf(collection string) Resource {
	switch collection {
	case "networks":
		return &Network{}
	case "subnets":
		return &Subnet{}
	case "network-interfaces":
		return &NetworkInterface{}
	case "vm-instances":
		return &VMInstance{}
	case "storage-buckets":
		return &StorageBucket{}
	}
	return nil
}

// --- concrete resource types ---------------------------------------------------

// Network is a virtual network (the root of the topology).
type Network struct {
	Meta
	Name      string `json:"name"`
	CIDRBlock string `json:"cidr_block"`
	Region    string `json:"region"`
}

func (n *Network) kind() Kind { return kinds[0] }
func (n *Network) attributes() map[string]string {
	return map[string]string{"name": n.Name, "cidr_block": n.CIDRBlock, "region": n.Region}
}
func (n *Network) refs() []ref { return nil }
func (n *Network) clone() Resource {
	c := *n
	return &c
}

// Subnet is a subnet within a Network.
type Subnet struct {
	Meta
	Name      string `json:"name"`
	NetworkID string `json:"network_id"`
	CIDRBlock string `json:"cidr_block"`
}

func (s *Subnet) kind() Kind { return kinds[1] }
func (s *Subnet) attributes() map[string]string {
	return map[string]string{"name": s.Name, "network_id": s.NetworkID, "cidr_block": s.CIDRBlock}
}
func (s *Subnet) refs() []ref {
	return []ref{{field: "network_id", collection: "networks", id: s.NetworkID}}
}
func (s *Subnet) clone() Resource {
	c := *s
	return &c
}

// NetworkInterface is a NIC attached to a Subnet.
type NetworkInterface struct {
	Meta
	Name      string `json:"name"`
	SubnetID  string `json:"subnet_id"`
	PrivateIP string `json:"private_ip"`
}

func (ni *NetworkInterface) kind() Kind { return kinds[2] }
func (ni *NetworkInterface) attributes() map[string]string {
	return map[string]string{"name": ni.Name, "subnet_id": ni.SubnetID, "private_ip": ni.PrivateIP}
}
func (ni *NetworkInterface) refs() []ref {
	return []ref{{field: "subnet_id", collection: "subnets", id: ni.SubnetID}}
}
func (ni *NetworkInterface) clone() Resource {
	c := *ni
	return &c
}

// VMInstance is a compute instance attached to a NetworkInterface.
type VMInstance struct {
	Meta
	Name        string `json:"name"`
	MachineType string `json:"machine_type"`
	NICID       string `json:"nic_id"`
}

func (vm *VMInstance) kind() Kind { return kinds[3] }
func (vm *VMInstance) attributes() map[string]string {
	return map[string]string{"name": vm.Name, "machine_type": vm.MachineType, "nic_id": vm.NICID}
}
func (vm *VMInstance) refs() []ref {
	return []ref{{field: "nic_id", collection: "network-interfaces", id: vm.NICID}}
}
func (vm *VMInstance) clone() Resource {
	c := *vm
	return &c
}

// StorageBucket is a standalone object store (no dependencies).
type StorageBucket struct {
	Meta
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (b *StorageBucket) kind() Kind { return kinds[4] }
func (b *StorageBucket) attributes() map[string]string {
	return map[string]string{"name": b.Name, "region": b.Region}
}
func (b *StorageBucket) refs() []ref { return nil }
func (b *StorageBucket) clone() Resource {
	c := *b
	return &c
}
