package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ResetS3Data() error {
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return fmt.Errorf("AWS_S3_BUCKET is not set in .env file")
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("aws.ResetS3Data() -> error loading default aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// List all objects
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("aws.ResetS3Data() -> error listing next page of paginator: %w", err)
		}

		if len(page.Contents) == 0 {
			continue
		}

		objects := make([]types.ObjectIdentifier, 0, len(page.Contents))
		for _, obj := range page.Contents {
			objects = append(objects, types.ObjectIdentifier{
				Key: obj.Key,
			})
		}

		// Delete all objects
		output, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: &bucket,
			Delete: &types.Delete{
				Objects: objects,
				Quiet:   aws.Bool(false),
			},
		})
		if err != nil {
			return fmt.Errorf("aws.ResetS3Data() -> error deleting bucket objects: %w", err)
		}

		for _, e := range output.Errors {
			fmt.Printf(
				"Key=%s Code=%s Message=%s\n",
				aws.ToString(e.Key),
				aws.ToString(e.Code),
				aws.ToString(e.Message),
			)
		}
	}
	return nil
}
