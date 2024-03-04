// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transformer

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/pandora/tools/data-api-sdk/v1/models"
	importerModels "github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

// NOTE: this file contains temporary glue code to enable refactoring this tool gradually, component-by-component.

func MapInternalTypesToDataAPISDKTypes(inputApiVersions []importerModels.AzureApiDefinition, resourceProvider, terraformPackageName *string, logger hclog.Logger) (*models.Service, error) {
	apiVersions := make(map[string]models.APIVersion)

	logger.Debug("Mapping API Versions..")
	for _, item := range inputApiVersions {
		logger.Trace(fmt.Sprintf("Mapping Service %q / API Version %q", item.ServiceName, item.ApiVersion))
		mapped, err := mapInternalAPIVersionTypeToDataAPISDKType(item)
		if err != nil {
			return nil, fmt.Errorf("mapping API Version %q: %+v", item.ApiVersion, err)
		}

		apiVersions[item.ApiVersion] = *mapped
	}

	output := models.Service{
		APIVersions:      apiVersions,
		Generate:         true,
		ResourceProvider: resourceProvider,
	}

	if terraformPackageName != nil {
		logger.Debug("Mapping Terraform Definition..")
		resources := make(map[string]models.TerraformResourceDefinition, 0)
		for _, apiVersion := range inputApiVersions {
			for _, apiResource := range apiVersion.Resources {
				if apiResource.Terraform == nil {
					continue
				}
				for key, value := range apiResource.Terraform.Resources {
					mapped, err := mapTerraformResourceDefinitionToSDKType(value)
					if err != nil {
						return nil, fmt.Errorf("mapping the Terraform Resource Definition for %q: %+v", key, err)
					}

					resources[key] = *mapped
				}
			}
		}

		output.TerraformDefinition = &models.TerraformDefinition{
			Resources:            resources,
			TerraformPackageName: *terraformPackageName,
		}
	}

	return &output, nil
}

func mapInternalAPIVersionTypeToDataAPISDKType(input importerModels.AzureApiDefinition) (*models.APIVersion, error) {
	resources := make(map[string]models.APIResource)

	for apiResource, apiResourceDetails := range input.Resources {
		resources[apiResource] = models.APIResource{
			Constants:   apiResourceDetails.Constants,
			Models:      apiResourceDetails.Models,
			Operations:  apiResourceDetails.Operations,
			ResourceIDs: apiResourceDetails.ResourceIds,
		}
	}

	return &models.APIVersion{
		Generate:  true,
		Preview:   input.IsPreviewVersion(),
		Resources: resources,
		Source:    models.AzureRestAPISpecsSourceDataOrigin,
	}, nil
}

func mapTerraformResourceDefinitionToSDKType(input resourcemanager.TerraformResourceDetails) (*models.TerraformResourceDefinition, error) {
	createMethod := mapTerraformMethodDefinitionToSDKType(input.CreateMethod)
	deleteMethod := mapTerraformMethodDefinitionToSDKType(input.DeleteMethod)
	readMethod := mapTerraformMethodDefinitionToSDKType(input.ReadMethod)
	tests := mapTerraformResourceTestsToSDKType(input.Tests)

	mappings, err := mapTerraformMappingsToSDKType(input.Mappings)
	if err != nil {
		return nil, fmt.Errorf("mapping Terraform Mappings: %+v", err)
	}

	schemaModels, err := mapTerraformSchemaModelsToSDKType(input.SchemaModels)
	if err != nil {
		return nil, fmt.Errorf("mapping Terraform Schema Models: %+v", err)
	}

	output := models.TerraformResourceDefinition{
		APIResource:  input.Resource,
		APIVersion:   input.ApiVersion,
		CreateMethod: createMethod,
		DeleteMethod: deleteMethod,
		Documentation: models.TerraformDocumentationDefinition{
			Category:        input.Documentation.Category,
			Description:     input.Documentation.Description,
			ExampleUsageHCL: input.Documentation.ExampleUsageHcl,
		},
		DisplayName:          input.DisplayName,
		Generate:             input.Generate,
		GenerateModel:        input.GenerateModel,
		GenerateIDValidation: input.GenerateIdValidation,
		GenerateSchema:       input.GenerateSchema,
		Mappings:             *mappings,
		ReadMethod:           readMethod,
		ResourceIDName:       input.ResourceIdName,
		ResourceName:         input.ResourceName,
		SchemaModelName:      input.SchemaModelName,
		SchemaModels:         *schemaModels,
		Tests:                tests,
		UpdateMethod:         nil,
	}

	if input.UpdateMethod != nil {
		updateMethod := mapTerraformMethodDefinitionToSDKType(*input.UpdateMethod)
		output.UpdateMethod = pointer.To(updateMethod)
	}

	return &output, nil
}

