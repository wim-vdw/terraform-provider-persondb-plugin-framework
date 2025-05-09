terraform {
  required_providers {
    persondb = {
      source = "hashicorp/persondb"
    }
  }
}

provider "persondb" {}

data "persondb_names" "test" {}

output "test" {
  value = data.persondb_names.test
}
