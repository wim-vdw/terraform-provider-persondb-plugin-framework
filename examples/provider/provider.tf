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
