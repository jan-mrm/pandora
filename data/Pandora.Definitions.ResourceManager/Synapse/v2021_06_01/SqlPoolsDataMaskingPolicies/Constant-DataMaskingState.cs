using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Synapse.v2021_06_01.SqlPoolsDataMaskingPolicies;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum DataMaskingStateConstant
{
    [Description("Disabled")]
    Disabled,

    [Description("Enabled")]
    Enabled,
}
