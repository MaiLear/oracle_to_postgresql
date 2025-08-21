package email

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	cockroachdbErrors "github.com/cockroachdb/errors"
	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/ports/out/services"
	"gopkg.in/mail.v2"
)

// SmtpService implementa el servicio de email usando SMTP
type SmtpService struct {
	dialer       *mail.Dialer
	fromEmail    string
	fromName     string
	templatePath string
}

// NewSmtpService crea una nueva instancia del servicio SMTP
func NewSmtpService() (*SmtpService, error) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, cockroachdbErrors.Wrap(err, "puerto SMTP inválido")
	}

	// Configuración del dialer
	d := mail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// Configuración según el ambiente y puerto
	environment := os.Getenv("ENVIRONMENT")
	if environment == "production" {
		// En producción no se usa autenticación
		d.SSL = false
		d.TLSConfig = nil
		d.Auth = nil
		d.StartTLSPolicy = mail.NoStartTLS
		fmt.Printf("📧 SMTP configurado para PRODUCCIÓN (sin autenticación)\n")
	} else {
		// En desarrollo, configurar según el puerto
		if port == 465 {
			// Puerto 465 requiere SSL
			d.SSL = true
			d.StartTLSPolicy = mail.NoStartTLS
			fmt.Printf("📧 SMTP configurado para DESARROLLO (puerto 465 con SSL)\n")
		} else {
			// Puerto 587 usa STARTTLS
			d.SSL = false
			d.StartTLSPolicy = mail.OpportunisticStartTLS
			fmt.Printf("📧 SMTP configurado para DESARROLLO (puerto %d con STARTTLS)\n", port)
		}
		d.Timeout = 30 * time.Second
	}

	fromEmail := os.Getenv("SMTP_FROM")
	if fromEmail == "" {
		return nil, cockroachdbErrors.New("SMTP_FROM no está configurado")
	}

	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "SENA Sofia Plus" // Valor por defecto
	}

	return &SmtpService{
		dialer:       d,
		fromEmail:    fromEmail,
		fromName:     fromName,
		templatePath: os.Getenv("EMAIL_TEMPLATE_PATH"),
	}, nil
}

// SendNotification implementa el envío de notificaciones genéricas
func (s *SmtpService) SendNotification(ctx context.Context, notification services.EmailNotification) error {
	// Determinar el email de destino según el ambiente
	destinationEmail := s.getDestinationEmail(notification.To)

	fmt.Printf("📧 Enviando notificación a: %s\n", destinationEmail)
	fmt.Printf("📧 Asunto: %s\n", notification.Subject)
	fmt.Printf("📧 Template: %s\n", notification.TemplateName)

	// Generar HTML desde template
	htmlBody, err := s.generateHTMLFromTemplate(notification.TemplateName, notification.Data)
	if err != nil {
		return cockroachdbErrors.Wrap(err, "error generando HTML desde template")
	}

	// Crear y enviar mensaje
	message := mail.NewMessage()
	message.SetHeader("From", message.FormatAddress(s.fromEmail, s.fromName))
	message.SetHeader("To", destinationEmail)
	message.SetHeader("Subject", notification.Subject)
	message.SetBody("text/html", htmlBody)

	// Enviar con timeout y manejo de errores detallado
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fmt.Printf("🔗 Intentando conexión SMTP a %s:%d...\n", s.dialer.Host, s.dialer.Port)

	if err := s.dialer.DialAndSend(message); err != nil {
		fmt.Printf("❌ Error detallado enviando email: %v\n", err)
		fmt.Printf("🔧 Configuración SMTP: Host=%s, Port=%d, User=%s\n",
			s.dialer.Host, s.dialer.Port, s.dialer.Username)
		return cockroachdbErrors.Wrap(err, "error enviando email por SMTP")
	}

	fmt.Printf("✅ Email enviado exitosamente a: %s\n", destinationEmail)
	return nil
}

