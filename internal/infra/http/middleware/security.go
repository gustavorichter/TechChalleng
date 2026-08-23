package middleware

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/techchalleng/oficina/internal/application/dto"
)

var (
	sqlInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union\s+select|insert\s+into|drop\s+table|delete\s+from|update\s+.+\s+set|;\s*--|'\s*or\s*'|'\s*or\s*1\s*=\s*1)`),
		regexp.MustCompile(`(?i)(exec\s*\(|xp_cmdshell|information_schema|pg_sleep\s*\()`),
	}

	commandInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(;\s*(rm|cat|wget|curl|bash|sh|cmd|powershell)\b)`),
		regexp.MustCompile(`(?i)(\|\||&&|\$\(|` + "`" + `)`),
		regexp.MustCompile(`(?i)(>\s*/dev/|/etc/passwd|\\\\windows\\\\system32)`),
	}
)

// InjectionGuard bloqueia padrões típicos de SQL Injection e Command Injection
// em query strings e parâmetros de rota antes de chegar aos handlers.
func InjectionGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if containsMaliciousInput(c.Request.URL.RawQuery) {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "entrada inválida: padrão suspeito detectado na requisição",
			})
			return
		}

		for _, param := range c.Params {
			if containsMaliciousInput(param.Value) {
				c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
					Error: "entrada inválida: parâmetro de rota rejeitado",
				})
				return
			}
		}

		c.Next()
	}
}

func containsMaliciousInput(value string) bool {
	if value == "" {
		return false
	}

	decoded, err := url.QueryUnescape(value)
	if err != nil {
		decoded = value
	}
	decoded = strings.ToLower(strings.TrimSpace(decoded))
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(decoded) {
			return true
		}
	}
	for _, pattern := range commandInjectionPatterns {
		if pattern.MatchString(decoded) {
			return true
		}
	}
	return false
}
