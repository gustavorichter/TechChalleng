terraform {
  required_version = ">= 1.5.0"

  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
    }
  }
}

# Cluster 100% local (Docker Desktop, minikube, kind, k3d).
# O provider lê o kubeconfig da máquina; não há credencial de nuvem.
provider "kubernetes" {
  config_path    = pathexpand(var.kubeconfig_path)
  config_context = var.kube_context != "" ? var.kube_context : null
}
