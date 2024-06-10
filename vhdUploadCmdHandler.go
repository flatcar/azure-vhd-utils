package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"gopkg.in/urfave/cli.v1"

	"github.com/flatcar/azure-vhd-utils/upload"
	"github.com/flatcar/azure-vhd-utils/upload/metadata"
	"github.com/flatcar/azure-vhd-utils/vhdcore/common"
	"github.com/flatcar/azure-vhd-utils/vhdcore/diskstream"
	"github.com/flatcar/azure-vhd-utils/vhdcore/validator"
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
				Usage: "Use managed identity.",
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

			ensureVHDSanity(localVHDPath)
			diskStream, err := diskstream.CreateNewDiskStream(localVHDPath)
			if err != nil {
				return err
			}
			defer diskStream.Close()

			serviceClient, err := createServiceClient(c, stgAccountName, stgAccountKey)
			if err != nil {
				return err
			}
			containerClient := serviceClient.NewContainerClient(containerName)
			pageblobClient := containerClient.NewPageBlobClient(blobName)
			blobClient := pageblobClient.BlobClient()

			_, err = containerClient.Create(context.TODO(), nil)
			if err != nil && !bloberror.HasCode(err, bloberror.ContainerAlreadyExists, bloberror.ResourceAlreadyExists) {
				return err
			}

			blobExists := true
			blobProperties, err := blobClient.GetProperties(context.TODO(), nil)
			if err != nil {
				if !bloberror.HasCode(err, bloberror.BlobNotFound, bloberror.ResourceNotFound) {
					return err
				}
				blobExists = false
			}

			resume := false
			var blobMetadata *metadata.Metadata
			if blobExists {
				if !overwrite {
					if len(blobProperties.ContentMD5) > 0 {
						log.Fatalf("VHD exists in blob storage with name '%s'. If you want to upload again, use the --overwrite option.", blobName)
					}
					blobMetadata, err = metadata.NewMetadataFromBlobMetadata(blobProperties.Metadata)
					if err != nil {
						return err
					}
					if blobMetadata == nil {
						log.Fatalf("There is no upload metadata associated with the existing blob '%s', so upload operation cannot be resumed, use --overwrite option.", blobName)
					}
					resume = true
					log.Printf("Blob with name '%s' already exists, checking upload can be resumed\n", blobName)
				}
			}

			localMetadata := getLocalVHDMetadata(localVHDPath)
			var rangesToSkip []*common.IndexRange
			if resume {
				if errs := metadata.CompareMetadata(blobMetadata, localMetadata); len(errs) > 0 {
					printErrorsAndFatal(errs)
				}
				rangesToSkip = getAlreadyUploadedBlobRanges(pageblobClient)
			} else {
				createBlob(pageblobClient, diskStream.GetSize(), localMetadata)
			}

			uploadableRanges, err := upload.LocateUploadableRanges(diskStream, rangesToSkip, PageBlobPageSize, PageBlobPageSetSize)
			if err != nil {
				return err
			}

			uploadableRanges, err = upload.DetectEmptyRanges(diskStream, uploadableRanges)
			if err != nil {
				return err
			}

			cxt := &upload.DiskUploadContext{
				VhdStream:             diskStream,
				AlreadyProcessedBytes: diskStream.GetSize() - common.TotalRangeLength(uploadableRanges),
				UploadableRanges:      uploadableRanges,
				PageblobClient:        pageblobClient,
				Parallelism:           parallelism,
				Resume:                resume,
			}

			err = upload.Upload(cxt)
			if err != nil {
				return err
			}

			setBlobMD5Hash(blobClient, localMetadata)
			fmt.Println("\nUpload completed")
			return nil
		},
	}
}

// printErrorsAndFatal prints the errors in a slice one by one and then exit
func printErrorsAndFatal(errs []error) {
	fmt.Println()
	for _, e := range errs {
		fmt.Println(e)
	}
	log.Fatal("Cannot continue due to above errors.")
}

// ensureVHDSanity ensure is VHD is valid for Azure.
func ensureVHDSanity(localVHDPath string) {
	if err := validator.ValidateVhd(localVHDPath); err != nil {
		log.Fatal(err)
	}

	if err := validator.ValidateVhdSize(localVHDPath); err != nil {
		log.Fatal(err)
	}
}

// getLocalVHDMetadata returns the metadata of a local VHD
func getLocalVHDMetadata(localVHDPath string) *metadata.Metadata {
	localMetadata, err := metadata.NewMetadataFromLocalVHD(localVHDPath)
	if err != nil {
		log.Fatal(err)
	}
	return localMetadata
}

// createBlob creates a page blob of specific size and sets custom metadata
// The parameter client is the Azure blob service client, parameter containerName is the name of an existing container
// in which the page blob needs to be created, parameter blobName is name for the new page blob, size is the size of
// the new page blob in bytes and parameter vhdMetadata is the custom metadata to be associacted with the page blob
func createBlob(client *pageblob.Client, size int64, vhdMetadata *metadata.Metadata) {
	m, err := vhdMetadata.ToMap()
	if err != nil {
		log.Fatal(err)
	}
	opts := pageblob.CreateOptions{
		Metadata: m,
	}
	_, err = client.Create(context.TODO(), size, &opts)
	if err != nil {
		log.Fatal(err)
	}
}

// setBlobMD5Hash sets MD5 hash of the blob in it's properties
func setBlobMD5Hash(client *blob.Client, vhdMetadata *metadata.Metadata) {
	if vhdMetadata.FileMetadata.MD5Hash == nil {
		return
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(vhdMetadata.FileMetadata.MD5Hash)))
	base64.StdEncoding.Encode(buf, vhdMetadata.FileMetadata.MD5Hash)
	blobHeaders := blob.HTTPHeaders{
		BlobContentMD5: buf,
	}
	_, err := client.SetHTTPHeaders(context.TODO(), blobHeaders, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// getAlreadyUploadedBlobRanges returns the range slice containing ranges of a page blob those are already uploaded.
// The parameter client is the Azure blob service client, parameter containerName is the name of an existing container
// in which the page blob resides, parameter blobName is name for the page blob
func getAlreadyUploadedBlobRanges(client *pageblob.Client) []*common.IndexRange {
	var (
		marker       *string
		rangesToSkip []*common.IndexRange
	)
	for {
		opts := pageblob.GetPageRangesOptions{
			Marker: marker,
		}
		pager := client.NewGetPageRangesPager(&opts)
		for pager.More() {
			response, err := pager.NextPage(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			tmpRanges := make([]*common.IndexRange, len(response.PageRange))
			for i, page := range response.PageRange {
				tmpRanges[i] = common.NewIndexRange(*page.Start, *page.End)
			}
			rangesToSkip = append(rangesToSkip, tmpRanges...)
			marker = response.NextMarker
		}
		if marker == nil || *marker == "" {
			break
		}
	}
	return rangesToSkip
}
