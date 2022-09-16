using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Synapse.v2021_06_01.Workspaces;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum DesiredStateConstant
{
    [Description("Disabled")]
    Disabled,

    [Description("Enabled")]
    Enabled,
}
