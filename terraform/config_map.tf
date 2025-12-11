resource "kubernetes_config_map_v1" "fence" {
  metadata {
    name = "fence-config"
  }
  data = {
    "fence.yaml" = "${file("${path.cwd}/fence.yaml")}"
  }
}
