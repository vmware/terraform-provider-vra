// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"testing"
)

func TestExpandImageMapping(t *testing.T) {
	c1 := map[string]interface{}{"mandatory": true, "expression": "Foo:Bar"}
	c2 := map[string]interface{}{"mandatory": false, "expression": "Env:Test"}

	imageMappings := make([]interface{}, 0)
	expandedImageMappings := expandImageMapping(imageMappings)

	if len(expandedImageMappings) != 0 {
		t.Errorf("error while expanding when there are no image mappings")
	}

	constraints := make([]interface{}, 0)
	constraints = append(constraints, c1)
	constraints = append(constraints, c2)

	im1 := make(map[string]interface{})
	im1["cloud_config"] = "runcmd sleep 60"
	im1["constraints"] = schema.NewSet(testSetFunc, constraints)
	im1["image_id"] = "vit-1234"
	im1["name"] = "Ubuntu"

	im2 := make(map[string]interface{})
	im2["cloud_config"] = "runcmd sleep 60"
	im2["constraints"] = schema.NewSet(testSetFunc, constraints)
	im2["image_name"] = "Template: CentOs"
	im2["name"] = "CentOs"

	imageMappings = append(imageMappings, im1)
	imageMappings = append(imageMappings, im2)

	expandedImageMappings = expandImageMapping(imageMappings)

	if len(expandedImageMappings) != 2 {
		t.Errorf("not all image mappings are expanded correctly")
	}

	if _, found := expandedImageMappings["Ubuntu"]; !found {
		t.Errorf("image mapping for Ubuntu not expanded correctly")
	}

	if _, found := expandedImageMappings["CentOs"]; !found {
		t.Errorf("image mapping for CentOs not expanded correctly")
	}

	id1 := expandedImageMappings["Ubuntu"]
	if id1.Name != "" || id1.ID != im1["image_id"] || id1.CloudConfig != im1["cloud_config"] || len(id1.Constraints) != im1["constraints"].(*schema.Set).Len() {
		t.Errorf("image description not created properly for Ubuntu image mapping")
	}

	id2 := expandedImageMappings["CentOs"]
	if id2.Name != im2["image_name"] || id2.ID != "" || id2.CloudConfig != im2["cloud_config"] || len(id2.Constraints) != im2["constraints"].(*schema.Set).Len() {
		t.Errorf("image description not created properly for CentOs image mapping")
	}
}

func TestFlattenImageMapping(t *testing.T) {
	imageDescriptions := make(map[string]models.ImageMappingDescription)
	flattenedImageMappings := flattenImageMappings(imageDescriptions)

	if len(flattenedImageMappings) != 0 {
		t.Errorf("error while flattening when there are no image mapping descriptions")
	}

	c1 := models.Constraint{Expression: withString("Foo:Bar"), Mandatory: withBool(true)}
	c2 := models.Constraint{Expression: withString("Env:Test"), Mandatory: withBool(false)}
	constraints := make([]*models.Constraint, 0)
	constraints = append(constraints, &c1)
	constraints = append(constraints, &c2)

	imd1 := models.ImageMappingDescription{
		CloudConfig:      "runcmd sleep 60",
		Constraints:      constraints,
		CreatedAt:        "2020-01-01",
		Description:      "description",
		ExternalID:       "vit-1234",
		ExternalRegionID: "Datacenter: Datacenter-2",
		ID:               withString("abc123"),
		IsPrivate:        true,
		Name:             "Template: Ubuntu",
		OrgID:            "abac",
		OsFamily:         "Linux",
		Owner:            "self",
		UpdatedAt:        "2020-02-10",
	}

	imd2 := models.ImageMappingDescription{
		CloudConfig:      "runcmd sleep 60",
		Constraints:      constraints,
		Description:      "description",
		ExternalID:       "vit-5678",
		ExternalRegionID: "Datacenter: Datacenter-2",
		ID:               withString("abc123"),
		IsPrivate:        true,
		Name:             "Template: CentOs",
		OrgID:            "abcd",
		OsFamily:         "Linux",
		Owner:            "self",
	}

	imageDescriptions["Ubuntu"] = imd1
	imageDescriptions["CentOs"] = imd2

	flattenedImageMappings = flattenImageMappings(imageDescriptions)

	if len(flattenedImageMappings) != 2 {
		t.Errorf("not all image mapping descriptions are flattened correctly")
	}

	im1 := flattenedImageMappings[0].(map[string]interface{})
	if im1["cloud_config"].(string) != imd1.CloudConfig ||
		len(im1["constraints"].([]interface{})) != len(imd1.Constraints) ||
		im1["description"].(string) != imd1.Description ||
		im1["external_id"].(string) != imd1.ExternalID ||
		im1["external_region_id"].(string) != imd1.ExternalRegionID ||
		im1["image_id"].(string) != *imd1.ID ||
		im1["private"].(bool) != imd1.IsPrivate ||
		im1["image_name"].(string) != imd1.Name ||
		im1["name"].(string) != "Ubuntu" ||
		im1["organization"].(string) != imd1.OrgID ||
		im1["os_family"].(string) != imd1.OsFamily ||
		im1["owner"].(string) != imd1.Owner {
		t.Errorf("image mapping descriptions 'Ubuntu' is not flattened correctly")
	}

	im2 := flattenedImageMappings[1].(map[string]interface{})
	if im2["cloud_config"].(string) != imd2.CloudConfig ||
		len(im2["constraints"].([]interface{})) != len(imd2.Constraints) ||
		im2["description"].(string) != imd2.Description ||
		im2["external_id"].(string) != imd2.ExternalID ||
		im2["external_region_id"].(string) != imd2.ExternalRegionID ||
		im2["image_id"].(string) != *imd2.ID ||
		im2["private"].(bool) != imd2.IsPrivate ||
		im2["image_name"].(string) != imd2.Name ||
		im2["name"].(string) != "CentOs" ||
		im2["organization"].(string) != imd2.OrgID ||
		im2["os_family"].(string) != imd2.OsFamily ||
		im2["owner"].(string) != imd2.Owner {
		t.Errorf("image mapping descriptions 'CentOs' is not flattened correctly")
	}
}

func TestFlattenImageMappingConstraints(t *testing.T) {
	constraints := make([]*models.Constraint, 0)
	flattenedImageMappingConstraints := flattenImageMappingConstraints(constraints)

	if len(flattenedImageMappingConstraints) != 0 {
		t.Errorf("error while flattening when there are no image mapping constraints")
	}

	c1 := models.Constraint{Expression: withString("Foo:Bar"), Mandatory: withBool(true)}
	c2 := models.Constraint{Expression: withString("Env:Test"), Mandatory: withBool(false)}

	constraints = append(constraints, &c1)
	constraints = append(constraints, &c2)

	flattenedImageMappingConstraints = flattenImageMappingConstraints(constraints)

	if len(flattenedImageMappingConstraints) != 2 {
		t.Errorf("not all image mapping constraints are flattened correctly")
	}

	imc1 := flattenedImageMappingConstraints[0].(map[string]interface{})
	imc2 := flattenedImageMappingConstraints[1].(map[string]interface{})
	if imc1["mandatory"].(bool) != *c1.Mandatory ||
		imc1["expression"].(string) != *c1.Expression ||
		imc2["mandatory"].(bool) != *c2.Mandatory ||
		imc2["expression"].(string) != *c2.Expression {
		t.Errorf("image mapping constraints are not flattened correctly")
	}
}
