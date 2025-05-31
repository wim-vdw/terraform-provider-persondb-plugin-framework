# List person with ID 1 in the database.
data "persondb_person" "wim" {
  person_id = "1"
}
