package secret

import (
	"fmt"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
	"github.com/solo-io/glooctl/pkg/client"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
)

const (
	// expected map identifiers for secrets
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"
)

func createAWS(storageOpts *client.StorageOptions, createOpts *CreateOptions) *cobra.Command {
	var useEnv bool
	var filename string
	var keyId string
	var secretKey string
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "create secret for upstream type AWS",
		Long: `
Creates a secret that can be used by upstream of type 'aws'.
By default, it will use credentials file. You can change the
location of the file using --filename flag. Alternatively,
use --env flag to use the default AWS environment variables
or provide them directly using --access-key-id and 
--secret-access-key flags.
		`,
		RunE: func(c *cobra.Command, a []string) error {
			name := createOpts.Name
			if name == "" {
				return fmt.Errorf("name for secret missing")
			}
			si, err := client.SecretClient(storageOpts)
			if err != nil {
				fmt.Println("Unable to get secret client:", err)
				return nil
			}

			id, key, err := idAndKey(useEnv, keyId, secretKey, filename)
			if err != nil {
				return errors.Wrap(err, "unable to get AWS credentials")
			}
			if err := runCreateAWS(si, name, id, key); err != nil {
				fmt.Printf("Unable to create secret %s: %q\n", name, err)
				return nil
			}
			fmt.Printf("Created secret %s for AWS upstream\n", name)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&useEnv, "env", "e", false,
		"use environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to create secret")
	flags.StringVarP(&filename, "filename", "f",
		"", "use credential file and default profile (defaults to ~/.aws/credentials)")
	flags.StringVar(&keyId, "access-key-id", "", "AWS access key ID")
	flags.StringVar(&secretKey, "secret-access-key", "", "AWS secret access key")
	cmd.MarkFlagFilename("filename")
	return cmd
}

func runCreateAWS(si secret.SecretInterface, name, id, key string) error {
	s := &secret.Secret{
		Name: name,
		Data: map[string][]byte{
			awsAccessKey: []byte(id),
			awsSecretKey: []byte(key),
		},
	}
	_, err := si.V1().Create(s)
	return err
}

func idAndKey(useEnv bool, keyId, secretKey, filename string) (string, string, error) {
	if keyId != "" || secretKey != "" {
		fmt.Println("Using values passed in CLI")
		if keyId != "" && secretKey != "" {
			return keyId, secretKey, nil
		}
		return "", "", fmt.Errorf("both access-key-id and secret-access-key must be provided")
	}

	var creds *credentials.Credentials
	if useEnv {
		creds = credentials.NewEnvCredentials()
	} else {
		//TODO: add a flag for profile
		creds = credentials.NewSharedCredentials(filename, "")
	}
	vals, err := creds.Get()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to retrieve aws credentials")
	}
	return vals.AccessKeyID, vals.SecretAccessKey, nil
}
