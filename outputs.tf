output "oci_state_par" {
  description = "State file PAR URI for this repository"
  value       = oci_objectstorage_preauthrequest.this.access_uri
}

output "cf_state_par" {
  description = "State file PAR URI for cf-infra"
  value       = oci_objectstorage_preauthrequest.cf_infra_state.access_uri
}