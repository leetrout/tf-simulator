package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

// NewNetworkResource creates the nimbus_network resource.
func NewNetworkResource() resource.Resource {
	return &baseResource{spec: resourceSpec{
		typeName:    "nimbus_network",
		collection:  "networks",
		description: "A virtual network (id prefix net-).",
		writable: []attrSpec{
			{name: "name", description: "Human-readable name."},
			{name: "cidr_block", description: "CIDR range for the network."},
			{name: "region", description: "Region the network lives in."},
		},
	}}
}
