# oci-infra

Everything is contained in a created compartment, created when bootstrapping.

- `OCI_COMPARTMENT_ID`: parent compartment

## Chicken or Egg?

Bootstrap this repository by running

```shell

OCI_COMPARTMENT_ID="ocid....." go run main.go

terraform init -backend-config="address="

terraform import oci_identity_compartment.this ${OUTPUT}

terraform import oci_objectstorage_bucket.state ${OUTPUT}

terraform refresh
```

## PAR urls

Managing PARs seem like a huge hassle.

## References

[OCI - terraform using object store](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformUsingObjectStore.htm)