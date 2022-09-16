using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Synapse.v2021_06_01.WorkspaceAzureADOnlyAuthentications;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum SensitivityLabelSourceConstant
{
    [Description("current")]
    Current,

    [Description("recommended")]
    Recommended,
}