func mapTerraformSchemaModelsToSDKType(input map[string]resourcemanager.TerraformSchemaModelDefinition) (*map[string]models.TerraformSchemaModel, error) {
	output := make(map[string]models.TerraformSchemaModel)

	for key, value := range input {
		fields, err := mapTerraformSchemaFieldsToSDKType(value.Fields)
		if err != nil {
			return nil, fmt.Errorf("mapping Terraform Schema Fields: %+v", err)
		}

		output[key] = models.TerraformSchemaModel{
			Fields: *fields,
		}
	}

	return &output, nil
}

func mapTerraformSchemaFieldsToSDKType(input map[string]resourcemanager.TerraformSchemaFieldDefinition) (*map[string]models.TerraformSchemaField, error) {
	output := make(map[string]models.TerraformSchemaField)

	for key, value := range input {
		objectDefinition, err := mapTerraformSchemaObjectDefinitionToSDKType(value.ObjectDefinition)
		if err != nil {
			return nil, fmt.Errorf("mapping ObjectDefinition for Field %q: %+v", key, err)
		}

		field := models.TerraformSchemaField{
			Computed: value.Computed,
			Documentation: models.TerraformSchemaFieldDocumentationDefinition{
				Markdown: value.Documentation.Markdown,
			},
			ForceNew:         value.ForceNew,
			HCLName:          value.HclName,
			ObjectDefinition: *objectDefinition,
			Optional:         value.Optional,
			Required:         value.Required,
			Validation:       nil,
		}

		if value.Validation != nil {
			mapped, err := mapTerraformSchemaFieldValidationToSDKType(*value.Validation)
			if err != nil {
				return nil, fmt.Errorf("mapping Validation for Field %q: %+v", key, err)
			}
			field.Validation = mapped
		}

		output[key] = field
	}

	return &output, nil
}

func mapTerraformSchemaObjectDefinitionToSDKType(input resourcemanager.TerraformSchemaFieldObjectDefinition) (*models.TerraformSchemaObjectDefinition, error) {
	val, ok := terraformObjectDefinitionToSDKType[input.Type]
	if !ok {
		return nil, fmt.Errorf("internal-error: missing mapping for Terraform Schema Object Definition Type %q", string(input.Type))
	}

	output := models.TerraformSchemaObjectDefinition{
		NestedObject:  nil,
		ReferenceName: input.ReferenceName,
		Type:          val,
	}

	if input.NestedObject != nil {
		nested, err := mapTerraformSchemaObjectDefinitionToSDKType(*input.NestedObject)
		if err != nil {
			return nil, fmt.Errorf("mapping NestedObject: %+v", err)
		}

		output.NestedObject = nested
	}

	return &output, nil
}

func mapTerraformSchemaFieldValidationToSDKType(input resourcemanager.TerraformSchemaValidationDefinition) (models.TerraformSchemaFieldValidationDefinition, error) {
	if input.Type == resourcemanager.TerraformSchemaValidationTypePossibleValues {
		mapped, ok := terraformSchemaValidationPossibleValuesTypeToSDKType[input.PossibleValues.Type]
		if !ok {
			return nil, fmt.Errorf("internal-error: missing mapping for Possible Values Type %q", string(input.PossibleValues.Type))
		}
		return models.TerraformSchemaFieldValidationPossibleValuesDefinition{
			PossibleValues: &models.TerraformSchemaFieldValidationPossibleValuesDefinitionImpl{
				Type:   mapped,
				Values: input.PossibleValues.Values,
			},
		}, nil
	}

	return nil, fmt.Errorf("internal-error: missing mapping for Schema Field Validation %q", string(input.Type))
}

func mapTerraformResourceTestsToSDKType(input resourcemanager.TerraformResourceTestsDefinition) models.TerraformResourceTestsDefinition {
	return models.TerraformResourceTestsDefinition{
		BasicConfiguration:          input.BasicConfiguration,
		RequiresImportConfiguration: input.RequiresImportConfiguration,
		CompleteConfiguration:       input.CompleteConfiguration,
		Generate:                    input.Generate,
		OtherTests:                  pointer.To(input.OtherTests),
		TemplateConfiguration:       input.TemplateConfiguration,
	}
}

func mapTerraformMappingsToSDKType(input resourcemanager.MappingDefinition) (*models.TerraformMappingDefinition, error) {
	fields, err := mapTerraformFieldMappingsToSDKType(input.Fields)
	if err != nil {
		return nil, fmt.Errorf("mapping Field Mappings: %+v", err)
	}
	modelToModels, err := mapTerraformModelToModelMappingsToSDKType(input.ModelToModels)
	if err != nil {
		return nil, fmt.Errorf("mapping ModelToModel Mappings: %+v", err)
	}
	resourceIds, err := mapTerraformResourceIDMappingsToSDKType(input.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("mapping ResourceId Mappings: %+v", err)
	}

	return &models.TerraformMappingDefinition{
		Fields:        *fields,
		ModelToModels: *modelToModels,
		ResourceID:    *resourceIds,
	}, nil
}

