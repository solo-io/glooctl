package secret

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/bootstrap/secretstorage"

	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/glooctl/pkg/secret"
	"github.com/spf13/cobra"
)

const (
	flagName     = "name"
	flagFilename = "filename"
)

func createCmd(opts *bootstrap.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a secret",
	}
	cmd.AddCommand(
		createAWS(opts),
		createGCF(opts),
		createCertificate(opts))

	return cmd
}

func createAWS(storageOpts *bootstrap.Options) *cobra.Command {
	opts := secret.AWSOptions{}
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
			si, err := secretstorage.Bootstrap(*storageOpts)
			if err != nil {
				fmt.Println("Unable to get secret client:", err)
				return nil
			}

			if err := secret.CreateAWS(si, &opts); err != nil {
				fmt.Printf("Unable to create secret %s: %q\n", opts.Name, err)
				return nil
			}
			fmt.Printf("Created secret %s for AWS\n", opts.Name)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.UseEnv, "env", "e", false,
		"use environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to create secret")
	flags.StringVar(&opts.Name, flagName, "", "name for secret")
	flags.StringVarP(&opts.Filename, flagFilename, "f",
		"", "use credential file and default profile (defaults to ~/.aws/credentials)")
	flags.StringVar(&opts.KeyID, "access-key-id", "", "AWS access key ID")
	flags.StringVar(&opts.SecretKey, "secret-access-key", "", "AWS secret access key")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagFilename(flagFilename)

	return cmd
}

func createGCF(storageOpts *bootstrap.Options) *cobra.Command {
	opts := secret.GoogleOptions{}
	cmd := &cobra.Command{
		Use:   "google",
		Short: "create secret for upstream type Google (Google Cloud Function)",
		RunE: func(c *cobra.Command, a []string) error {
			si, err := secretstorage.Bootstrap(*storageOpts)
			if err != nil {
				fmt.Println("Unable to get secret client:", err)
				return nil
			}
			if err := secret.CreateGoogle(si, &opts); err != nil {
				fmt.Printf("Unable to create secret %s: %q\n", opts.Name, err)
				return nil
			}
			fmt.Printf("Created secret %s for Google Cloud Function\n", opts.Name)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.Name, flagName, "", "name for secret")
	cmd.MarkFlagRequired(flagName)
	flags.StringVarP(&opts.Filename, flagFilename, "f", "", "service account key file")
	cmd.MarkFlagFilename(flagFilename)
	cmd.MarkFlagRequired(flagFilename)
	return cmd
}

func createCertificate(storageOpts *bootstrap.Options) *cobra.Command {
	opts := secret.CertificateOptions{}
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "create a secret for certificate",
		Run: func(c *cobra.Command, args []string) {
			si, err := secretstorage.Bootstrap(*storageOpts)
			if err != nil {
				fmt.Println("Unable to get secret client:", err)
				return
			}
			if err := secret.CreateCertificate(si, &opts); err != nil {
				fmt.Printf("Unable to create secret %s: %q\n", opts.Name, err)
				return
			}
			fmt.Printf("Created secret %s for certificate\n", opts.Name)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.Name, flagName, "", "name for secret")
	cmd.MarkFlagRequired(flagName)
	flags.StringVarP(&opts.CAChain, "ca-chain", "c", "", "certificate authority chain certificate")
	cmd.MarkFlagFilename("ca-chain")
	cmd.MarkFlagRequired("ca-chain")
	flags.StringVarP(&opts.PrivateKey, "private-key", "p", "", "private key file")
	cmd.MarkFlagFilename("private-key")
	cmd.MarkFlagRequired("private-key")

	return cmd
}
