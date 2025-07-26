package utils

import (
	"fmt"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
)

func _(nome, errMsg string) {
	smtpCfg := config.Current.SMTP
	if !smtpCfg.Enabled {
		return
	}

	email := EmailConfig{
		SMTPHost: smtpCfg.SMTPHost,
		SMTPPort: smtpCfg.SMTPPort,
		Username: smtpCfg.Username,
		Password: smtpCfg.Password,
		From:     smtpCfg.From,
		To:       smtpCfg.To,
	}

	subject := fmt.Sprintf("Teste de e-mail Agente %s", nome)
	body := fmt.Sprintf("E-mail enviado '%s':\n\n%s", nome, errMsg)

	if err := SendEmail(email, subject, body); err != nil {
		logger.Log.Errorf("Erro ao enviar e-mail: %v", err)
	}
}
