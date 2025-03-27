data "external" "assets_from_k8s" {
  program = ["./kservices"]
  
  query = {
    ns = "argo-cd"
  }
}

variable "assets" {
  type = set(string)
  default = []
}

locals {
  assets = tomap(data.external.assets_from_k8s.result)
}

resource "panop_zone" "zone1" {
  zone_name = "sks"
  zone_type = "internal"
}

resource "panop_asset" "asset" {
  for_each   = local.assets
  asset_name = each.key
  asset_type = "dns"
  zone_id    = panop_zone.zone1.id
}
