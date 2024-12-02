// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// snapshotsSchema returns the schema to use for the snapshots property
func snapshotsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_at": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date when the block-device snapshot was created. The date is in ISO 8601 and UTC.",
				},
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A human-friendly description.",
				},
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The id of this block-device snapshot.",
				},
				"is_current": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether this snapshot is the current snapshot on the block-device.",
				},
				"links": linksSchema(),
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A human-friendly name for the block-device snapshot.",
				},
				"org_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The id of the organization this block-device snapshot belongs to.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Email of the user user that owns the block-device snapshot.",
				},
				"updated_at": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date when the entity was last updated. The date is in ISO 8601 and UTC.",
				},
			},
		},
	}
}

func flattenSnapshots(diskSnapshots []*models.DiskSnapshot) []map[string]interface{} {
	if len(diskSnapshots) == 0 {
		return make([]map[string]interface{}, 0)
	}

	snapshots := make([]map[string]interface{}, 0, len(diskSnapshots))

	for _, diskSnapshot := range diskSnapshots {
		helper := make(map[string]interface{})

		helper["created_at"] = diskSnapshot.CreatedAt
		helper["description"] = diskSnapshot.Desc
		helper["name"] = diskSnapshot.Name
		helper["id"] = diskSnapshot.ID
		helper["org_id"] = diskSnapshot.OrgID
		helper["owner"] = diskSnapshot.Owner
		helper["updated_at"] = diskSnapshot.UpdatedAt
		helper["links"] = flattenLinks(diskSnapshot.Links)

		helper["is_current"] = false
		if isCurrent, ok := diskSnapshot.SnapshotProperties["isCurrent"]; ok && isCurrent == "true" {
			helper["is_current"] = true
		}

		snapshots = append(snapshots, helper)
	}

	return snapshots
}
