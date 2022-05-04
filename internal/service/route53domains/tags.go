// Custom Route 53 Domains tag service update functions using the same format as generated code.
// Modified to support AWS Go SDK v2.

package route53domains

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/aws/aws-sdk-go-v2/service/route53domains/types"
	tftags "github.com/PixarV/terraform-provider-ritt/internal/tags"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

// GetTag fetches an individual route53domains service tag for a resource.
// Returns whether the key value and any errors. A NotFoundError is used to signal that no value was found.
// This function will optimise the handling over ListTags, if possible.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func GetTag(ctx context.Context, conn *route53domains.Client, identifier string, key string) (*string, error) {
	listTags, err := ListTags(ctx, conn, identifier)

	if err != nil {
		return nil, err
	}

	if !listTags.KeyExists(key) {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	return listTags.KeyValue(key), nil
}

// ListTags lists route53domains service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(ctx context.Context, conn *route53domains.Client, identifier string) (tftags.KeyValueTags, error) {
	input := &route53domains.ListTagsForDomainInput{
		DomainName: aws.String(identifier),
	}

	output, err := conn.ListTagsForDomain(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.TagList), nil
}

// []*SERVICE.Tag handling

// Tags returns route53domains service tags.
func Tags(tags tftags.KeyValueTags) []types.Tag {
	result := make([]types.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from route53domains service tags.
func KeyValueTags(tags []types.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.ToString(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates route53domains service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(ctx context.Context, conn *route53domains.Client, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &route53domains.DeleteTagsForDomainInput{
			DomainName:   aws.String(identifier),
			TagsToDelete: removedTags.IgnoreAWS().Keys(),
		}

		_, err := conn.DeleteTagsForDomain(ctx, input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &route53domains.UpdateTagsForDomainInput{
			DomainName:   aws.String(identifier),
			TagsToUpdate: Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.UpdateTagsForDomain(ctx, input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
