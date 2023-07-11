package utilities

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &moduleTemplateDataSource{}
)

func NewModuleTemplateDataSource() datasource.DataSource {
	return &moduleTemplateDataSource{}
}

type moduleTemplateDataSource struct{}

type moduleTemplateDataSourceModel struct {
	ModuleInfo types.Map `tfsdk:"module_info"`
	ModuleTmpl types.Map `tfsdk:"module_tmpl"`
}

func (d *moduleTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_tmpl"
}

func (d *moduleTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"module_info": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"module_tmpl": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (d *moduleTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var plan, state moduleTemplateDataSourceModel

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state = plan

	moduleInfoInput := make(map[string]string)
	convertModuleInfoDiags := plan.ModuleInfo.ElementsAs(ctx, &moduleInfoInput, false)
	resp.Diagnostics.Append(convertModuleInfoDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	moduleTmplInput := make(map[string]string)
	convertModuleTmplDiags := plan.ModuleTmpl.ElementsAs(ctx, &moduleTmplInput, false)
	resp.Diagnostics.Append(convertModuleTmplDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	re := regexp.MustCompile(`{(.*?)}`)
	moduleTmplResult := make(map[string]attr.Value)
	for tmplKey, tmplValue := range moduleTmplInput {
		for infoKey, infoValue := range moduleInfoInput {
			tmplValue = strings.Replace(tmplValue, fmt.Sprintf("{%s}", infoKey), infoValue, 1)
		}
		containIllegal := re.FindAllStringSubmatch(tmplValue, -1)
		if len(containIllegal) > 0 {
			for _, x := range containIllegal {
				resp.Diagnostics.AddError(
					"[Error] Template contains undefine key from module_info key",
					fmt.Sprintf("Unknown key: %s", x[1]),
				)
			}
			break
		}
		moduleTmplResult[tmplKey] = types.StringValue(tmplValue)
	}

	state.ModuleTmpl = types.MapValueMust(types.StringType, moduleTmplResult)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
