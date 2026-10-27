# Estado gravado no disco local. Nenhum backend remoto (S3, GCS, AzureRM, Terraform Cloud).
terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}