// SendEnrollmentConfirmation envía confirmación de inscripción
func (s *SmtpService) SendEnrollmentConfirmation(ctx context.Context, userEmail, userName, courseName, formationLevel string, enrollmentId int) error {
	fmt.Printf("📧 Enviando confirmación de inscripción a: %s\n", userEmail)

	notification := services.EmailNotification{
		To:           userEmail,
		Subject:      "SENA - Confirmación de Inscripción programa SOFIA Plus – Betowa",
		TemplateName: "enrollment_confirmation",
		Data: map[string]interface{}{
			"userName":       userName,
			"courseName":     courseName,
			"formationLevel": formationLevel,
			"enrollmentId":   enrollmentId,
			"date":           s.formatSpanishDate(time.Now()),
			"currentYear":    time.Now().Year(),
		},
	}

	return s.SendNotification(ctx, notification)
}

// SendRegistrationConfirmation envía confirmación de registro
func (s *SmtpService) SendRegistrationConfirmation(ctx context.Context, userEmail, userName, documentType, documentNumber string, nis int) error {
	fmt.Printf("📧 Enviando confirmación de registro a: %s\n", userEmail)

	notification := services.EmailNotification{
		To:           userEmail,
		Subject:      "Bienvenido a SENA SOFIA Plus y Betowa - Registro Exitoso",
		TemplateName: "registration_confirmation",
		Data: map[string]interface{}{
			"userName":       userName,
			"documentType":   documentType,
			"documentNumber": documentNumber,
			"nis":            nis,
			"date":           s.formatSpanishDate(time.Now()),
			"currentYear":    time.Now().Year(),
		},
	}

	return s.SendNotification(ctx, notification)
}

// SendRegistrationConfirmation envía confirmación de registro
func (s *SmtpService) SendEnrollmentError(ctx context.Context, userEmail, userName, courseName, formationLevel string, courseId int,errorMessage string) error {
	fmt.Printf("📧 Enviando confirmación de registro a: %s\n", userEmail)

	notification := services.EmailNotification{
		To:           userEmail,
		Subject:      "Error de Inscripción programa SOFIA Plus – Betowa",
		TemplateName: "enrollment_error",
		Data: map[string]any{
			"userName":       userName,
			"courseName":     courseName,
			"formationLevel": formationLevel,
			"courseId":   courseId,
			"errorMessage": errorMessage,
			"date":           s.formatSpanishDate(time.Now()),
			"currentYear":    time.Now().Year(),
		},
	}

	return s.SendNotification(ctx, notification)
}



// generateHTMLFromTemplate genera el HTML usando templates
func (s *SmtpService) generateHTMLFromTemplate(templateName string, data map[string]interface{}) (string, error) {
	switch templateName {
	case "enrollment_confirmation":
		return s.enrollmentConfirmationTemplate(data)
	case "registration_confirmation":
		return s.registrationConfirmationTemplate(data)
	case "enrollment_error":
		return s.enrollmentErrorTemplate(data)
	default:

		return "", fmt.Errorf("template no encontrado: %s", templateName)
	}
}

// formatSpanishDate formatea la fecha en español con zona horaria de Colombia
func (s *SmtpService) formatSpanishDate(t time.Time) string {
	months := []string{
		"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
	}

	// Convertir a zona horaria de Colombia (UTC-5)
	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		// Si falla, usar UTC-5 manualmente
		location = time.FixedZone("COT", -5*60*60) // Colombia Time (UTC-5)
	}

	colombiaTime := t.In(location)

	return fmt.Sprintf("%d de %s de %d a las %s",
		colombiaTime.Day(),
		months[colombiaTime.Month()-1],
		colombiaTime.Year(),
		colombiaTime.Format("15:04"),
	)
}

// getDestinationEmail determina el email de destino según el ambiente
//
// :param originalEmail: email original del usuario
// :return: email de destino según el ambiente
func (s *SmtpService) getDestinationEmail(originalEmail string) string {
	environment := os.Getenv("ENVIRONMENT")

	// En ambiente de desarrollo, siempre enviar al email configurado en SMTP_TO
	if environment != "production" {
		smtpTo := os.Getenv("SMTP_TO")
		if smtpTo == "" {
			smtpTo = "betowa@cristhiancano.com" // Valor por defecto si SMTP_TO no está configurado
		}
		fmt.Printf("🧪 AMBIENTE DE DESARROLLO: Redirigiendo email de %s a %s\n", originalEmail, smtpTo)
		return smtpTo
	}

	// En producción, enviar al email original del usuario
	return originalEmail
}
