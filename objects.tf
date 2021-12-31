# State files, PARs

resource "oci_objectstorage_preauthrequest" "this" {
  object_name  = data.oci_objectstorage_object.state.object
  access_type  = "ObjectReadWrite"
  bucket       = oci_objectstorage_bucket.state.name
  name         = "bootstrap"
  namespace    = data.oci_objectstorage_namespace.ns.namespace
  time_expires = "2022-06-30T23:59:59.000Z"
}
