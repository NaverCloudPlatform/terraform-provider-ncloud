// Copyright (c) 2017, 2020, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package ncloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

// Get the schema for a nested DataSourceSchema generated from the ResourceSchema
func GetDataSourceItemSchema(resourceSchema *schema.Resource) *schema.Resource {
	if _, idExists := resourceSchema.Schema["id"]; !idExists {
		resourceSchema.Schema["id"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}

	// Ensure Create/Read are not set for nested sub-resource schemas. Otherwise, terraform will validate them
	// as though they were resources.
	resourceSchema.Create = nil
	resourceSchema.Read = nil

	return convertResourceFieldsToDatasourceFields(resourceSchema)
}

// Get the Singular DataSource Schema from Resource Schema with additional fields and Read Function
func GetSingularDataSourceItemSchema(resourceSchema *schema.Resource, addFieldMap map[string]*schema.Schema, readFunc schema.ReadFunc) *schema.Resource {
	if _, idExists := resourceSchema.Schema["id"]; !idExists {
		resourceSchema.Schema["id"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}

	// Ensure Create,Read, Update and Delete are not set for data source schemas. Otherwise, terraform will validate them
	// as though they were resources.
	resourceSchema.Create = nil
	resourceSchema.Update = nil
	resourceSchema.Delete = nil
	resourceSchema.Read = readFunc
	resourceSchema.Importer = nil
	resourceSchema.Timeouts = nil
	resourceSchema.CustomizeDiff = nil

	var dataSourceSchema *schema.Resource = convertResourceFieldsToDatasourceFields(resourceSchema)

	for key, value := range addFieldMap {
		dataSourceSchema.Schema[key] = value
	}

	return dataSourceSchema
}

// Get the Singular DataSource Schema from Resource Schema with additional fields and ReadContext Function
func GetSingularDataSourceItemSchemaContext(resourceSchema *schema.Resource, addFieldMap map[string]*schema.Schema, readFunc schema.ReadContextFunc) *schema.Resource {
	if _, idExists := resourceSchema.Schema["id"]; !idExists {
		resourceSchema.Schema["id"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}

	// Ensure Create,Read, Update and Delete are not set for data source schemas. Otherwise, terraform will validate them
	// as though they were resources.
	resourceSchema.CreateContext = nil
	resourceSchema.ReadContext = readFunc
	resourceSchema.UpdateContext = nil
	resourceSchema.DeleteContext = nil
	resourceSchema.Create = nil
	resourceSchema.Update = nil
	resourceSchema.Delete = nil
	resourceSchema.Read = nil
	resourceSchema.Importer = nil
	resourceSchema.Timeouts = nil
	resourceSchema.CustomizeDiff = nil

	var dataSourceSchema *schema.Resource = convertResourceFieldsToDatasourceFields(resourceSchema)

	for key, value := range addFieldMap {
		dataSourceSchema.Schema[key] = value
	}

	return dataSourceSchema
}

// This is mainly used to ensure that fields of a datasource item are compliant with Terraform schema validation
// All datasource return items should have computed-only fields; and not require Diff, Validation, or Default settings.
func convertResourceFieldsToDatasourceFields(resourceSchema *schema.Resource) *schema.Resource {
	resultSchema := map[string]*schema.Schema{}
	for k, fieldSchema := range resourceSchema.Schema {
		isComputed := fieldSchema.Required || fieldSchema.Computed
		fieldSchema.Computed = true
		fieldSchema.Required = false
		fieldSchema.Optional = false
		fieldSchema.DiffSuppressFunc = nil
		fieldSchema.ValidateFunc = nil
		fieldSchema.ValidateDiagFunc = nil
		fieldSchema.ConflictsWith = nil
		fieldSchema.Default = nil
		fieldSchema.MaxItems = 0
		fieldSchema.MaxItems = 0
		if fieldSchema.Type == schema.TypeSet {
			fieldSchema.Type = schema.TypeList
			fieldSchema.Set = nil
		}

		if fieldSchema.Elem != nil {
			if resource, ok := fieldSchema.Elem.(*schema.Resource); ok {
				fieldSchema.Elem = convertResourceFieldsToDatasourceFields(resource)
			}
		}

		if isComputed {
			resultSchema[k] = fieldSchema
		}
	}

	resourceSchema.Schema = resultSchema
	return resourceSchema
}

// SetSingularResourceDataFromMap Set the Singular DataSource from Map
func SetSingularResourceDataFromMap(d *schema.ResourceData, resources map[string]interface{}) {
	for k, v := range resources {
		if k == "id" {
			d.SetId(v.(string))
			continue
		}
		d.Set(k, v)
	}
}

// SetSingularResourceDataFromMapSchema Set the Singular DataSource from Map
func SetSingularResourceDataFromMapSchema(resourceSchema *schema.Resource, d *schema.ResourceData, resources map[string]interface{}) {
	for k, fieldSchema := range resourceSchema.Schema {
		if resources[k] != nil {
			if fieldSchema.Computed || fieldSchema.Required {
				if k == "id" {
					d.SetId(resources[k].(string))
				} else {
					d.Set(k, resources[k])
				}
			} else {
				log.Printf("[TRACE] SetSingularResourceDataFromMapSchema >> [%s] is not nil but Not Computed", k)
			}
		} else {
			log.Printf("[TRACE] SetSingularResourceDataFromMapSchema >> [%s] is nil", k)
		}
	}

	for k, _ := range resources {
		if k == "id" {
			continue
		}

		if resourceSchema.Schema[k] == nil {
			log.Printf("[TRACE] SetSingularResourceDataFromMapSchema >> [%s] is doesn't have schema", k)
		}
	}
}

// GetValueClassicOrVPC get value classic for vpc
func GetValueClassicOrVPC(config *ProviderConfig, classicValue, vpcValue string) string {
	if config.SupportVPC {
		return vpcValue
	}
	return classicValue
}
