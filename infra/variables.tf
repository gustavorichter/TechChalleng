variable "kubeconfig_path" {
  description = "Caminho do kubeconfig do cluster local. O padrão é o arquivo de contexto do usuário."
  type        = string
  default     = "~/.kube/config"
}

variable "kube_context" {
  description = "Contexto do kubeconfig. Vazio usa o current-context."
  type        = string
  default     = ""
}

variable "namespace" {
  description = "Namespace Kubernetes da oficina."
  type        = string
  default     = "oficina"
}

variable "app_name" {
  description = "Nome dos recursos da API (Deployment, Service, HPA)."
  type        = string
  default     = "oficina-api"
}

variable "api_image" {
  description = "Imagem local da API. Precisa existir no runtime do cluster (docker build -t oficina-api:local .)."
  type        = string
  default     = "oficina-api:local"
}

variable "api_container_port" {
  description = "Porta HTTP do processo da API dentro do container."
  type        = number
  default     = 8080
}

variable "api_service_port" {
  description = "Porta do Service da API."
  type        = number
  default     = 80
}

variable "api_node_port" {
  description = "NodePort exposta no cluster local (30000-32767)."
  type        = number
  default     = 30080

  validation {
    condition     = var.api_node_port >= 30000 && var.api_node_port <= 32767
    error_message = "api_node_port precisa estar na faixa de NodePort (30000-32767)."
  }
}

variable "api_replicas" {
  description = "Réplicas iniciais da API. O HPA assume o controle depois do apply (ignore_changes em replicas)."
  type        = number
  default     = 2
}

variable "api_cpu_request" {
  description = "CPU request do container da API."
  type        = string
  default     = "50m"
}

variable "api_cpu_limit" {
  description = "CPU limit do container da API."
  type        = string
  default     = "250m"
}

variable "api_memory_request" {
  description = "Memory request do container da API."
  type        = string
  default     = "64Mi"
}

variable "api_memory_limit" {
  description = "Memory limit do container da API."
  type        = string
  default     = "128Mi"
}

variable "hpa_min_replicas" {
  description = "Mínimo de réplicas do HPA."
  type        = number
  default     = 2
}

variable "hpa_max_replicas" {
  description = "Máximo de réplicas do HPA."
  type        = number
  default     = 6
}

variable "hpa_cpu_target" {
  description = "Utilização média de CPU (%) que dispara o scale up."
  type        = number
  default     = 70
}

variable "hpa_memory_target" {
  description = "Utilização média de memória (%) que dispara o scale up."
  type        = number
  default     = 80
}

variable "db_host" {
  description = "Host do PostgreSQL visto pelos pods. O padrão é o Service postgres no mesmo namespace."
  type        = string
  default     = "postgres"
}

variable "db_port" {
  description = "Porta do PostgreSQL."
  type        = string
  default     = "5432"
}

variable "db_name" {
  description = "Nome do banco."
  type        = string
  default     = "oficina_db"
}

variable "db_sslmode" {
  description = "Modo SSL do cliente PostgreSQL. Em cluster local o padrão é disable."
  type        = string
  default     = "disable"
}

variable "db_user" {
  description = "Usuário do PostgreSQL. Placeholder de desenvolvimento, igual ao .env.example."
  type        = string
  sensitive   = true
  default     = "oficina"
}

variable "db_password" {
  description = "Senha do PostgreSQL. Placeholder de desenvolvimento, igual ao .env.example."
  type        = string
  sensitive   = true
  default     = "oficina_secret"
}

variable "jwt_secret" {
  description = "Segredo JWT da API. Placeholder de desenvolvimento, igual ao .env.example."
  type        = string
  sensitive   = true
  default     = "super-secret-jwt-key-change-me"
}

variable "admin_username" {
  description = "Usuário administrativo inicial da API."
  type        = string
  sensitive   = true
  default     = "admin"
}

variable "admin_password" {
  description = "Senha administrativa inicial da API. Placeholder de desenvolvimento."
  type        = string
  sensitive   = true
  default     = "admin123"
}

variable "postgres_image" {
  description = "Imagem do PostgreSQL."
  type        = string
  default     = "postgres:16-alpine"
}

variable "postgres_storage" {
  description = "Tamanho do volume persistente do PostgreSQL."
  type        = string
  default     = "1Gi"
}

variable "postgres_storage_class" {
  description = "StorageClass do PVC. Vazio usa a classe padrão do cluster local."
  type        = string
  default     = ""
}
