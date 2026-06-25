package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

// NewNetworkInterfaceResource creates the nimbus_network_interface resource.
func NewNetworkInterfaceResource() resource.Resource {
	return &baseResource{spec: resourceSpec{
		typeName:    "nimbus_network_interface",
		collection:  "network-interfaces",
		description: "A network interface attached to a subnet (id prefix nic-).",
		writable: []attrSpec{
			{name: "name", description: "Human-readable name."},
			{name: "subnet_id", description: "ID of the parent subnet.", requiresReplace: true},
			{name: "private_ip", description: "Private IP address."},
		},
	}}
}
