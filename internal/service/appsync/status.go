package appsync

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

func StatusApiCache(conn *appsync.AppSync, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindApiCacheByID(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusDomainNameApiAssociation(conn *appsync.AppSync, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindDomainNameApiAssociationByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.AssociationStatus), nil
	}
}
