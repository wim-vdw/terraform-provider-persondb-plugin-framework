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

data "persondb_names" "wim" {
  person_id = "1"
}

output "test" {
  value = data.persondb_names.wim
}
