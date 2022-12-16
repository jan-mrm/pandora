using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.StreamAnalytics.v2021_10_01_preview.Outputs;

[ValueForType("Microsoft.EventHub/EventHub")]
internal class EventHubV2OutputDataSourceModel : OutputDataSourceModel
{
    [JsonPropertyName("properties")]
    public EventHubOutputDataSourcePropertiesModel? Properties { get; set; }
}
