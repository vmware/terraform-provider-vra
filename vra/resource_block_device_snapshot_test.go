// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/disk"
)

func TestBlockDeviceSnapshotCreate(t *testing.T) {
	for _, tc := range []struct {
		name       string
		statusCode int
		body       string
		wantID     string
		wantError  string
	}{
		{
			name:       "empty snapshot list",
			statusCode: http.StatusOK,
			body:       `[]`,
			wantError:  "failed to find the created snapshot",
		},
		{
			name:       "no current snapshot",
			statusCode: http.StatusOK,
			body:       `[{"id":"old-snapshot","snapshotProperties":{"isCurrent":"false"}}]`,
			wantError:  "failed to find the created snapshot",
		},
		{
			name:       "snapshot list API error",
			statusCode: http.StatusForbidden,
			body:       `{"message":"snapshot access denied","statusCode":403}`,
			wantError:  "snapshot access denied",
		},
		{
			name:       "current snapshot found",
			statusCode: http.StatusOK,
			body:       `[{"id":"old-snapshot","snapshotProperties":{"isCurrent":"false"}},{"id":"test-snapshot","snapshotProperties":{"isCurrent":"true"}}]`,
			wantID:     "test-snapshot",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			meta := newBlockDeviceSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method + " " + r.URL.Path {
				case "POST /iaas/api/block-devices/test-disk/operations/snapshots":
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprint(w, `{"id":"test-request","status":"INPROGRESS"}`)
				case "GET /iaas/api/request-tracker/test-request":
					fmt.Fprint(w, `{"id":"test-request","status":"FINISHED"}`)
				case "GET /iaas/api/block-devices/test-disk/snapshots":
					w.WriteHeader(tc.statusCode)
					fmt.Fprint(w, tc.body)
				case "GET /iaas/api/block-devices/test-disk/snapshots/test-snapshot":
					fmt.Fprint(w, `{"id":"test-snapshot","desc":"test snapshot","snapshotProperties":{"isCurrent":"true"}}`)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.NotFound(w, r)
				}
			})
			d := schema.TestResourceDataRaw(t, resourceBlockDeviceSnapshot().Schema, map[string]interface{}{
				"block_device_id": "test-disk",
				"description":     "test snapshot",
			})

			diags := resourceBlockDeviceSnapshotCreate(context.Background(), d, meta)
			if tc.wantError != "" {
				if !diags.HasError() {
					t.Fatalf("expected an error containing %q, got %v", tc.wantError, diags)
				}
				if !strings.Contains(diags[0].Summary, tc.wantError) || !strings.Contains(diags[0].Summary, "test-disk") {
					t.Errorf("expected diagnostic to include %q and the block device ID, got %v", tc.wantError, diags)
				}
			} else {
				if diags.HasError() {
					t.Fatalf("unexpected create error: %v", diags)
				}
				if current := d.Get("is_current"); current != true {
					t.Errorf("expected current snapshot to be read into state, got is_current=%v", current)
				}
				if description := d.Get("description"); description != "test snapshot" {
					t.Errorf("expected snapshot description to be read into state, got %q", description)
				}
			}
			if d.Id() != tc.wantID {
				t.Errorf("expected snapshot ID %q, got %q", tc.wantID, d.Id())
			}
		})
	}
}

func TestFindCreatedBlockDeviceSnapshotAPIError(t *testing.T) {
	meta := newBlockDeviceSnapshotTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/iaas/api/block-devices/test-disk/snapshots" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"message":"snapshot access denied","statusCode":403}`)
	})

	id, err := findCreatedBlockDeviceSnapshot("test-disk", meta)
	if err == nil {
		t.Fatal("expected snapshot lookup to fail")
	}
	if id != "" {
		t.Errorf("expected an empty snapshot ID, got %q", id)
	}
	if !strings.Contains(err.Error(), "test-disk") || !strings.Contains(err.Error(), "snapshot access denied") {
		t.Errorf("expected block device ID and API error details, got %v", err)
	}
	var forbidden *disk.GetDiskSnapshotsForbidden
	if !errors.As(err, &forbidden) {
		t.Errorf("expected wrapped GetDiskSnapshotsForbidden error, got %T: %v", err, err)
	}
}

func newBlockDeviceSnapshotTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	transport := httptransport.New(strings.TrimPrefix(server.URL, "http://"), "/", []string{"http"})
	transport.SetDebug(false)
	return &Client{apiClient: client.New(transport, strfmt.Default)}
}

func TestAccVRABlockDeviceSnapshotBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBlockDeviceSnapshotResource(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRABlockDeviceSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRABlockDeviceSnapshotConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRABlockDeviceSnapshotExists("vra_block_device_snapshot.this"),
					resource.TestCheckResourceAttr(
						"vra_block_device_snapshot.this", "description", "terraform block device snapshot"),
					resource.TestCheckResourceAttrSet(
						"vra_block_device_snapshot.this", "is_current"),
				),
			},
		},
	})
}

func testAccCheckVRABlockDeviceSnapshotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no block_device_snapshot ID is set")
		}

		return nil
	}
}

func testAccCheckVRABlockDeviceSnapshotDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "vra_block_device" {
			_, err := apiClient.Disk.GetBlockDevice(disk.NewGetBlockDeviceParams().WithID(rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Resource 'vra_block_device' still exists with id %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckVRABlockDeviceSnapshotConfig(rInt int) string {
	return testAccCheckVRABlockDeviceBasicConfig(rInt) + `
	resource "vra_block_device_snapshot" "this" {
	  block_device_id = vra_block_device.disk1.id
	  description = "terraform block device snapshot"
	}`
}

func testAccCheckVRABlockDeviceBasicConfig(rInt int) string {
	// Need valid credentials since this is creating a real cloud account
	name := os.Getenv("VRA_VSPHERE_CLOUD_ACCOUNT_NAME")
	region := os.Getenv("VRA_REGION")
	projectName := os.Getenv("VRA_PROJECT")
	datastoreName := os.Getenv("VRA_VSPHERE_DATASTORE_NAME")
	storagePolicyName := os.Getenv("VRA_VSPHERE_STORAGE_POLICY_NAME")

	return fmt.Sprintf(`
	data "vra_cloud_account_vsphere" "this" {
		name = "%s"
	  }

	data "vra_region" "this" {
		cloud_account_id = data.vra_cloud_account_vsphere.this.id
		region = "%s"
	}

	data "vra_project" "this" {
		name = "%s"
	 }

	# Lookup vSphere fabric datastore using its name
	data "vra_fabric_datastore_vsphere" "this" {
	  filter = "name eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	# Lookup vSphere fabric storage policy using its name
	data "vra_fabric_storage_policy_vsphere" "this" {
	  filter = "name eq '%s' and cloudAccountId eq '${data.vra_cloud_account_vsphere.this.id}'"
	}

	resource "vra_storage_profile" "this" {
	  name         = "vSphere-first-class-disk"
	  description  = "vSphere Storage Profile with first class disk."
	  region_id    = data.vra_region.this.id
	  default_item = true

	  disk_properties = {
		diskType         = "firstClass"
		provisioningType = "thin" // Supported values: "thin", "thick", "eagerZeroedThick"
	  }

	  disk_target_properties = {
		datastoreId     = data.vra_fabric_datastore_vsphere.this.id
		storagePolicyId = data.vra_fabric_storage_policy_vsphere.this.id
	  }

	  tags {
		key   = "foo"
		value = "bar"
	  }
	}

	resource "vra_block_device" "disk1" {
	  capacity_in_gb = 10
      name = "block-device-%d"
      description = "terraform block device for snapshot"
	  project_id = data.vra_project.this.id
      depends_on = [vra_storage_profile.this]
      persistent = true
      purge = true

	  constraints {
		mandatory  = true
		expression = "foo:bar"
	  }
	}
`, name, region, projectName, datastoreName, storagePolicyName, rInt)
}
