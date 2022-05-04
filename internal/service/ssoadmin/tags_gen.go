// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ssoadmin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	tftags "github.com/PixarV/terraform-provider-ritt/internal/tags"
)

// ListTags lists ssoadmin service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn *ssoadmin.SSOAdmin, identifier string, resourceType string) (tftags.KeyValueTags, error) {
	input := &ssoadmin.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
		InstanceArn: aws.String(resourceType),
	}

	output, err := conn.ListTagsForResource(input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns ssoadmin service tags.
func Tags(tags tftags.KeyValueTags) []*ssoadmin.Tag {
	result := make([]*ssoadmin.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &ssoadmin.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from ssoadmin service tags.
func KeyValueTags(tags []*ssoadmin.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates ssoadmin service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn *ssoadmin.SSOAdmin, identifier string, resourceType string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &ssoadmin.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			InstanceArn: aws.String(resourceType),
			TagKeys:     aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResource(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &ssoadmin.TagResourceInput{
			ResourceArn: aws.String(identifier),
			InstanceArn: aws.String(resourceType),
			Tags:        Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResource(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
