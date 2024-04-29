#########################################################################################################################
# Providers
#########################################################################################################################

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
    template = {
      source = "hashicorp/template"
      version = "2.2.0"
    }
  }
}


terraform {
  required_version = ">= 0.13.5"
}