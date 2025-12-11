resource "kubernetes_service_v1" "fence-agent" {
  wait_for_load_balancer = false
  metadata {
    name = "fence-agent"
  }
  spec {
    selector = {
      app = "fence-agent"
    }
    port {
      port        = 8081
    }

  }
}
