package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePlaybackConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybackConfigurationCreate,
		UpdateContext: resourcePlaybackConfigurationUpdate,
		ReadContext:   resourcePlaybackConfigurationRead,
		DeleteContext: resourcePlaybackConfigurationDelete,
		// schema based on: https://docs.aws.amazon.com/mediatailor/latest/apireference/playbackconfiguration.html#playbackconfiguration-prop-putplaybackconfigurationrequest-personalizationthresholdseconds
		// and https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor#PutPlaybackConfigurationInput
		Schema: map[string]*schema.Schema{
			"ad_decision_server_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cdn_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_segment_url_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_segment_url_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"dash_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mpd_location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"origin_manifest_type": {
							//SINGLE_PERIOD | MULTI_PERIOD
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slate_ad_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"transcode_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"video_content_source_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePlaybackConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics

	input := getPlaybackConfigurationInput(d)

	_, err := client.PutPlaybackConfiguration(&input)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePlaybackConfigurationUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourcePlaybackConfigurationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourcePlaybackConfigurationDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func getPlaybackConfigurationInput(d *schema.ResourceData) mediatailor.PutPlaybackConfigurationInput {
	input := mediatailor.PutPlaybackConfigurationInput{}
	if v, ok := d.GetOk("ad_decision_server_url"); ok {
		val := v.(string)
		input.AdDecisionServerUrl = &val
	}
	if v, ok := d.GetOk("cdn_configuration"); ok {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.CdnConfiguration{}
		if str, ok := val["ad_segment_url_prefix"]; ok {
			converted := str.(string)
			output.AdSegmentUrlPrefix = &converted
		}
		if str, ok := val["content_segment_url_prefix"]; ok {
			converted := str.(string)
			output.ContentSegmentUrlPrefix = &converted
		}
		input.CdnConfiguration = &output
	}
	if v, ok := d.GetOk("dash_configuration"); ok {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.DashConfigurationForPut{}
		if str, ok := val["mpd_location"]; ok {
			converted := str.(string)
			output.MpdLocation = &converted
		}
		if str, ok := val["origin_manifest_type"]; ok {
			converted := str.(string)
			output.OriginManifestType = &converted
		}
		input.DashConfiguration = &output
	}
	if v, ok := d.GetOk("name"); ok {
		val := v.(string)
		input.Name = &val
	}
	if v, ok := d.GetOk("slate_ad_url"); ok {
		val := v.(string)
		input.SlateAdUrl = &val
	}
	if v, ok := d.GetOk("tags"); ok {
		outputMap := make(map[string]*string)
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
		input.Tags = outputMap
	}
	if v, ok := d.GetOk("transcode_profile_name"); ok {
		val := v.(string)
		input.TranscodeProfileName = &val
	}
	if v, ok := d.GetOk("video_content_source_url"); ok {
		val := v.(string)
		input.VideoContentSourceUrl = &val
	}
	return input
}
