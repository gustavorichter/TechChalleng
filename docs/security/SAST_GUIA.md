# Guia de Análise Estática de Segurança (SAST) — Tech Challenge

Este guia descreve como executar localmente as ferramentas de segurança exigidas para o back-end monolítico em Go da **Oficina Mecânica (Tech Challenge)**.

## Pré-requisitos

| Ferramenta | Versão mínima | Instalação |
|------------|---------------|------------|
| Go | 1.25+ | [go.dev/dl](https://go.dev/dl/) |
| Docker Desktop | 4.x+ | [docker.com](https://docs.docker.com/get-docker/) |
| Git | 2.x+ | Já instalado no ambiente de desenvolvimento |

Opcional (varredura de imagem):

| Ferramenta | Instalação |
|------------|------------|
| Snyk CLI | `npm install -g snyk` + `snyk auth` |
| Docker Scout | Incluso no Docker Desktop 4.17+ |

---

## 1. gosec — Go Security Checker (SAST)

Analisa o código-fonte Go em busca de vulnerabilidades estáticas: injeções, criptografia fraca, concorrência insegura, configurações HTTP inseguras, etc.

### Instalação

```powershell
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

Certifique-se de que `%USERPROFILE%\go\bin` está no `PATH`.

### Execução

Na raiz do repositório:

```powershell
# Relatório no terminal
gosec ./...

# Relatório JSON (para CI/CD ou anexo do documento)
mkdir docs\security -Force
gosec -fmt json -out docs\security\gosec-report.json ./...

# Relatório HTML (visual)
gosec -fmt html -out docs\security\gosec-report.html ./...
```

### Interpretação

- **Exit code 0**: nenhum issue encontrado (ou apenas suprimidos com `#nosec`).
- **Exit code 1**: issues detectados — revisar severidade (HIGH, MEDIUM, LOW).
- Regras comuns: `G101` (hardcoded credentials), `G112` (Slowloris), `G115` (integer overflow), `G706` (log injection).

---

## 2. govulncheck — Vulnerabilidades em Dependências

Ferramenta **oficial** do time Go. Cruza o grafo de importação do projeto com o [Go Vulnerability Database](https://vuln.go.dev/).

### Instalação

```powershell
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### Execução

```powershell
# Scan padrão (mostra apenas vulns alcançáveis pelo código)
govulncheck ./...

# Detalhes completos (inclui dependências transitivas não usadas)
govulncheck -show verbose ./...

# Relatório JSON
govulncheck -json ./... > docs\security\govulncheck-report.json
```

### Mitigação típica

```powershell
# Atualizar dependências diretas
go get github.com/jackc/pgx/v5@v5.9.2
go get github.com/golang-jwt/jwt/v5@v5.2.2
go get golang.org/x/text@latest
go get golang.org/x/net@latest
go mod tidy
govulncheck ./...
```

Para vulnerabilidades da **stdlib Go**, atualize o toolchain:

```powershell
go install golang.org/dl/go1.26.6@latest
go1.26.6 download
go1.26.6 version
```

---

## 3. Varredura de Imagem Docker (Snyk ou Docker Scout)

Após construir a imagem da API:

```powershell
docker build -t techchalleng-oficina:latest .
```

### Opção A — Docker Scout (recomendado, sem conta extra)

```powershell
# Habilitar Scout no Docker Desktop (Settings > Docker Scout)

docker scout quickview techchalleng-oficina:latest
docker scout cves techchalleng-oficina:latest
docker scout recommendations techchalleng-oficina:latest

# Exportar SBOM
docker scout sbom techchalleng-oficina:latest --output docs\security\sbom.json
```

### Opção B — Snyk CLI

```powershell
snyk auth
snyk container test techchalleng-oficina:latest --file=Dockerfile
snyk container test techchalleng-oficina:latest --json > docs\security\snyk-container-report.json
```

---

## 4. Script automatizado (Windows)

Execute todos os scans de uma vez:

```powershell
.\scripts\security-scan.ps1
```

Os artefatos são gravados em `docs/security/`.

---

## 5. Integração CI/CD (referência)

Exemplo de job GitHub Actions:

```yaml
name: Security Scan
on: [push, pull_request]
jobs:
  sast:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: go install github.com/securego/gosec/v2/cmd/gosec@latest
      - run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - run: gosec ./...
      - run: govulncheck ./...
      - run: docker build -t oficina:ci .
      - run: docker scout cves oficina:ci
```

---

## 6. Checklist de entrega (Fase 1)

- [ ] `gosec ./...` executado e issues documentados
- [ ] `govulncheck ./...` executado e plano de atualização registrado
- [ ] Imagem Docker escaneada (Scout ou Snyk)
- [ ] Relatório técnico preenchido (`docs/security/RELATORIO_VULNERABILIDADES.md`)
- [ ] Mitigações implementadas (middleware, timeouts HTTP, usuário não-root)
