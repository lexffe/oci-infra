# The following resources are bootstrapped and imported.

resource "oci_identity_compartment" "this" {
  compartment_id = var.parent_compartment_ocid
  description    = "IaC Managed infrastructure"
  name           = "Terraformed Compartment"
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

resource "oci_objectstorage_object" "state" {
  bucket = oci_objectstorage_bucket.state.name
  source_uri_details {
    region = var.region
    namespace = data.oci_objectstorage_namespace.ns.namespace
    bucket = oci_objectstorage_bucket.state
  }
}

# this PAU will always expire "now", i.e. whenever tf apply runs, this PAU will never work. 
# (this resource should only be used for bootstrapping)
resource "oci_objectstorage_preauthrequest" "bootstrap" {
  object_name = ""
  access_type = "ObjectReadWrite"
  bucket = oci_objectstorage_bucket.state.name
  name = "bootstrap"
  namespace = data.oci_objectstorage_namespace.ns.namespace
  time_expires = timestamp()
}