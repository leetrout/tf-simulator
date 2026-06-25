package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

// NewSubnetResource creates the nimbus_subnet resource.
func NewSubnetResource() resource.Resource {
	return &baseResource{spec: resourceSpec{
		typeName:    "nimbus_subnet",
		collection:  "subnets",
		description: "A subnet within a network (id prefix subnet-).",
		writable: []attrSpec{
			{name: "name", description: "Human-readable name."},
			{name: "network_id", description: "ID of the parent network.", requiresReplace: true},
			{name: "cidr_block", description: "CIDR range for the subnet."},
		},
	}}
}
