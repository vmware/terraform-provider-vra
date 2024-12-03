// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vra

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// routesSchema returns the schema to use for the routes property
func routesSchema(_ bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"health_check_configuration": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"healthy_threshold": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"interval_seconds": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"port": {
								Type:     schema.TypeString,
								Required: true,
							},
							"protocol": {
								Type:     schema.TypeString,
								Required: true,
							},
							"timeout_seconds": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"unhealthy_threshold": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"url_path": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"member_port": {
					Type:     schema.TypeString,
					Required: true,
				},
				"member_protocol": {
					Type:     schema.TypeString,
					Required: true,
				},
				"port": {
					Type:     schema.TypeString,
					Required: true,
				},
				"protocol": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func expandRoutes(configRoutes []interface{}) []*models.RouteConfiguration {
	routes := make([]*models.RouteConfiguration, 0, len(configRoutes))

	for _, configRoutes := range configRoutes {
		routeMap := configRoutes.(map[string]interface{})

		route := models.RouteConfiguration{
			MemberPort:     withString(routeMap["member_port"].(string)),
			MemberProtocol: withString(routeMap["member_protocol"].(string)),
			Port:           withString(routeMap["port"].(string)),
			Protocol:       withString(routeMap["protocol"].(string)),
		}

		if v, ok := routeMap["health_check_configuration"].(map[string]interface{}); ok && len(v) == 1 {
			healthCheckConfigMap := v

			healthCheckConfiguration := models.HealthCheckConfiguration{
				Protocol: healthCheckConfigMap["protocol"].(string),
				Port:     healthCheckConfigMap["port"].(string),
			}

			if v, ok := healthCheckConfigMap["url_path"].(string); ok && v != "" {
				healthCheckConfiguration.URLPath = v
			}

			if v, ok := healthCheckConfigMap["interval_seconds"].(string); ok && v != "" {
				healthCheckConfiguration.IntervalSeconds = healthCheckConfigMap["interval_seconds"].(int32)
			}

			if v, ok := healthCheckConfigMap["timeout_seconds"].(string); ok && v != "" {
				healthCheckConfiguration.TimeoutSeconds = healthCheckConfigMap["timeout_seconds"].(int32)
			}

			if v, ok := healthCheckConfigMap["unhealthy_threshold"].(string); ok && v != "" {
				healthCheckConfiguration.UnhealthyThreshold = healthCheckConfigMap["unhealthy_threshold"].(int32)
			}

			if v, ok := healthCheckConfigMap["healthy_threshold"].(string); ok && v != "" {
				healthCheckConfiguration.HealthyThreshold = healthCheckConfigMap["healthy_threshold"].(int32)
			}

			route.HealthCheckConfiguration = &healthCheckConfiguration
		}

		routes = append(routes, &route)
	}
	return routes
}

func flattenRoutes(routes []*models.RouteConfiguration) []map[string]interface{} {
	if len(routes) == 0 {
		return make([]map[string]interface{}, 0)
	}

	configRoutes := make([]map[string]interface{}, 0, len(routes))

	for _, route := range routes {
		helper := make(map[string]interface{})
		helper["member_port"] = route.MemberPort
		helper["member_protocol"] = route.MemberProtocol
		helper["port"] = route.Port
		helper["protocol"] = route.Protocol

		if route.HealthCheckConfiguration != nil {
			healthCheckConfigMap := make(map[string]interface{})
			healthCheckConfigMap["healthy_threshold"] = strconv.Itoa(int(route.HealthCheckConfiguration.HealthyThreshold))
			healthCheckConfigMap["interval_seconds"] = strconv.Itoa(int(route.HealthCheckConfiguration.IntervalSeconds))
			healthCheckConfigMap["port"] = route.HealthCheckConfiguration.Port
			healthCheckConfigMap["protocol"] = route.HealthCheckConfiguration.Protocol
			healthCheckConfigMap["timeout_seconds"] = strconv.Itoa(int(route.HealthCheckConfiguration.TimeoutSeconds))
			healthCheckConfigMap["unhealthy_threshold"] = strconv.Itoa(int(route.HealthCheckConfiguration.UnhealthyThreshold))
			healthCheckConfigMap["url_path"] = route.HealthCheckConfiguration.URLPath

			helper["health_check_configuration"] = healthCheckConfigMap
		}

		configRoutes = append(configRoutes, helper)
	}

	return configRoutes
}
