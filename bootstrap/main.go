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
		log.Fatal(err)
	}

	idClient, err := identity.NewIdentityClientWithConfigurationProvider(config)

	if err != nil {
		log.Fatal(err)
	}

	// Create compartment

	compId := createCompartment(ctx, &idClient, common.String(parentCompartment), common.String("terraform-compartment"), common.String("IaC Managed infrastructure"), commonTags)

	log.Println("sleeping 2 minutes for compartment to initialise")
	time.Sleep(time.Minute * 2)

	// Get OS NS
	// https://docs.oracle.com/en-us/iaas/Content/Object/Tasks/understandingnamespaces.htm#Understanding_Object_Storage_Namespaces
	// "The namespace spans all compartments within a region."

	ns := getNS(ctx, &osClient, common.String(parentCompartment), commonTags)

	// Create OS bucket in new compartment

	bucketName := "terraform-state-bucket"
	createBucket(ctx, &osClient, ns, compId, common.String(bucketName), commonTags)

	// upload empty file to bucket as state file

	filename := "oci.tfstate"
	empty := strings.NewReader("")
	uploadToBucket(ctx, &osClient, ns, common.String(bucketName), common.String(filename), empty)

	// Generate pre-authenticated request URL

	generatePar(ctx, &osClient, ns, common.String(bucketName), common.String(filename), common.String("bootstrap"))

}

func createCompartment(ctx context.Context, c *identity.IdentityClient, parent *string, name *string, desc *string, tags kv) *string {

	log.Printf("creating compartment %v\n", *name)

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
		log.Fatal(err)
	}

	log.Printf("Compartment %v created (%v)\n", *r.Compartment.Name, *r.Compartment.Id)

	return r.Compartment.Id
}

func getNS(ctx context.Context, c *objectstorage.ObjectStorageClient, compId *string, tags kv) *string {

	log.Printf("getting ns for compartment %v\n", *compId)

	req := objectstorage.GetNamespaceRequest{
		CompartmentId: compId,
	}

	r, err := c.GetNamespace(ctx, req)

	if err != nil {
		log.Fatal(err)
	}

	return r.Value
}

func createBucket(ctx context.Context, c *objectstorage.ObjectStorageClient, ns *string, compId *string, bucket *string, tags kv) {

	log.Printf("creating bucket %v\n", *bucket)

	req := objectstorage.CreateBucketRequest{
		NamespaceName: ns,
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			Name:          bucket,
			CompartmentId: compId,
			FreeformTags:  tags,
			Versioning:    objectstorage.CreateBucketDetailsVersioningEnabled,
		},
	}

	r, err := c.CreateBucket(ctx, req)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bucket %v created (%v)\n", *r.Bucket.Name, *r.Bucket.Id)
}

func uploadToBucket(ctx context.Context, c *objectstorage.ObjectStorageClient, ns *string, bucket *string, name *string, st io.Reader) {

	log.Printf("uploading to bucket %v\n", *bucket)

	manager := transfer.NewUploadManager()

	req := transfer.UploadStreamRequest{
		UploadRequest: transfer.UploadRequest{
			ObjectStorageClient: c,
			NamespaceName:       ns,
			BucketName:          bucket,
			ObjectName:          name,
		},
		StreamReader: st,
	}

	r, err := manager.UploadStream(ctx, req)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("n/%v/b/%v/o/%v \t etag: %v\n", *ns, *bucket, *name, *r.SinglepartUploadResponse.PutObjectResponse.ETag)

}

func generatePar(ctx context.Context, c *objectstorage.ObjectStorageClient, ns *string, bucket *string, object *string, parName *string) {

	log.Printf("generating PAR %v for %v\n", *parName, *object)

	expiryTime := time.Now().Add(time.Hour)

	sdkTime := common.SDKTime{
		Time: expiryTime,
	}

	req := objectstorage.CreatePreauthenticatedRequestRequest{
		NamespaceName: ns,
		BucketName:    bucket,
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:        parName,
			ObjectName:  object,
			AccessType:  objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectreadwrite,
			TimeExpires: &sdkTime,
		},
	}

	r, err := c.CreatePreauthenticatedRequest(ctx, req)

	if err != nil {
		log.Fatal(err)
	}
	// r.AccessUri

	timeStr, _ := expiryTime.MarshalText()

	log.Printf("Generated PAR %v, expiry date in %v\n", *r.AccessUri, string(timeStr))
}
