provider "kubernetes" {
  config_path = "~/.kube/config"
}
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.0.1"
    }
  }
}
