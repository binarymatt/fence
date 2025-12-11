resource "kubernetes_deployment_v1" "fence-agent" {
  wait_for_rollout = false
  metadata {
    name = "fence-agent"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "fence-agent"
      }
    }

    template {
      metadata {
        labels = {
          app = "fence-agent"
        }
      }

      spec {
        automount_service_account_token = false
        enable_service_links = false
        container {
          name  = "fence-agent"
          image = "registry.homelab.lan/fence-agent:58833f5"
          port {
            container_port =  8081
          }
          volume_mount {
            name = "storage"
            mount_path = "/fence"
          }
          volume_mount {
            name = "config"
            mount_path = "/config"
            read_only = true
          }
          env {
            name = "FENCE_CONFIG"
            value = "/config/fence.yaml"
          }

        }
        volume {
          name = "config"
          config_map {
            name = "fence-config"
            items {
              key = "fence.yaml"
              path = "fence.yaml"
            }
          }
        }
        volume {
          name = "storage"
          persistent_volume_claim {
            claim_name = "fence-agent-pvc"
          }
        }
      }
    }
  }
}
