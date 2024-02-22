// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceids

import (
	"testing"

	importerModels "github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func TestCommonResourceID_ResourceGroup(t *testing.T) {
	valid := importerModels.ParsedResourceId{
		Constants: map[string]resourcemanager.ConstantDetails{},
		Segments: []resourcemanager.ResourceIdSegment{
			importerModels.StaticResourceIDSegment("subscriptions", "subscriptions"),
			importerModels.SubscriptionIDResourceIDSegment("subscriptionId"),
			importerModels.StaticResourceIDSegment("resourceGroups", "resourceGroups"),
			importerModels.ResourceGroupResourceIDSegment("resourceGroupName"),
		},
	}
	invalid := importerModels.ParsedResourceId{
		Constants: map[string]resourcemanager.ConstantDetails{},
		Segments: []resourcemanager.ResourceIdSegment{
			importerModels.StaticResourceIDSegment("subscriptions", "subscriptions"),
			importerModels.SubscriptionIDResourceIDSegment("subscriptionId"),
			importerModels.StaticResourceIDSegment("resourceGroups", "resourceGroups"),
			importerModels.ResourceGroupResourceIDSegment("resourceGroupName"),
			importerModels.StaticResourceIDSegment("someResource", "someResource"),
			importerModels.UserSpecifiedResourceIDSegment("resourceName"),
		},
	}
	input := []importerModels.ParsedResourceId{
		valid,
		invalid,
	}
	output := switchOutCommonResourceIDsAsNeeded(input)
	for _, actual := range output {
		if normalizedResourceId(actual.Segments) == normalizedResourceId(valid.Segments) {
			if actual.CommonAlias == nil {
				t.Fatalf("Expected `valid` to have the CommonAlias `ResourceGroup` but got nil")
			}
			if *actual.CommonAlias != "ResourceGroup" {
				t.Fatalf("Expected `valid` to have the CommonAlias `ResourceGroup` but got %q", *actual.CommonAlias)
			}

			continue
		}

		if normalizedResourceId(actual.Segments) == normalizedResourceId(invalid.Segments) {
			if actual.CommonAlias != nil {
				t.Fatalf("Expected `invalid` to have no CommonAlias but got %q", *actual.CommonAlias)
			}
			continue
		}

		t.Fatalf("unexpected Resource ID %q", normalizedResourceId(actual.Segments))
	}
}

func TestCommonResourceID_ResourceGroupIncorrectSegment(t *testing.T) {
	input := []importerModels.ParsedResourceId{
		{
			Constants: map[string]resourcemanager.ConstantDetails{},
			Segments: []resourcemanager.ResourceIdSegment{
				importerModels.StaticResourceIDSegment("subscriptions", "subscriptions"),
				importerModels.SubscriptionIDResourceIDSegment("subscriptionId"),
				importerModels.StaticResourceIDSegment("resourceGroups", "resourceGroups"),
				importerModels.ResourceGroupResourceIDSegment("resourceGroupName"),
			},
		},
		{
			Constants: map[string]resourcemanager.ConstantDetails{},
			Segments: []resourcemanager.ResourceIdSegment{
				importerModels.StaticResourceIDSegment("subscriptions", "subscriptions"),
				importerModels.SubscriptionIDResourceIDSegment("subscriptionId"),
				importerModels.StaticResourceIDSegment("resourceGroups", "resourceGroups"),
				importerModels.ResourceGroupResourceIDSegment("sourceResourceGroupName"),
			},
		},
	}
	output := switchOutCommonResourceIDsAsNeeded(input)
	for i, actual := range output {
		t.Logf("testing %d", i)
		if actual.CommonAlias == nil || *actual.CommonAlias != "ResourceGroup" {
			t.Fatalf("expected item %d to be detected as a ResourceGroup but it wasn't", i)
		}
	}
}
