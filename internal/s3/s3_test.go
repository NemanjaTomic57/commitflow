package s3

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/joho/godotenv"
)

func TestSTSIntegration(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	client := sts.NewFromConfig(cfg)

	result, err := client.GetCallerIdentity(
		ctx,
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.Account == nil || *result.Account == "" {
		t.Fatal("account id is empty")
	}

	t.Logf("Account: %s", *result.Account)
	t.Logf("ARN: %s", *result.Arn)
	t.Logf("UserID: %s", *result.UserId)
}
