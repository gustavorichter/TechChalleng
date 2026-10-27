# Infraestrutura local (Terraform)

Provisiona a stack da oficina em um cluster Kubernetes **somente local**. O estado fica em `terraform.tfstate` nesta pasta. Não há backend remoto nem provider de nuvem.

O provider usado é o [`hashicorp/kubernetes`](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs). Helm e Docker não entram aqui: os recursos espelham os manifestos de `k8s/` e o PostgreSQL do `docker-compose.yml`, sem chart e sem criar o cluster.

## Recursos

| Recurso | Nome | Função |
|---------|------|--------|
| Namespace | `oficina` | Isola a stack no cluster local |
| ConfigMap | `oficina-config` | `PORT`, host, porta, nome e SSL do banco |
| Secret | `oficina-secret` | Usuário e senha do banco, JWT e admin |
| PVC | `postgres-data` | Volume do PostgreSQL (StorageClass padrão do cluster) |
| Deployment | `postgres` | PostgreSQL 16 (`postgres:16-alpine`), 1 réplica |
| Service | `postgres` | ClusterIP na porta 5432, visível só dentro do namespace |
| Deployment | `oficina-api` | API (`oficina-api:local`), probes em `/health`, root filesystem somente leitura |
| Service | `oficina-api` | NodePort **30080** → porta 80 do Service → 8080 do container |
| HPA | `oficina-api` | 2 a 6 réplicas, CPU 70% e memória 80% |

Os defaults de credencial são os mesmos do `.env.example` e de `k8s/secret.yaml`. Servem só para desenvolvimento local.

O Deployment da API ignora mudanças posteriores em `spec.replicas` para o HPA poder escalar sem o próximo `terraform apply` devolver a contagem inicial.

## Pré-requisitos

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.5
- Cluster Kubernetes local com contexto em `~/.kube/config`
  - Docker Desktop: Settings → Kubernetes → Enable Kubernetes
  - ou minikube, kind, k3d
- `kubectl` apontando para esse contexto (`kubectl config current-context`)
- Imagem da API presente no runtime do cluster:

```powershell
docker build -t oficina-api:local .
```

No Docker Desktop a imagem local fica disponível para o Kubernetes. Em outros runtimes, carregue a imagem antes do apply:

```powershell
# minikube
minikube image load oficina-api:local

# kind (troque o nome do cluster se não for "kind")
kind load docker-image oficina-api:local --name kind
```

O HPA precisa do [metrics-server](https://github.com/kubernetes-sigs/metrics-server). Docker Desktop e `minikube addons enable metrics-server` já entregam isso. Sem o metrics-server o HPA é criado, mas não escala.

## Passo a passo

Na pasta `infra`:

```powershell
cd infra

terraform init
terraform plan
terraform apply
```

`terraform init` baixa o provider Kubernetes e inicializa o backend local (`terraform.tfstate`).

`terraform plan` mostra o que será criado no contexto atual do kubeconfig. Nada é aplicado ainda.

`terraform apply` pede confirmação (`yes`) e cria os recursos. Para pular o prompt:

```powershell
terraform apply -auto-approve
```

Conferir o status:

```powershell
terraform output
kubectl -n oficina get pods,svc,hpa,pvc
kubectl -n oficina rollout status deployment/oficina-api
```

Health check (Docker Desktop e kind com acesso ao NodePort no localhost):

```powershell
curl http://localhost:30080/health
```

Resposta esperada: `{"status":"ok"}`.

No minikube o NodePort não escuta em `localhost` do host. Use:

```powershell
minikube service oficina-api -n oficina --url
```

Swagger, com a API no ar: `http://localhost:30080/swagger/index.html`.

## Parametrizar

Os defaults estão em `variables.tf`. Para outro contexto, outra NodePort ou outras credenciais, crie `infra/terraform.tfvars` (este arquivo está no `.gitignore`):

```hcl
kube_context  = "docker-desktop"
namespace     = "oficina"
api_node_port = 30080
db_password   = "oficina_secret"
jwt_secret    = "super-secret-jwt-key-change-me"
```

Ou passe na linha de comando:

```powershell
terraform apply -var "kube_context=docker-desktop" -var "api_node_port=30081"
```

`kubeconfig_path` padrão é `~/.kube/config`. `kube_context` vazio usa o `current-context`.

## Saídas

Depois do apply, `terraform output` mostra namespace, nome do Deployment, NodePort, URL da API, faixa do HPA e o DNS interno do PostgreSQL. Valores de Secret não são impressos.

## Destruir

```powershell
terraform destroy
```

O PVC é removido junto com o namespace. Os dados do banco local deixam de existir.

## Estado

O arquivo `infra/terraform.tfstate` contém o plano aplicado e os valores sensíveis do Secret. Ele não deve ir para o Git (já coberto pelo `.gitignore`). Não há lock remoto: evite dois `terraform apply` ao mesmo tempo nesta pasta.
