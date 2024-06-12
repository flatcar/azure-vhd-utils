package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"gopkg.in/urfave/cli.v1"

	"github.com/flatcar/azure-vhd-utils/op"
)

func createServiceClient(c *cli.Context, account, key string) (*service.Client, error) {
	var (
		client *service.Client
		err    error
	)
	accountURL := fmt.Sprintf("https://%s.blob.core.windows.net", url.PathEscape(account))

	if key != "" {
		skc, err := service.NewSharedKeyCredential(account, key)
		if err != nil {
			return nil, fmt.Errorf("Failed to create shared key credential: %w", err)
		}
		client, err = service.NewClientWithSharedKeyCredential(accountURL, skc, nil)
	} else {
		opts := azidentity.DefaultAzureCredentialOptions{
			DisableInstanceDiscovery: c.Bool("disableinstancediscovery"),
			TenantID:                 c.String("tenantid"),
		}
		creds, err := azidentity.NewDefaultAzureCredential(&opts)
		if err != nil {
			return nil, fmt.Errorf("Failed to create default Azure credential: %w", err)
		}
		client, err = service.NewClient(accountURL, creds, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to create storage service client: %w", err)
	}

	return client, nil
}

func vhdUploadCmdHandler() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "Upload a local VHD to Azure storage as page blob",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "localvhdpath",
				Usage: "Path to source VHD in the local machine.",
			},
			cli.StringFlag{
				Name:  "stgaccountname",
				Usage: "Azure storage account name.",
			},
			cli.StringFlag{
				Name:  "stgaccountkey",
				Usage: "Azure storage account key (optional).",
			},
			cli.StringFlag{
				Name:  "tenantid",
				Usage: "Azure Tenant ID.",
			},
			cli.BoolFlag{
				Name:  "disableinstancediscovery",
				Usage: "Skip the request to Microsoft Entra before authenticating.",
			},
			cli.StringFlag{
				Name:  "containername",
				Usage: "Name of the container holding destination page blob. (Default: vhds)",
			},
			cli.StringFlag{
				Name:  "blobname",
				Usage: "Name of the destination page blob.",
			},
			cli.StringFlag{
				Name:  "parallelism",
				Usage: "Number of concurrent goroutines to be used for upload",
			},
			cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Overwrite the blob if already exists.",
			},
		},
		Action: func(c *cli.Context) error {
			const PageBlobPageSize int64 = 512
			const PageBlobPageSetSize int64 = 4 * 1024 * 1024

			localVHDPath := c.String("localvhdpath")
			if localVHDPath == "" {
				return errors.New("Missing required argument --localvhdpath")
			}

			stgAccountName := c.String("stgaccountname")
			if stgAccountName == "" {
				return errors.New("Missing required argument --stgaccountname")
			}

			// account key is optional, if not passed,
			// then we expect that the required storage
			// blob roles for storage account are already
			// assigned to azure account
			stgAccountKey := c.String("stgaccountkey")

			containerName := c.String("containername")
			if containerName == "" {
				containerName = "vhds"
				log.Println("Using default container 'vhds'")
			}

			blobName := c.String("blobname")
			if blobName == "" {
				return errors.New("Missing required argument --blobname")
			}

			if !strings.HasSuffix(strings.ToLower(blobName), ".vhd") {
				blobName = blobName + ".vhd"
			}

			parallelism := int(0)
			if c.IsSet("parallelism") {
				p, err := strconv.ParseUint(c.String("parallelism"), 10, 32)
				if err != nil {
					return fmt.Errorf("invalid index value --parallelism: %s", err)
				}
				parallelism = int(p)
			} else {
				parallelism = 8 * runtime.NumCPU()
				log.Printf("Using default parallelism [8*NumCPU] : %d\n", parallelism)
			}

			overwrite := c.IsSet("overwrite")

			serviceClient, err := createServiceClient(c, stgAccountName, stgAccountKey)
			if err != nil {
				return err
			}

			uopts := op.UploadOptions{
				Overwrite:   overwrite,
				Parallelism: parallelism,
				Logger: func(s string) {
					log.Println(s)
				},
			}
			err = op.Upload(context.TODO(), serviceClient, containerName, blobName, localVHDPath, &uopts)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
}
