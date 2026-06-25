package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

// NewVMInstanceResource creates the nimbus_vm_instance resource.
func NewVMInstanceResource() resource.Resource {
	return &baseResource{spec: resourceSpec{
		typeName:    "nimbus_vm_instance",
		collection:  "vm-instances",
		description: "A virtual machine instance (id prefix vm-).",
		writable: []attrSpec{
			{name: "name", description: "Human-readable name."},
			{name: "machine_type", description: "Machine type / size."},
			{name: "nic_id", description: "ID of the attached network interface.", requiresReplace: true},
		},
	}}
}
