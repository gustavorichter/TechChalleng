locals {
  labels = {
    "app.kubernetes.io/name"       = var.app_name
    "app.kubernetes.io/component"  = "api"
    "app.kubernetes.io/part-of"    = "oficina"
    "app.kubernetes.io/managed-by" = "terraform"
  }

  postgres_labels = {
    "app.kubernetes.io/name"       = "postgres"
    "app.kubernetes.io/component"  = "database"
    "app.kubernetes.io/part-of"    = "oficina"
    "app.kubernetes.io/managed-by" = "terraform"
  }
}

resource "kubernetes_namespace_v1" "oficina" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "oficina"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_config_map_v1" "oficina" {
  metadata {
    name      = "oficina-config"
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.labels
  }

  data = {
    PORT       = tostring(var.api_container_port)
    DB_HOST    = var.db_host
    DB_PORT    = var.db_port
    DB_NAME    = var.db_name
    DB_SSLMODE = var.db_sslmode
  }
}

resource "kubernetes_secret_v1" "oficina" {
  metadata {
    name      = "oficina-secret"
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.labels
  }

  type = "Opaque"

  data = {
    DB_USER        = var.db_user
    DB_PASSWORD    = var.db_password
    JWT_SECRET     = var.jwt_secret
    ADMIN_USERNAME = var.admin_username
    ADMIN_PASSWORD = var.admin_password
  }
}

resource "kubernetes_persistent_volume_claim_v1" "postgres" {
  metadata {
    name      = "postgres-data"
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.postgres_labels
  }

  spec {
    access_modes       = ["ReadWriteOnce"]
    storage_class_name = var.postgres_storage_class != "" ? var.postgres_storage_class : null

    resources {
      requests = {
        storage = var.postgres_storage
      }
    }
  }
}

resource "kubernetes_deployment_v1" "postgres" {
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.postgres_labels
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name" = "postgres"
      }
    }

    strategy {
      type = "Recreate"
    }

    template {
      metadata {
        labels = local.postgres_labels
      }

      spec {
        container {
          name              = "postgres"
          image             = var.postgres_image
          image_pull_policy = "IfNotPresent"

          port {
            name           = "postgres"
            container_port = tonumber(var.db_port)
            protocol       = "TCP"
          }

          env {
            name = "POSTGRES_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret_v1.oficina.metadata[0].name
                key  = "DB_USER"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret_v1.oficina.metadata[0].name
                key  = "DB_PASSWORD"
              }
            }
          }

          env {
            name  = "POSTGRES_DB"
            value = var.db_name
          }

          env {
            name  = "PGDATA"
            value = "/var/lib/postgresql/data/pgdata"
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }

          readiness_probe {
            exec {
              command = ["pg_isready", "-U", var.db_user, "-d", var.db_name]
            }
            initial_delay_seconds = 5
            period_seconds        = 5
            timeout_seconds       = 3
            failure_threshold     = 12
          }

          liveness_probe {
            exec {
              command = ["pg_isready", "-U", var.db_user, "-d", var.db_name]
            }
            initial_delay_seconds = 20
            period_seconds        = 15
            timeout_seconds       = 3
            failure_threshold     = 3
          }

          volume_mount {
            name       = "data"
            mount_path = "/var/lib/postgresql/data"
          }
        }

        volume {
          name = "data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim_v1.postgres.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "postgres" {
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.postgres_labels
  }

  spec {
    type = "ClusterIP"

    selector = {
      "app.kubernetes.io/name" = "postgres"
    }

    port {
      name        = "postgres"
      port        = tonumber(var.db_port)
      target_port = "postgres"
      protocol    = "TCP"
    }
  }
}

resource "kubernetes_deployment_v1" "api" {
  metadata {
    name      = var.app_name
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.labels
  }

  spec {
    replicas = var.api_replicas

    selector {
      match_labels = {
        "app.kubernetes.io/name" = var.app_name
      }
    }

    strategy {
      type = "RollingUpdate"
      rolling_update {
        max_unavailable = "0"
        max_surge       = "1"
      }
    }

    template {
      metadata {
        labels = local.labels
      }

      spec {
        termination_grace_period_seconds = 15

        security_context {
          run_as_non_root = true
          run_as_user     = 10001
          run_as_group    = 10001
          fs_group        = 10001
          seccomp_profile {
            type = "RuntimeDefault"
          }
        }

        container {
          name              = "api"
          image             = var.api_image
          image_pull_policy = "IfNotPresent"

          port {
            name           = "http"
            container_port = var.api_container_port
            protocol       = "TCP"
          }

          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.oficina.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret_v1.oficina.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = var.api_cpu_request
              memory = var.api_memory_request
            }
            limits = {
              cpu    = var.api_cpu_limit
              memory = var.api_memory_limit
            }
          }

          startup_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            initial_delay_seconds = 5
            period_seconds        = 5
            timeout_seconds       = 3
            failure_threshold     = 24
            success_threshold     = 1
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            period_seconds    = 15
            timeout_seconds   = 3
            failure_threshold = 3
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            period_seconds    = 10
            timeout_seconds   = 3
            failure_threshold = 3
          }

          security_context {
            allow_privilege_escalation = false
            read_only_root_filesystem  = true
            capabilities {
              drop = ["ALL"]
            }
          }

          volume_mount {
            name       = "tmp"
            mount_path = "/tmp"
          }
        }

        volume {
          name = "tmp"
          empty_dir {}
        }
      }
    }
  }

  depends_on = [kubernetes_service_v1.postgres]

  lifecycle {
    ignore_changes = [spec[0].replicas]
  }
}

resource "kubernetes_service_v1" "api" {
  metadata {
    name      = var.app_name
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.labels
  }

  spec {
    type = "NodePort"

    selector = {
      "app.kubernetes.io/name" = var.app_name
    }

    port {
      name        = "http"
      port        = var.api_service_port
      target_port = "http"
      node_port   = var.api_node_port
      protocol    = "TCP"
    }
  }
}

resource "kubernetes_horizontal_pod_autoscaler_v2" "api" {
  metadata {
    name      = var.app_name
    namespace = kubernetes_namespace_v1.oficina.metadata[0].name
    labels    = local.labels
  }

  spec {
    min_replicas = var.hpa_min_replicas
    max_replicas = var.hpa_max_replicas

    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.api.metadata[0].name
    }

    metric {
      type = "Resource"
      resource {
        name = "cpu"
        target {
          type                = "Utilization"
          average_utilization = var.hpa_cpu_target
        }
      }
    }

    metric {
      type = "Resource"
      resource {
        name = "memory"
        target {
          type                = "Utilization"
          average_utilization = var.hpa_memory_target
        }
      }
    }

    behavior {
      scale_up {
        stabilization_window_seconds = 30
        select_policy                = "Max"
        policy {
          type           = "Pods"
          value          = 2
          period_seconds = 60
        }
      }

      scale_down {
        stabilization_window_seconds = 120
        select_policy                = "Min"
        policy {
          type           = "Pods"
          value          = 1
          period_seconds = 60
        }
      }
    }
  }
}
