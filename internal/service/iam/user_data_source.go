package iam

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"identity_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_login_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"otp_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions_boundary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	userName := d.Get("user_name").(string)
	user, err := FindUserByName(conn, userName)

	if err != nil {
		return fmt.Errorf("error reading IAM User: %w", err)
	}

	d.SetId(aws.StringValue(user.UserId))
	d.Set("arn", user.UserArn)
	d.Set("display_name", user.DisplayName)
	d.Set("email", user.Email)
	d.Set("enabled", user.Enabled)
	d.Set("identity_provider", user.IdentityProvider)

	if user.LastLoginDate != nil {
		d.Set("last_login_date", aws.TimeValue(user.LastLoginDate).Format(time.RFC3339))
	} else {
		d.Set("last_login_date", nil)
	}

	d.Set("login", user.Login)
	d.Set("otp_required", user.OtpRequired)
	d.Set("path", user.Path)

	d.Set("permissions_boundary", "")
	if user.PermissionsBoundary != nil {
		d.Set("permissions_boundary", user.PermissionsBoundary.PermissionsBoundaryArn)
	}

	d.Set("phone", user.Phone)

	if user.UpdateDate != nil {
		d.Set("update_date", aws.TimeValue(user.UpdateDate).Format(time.RFC3339))
	} else {
		d.Set("update_date", nil)
	}

	d.Set("user_name", user.UserName)
	d.Set("user_id", user.UserId)

	tags := KeyValueTags(user.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	// lintignore:AWSR002
	if err := d.Set("tags", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
