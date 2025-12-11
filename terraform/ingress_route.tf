resource "kubernetes_manifest" "ingress_route" {
  manifest = {
    apiVersion = "traefik.io/v1alpha1"
    kind = "IngressRoute"
    metadata = {
      name = "fence-agent"
    namespace = "default"
    }
    spec = {
      entryPoints = ["websecure"]
      routes = [{
        kind = "Rule"
        match = "Host(`fence.homelab.lan`)"
        services = [{
          name = "fence-agent"
          port = 8081
        }]
      }]
    }
  }
}

