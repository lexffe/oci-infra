# oci-infra

Everything is contained in a created compartment, created when bootstrapping.

- `OCI_COMPARTMENT_ID`: parent compartment

## Chicken or Egg?

Bootstrap this repository by running

```shell

OCI_COMPARTMENT_ID="ocid....." go run main.go

terraform init -backend-config="address=PAR_URI"

## compartment id
terraform import oci_identity_compartment.this ${OUTPUT}

## id: n/{namespaceName}/b/{bucketName}
terraform import oci_objectstorage_bucket.state ${OUTPUT}

## id: n/{namespaceName}/b/{bucketName}/p/{parId}
terraform import oci_objectstorage_preauthrequest.bootstrap ${OUTPUT}

terraform refresh
```

## Backend

It is possible to use object storage's S3-compat endpoint as terraform backend, but that introduces some complexity.

## PAR urls

Managing PARs seem like a huge hassle.

Every time terraform apply runs, the PAR url expiry date for the state file should extend (?)

## References

[OCI - terraform using object store](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformUsingObjectStore.htm)