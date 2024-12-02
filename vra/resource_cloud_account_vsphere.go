// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vmware/vra-sdk-go/pkg/client"
	"github.com/vmware/vra-sdk-go/pkg/client/cloud_account"
	"github.com/vmware/vra-sdk-go/pkg/client/request"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudAccountVsphere() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudAccountVsphereCreate,
		ReadContext:   resourceCloudAccountVsphereRead,
		UpdateContext: resourceCloudAccountVsphereUpdate,
		DeleteContext: resourceCloudAccountVsphereDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required arguments
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address or FQDN of the vCenter Server.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this resource instance.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the vCenter Server.",
			},
			"regions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The set of region ids that will be enabled for this cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the vCenter Server.",
			},

			// Optional arguments
			"accept_self_signed_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to accept self signed certificate when connecting to the vCenter Server.",
			},
			"associated_cloud_account_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "NSX-V or NSX-T account ids to associate with this vSphere cloud account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of a data collector vm deployed in the on premise infrastructure.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-friendly description.",
			},
			"tags": tagsSchema(),

			// Computed attributes
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was created. The date is in ISO 8601 and UTC.",
			},
			"links": linksSchema(),
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization this entity belongs to.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user that owns the entity.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the entity was last updated. The date is ISO 8601 and UTC.",
			},
		},
	}
}

func resourceCloudAccountVsphereCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var associatedCloudAccountIDs []string
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("associated_cloud_account_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified associated cloud account ids are not unique"))
		}
		associatedCloudAccountIDs = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	createResp, err := apiClient.CloudAccount.CreateVSphereCloudAccountAsync(
		cloud_account.NewCreateVSphereCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithBody(&models.CloudAccountVsphereSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				AssociatedCloudAccountIds:   associatedCloudAccountIDs,
				CreateDefaultZones:          false,
				Dcid:                        d.Get("dc_id").(string),
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
				Name:                        withString(d.Get("name").(string)),
				Password:                    d.Get("password").(string),
				Regions:                     regions,
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				Username:                    d.Get("username").(string),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountVsphereStateRefreshFunc(*apiClient, *createResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 5 * time.Second,
	}

	resourceIDs, err := stateChangeFunc.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudAccountVsphere := (resourceIDs.([]string))[0]

	d.SetId(cloudAccountVsphere)

	return resourceCloudAccountVsphereRead(ctx, d, m)
}

func resourceCloudAccountVsphereRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	ret, err := apiClient.CloudAccount.GetVSphereCloudAccount(cloud_account.NewGetVSphereCloudAccountParams().WithID(id))
	if err != nil {
		switch err.(type) {
		case *cloud_account.GetVSphereCloudAccountNotFound:
			d.SetId("")
			return diag.Errorf("vsphere cloud account '%s' not found", id)
		}
		return diag.FromErr(err)
	}

	vsphereAccount := *ret.Payload
	d.Set("associated_cloud_account_ids", flattenAssociatedCloudAccountIDs(vsphereAccount.Links))
	d.Set("created_at", vsphereAccount.CreatedAt)
	d.Set("dc_id", vsphereAccount.Dcid)
	d.Set("description", vsphereAccount.Description)
	d.Set("hostname", vsphereAccount.HostName)
	d.Set("name", vsphereAccount.Name)
	d.Set("org_id", vsphereAccount.OrgID)
	d.Set("owner", vsphereAccount.Owner)
	d.Set("updated_at", vsphereAccount.UpdatedAt)
	d.Set("username", vsphereAccount.Username)

	if err := d.Set("links", flattenLinks(vsphereAccount.Links)); err != nil {
		return diag.Errorf("error setting cloud_account_vsphere links - error: %#v", err)
	}
	if err := d.Set("regions", extractIDsFromRegion(vsphereAccount.EnabledRegions)); err != nil {
		return diag.Errorf("error setting cloud_account_vsphere regions - error: %#v", err)
	}

	if err := d.Set("tags", flattenTags(vsphereAccount.Tags)); err != nil {
		return diag.Errorf("Error setting cloud_account_vsphere tags - error: %#v", err)
	}

	return nil
}

func resourceCloudAccountVsphereUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var associatedCloudAccountIDs []string
	var regions []*models.RegionSpecification

	apiClient := m.(*Client).apiClient

	if v, ok := d.GetOk("associated_cloud_account_ids"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified associated cloud account ids are not unique"))
		}
		associatedCloudAccountIDs = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("regions"); ok {
		if !compareUnique(v.(*schema.Set).List()) {
			return diag.FromErr(errors.New("specified regions are not unique"))
		}
		regions = expandRegionSpecificationList(v.(*schema.Set).List())
	}

	id := d.Id()
	updateResp, err := apiClient.CloudAccount.UpdateVSphereCloudAccountAsync(
		cloud_account.NewUpdateVSphereCloudAccountAsyncParams().
			WithAPIVersion(IaaSAPIVersion).
			WithID(id).
			WithBody(&models.UpdateCloudAccountVsphereSpecification{
				AcceptSelfSignedCertificate: d.Get("accept_self_signed_cert").(bool),
				AssociatedCloudAccountIds:   associatedCloudAccountIDs,
				CreateDefaultZones:          false,
				Description:                 d.Get("description").(string),
				HostName:                    withString(d.Get("hostname").(string)),
				Name:                        d.Get("name").(string),
				Password:                    d.Get("password").(string),
				Regions:                     regions,
				Tags:                        expandTags(d.Get("tags").(*schema.Set).List()),
				Username:                    d.Get("username").(string),
			}))
	if err != nil {
		return diag.FromErr(err)
	}

	stateChangeFunc := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{models.RequestTrackerStatusINPROGRESS},
		Refresh:    resourceCloudAccountVsphereStateRefreshFunc(*apiClient, *updateResp.Payload.ID),
		Target:     []string{models.RequestTrackerStatusFINISHED},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateChangeFunc.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudAccountVsphereRead(ctx, d, m)
}

func resourceCloudAccountVsphereDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*Client).apiClient

	id := d.Id()
	if _, _, err := apiClient.CloudAccount.DeleteVSphereCloudAccount(cloud_account.NewDeleteVSphereCloudAccountParams().WithID(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceCloudAccountVsphereStateRefreshFunc(apiClient client.API, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ret, err := apiClient.Request.GetRequestTracker(request.NewGetRequestTrackerParams().WithID(id))
		if err != nil {
			return "", models.RequestTrackerStatusFAILED, err
		}

		status := ret.Payload.Status
		switch *status {
		case models.RequestTrackerStatusFAILED:
			return []string{""}, *status, errors.New(ret.Payload.Message)
		case models.RequestTrackerStatusINPROGRESS:
			return [...]string{id}, *status, nil
		case models.RequestTrackerStatusFINISHED:
			cloudAccountIDs := make([]string, len(ret.Payload.Resources))
			for i, r := range ret.Payload.Resources {
				cloudAccountIDs[i] = strings.TrimPrefix(r, "/iaas/api/cloud-accounts/")
			}
			return cloudAccountIDs, *status, nil
		default:
			return [...]string{id}, ret.Payload.Message, fmt.Errorf("resourceCloudAccountVsphereStateRefreshFunc: unknown status %v", *status)
		}
	}
}
