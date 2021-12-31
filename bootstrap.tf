# The following resources are bootstrapped and imported.

resource "oci_identity_compartment" "this" {
  compartment_id = var.parent_compartment_ocid
  description    = "IaC Managed infrastructure"
  name           = "terraform-compartment"
}

data "oci_objectstorage_namespace" "ns" {
  compartment_id = oci_identity_compartment.this.id
}

resource "oci_objectstorage_bucket" "state" {
  compartment_id = oci_identity_compartment.this.id
  name           = "terraform-state-bucket"
  namespace      = data.oci_objectstorage_namespace.ns.namespace
  versioning     = "Enabled"
}

# Object managed out of band
data "oci_objectstorage_object" "state" {
  bucket    = oci_objectstorage_bucket.state.name
  namespace = data.oci_objectstorage_namespace.ns.namespace
  object    = "oci.tfstate"
}

# this PAU will always change.
resource "oci_objectstorage_preauthrequest" "bootstrap" {
  object_name  = data.oci_objectstorage_object.state.object
  access_type  = "ObjectReadWrite"
  bucket       = oci_objectstorage_bucket.state.name
  name         = "bootstrap"
  namespace    = data.oci_objectstorage_namespace.ns.namespace
  time_expires = timeadd(timestamp(), "336h")
}