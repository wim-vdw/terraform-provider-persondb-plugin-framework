# Test provider based on Terraform Plugin Framework

This repository demonstrates the implementation of a custom Terraform provider using the Terraform Plugin Framework.  
The provider interacts with a simple SQLite database to manage `person` resources, showcasing key features of Terraform
provider development.

## Features

- **CRUD Operations**: Full support for Create, Read, Update, and Delete operations on `person` resources.
- **Drift Detection**: Automatically detects and reconciles changes made outside of Terraform to ensure the state
  remains consistent.
- **Resource Recreation**: Supports resource recreation when required, such as when a resource is forcefully replaced.
- **Import Functionality**: Allows importing existing resources into Terraform state for management.

This project is designed for learning purposes and provides a hands-on example of how to build and test a Terraform
provider.

## Create a development build for local testing of the provider

Create directory `local_dev_build` in the `root` of this repository when it does not exist.  
The directory `local_dev_build` is used to store the local build artifacts and is excluded via `.gitignore`.

Create a local development build:

```bash
go mod tidy

# Linux/MacOS
go build -o local_dev_build/terraform-provider-persondb

# Windows (go-sqlite3 package requires CGO to be enabled -> gcc.exe must be installed)
$env:CGO_ENABLED=1; go build -o local_dev_build/terraform-provider-persondb.exe
```

## Run the local development tests with the Terraform CLI

Make sure you are located in the `local_tests` directory.  
File `terraformrc-local-dev` contains the Terraform CLI dev configuration overrides for the local development tests.  
File `main.tf` contains the Terraform configuration for the local development tests.

Activate the local development build tests by setting the `TF_CLI_CONFIG_FILE` environment variable to point to the
`terraformrc-local-dev` file:

```bash
# Linux
export TF_CLI_CONFIG_FILE="terraformrc-local-dev"

# Windows
$env:TF_CLI_CONFIG_FILE="terraformrc-local-dev"

```

> **ATTENTION:** You do not need to run `terraform init` as the local development build will be used automatically now.

Example:

```hcl
provider "persondb" {
  database_filename = "persons.db"
}

resource "persondb_person" "wim" {
  person_id  = "1"
  last_name  = "Van den Wyngaert"
  first_name = "Wim"
}
```

Run the local development tests (data will be persisted in the SQLite database `persons.db`):

```bash
# Run a plan to see the changes
terraform plan

# Create the resource and check the results
terraform apply -auto-approve
terraform show
terraform state list

# After removing the item from state and keeping the resource in the SQLite database, you can re-import it
terraform state rm persondb_person.wim
terraform plan
terraform apply -auto-approve
terraform import 'persondb_person.wim' '/person/1'
terraform show
terraform state list

# Destroy the resource
terraform destroy
```

Changing the `person_id` in the `main.tf` file will trigger a recreation of the resource.  
You can also change the `last_name` or `first_name` attributes to see how the provider handles updates.  
You can also make changes directly in the SQLite database to see how the provider handles drift detection.
