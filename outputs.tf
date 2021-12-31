output "oci_state_par" {
  description = "State file PAR URI for this repository"
  value = oci_objectstorage_preauthrequest.this.access_uri
}