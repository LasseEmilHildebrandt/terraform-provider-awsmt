package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

func init() {
	resource.AddTestSweepers("test_playback_configuration", &resource.Sweeper{
		Name: "test_playback_configuration",
		F: func(region string) error {
			client, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}
			conn := client.(*mediatailor.MediaTailor)
			name := "test-playback-configuration-awsmt"
			_, err = conn.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
			if err != nil {
				return err
			}
			return nil
		},
	})
}

func TestAccPlaybackConfigurationResourceBasic(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r1"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckPlaybackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					resource.TestMatchResourceAttr(resourceName, "playback_configuration_arn", regexp.MustCompile(`arn:aws:mediatailor`)),
				),
			},
			{
				Config: testAccPlaybackConfigurationUpdateResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr(resourceName, "slate_ad_url", "https://exampleurl.com/updated"),
				),
			},
			{
				Config: testAccPlaybackConfigurationUpdateResourceName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt-changed"),
					resource.TestCheckResourceAttr(resourceName, "slate_ad_url", "https://exampleurl.com/updated"),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationResourceImport(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r2"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { importPreCheck(t, "eu-central-1") },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckPlaybackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationImportResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					resource.TestMatchResourceAttr(resourceName, "playback_configuration_arn", regexp.MustCompile(`arn:aws:mediatailor`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPlaybackConfigurationResourceTaint(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r2"
	firstARN := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckPlaybackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationImportResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					resource.TestMatchResourceAttr(resourceName, "playback_configuration_arn", regexp.MustCompile(`arn:aws:mediatailor`)),
					testAccAssignARN(resourceName, &firstARN, true),
				),
			},
			{
				Config: testAccPlaybackConfigurationImportResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					resource.TestMatchResourceAttr(resourceName, "playback_configuration_arn", regexp.MustCompile(`arn:aws:mediatailor`)),
					testAccCheckARN(resourceName, &firstARN, true),
				),
				Taint: []string{resourceName},
			},
		},
	})
}

func TestAccPlaybackConfigurationRemoveResourceTag(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r1"
	firstARN := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckPlaybackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					testAccAssignARN(resourceName, &firstARN, false),
				),
			},
			{
				Config: testAccPlaybackConfigurationRemoveResourceTag(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
					testAccCheckARN(resourceName, &firstARN, false),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAssignARN(resourceName string, ARNVariable *string, taintArn bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if taintArn {
			rs.Primary.Attributes["playback_configuration_arn"] = "taintedARN"
		}
		*ARNVariable = rs.Primary.Attributes["playback_configuration_arn"]
		return nil
	}
}

func testAccCheckARN(resourceName string, ARNVariable *string, expectDifferent bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if expectDifferent && *ARNVariable == rs.Primary.Attributes["playback_configuration_arn"] {
			return fmt.Errorf("same ARN (%s), resource not recreated", *ARNVariable)
		} else if !expectDifferent && *ARNVariable != rs.Primary.Attributes["playback_configuration_arn"] {
			return fmt.Errorf("different ARN (%s) != (%s), resource recreated", *ARNVariable, rs.Primary.Attributes["playback_configuration_arn"])
		}
		return nil
	}
}

func importPreCheck(t *testing.T, region string) {
	testAccPreCheck(t)

	client, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("error getting client: %s", err)
	}
	conn := client.(*mediatailor.MediaTailor)

	exampleUrl := "https://exampleurl.com/"
	mpdLocation := "DISABLED"
	manifestType := "SINGLE_PERIOD"
	name := "test-playback-configuration-awsmt"
	env := "dev"
	input := mediatailor.PutPlaybackConfigurationInput{
		AdDecisionServerUrl:   &exampleUrl,
		CdnConfiguration:      &mediatailor.CdnConfiguration{AdSegmentUrlPrefix: &exampleUrl},
		DashConfiguration:     &mediatailor.DashConfigurationForPut{MpdLocation: &mpdLocation, OriginManifestType: &manifestType},
		Name:                  &name,
		VideoContentSourceUrl: &exampleUrl,
		Tags:                  map[string]*string{"Environment": &env},
	}

	_, err = conn.PutPlaybackConfiguration(&input)
	if err != nil {
		t.Fatal(diag.FromErr(err))
	}
}

func testAccCheckPlaybackConfigurationDestroy(_ *terraform.State) error {
	c := testAccProvider.Meta().(*mediatailor.MediaTailor)
	name := "test-playback-configuration-awsmt"
	_, err := c.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
	if err != nil {
		return err
	}
	return nil
}

func testAccPlaybackConfigurationResource() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression {
   mode = "OFF"
  }
  bumper {}
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  live_pre_roll_configuration {
	max_duration_seconds = 1
  }
  manifest_processing_rules {
	ad_marker_passthrough {
	  enabled = true
	}
  }
  name = "test-playback-configuration-awsmt"
  personalization_threshold_seconds = 2
  slate_ad_url = "https://exampleurl.com/"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}

`
}

func testAccPlaybackConfigurationRemoveResourceTag() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression {
   mode = "OFF"
  }
  bumper {}
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  live_pre_roll_configuration {
	max_duration_seconds = 1
  }
  manifest_processing_rules {
	ad_marker_passthrough {
	  enabled = true
	}
  }
  name = "test-playback-configuration-awsmt"
  personalization_threshold_seconds = 2
  slate_ad_url = "https://exampleurl.com/"
  video_content_source_url = "https://exampleurl.com/"
}
`
}

func testAccPlaybackConfigurationUpdateResource() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression {
   mode = "OFF"
  }
  bumper {}
  cdn_configuration {
    ad_segment_url_prefix = "test-updated"
    content_segment_url_prefix = "test-updated"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  live_pre_roll_configuration {
	max_duration_seconds = 1
  }
  manifest_processing_rules {
	ad_marker_passthrough {
	  enabled = true
	}
  }
  name = "test-playback-configuration-awsmt"
  personalization_threshold_seconds = 2
  slate_ad_url = "https://exampleurl.com/updated"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/updated"
}

`
}
func testAccPlaybackConfigurationUpdateResourceName() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression {
   mode = "OFF"
  }
  bumper {}
  cdn_configuration {
    ad_segment_url_prefix = "test-updated"
    content_segment_url_prefix = "test-updated"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  live_pre_roll_configuration {
	max_duration_seconds = 1
  }
  manifest_processing_rules {
	ad_marker_passthrough {
	  enabled = true
	}
  }
  name = "test-playback-configuration-awsmt-changed"
  personalization_threshold_seconds = 2
  slate_ad_url = "https://exampleurl.com/updated"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/updated"
}

`
}
func testAccPlaybackConfigurationImportResource() string {
	return `
resource "awsmt_playback_configuration" "r2" {
  ad_decision_server_url = "https://exampleurl.com/"
  bumper {
	end_url = "https://wxample.com/endbumper"
    start_url = "https://wxample.com/startbumper"
  }
  cdn_configuration {
    ad_segment_url_prefix = "https://exampleurl.com/"
  }
  dash_configuration {
    mpd_location = "DISABLED"
	origin_manifest_type = "SINGLE_PERIOD"
  }
  live_pre_roll_configuration {
	max_duration_seconds = 1
  }
  manifest_processing_rules {
	ad_marker_passthrough {
	  enabled = true
	}
  }
  name = "test-playback-configuration-awsmt"
  personalization_threshold_seconds = 2
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}
`
}
