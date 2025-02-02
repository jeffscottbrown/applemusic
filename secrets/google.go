package secrets

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretspb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log/slog"
	"os"
)

var projectId = os.Getenv("PROJECT_ID")

func RetrieveSecret(secretName string) (string, error) {
	if projectId == "" {
		slog.Error("PROJECT_ID environment variable not set")
		return "", fmt.Errorf("PROJECT_ID environment variable not set")
	}
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	secretResource := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectId, secretName)

	req := &secretspb.AccessSecretVersionRequest{
		Name: secretResource,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		slog.Error("Error retrieving client secret", "secretName", secretName, "error", err)
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}
