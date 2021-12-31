# oci-infra

Everything is contained in a created compartment, created when bootstrapping.

- `OCI_COMPARTMENT_ID`: parent compartment

## Chicken or Egg?

Bootstrap this repository by running

```shell

## get parent compartment
oci iam compartment list # optionally --include-root -> tenancy id

## populate 

echo -e "parent_compartment_ocid = \"ocid...\"" >> terraform.tfvars

OCI_COMPARTMENT_ID="ocid....." go run main.go

echo -e "bootstrap_par_time = \"2021-12-31T18:00:00.123Z\"" >> terraform.tfvars

terraform init -backend-config="address=https://objectstorage.{region}.oraclecloud.com/p/XYZ/n/{ns}/b/{bucket}/o/oci.tfstate"

## compartment id
terraform import oci_identity_compartment.this ${new_compartment_id}

## id: n/{namespaceName}/b/{bucketName}
terraform import oci_objectstorage_bucket.state ${new_bucket_path}

## id: n/{namespaceName}/b/{bucketName}/p/{parId}
terraform import oci_objectstorage_preauthrequest.bootstrap ${new_par_id}

terraform refresh
```

You should create a new PAR with a later expiry date, and re-initialise the backend.

```shell
terraform init -reconfigure -backend-config="address=https://objectstorage.{region}.oraclecloud.com/p/parId/n/ns/b/bucket/o/oci.tfstate"
```

## Backend

It is possible to use object storage's S3-compat endpoint as terraform backend, but that introduces some complexity.

## PAR urls

Managing PARs seem like a huge hassle.

Every time terraform apply runs, the PAR url expiry date for the state file should extend (?)

## References

[OCI - terraform using object store](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformUsingObjectStore.htm)