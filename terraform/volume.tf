resource "kubernetes_persistent_volume_claim_v1" "fence-pvc" {
  metadata {
    name = "fence-agent-pvc"
    labels = {
      app = "fence-agent"
    }
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    storage_class_name = "local-path"
    resources {
      requests = {
        storage = "1Gi"
      }
    }
  }
}

