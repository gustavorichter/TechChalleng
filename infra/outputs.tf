output "namespace" {
  description = "Namespace em que a stack foi criada."
  value       = kubernetes_namespace_v1.oficina.metadata[0].name
}

output "api_deployment" {
  description = "Nome do Deployment da API."
  value       = kubernetes_deployment_v1.api.metadata[0].name
}

output "api_replicas_declaradas" {
  description = "Réplicas declaradas no primeiro apply. Depois disso o HPA controla a escala."
  value       = var.api_replicas
}

output "api_service" {
  description = "Nome do Service NodePort da API."
  value       = kubernetes_service_v1.api.metadata[0].name
}

output "api_node_port" {
  description = "NodePort da API no cluster local."
  value       = kubernetes_service_v1.api.spec[0].port[0].node_port
}

output "api_endpoint" {
  description = "URL da API no Docker Desktop / kind com port-map. No minikube use minikube service."
  value       = "http://localhost:${kubernetes_service_v1.api.spec[0].port[0].node_port}"
}

output "api_health_endpoint" {
  description = "Health check da API via NodePort."
  value       = "http://localhost:${kubernetes_service_v1.api.spec[0].port[0].node_port}/health"
}

output "hpa_name" {
  description = "Nome do HorizontalPodAutoscaler da API."
  value       = kubernetes_horizontal_pod_autoscaler_v2.api.metadata[0].name
}

output "hpa_range" {
  description = "Faixa de réplicas do HPA (mínimo-máximo)."
  value       = "${var.hpa_min_replicas}-${var.hpa_max_replicas}"
}

output "postgres_service" {
  description = "DNS interno do PostgreSQL (ClusterIP, sem exposição no host)."
  value       = "${kubernetes_service_v1.postgres.metadata[0].name}.${kubernetes_namespace_v1.oficina.metadata[0].name}.svc.cluster.local:${var.db_port}"
}

output "postgres_pvc" {
  description = "PersistentVolumeClaim do PostgreSQL."
  value       = kubernetes_persistent_volume_claim_v1.postgres.metadata[0].name
}

output "config_map" {
  description = "ConfigMap com a configuração não sensível da API."
  value       = kubernetes_config_map_v1.oficina.metadata[0].name
}

output "secret" {
  description = "Secret com credenciais. O conteúdo não é exibido nos outputs."
  value       = kubernetes_secret_v1.oficina.metadata[0].name
}
