// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataapigeneratorjson

import (
	"encoding/json"
	"fmt"
	
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/dataapigeneratorjson/transforms"
	importerModels "github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
)

func codeForResourceId(name string, input importerModels.ParsedResourceId) ([]byte, error) {
	resourceId, err := transforms.MapResourceIDToRepository(name, input)
	if err != nil {
		return nil, fmt.Errorf("mapping ResourceID %q: %+v", name, err)
	}

	data, err := json.MarshalIndent(resourceId, "", " ")
	if err != nil {
		return nil, err
	}

	return data, nil
}
