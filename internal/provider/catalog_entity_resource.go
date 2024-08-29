package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogEntityResource{}
var _ resource.ResourceWithImportState = &CatalogEntityResource{}

func NewCatalogEntityResource() resource.Resource {
	return &CatalogEntityResource{}
}

func NewCatalogEntityResourceModel() CatalogEntityResourceModel {
	return CatalogEntityResourceModel{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// CatalogEntityResource defines the resource implementation.
type CatalogEntityResource struct {
	client *cortex.HttpClient
}

func (r *CatalogEntityResource) toUpsertRequest(ctx context.Context, diagnostics *diag.Diagnostics, data *CatalogEntityResourceModel) cortex.UpsertCatalogEntityRequest {
	return cortex.UpsertCatalogEntityRequest{
		Info: data.ToApiModel(ctx, diagnostics),
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *CatalogEntityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity"
}

func (r *CatalogEntityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Catalog Entity",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name for the entity",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog. Markdown is supported.",
				Optional:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the entity. Corresponds to the x-cortex-tag field in the entity descriptor.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Set when the entity is a Resource or Team. This must match a tag of a valid Resource Definition or be \"team\" or \"domain\".",
				Optional:            true,
			},
			"definition": schema.StringAttribute{
				MarkdownDescription: "Set when the entity is a Resource. These are the properties defined by the Resource Definition, in JSON format in a string (use the `jsonencode` function to convert a JSON object to a string).",
				Optional:            true,
			},

			// Optional attributes
			"owners": schema.ListNestedAttribute{
				MarkdownDescription: "List of owners for the entity. Owners can be users, groups, or Slack channels.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of owner. Valid values are `EMAIL`, `GROUP`, or `SLACK`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("EMAIL", "GROUP", "SLACK"),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the owner. Only required for `user` or `group` types.",
							Optional:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email of the owner. Only allowed if `type` is `user`.",
							Optional:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the owner. Optional.",
							Optional:            true,
						},
						"provider": schema.StringAttribute{
							MarkdownDescription: "Provider of the owner. Only allowed if `type` is `group`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ACTIVE_DIRECTORY", "BAMBOO_HR", "CORTEX", "GITHUB", "GITLAB", "GOOGLE", "OKTA", "OPSGENIE", "SERVICE_NOW", "WORKDAY"),
							},
						},
						"channel": schema.StringAttribute{
							MarkdownDescription: "Channel of the owner. Only allowed if `type` is `slack`. Omit the #.",
							Optional:            true,
						},
						"notifications_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether Slack notifications are enabled for all owners of this service. Only allowed if `type` is `slack`.",
							Optional:            true,
						},
					},
				},
			},
			"children": schema.ListNestedAttribute{
				MarkdownDescription: "List of child entities for the entity. Only used for entities of type `TEAM` or `DOMAIN`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the child entity.",
							Required:            true,
						},
					},
				},
			},
			"parents": schema.ListNestedAttribute{
				MarkdownDescription: "List of parents for the entity. The list of parents can only be entities of type `DOMAIN`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the parent domain.",
							Required:            true,
						},
					},
				},
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "List of groups related to the entity.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"links": schema.ListNestedAttribute{
				MarkdownDescription: "List of links related to the entity.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the link.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the link. Valid values are `runbook`, `documentation`, `logs`, `dashboard`, `metrics`, `healthcheck`, `OPENAPI`, `ASYNC_API`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("runbook", "documentation", "logs", "dashboard", "metrics", "healthcheck", "OPENAPI", "ASYNC_API"),
							},
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL of the link.",
							Required:            true,
						},
					},
				},
			},
			"ignore_metadata": schema.BoolAttribute{
				MarkdownDescription: "Whether the entity's custom metadata is managed by Terraform. Defaults to `false`. If set to `true`, the provider will ignore any metadata on the Entity and not persist it to state.",
				Optional:            true,
			},
			"metadata": schema.StringAttribute{
				MarkdownDescription: "Custom metadata for the entity, in JSON format in a string. (Use the `jsonencode` function to convert a JSON object to a string.)",
				Optional:            true,
			},
			"dependencies": schema.ListNestedAttribute{
				MarkdownDescription: "List of dependencies for the entity.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the dependency.",
							Required:            true,
						},
						"method": schema.StringAttribute{
							MarkdownDescription: "HTTP method if depending on a specific endpoint.",
							Optional:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "The actual endpoint this dependency refers to.",
							Optional:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the dependency.",
							Optional:            true,
						},
						"metadata": schema.StringAttribute{
							MarkdownDescription: "Custom metadata for the dependency, in JSON format in a string. (Use the `jsonencode` function to convert a JSON object to a string.)",
							Optional:            true,
						},
					},
				},
			},
			"alerts": schema.ListNestedAttribute{
				MarkdownDescription: "List of alerts for the entity.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of alert. Valid values are `opsgenie`",
							Required:            true,
						},
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the alert.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the alert.",
							Optional:            true,
						},
					},
				},
			},
			"git": schema.SingleNestedAttribute{
				MarkdownDescription: "Git configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"github": schema.SingleNestedAttribute{
						MarkdownDescription: "GitHub configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"repository": schema.StringAttribute{
								MarkdownDescription: "GitHub repository for the entity.",
								Required:            true,
							},
							"base_path": schema.StringAttribute{
								MarkdownDescription: "Base path if not /",
								Optional:            true,
							},
						},
					},
					"gitlab": schema.SingleNestedAttribute{
						MarkdownDescription: "GitLab configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"repository": schema.StringAttribute{
								MarkdownDescription: "GitLab repository for the entity.",
								Required:            true,
							},
							"base_path": schema.StringAttribute{
								MarkdownDescription: "Base path if not /",
								Optional:            true,
							},
						},
					},
					"azure": schema.SingleNestedAttribute{
						MarkdownDescription: "Azure configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"project": schema.StringAttribute{
								MarkdownDescription: "Azure project for the entity.",
								Required:            true,
							},
							"repository": schema.StringAttribute{
								MarkdownDescription: "Azure repository for the entity.",
								Required:            true,
							},
							"base_path": schema.StringAttribute{
								MarkdownDescription: "Base path if not /",
								Optional:            true,
							},
						},
					},
					"bitbucket": schema.SingleNestedAttribute{
						MarkdownDescription: "BitBucket configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"repository": schema.StringAttribute{
								MarkdownDescription: "BitBucket repository for the entity.",
								Required:            true,
							},
						},
					},
				},
			},
			"issues": schema.SingleNestedAttribute{
				MarkdownDescription: "Issue tracking configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"jira": schema.SingleNestedAttribute{
						MarkdownDescription: "Jira configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"default_jql": schema.StringAttribute{
								MarkdownDescription: "Default JQL to surface issues for the entity.",
								Optional:            true,
							},
							"projects": schema.ListAttribute{
								MarkdownDescription: "List of Jira projects for the entity.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"components": schema.ListAttribute{
								MarkdownDescription: "List of Jira components for the entity.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"labels": schema.ListAttribute{
								MarkdownDescription: "List of Jira labels for the entity.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
				},
			},
			"on_call": schema.SingleNestedAttribute{
				MarkdownDescription: "On-call configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"pager_duty": schema.SingleNestedAttribute{
						MarkdownDescription: "PagerDuty configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "PagerDuty Service, Schedule, or Escalation Policy ID.",
								Required:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "Type. Valid values are `SERVICE`, `SCHEDULE`, or `ESCALATION_POLICY`.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("SERVICE", "SCHEDULE", "ESCALATION_POLICY"),
								},
							},
						},
					},
					"ops_genie": schema.SingleNestedAttribute{
						MarkdownDescription: "OpsGenie configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "OpsGenie Schedule ID.",
								Required:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "Type. Valid values are `SCHEDULE`.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("SCHEDULE"),
								},
							},
						},
					},
					"victor_ops": schema.SingleNestedAttribute{
						MarkdownDescription: "VictorOps configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "VictorOps Schedule ID.",
								Required:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "Type. Valid values are `SCHEDULE`.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("SCHEDULE"),
								},
							},
						},
					},
					"xmatters": schema.SingleNestedAttribute{
						MarkdownDescription: "XMatters configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "XMatters Group Name.",
								Required:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "Type. Valid values are `SERVICE`.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("SERVICE"),
								},
							},
						},
					},
				},
			},
			"apm": schema.SingleNestedAttribute{
				MarkdownDescription: "APM configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"data_dog": schema.SingleNestedAttribute{
						MarkdownDescription: "DataDog configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"monitors": schema.ListAttribute{
								MarkdownDescription: "List of DataDog monitors for the entity.",
								Optional:            true,
								ElementType:         types.Int64Type,
							},
						},
					},
					"dynatrace": schema.SingleNestedAttribute{
						MarkdownDescription: "Dynatrace configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"entity_ids": schema.ListAttribute{
								MarkdownDescription: "List of Dynatrace entity IDs for the entity.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"entity_name_matchers": schema.ListAttribute{
								MarkdownDescription: "List of Dynatrace entity name matchers for the entity.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"new_relic": schema.ListNestedAttribute{
						MarkdownDescription: "NewRelic configurations for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"application_id": schema.Int64Attribute{
									MarkdownDescription: "NewRelic application ID for the entity.",
									Optional:            true,
								},
								"alias": schema.StringAttribute{
									MarkdownDescription: "Alias for the service. Only used if opted into multi-account support in New Relic.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"dashboards": schema.SingleNestedAttribute{
				MarkdownDescription: "Dashboards configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"embeds": schema.ListNestedAttribute{
						MarkdownDescription: "List of dashboard embeds for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type of embed. Valid values are `grafana`, `datadog` or `newrelic`.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("grafana", "datadog", "newrelic"),
									},
								},
								"url": schema.StringAttribute{
									MarkdownDescription: "URL of the embed.",
									Required:            true,
								},
							},
						},
					},
				},
			},
			"slos": schema.SingleNestedAttribute{
				MarkdownDescription: "Service-level Objectives configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"data_dog": schema.ListNestedAttribute{
						MarkdownDescription: "DataDog SLO configuration for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "DataDog SLO ID.",
									Required:            true,
								},
							},
						},
					},
					"dynatrace": schema.ListNestedAttribute{
						MarkdownDescription: "Dynatrace SLO configuration for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Dynatrace SLO ID.",
									Required:            true,
								},
							},
						},
					},
					"lightstep": schema.SingleNestedAttribute{
						MarkdownDescription: "Lightstep SLO configuration for the entity.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"streams": schema.ListNestedAttribute{
								MarkdownDescription: "List of Lightstep streams for the entity.",
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"stream_id": schema.StringAttribute{
											MarkdownDescription: "Lightstep stream ID.",
											Required:            true,
										},
										"targets": schema.SingleNestedAttribute{
											MarkdownDescription: "List of target latencies and error rates for the stream.",
											Optional:            true,
											Attributes: map[string]schema.Attribute{
												"latencies": schema.ListNestedAttribute{
													MarkdownDescription: "List of latency targets for the stream.",
													Optional:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"percentile": schema.Float64Attribute{
																MarkdownDescription: "Percentile latency for your given streamId, out of 1.",
																Required:            true,
															},
															"target": schema.Int64Attribute{
																MarkdownDescription: "Target latency, in ms.",
																Required:            true,
															},
															"slo": schema.Float64Attribute{
																MarkdownDescription: "SLO percentile, out of 1.",
																Required:            true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"prometheus": schema.ListNestedAttribute{
						MarkdownDescription: "Prometheus SLO configuration for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"error_query": schema.StringAttribute{
									MarkdownDescription: "Query that indicates error events for your metric.",
									Required:            true,
								},
								"total_query": schema.StringAttribute{
									MarkdownDescription: "Query that indicates all events to be considered for your metric.",
									Required:            true,
								},
								"slo": schema.Float64Attribute{
									MarkdownDescription: "Target number for SLO.",
									Required:            true,
								},
								"alias": schema.StringAttribute{
									MarkdownDescription: "Ties the SLO registration to a Prometheus instance listed under Settings → Prometheus. The alias parameter is optional, but if not provided the SLO will use the default configuration under Settings → Prometheus.",
									Optional:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "The SLO's name in Prometheus. The name parameter is optional.",
									Optional:            true,
								},
							},
						},
					},
					"signal_fx": schema.ListNestedAttribute{
						MarkdownDescription: "SignalFx SLO configuration for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"query": schema.StringAttribute{
									MarkdownDescription: "Elasticsearch query for your metric. Filter by metric with `sf_metric` and add additional dimensions to narrow the search. Queries resulting in multiple datasets will be rolled up according to `rollup`.",
									Required:            true,
								},
								"rollup": schema.StringAttribute{
									MarkdownDescription: "SUM or AVERAGE.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("SUM", "AVERAGE"),
									},
								},
								"target": schema.Int64Attribute{
									MarkdownDescription: "Target number for SLO.",
									Required:            true,
								},
								"lookback": schema.StringAttribute{
									MarkdownDescription: "ISO-8601 duration `(P[n]Y[n]M[n]DT[n]H[n]M[n]S)`.",
									Required:            true,
								},
								"operation": schema.StringAttribute{
									MarkdownDescription: "> / < / = / <=, >=.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(">", "<", "=", "<=", ">="),
									},
								},
							},
						},
					},
					"sumo_logic": schema.ListNestedAttribute{
						MarkdownDescription: "SumoLogic SLO configuration for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "SumoLogic SLO ID.",
									Required:            true,
								},
							},
						},
					},
				},
			},
			"static_analysis": schema.SingleNestedAttribute{
				MarkdownDescription: "Static analysis configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"code_cov": schema.SingleNestedAttribute{
						MarkdownDescription: "Code coverage configuration.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"repository": schema.StringAttribute{
								MarkdownDescription: "Git repository, in `organization/repository` format.",
								Required:            true,
							},
							"provider": schema.StringAttribute{
								MarkdownDescription: "Git provider. One of: `GITHUB`, `GITLAB`, or `BITBUCKET`.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("GITHUB", "GITLAB", "BITBUCKET"),
								},
							},
							"owner": schema.StringAttribute{
								MarkdownDescription: "Username from service provider (optional).",
								Optional:            true,
							},
							"flag": schema.StringAttribute{
								MarkdownDescription: "Pulls from isolated & categorized coverage reports (optional).",
								Optional:            true,
							},
						},
					},
					"mend": schema.SingleNestedAttribute{
						MarkdownDescription: "Mend static analysis configuration.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"application_ids": schema.ListAttribute{
								MarkdownDescription: "Mend application IDs, found in the Mend SAST web interface.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"project_ids": schema.ListAttribute{
								MarkdownDescription: "Mend project IDs, found in the Mend SCA web interface.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"sonar_qube": schema.SingleNestedAttribute{
						MarkdownDescription: "SonarQube static analysis configuration.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"project": schema.StringAttribute{
								MarkdownDescription: "SonarQube project key.",
								Required:            true,
							},
							"alias": schema.StringAttribute{
								MarkdownDescription: "Ties the SonarQube registration to a SonarQube instance listed under Settings → SonarQube. The alias parameter is optional, but if not provided the SonarQube registration will use the default configuration under Settings → SonarQube.",
								Optional:            true,
							},
						},
					},
					"veracode": schema.SingleNestedAttribute{
						MarkdownDescription: "Veracode static analysis configuration.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"application_names": schema.ListAttribute{
								MarkdownDescription: "Veracode application names.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"sandboxes": schema.ListNestedAttribute{
								MarkdownDescription: "Veracode sandboxes.",
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"application_name": schema.StringAttribute{
											MarkdownDescription: "Veracode application name.",
											Required:            true,
										},
										"sandbox_name": schema.StringAttribute{
											MarkdownDescription: "Veracode sandbox name.",
											Required:            true,
										},
									},
								},
							},
						},
					},
				},
			},
			"ci_cd": schema.SingleNestedAttribute{
				MarkdownDescription: "CI/CD configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"buildkite": schema.SingleNestedAttribute{
						MarkdownDescription: "Buildkite CI/CD configuration.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"pipelines": schema.SetNestedAttribute{
								MarkdownDescription: "Buildkite pipelines.",
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"slug": schema.StringAttribute{
											MarkdownDescription: "Buildkite pipeline slug.",
											Required:            true,
										},
									},
								},
							},
							"tags": schema.SetNestedAttribute{
								MarkdownDescription: "Buildkite pipelines.",
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"tag": schema.StringAttribute{
											MarkdownDescription: "Buildkite tag.",
											Required:            true,
										},
									},
								},
							},
						},
					},
				},
			},
			"bug_snag": schema.SingleNestedAttribute{
				MarkdownDescription: "BugSnag configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"project": schema.StringAttribute{
						MarkdownDescription: "BugSnag project ID for the entity.",
						Required:            true,
					},
				},
			},
			"checkmarx": schema.SingleNestedAttribute{
				MarkdownDescription: "Checkmarx configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"projects": schema.ListNestedAttribute{
						MarkdownDescription: "List of Checkmarx projects for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									MarkdownDescription: "Checkmarx project ID. Required if Name is not set.",
									Optional:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Checkmarx project name. Required if ID is not set.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"coralogix": schema.SingleNestedAttribute{
				MarkdownDescription: "Coralogix configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"applications": schema.ListNestedAttribute{
						MarkdownDescription: "List of Coralogix applications for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Coralogix application Name.",
									Required:            true,
								},
								"alias": schema.StringAttribute{
									MarkdownDescription: "Alias is optional and only relevant if you have opted into multi account support.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"firehydrant": schema.SingleNestedAttribute{
				MarkdownDescription: "FireHydrant configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"services": schema.ListNestedAttribute{
						MarkdownDescription: "List of FireHydrant services for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "FireHydrant service identifier. Either the ID or the Slug, based on the value of Type below.",
									Required:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "FireHydrant service identifier type. Either `ID` or `SLUG`.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("ID", "SLUG"),
									},
								},
							},
						},
					},
				},
			},
			"k8s": schema.SingleNestedAttribute{
				MarkdownDescription: "Kubernetes configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"deployments": schema.ListNestedAttribute{
						MarkdownDescription: "List of K8s deployment configurations for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									MarkdownDescription: "Deployment identifier. `namespace/name` as found in Kubernetes.",
									Required:            true,
								},
								"cluster": schema.StringAttribute{
									MarkdownDescription: "Optional. Kubernetes cluster name.",
									Optional:            true,
								},
							},
						},
					},
					"argo_rollouts": schema.ListNestedAttribute{
						MarkdownDescription: "List of K8s argo rollout configurations for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									MarkdownDescription: "Argo Rollout identifier. `namespace/name` as found in Kubernetes.",
									Required:            true,
								},
								"cluster": schema.StringAttribute{
									MarkdownDescription: "Optional. Kubernetes cluster name.",
									Optional:            true,
								},
							},
						},
					},
					"stateful_sets": schema.ListNestedAttribute{
						MarkdownDescription: "List of K8s stateful set configurations for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									MarkdownDescription: "Stateful Set identifier. `namespace/name` as found in Kubernetes.",
									Required:            true,
								},
								"cluster": schema.StringAttribute{
									MarkdownDescription: "Optional. Kubernetes cluster name.",
									Optional:            true,
								},
							},
						},
					},
					"cron_jobs": schema.ListNestedAttribute{
						MarkdownDescription: "List of K8s cron job configurations for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									MarkdownDescription: "Cron Job identifier. `namespace/name` as found in Kubernetes.",
									Required:            true,
								},
								"cluster": schema.StringAttribute{
									MarkdownDescription: "Optional. Kubernetes cluster name.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"launch_darkly": schema.SingleNestedAttribute{
				MarkdownDescription: "LaunchDarkly configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"projects": schema.ListNestedAttribute{
						MarkdownDescription: "List of LaunchDarkly projects for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "LaunchDarkly project ID.",
									Required:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Type of project identifier to use to map to the LaunchDarkly project. Valid values are `ID` or `KEY`.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("ID", "KEY"),
									},
								},
								"environments": schema.ListNestedAttribute{
									MarkdownDescription: "List of LaunchDarkly environments for the project.",
									Optional:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "LaunchDarkly environment name.",
												Required:            true,
											},
										},
									},
								},
								"alias": schema.StringAttribute{
									MarkdownDescription: "Alias is optional and only relevant if you have opted into multi-account support.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"microsoft_teams": schema.ListNestedAttribute{
				MarkdownDescription: "List of Microsoft Teams channels for the entity.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Exact match name of the channel.",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the channel.",
							Optional:            true,
						},
						"notifications_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the channel should receive notifications.",
							Optional:            true,
						},
					},
				},
			},
			"rollbar": schema.SingleNestedAttribute{
				MarkdownDescription: "Rollbar configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"project": schema.StringAttribute{
						MarkdownDescription: "Rollbar project ID for the entity.",
						Required:            true,
					},
				},
			},
			"sentry": schema.SingleNestedAttribute{
				MarkdownDescription: "Sentry configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"project": schema.StringAttribute{
						MarkdownDescription: "Sentry project ID for the entity.",
						Required:            true,
					},
				},
			},
			"service_now": schema.SingleNestedAttribute{
				MarkdownDescription: "ServiceNow configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"services": schema.ListNestedAttribute{
						MarkdownDescription: "List of ServiceNow services for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									MarkdownDescription: "ServiceNow service ID.",
									Optional:            true,
								},
								"table_name": schema.StringAttribute{
									MarkdownDescription: "ServiceNow table name.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"slack": schema.SingleNestedAttribute{
				MarkdownDescription: "Slack configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"channels": schema.SetNestedAttribute{
						MarkdownDescription: "List of Slack channels for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Slack channel name.",
									Required:            true,
								},
								"notifications_enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether the slack channel should receive notifications.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"snyk": schema.SingleNestedAttribute{
				MarkdownDescription: "Snyk configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"projects": schema.ListNestedAttribute{
						MarkdownDescription: "List of Snyk projects for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"organization": schema.StringAttribute{
									MarkdownDescription: "Snyk organization ID.",
									Required:            true,
								},
								"project_id": schema.StringAttribute{
									MarkdownDescription: "Snyk project ID.",
									Required:            true,
								},
								"source": schema.StringAttribute{
									MarkdownDescription: "Type of Snyk product. Valid values are `CODE` or `OPEN_SOURCE`.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("CODE", "OPEN_SOURCE"),
									},
								},
							},
						},
					},
				},
			},
			"wiz": schema.SingleNestedAttribute{
				MarkdownDescription: "Wiz configuration for the entity.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"projects": schema.ListNestedAttribute{
						MarkdownDescription: "List of Wiz projects for the entity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"project_id": schema.StringAttribute{
									MarkdownDescription: "Wiz project ID.",
									Required:            true,
								},
							},
						},
					},
				},
			},
			// Team attributes
			"team": schema.SingleNestedAttribute{
				MarkdownDescription: "Team configuration for the entity. Only used for entities of type `TEAM`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"members": schema.ListNestedAttribute{
						MarkdownDescription: "List of additional team members for the team not derived from any groups.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"email": schema.StringAttribute{
									MarkdownDescription: "Email address of the team member.",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name of the team member.",
									Required:            true,
								},
								"role": schema.StringAttribute{
									MarkdownDescription: "Role of the team member. Optional.",
									Optional:            true,
								},
								"notifications_enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether the team member should receive notifications.",
									Optional:            true,
								},
							},
						},
					},
					"groups": schema.ListNestedAttribute{
						MarkdownDescription: "List of groups to derive team members from.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Name of the group.",
									Required:            true,
								},
								"provider": schema.StringAttribute{
									MarkdownDescription: "Provider of the group.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("ACTIVE_DIRECTORY", "BAMBOO_HR", "CORTEX", "GITHUB", "GOOGLE", "OKTA", "OPSGENIE", "WORKDAY"),
									},
								},
							},
						},
					},
				},
			},

			//Computed
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CatalogEntityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *CatalogEntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewCatalogEntityResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	oldMetadata := data.Metadata

	// Parse configuration into an upsert entity
	upsertRequest := r.toUpsertRequest(ctx, &resp.Diagnostics, &data)
	if resp.Diagnostics.HasError() {
		// If we have an issue parsing the TF configuration, bail out early
		return
	}

	// Issue API request
	entity, err := r.client.CatalogEntities().Upsert(ctx, upsertRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entity, got error: %s", err))
		return
	}

	// Set computed attributes
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if data.IgnoreMetadata.ValueBool() {
		data.Metadata = oldMetadata
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewCatalogEntityResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	oldMetadata := data.Metadata

	// Issue API request
	entity, err := r.client.CatalogEntities().GetFromDescriptor(ctx, data.Tag.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// Set attributes from API response
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if data.IgnoreMetadata.ValueBool() {
		data.Metadata = oldMetadata
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewCatalogEntityResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldMetadata := data.Metadata

	// Parse configuration into API entity
	upsertRequest := r.toUpsertRequest(ctx, &resp.Diagnostics, &data)
	if resp.Diagnostics.HasError() {
		// If we have an issue parsing the TF configuration, bail out early
		return
	}

	// Issue API request to Cortex
	entity, err := r.client.CatalogEntities().Upsert(ctx, upsertRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entity, got error: %s", err))
		return
	}

	// Set computed attributes
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if data.IgnoreMetadata.ValueBool() {
		data.Metadata = oldMetadata
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewCatalogEntityResourceModel()

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CatalogEntities().Delete(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete catalog entity, got error: %s", err))
		return
	}
}

func (r *CatalogEntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
