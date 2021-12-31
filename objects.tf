# State files, PARs

resource "oci_objectstorage_preauthrequest" "this" {
  object_name  = data.oci_objectstorage_object.state.object
  access_type  = "ObjectReadWrite"
  bucket       = oci_objectstorage_bucket.state.name
  name         = "bootstrap"
  namespace    = data.oci_objectstorage_namespace.ns.namespace
  time_expires = "2022-06-30T23:59:59.000Z"
}

resource "oci_objectstorage_object" "cf_infra_state" {
  bucket    = oci_objectstorage_bucket.state.name
  namespace = data.oci_objectstorage_namespace.ns.namespace
  object    = "cf-infra.tfstate"

  content = "" # content managed by external world
}

resource "oci_objectstorage_preauthrequest" "cf_infra_state" {
  object_name  = oci_objectstorage_object.cf_infra_state.object
  access_type  = "ObjectReadWrite"
  bucket       = oci_objectstorage_bucket.state.name
  name         = "bootstrap"
  namespace    = data.oci_objectstorage_namespace.ns.namespace
  time_expires = "2022-06-30T23:59:59.000Z"
}
