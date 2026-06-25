package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

// NewStorageBucketResource creates the nimbus_storage_bucket resource.
func NewStorageBucketResource() resource.Resource {
	return &baseResource{spec: resourceSpec{
		typeName:    "nimbus_storage_bucket",
		collection:  "storage-buckets",
		description: "An object storage bucket (id prefix bkt-).",
		writable: []attrSpec{
			{name: "name", description: "Human-readable name."},
			{name: "region", description: "Region the bucket lives in."},
		},
	}}
}
