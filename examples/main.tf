terraform {
  required_providers {
    persondb = {
      source = "hashicorp/persondb"
    }
  }
}

provider "persondb" {
  database_filename = "persons.db"
}

resource "persondb_person" "wim" {
  person_id  = "1"
  last_name  = "Van den Wyngaert"
  first_name = "Wim"
}
