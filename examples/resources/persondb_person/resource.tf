# Create a person in the database.
resource "persondb_person" "wim" {
  person_id  = "1"
  last_name  = "Van den Wyngaert"
  first_name = "Wim"
}
