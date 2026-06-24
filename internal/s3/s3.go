package s3

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var bucketName = os.Getenv("AWS_S3_BUCKET")

func InitializeS3(ctx context.Context) {
	client := getS3Client()
	createBucket(ctx, client)
	deleteAllObjects(ctx, client)
}

func getS3Client() *s3.Client {
	endpoint := os.Getenv("MINIO_URL")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	opts := s3.Options{
		Region: os.Getenv("AWS_REGION"),

		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),

		BaseEndpoint: aws.String(endpoint),

		UsePathStyle: true,
	}

	client := s3.New(opts)

	return client
}

func createBucket(ctx context.Context, client *s3.Client) {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			log.Fatalf("aws.CreateBucket() -> ERROR: create bucket failed: %v", err)
		}
	}
}

func deleteAllObjects(ctx context.Context, client *s3.Client) {
	// List all objects
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Printf("aws.ResetS3Data() -> error listing next page of paginator: %v", err)
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
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{
				Objects: objects,
				Quiet:   aws.Bool(false),
			},
		})
		if err != nil {
			log.Printf("aws.ResetS3Data() -> error deleting bucket objects: %v", err)
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
}