func mapTerraformFieldMappingsToSDKType(input []resourcemanager.FieldMappingDefinition) (*[]models.TerraformFieldMappingDefinition, error) {
	output := make([]models.TerraformFieldMappingDefinition, 0)

	for _, item := range input {
		if item.Type == resourcemanager.DirectAssignmentMappingDefinitionType {
			output = append(output, models.TerraformDirectAssignmentFieldMappingDefinition{
				DirectAssignment: models.TerraformDirectAssignmentFieldMappingDefinitionImpl{
					TerraformSchemaModelName: item.DirectAssignment.SchemaModelName,
					TerraformSchemaFieldName: item.DirectAssignment.SchemaFieldPath,
					SDKModelName:             item.DirectAssignment.SdkModelName,
					SDKFieldName:             item.DirectAssignment.SdkFieldPath,
				},
			})
			continue
		}

		if item.Type == resourcemanager.ModelToModelMappingDefinitionType {
			output = append(output, models.TerraformModelToModelFieldMappingDefinition{
				ModelToModel: models.TerraformModelToModelFieldMappingDefinitionImpl{
					TerraformSchemaModelName: item.ModelToModel.SchemaModelName,
					SDKModelName:             item.ModelToModel.SdkModelName,
					SDKFieldName:             item.ModelToModel.SdkFieldName,
				},
			})
			continue
		}

		return nil, fmt.Errorf("internal-error: missing mapping for Mapping Type %q", string(item.Type))
	}

	return &output, nil
}

func mapTerraformModelToModelMappingsToSDKType(input []resourcemanager.ModelToModelMappingDefinition) (*[]models.TerraformModelToModelMappingDefinition, error) {
	output := make([]models.TerraformModelToModelMappingDefinition, 0)

	for _, item := range input {
		output = append(output, models.TerraformModelToModelMappingDefinition{
			SDKModelName:             item.SdkModelName,
			TerraformSchemaModelName: item.SchemaModelName,
		})
	}

	return &output, nil
}

func mapTerraformResourceIDMappingsToSDKType(input []resourcemanager.ResourceIdMappingDefinition) (*[]models.TerraformResourceIDMappingDefinition, error) {
	output := make([]models.TerraformResourceIDMappingDefinition, 0)

	for _, item := range input {
		output = append(output, models.TerraformResourceIDMappingDefinition{
			ParsedFromParentID:       item.ParsedFromParentID,
			SegmentName:              item.SegmentName,
			TerraformSchemaFieldName: item.SchemaFieldName,
		})
	}

	return &output, nil
}

func mapTerraformMethodDefinitionToSDKType(input resourcemanager.MethodDefinition) models.TerraformMethodDefinition {
	return models.TerraformMethodDefinition{
		Generate:         input.Generate,
		SDKOperationName: input.MethodName,
		TimeoutInMinutes: input.TimeoutInMinutes,
	}
}

var terraformObjectDefinitionToSDKType = map[resourcemanager.TerraformSchemaFieldType]models.TerraformSchemaFieldType{
	resourcemanager.TerraformSchemaFieldTypeBoolean:    models.BooleanTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeDateTime:   models.DateTimeTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeDictionary: models.DictionaryTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeFloat:      models.FloatTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeInteger:    models.IntegerTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeList:       models.ListTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeReference:  models.ReferenceTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeSet:        models.SetTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeString:     models.StringTerraformSchemaObjectDefinitionType,

	resourcemanager.TerraformSchemaFieldTypeEdgeZone:                      models.EdgeZoneTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeIdentitySystemAssigned:        models.SystemAssignedIdentityTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeIdentitySystemAndUserAssigned: models.SystemAndUserAssignedIdentityTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeIdentitySystemOrUserAssigned:  models.SystemOrUserAssignedIdentityTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeIdentityUserAssigned:          models.UserAssignedIdentityTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeLocation:                      models.LocationTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeResourceGroup:                 models.ResourceGroupTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeTags:                          models.TagsTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeZone:                          models.ZoneTerraformSchemaObjectDefinitionType,
	resourcemanager.TerraformSchemaFieldTypeZones:                         models.ZonesTerraformSchemaObjectDefinitionType,
	// NOTE: sku is intentionally not mapped
}

var terraformSchemaValidationPossibleValuesTypeToSDKType = map[resourcemanager.TerraformSchemaValidationPossibleValueType]models.TerraformSchemaFieldValidationPossibleValuesType{
	resourcemanager.TerraformSchemaValidationPossibleValueTypeFloat:  models.FloatTerraformSchemaFieldValidationPossibleValuesType,
	resourcemanager.TerraformSchemaValidationPossibleValueTypeInt:    models.IntegerTerraformSchemaFieldValidationPossibleValuesType,
	resourcemanager.TerraformSchemaValidationPossibleValueTypeString: models.StringTerraformSchemaFieldValidationPossibleValuesType,
}
