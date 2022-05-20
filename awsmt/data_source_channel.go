package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceChannel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelRead,
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"channel_name":  &requiredString,
			"channel_state": &computedString,
			"creation_time": &computedString,
			"filler_slate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_location_name": &computedString,
						"vod_source_name":      &computedString,
					},
				},
			},
			"last_modified_time": &computedString,
			"outputs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dash_manifest_windows_seconds":             &computedInt,
						"dash_min_buffer_time_seconds":              &computedInt,
						"dash_min_update_period_seconds":            &computedInt,
						"dash_suggested_presentation_delay_seconds": &computedInt,
						"hls_manifest_windows_seconds":              &computedInt,
						"manifest_name":                             &computedString,
						"playback_url":                              &computedString,
						"source_group":                              &computedString,
					},
				},
			},
			"playback_mode": &computedString,
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tier": &computedString,
		},
	}
}

func dataSourceChannelRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)

	name := d.Get("channel_name").(string)
	if name == "" {
		return diag.Errorf("`channel_name` parameter required")
	}
	res, err := client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: &name})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the channel: %w", err))
	}

	d.SetId(aws.StringValue(res.ChannelName))

	err = returnChannel(res, d)
	if err != nil {
		diag.FromErr(err)
	}

	return nil
}