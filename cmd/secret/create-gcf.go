package secret

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	secret "github.com/solo-io/gloo-secret"
	"github.com/solo-io/glooctl/pkg/secrets"
	"github.com/spf13/cobra"
)

const (
	serviceAccountJsonKeyFile = "json_key_file"
)

func createGCF() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "google",
		Short: "create secret for upstream type Google (Google Cloud Function)",
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
			if err := runCreateGCF(si, name, filename); err != nil {
				fmt.Printf("Unable to create secret %s: %q\n", name, err)
				return nil
			}
			fmt.Printf("Created secret %s for Google Cloud Function upstream\n", name)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, "filename", "f", "", "service account key file")
	cmd.MarkFlagFilename("filename")
	cmd.MarkFlagRequired("filename")
	return cmd
}

func runCreateGCF(si secret.SecretInterface, name, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "unable to read service account key file %s", filename)
	}
	s := &secret.Secret{
		Name: name,
		Data: map[string][]byte{
			serviceAccountJsonKeyFile: b,
		},
	}
	_, err = si.V1().Create(s)
	return err
}
