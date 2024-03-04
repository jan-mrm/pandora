// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"strings"

	"github.com/hashicorp/pandora/tools/data-api-sdk/v1/models"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/parser/internal"
)

var _ customFieldMatcher = systemOrUserAssignedIdentityMapMatcher{}

type systemOrUserAssignedIdentityMapMatcher struct{}

func (systemOrUserAssignedIdentityMapMatcher) ReplacementObjectDefinition() models.SDKObjectDefinition {
	return models.SDKObjectDefinition{
		Type: models.SystemOrUserAssignedIdentityMapSDKObjectDefinitionType,
	}
}

func (systemOrUserAssignedIdentityMapMatcher) IsMatch(field models.SDKField, known internal.ParseResult) bool {
	if field.ObjectDefinition.Type != models.ReferenceSDKObjectDefinitionType {
		return false
	}

	// retrieve the model from the reference
	model, ok := known.Models[*field.ObjectDefinition.ReferenceName]
	if !ok {
		return false
	}

	hasUserAssignedIdentities := false
	hasMatchingType := false
	hasPrincipalId := false
	hasTenantId := false

	for fieldName, fieldVal := range model.Fields {
		if strings.EqualFold(fieldName, "PrincipalId") {
			hasPrincipalId = true
			continue
		}

		if strings.EqualFold(fieldName, "TenantId") {
			hasTenantId = true
			continue
		}

		if strings.EqualFold(fieldName, "UserAssignedIdentities") {
			// this should be a Map of an Object containing ClientId/PrincipalId
			if fieldVal.ObjectDefinition.Type != models.DictionarySDKObjectDefinitionType {
				continue
			}
			if fieldVal.ObjectDefinition.NestedItem == nil || fieldVal.ObjectDefinition.NestedItem.Type != models.ReferenceSDKObjectDefinitionType {
				continue
			}

			inlinedModel, ok := known.Models[*fieldVal.ObjectDefinition.NestedItem.ReferenceName]
			if !ok {
				continue
			}

			innerHasClientId := false
			innerHasPrincipalId := false
			for innerName, innerVal := range inlinedModel.Fields {
				if strings.EqualFold(innerName, "ClientId") {
					if innerVal.ObjectDefinition.Type != models.StringSDKObjectDefinitionType {
						continue
					}

					innerHasClientId = true
					continue
				}

				if strings.EqualFold(innerName, "PrincipalId") {
					if innerVal.ObjectDefinition.Type != models.StringSDKObjectDefinitionType {
						continue
					}

					innerHasPrincipalId = true
					continue
				}

				// if extra fields are returned within the UAI properties block then we ignore them for now
			}

			hasUserAssignedIdentities = innerHasClientId && innerHasPrincipalId
			continue
		}

		if strings.EqualFold(fieldName, "Type") {
			if fieldVal.ObjectDefinition.Type != models.ReferenceSDKObjectDefinitionType {
				continue
			}
			constant, ok := known.Constants[*fieldVal.ObjectDefinition.ReferenceName]
			if !ok {
				continue
			}
			expected := map[string]string{
				"SystemAssigned": "SystemAssigned",
				"UserAssigned":   "UserAssigned",
			}
			hasMatchingType = validateIdentityConstantValues(constant, expected)
			continue
		}

		// Per the API Definition, the `DelegatedResources` field is for `internal use only` - therefore we should ignore this if it's present
		if strings.EqualFold(fieldName, "DelegatedResources") {
			continue
		}

		// other fields
		return false
	}

	return hasUserAssignedIdentities && hasMatchingType && hasPrincipalId && hasTenantId
}
