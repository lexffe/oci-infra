package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v54/common"
	"github.com/oracle/oci-go-sdk/v54/identity"
	"github.com/oracle/oci-go-sdk/v54/objectstorage"
	"github.com/oracle/oci-go-sdk/v54/objectstorage/transfer"
)

type kv map[string]string

/*
https://github.com/oracle/oci-go-sdk/blob/master/example/example_objectstorage_test.go
*/

func main() {

	// Main routine: create new compartment, create bucket, upload an empty file, return PAR URL in log, to bootstrap terraform state file

	// Parent compartment Id

	parentCompartment := os.Getenv("OCI_COMPARTMENT_ID")

	commonTags := kv{"Type": "Bootstrap"}

	ctx := context.Background()
	config := common.DefaultConfigProvider()

	osClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(config)

	if err != nil {
		log.Fatalln("%v", err)
	}

	idClient, err := identity.NewIdentityClientWithConfigurationProvider(config)

	if err != nil {
		log.Fatalln("%v", err)
	}

	// Create compartment

	comp := createCompartment(ctx, &idClient, common.String(parentCompartment), common.String("Terraformed Compartment"), common.String("IaC Managed infrastructure"), commonTags)

	// Get OS NS

	ns := getNS(ctx, &osClient, comp, commonTags)

	// Create OS bucket

	bucketName := "terraform-state-bucket"
	createBucket(ctx, &osClient, ns, common.String(parentCompartment), common.String(bucketName), commonTags)

	// upload empty file to bucket as state file

	filename := "oci.tfstate"
	empty := strings.NewReader("")
	uploadToBucket(ctx, &osClient, ns, common.String(bucketName), common.String(filename), empty)

	// Generate pre-authenticated request URL

	generatePar(ctx, &osClient, ns, common.String(bucketName), common.String(filename), common.String("bootstrap"))

}

func createCompartment(ctx context.Context, c *identity.IdentityClient, parent *string, name *string, desc *string, tags kv) *string {

	req := identity.CreateCompartmentRequest{
		CreateCompartmentDetails: identity.CreateCompartmentDetails{
			CompartmentId: parent,
			Name:          name,
			Description:   desc,
			FreeformTags:  tags,
		},
	}

	r, err := c.CreateCompartment(ctx, req)

	if err != nil {
		log.Fatalln("%v", err)
	}

	return r.Compartment.Id
}

func getNS(ctx context.Context, c *objectstorage.ObjectStorageClient, comp *string, tags kv) *string {

	req := objectstorage.GetNamespaceRequest {
		CompartmentId: comp,
	}

	r, err := c.GetNamespace(ctx, req)

	if err != nil {
		log.Fatalln("%v", err)
	}

	return r.Value

}

func createBucket(ctx context.Context, c *objectstorage.ObjectStorageClient, ns *string, comp *string, bn *string, tags kv) {
	
	req := objectstorage.CreateBucketRequest{
		NamespaceName: ns,
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			Name:          bn,
			CompartmentId: comp,
			FreeformTags:  tags,
		},
	}

	r, err := c.CreateBucket(ctx, req)

	if err != nil {
		log.Fatalln("%v", err)
	}

}

func uploadToBucket(ctx context.Context, c *objectstorage.ObjectStorageClient, ns *string, bucket *string, name *string, st io.Reader) {
	
	manager := transfer.NewUploadManager()

	req := transfer.UploadStreamRequest {
		UploadRequest: transfer.UploadRequest {
			ObjectStorageClient: c,
			NamespaceName: ns,
			BucketName: bucket,
			ObjectName: name,
		},
		StreamReader: st,
	}

	r, err := manager.UploadStream(ctx, req)

	if err != nil {
		log.Fatalln("%v", err)
	}

	log.Println("n/%v/b/%v/o/%v \t etag: %v", ns, bucket, name, r.SinglepartUploadResponse.PutObjectResponse.ETag)

}

func generatePar(ctx context.Context, c *objectstorage.ObjectStorageClient, ns *string, bucket *string, object *string, parName *string) {

	expiryTime := time.Now().Add(time.Hour)

	sdkTime := common.SDKTime {
		Time: expiryTime,
	}

	req := objectstorage.CreatePreauthenticatedRequestRequest {
		NamespaceName: ns,
		BucketName: bucket,
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails { 
			Name: parName,
			ObjectName: object,
			AccessType: objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectreadwrite,
			TimeExpires: &sdkTime,
		},
	}

	r, err := c.CreatePreauthenticatedRequest(ctx, req)
	
	if err != nil {
		log.Fatalln("%v", err)
	}
	// r.AccessUri

	log.Println("Generated PAR %v, expiry date in %v", r.AccessUri, expiryTime)
}
