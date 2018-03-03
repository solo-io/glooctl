package secret

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/glooctl/pkg/secrets"
	"github.com/solo-io/glooctl/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/spf13/cobra"
)

const (
	// expected map identifiers for secrets
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"
)

func createAWS() *cobra.Command {
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
			flags := c.InheritedFlags()
			name, _ := flags.GetString("name")
			if name == "" {
				return fmt.Errorf("name for secret missing")
			}
			si, err := secrets.GetSecretClient(c)
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

func runCreateAWS(si corev1.SecretInterface, name, id, key string) error {
	s := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: map[string]string{
			awsAccessKey: id,
			awsSecretKey: key,
		},
	}
	_, err := si.Create(s)
	return err
}

func idAndKey(useEnv bool, keyId, secretKey, filename string) (string, string, error) {
	if useEnv {
		fmt.Println("Using environment variables")
		keyId = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		if keyId == "" || secretKey == "" {
			return "", "", fmt.Errorf("Both environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY should be set")
		}
		return keyId, secretKey, nil
	}

	if keyId != "" || secretKey != "" {
		fmt.Println("Using values passed in CLI")
		if keyId != "" && secretKey != "" {
			return keyId, secretKey, nil
		}
		return "", "", fmt.Errorf("both access-key-id and secret-access-key must be provided")
	}

	if filename == "" {
		filename = filepath.Join(util.HomeDir(), ".aws", "credentials")
	}
	fmt.Println("Using the file", filename)
	// TODO - use a toml parser? https://github.com/pelletier/go-toml has
	// errors parsing the file
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", err
	}
	lines := strings.Split(string(b), "\n")
	for i, l := range lines {
		if strings.TrimSpace(l) == "[default]" {
			if i+2 >= len(lines) {
				return "", "", fmt.Errorf("unable to find key in the file")
			}
			parts := strings.SplitN(lines[i+1], "=", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "aws_access_key_id" {
				return "", "", fmt.Errorf("unable to find aws access key id in the file")
			}
			id := strings.TrimSpace(parts[1])

			parts = strings.SplitN(lines[i+2], "=", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "aws_secret_access_key" {
				return "", "", fmt.Errorf("unable to find aws secret access key in the file")
			}
			key := strings.TrimSpace(parts[1])
			return id, key, nil
		}
	}
	return "", "", fmt.Errorf("unable to find the key in the file")
}
