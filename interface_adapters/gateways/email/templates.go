package email

import (
	"fmt"
	"strings"
)

// enrollmentConfirmationTemplate genera el HTML para confirmación de inscripción
func (s *SmtpService) enrollmentConfirmationTemplate(data map[string]interface{}) (string, error) {
	userName := data["userName"].(string)
	courseName := data["courseName"].(string)
	formationLevel := data["formationLevel"].(string)
	enrollmentId := data["enrollmentId"].(int)
	date := data["date"].(string)
	currentYear := data["currentYear"].(int)

	html := fmt.Sprintf(`<body style="font-family: 'Arial', sans-serif; background-color: #f4f4f4; margin: 0; padding: 0;">
		<div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
			<div style="text-align: center;">
				<img src="https://betowa.sena.edu.co/assets/banner/banner-email.jpg" 
					alt="Logo SOFIA Plus" 
					style="width: 100%%; height: auto; display: block;"> 
			</div>
			<div style="padding: 30px 20px;">

				<h1 style="color: #333333; font-size: 24px; margin-bottom: 20px; text-align: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px;">✅ Confirmación de Inscripción</h1>
				<div style="color: #555555; font-size: 16px; line-height: 1.6;">
					<p>Hola <strong>%s</strong>,</p>

					<p>¡Excelente noticia! Tu inscripción ha sido realizada de manera correcta y exitosa.</p>

					<div style="background-color: #e8f5e8; border-left: 4px solid #28a745; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h3 style="color: #28a745; margin: 0 0 10px 0;">📚 Detalles de tu inscripción:</h3>
						<p style="margin: 5px 0;"><strong>Programa:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Nivel de Formación:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>ID de Inscripción:</strong> #%d</p>
						<p style="margin: 5px 0;"><strong>Fecha de inscripción:</strong> %s</p>
					</div>

					<p>Tu inscripción ha sido procesada y registrada en nuestro sistema. En los próximos días recibirás más información sobre el inicio del programa de formación.</p>

					<div style="background-color: #e3f2fd; border-left: 4px solid #2196f3; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h4 style="color: #1976d2; margin: 0 0 10px 0;">⏱️ Importante - Sincronización:</h4>
						<p style="margin: 5px 0; color: #1976d2; font-weight: bold;">Tu inscripción puede tardar hasta 30 minutos en aparecer en el portal de SOFIA Plus una vez recibida esta notificación.</p>
					</div>

					<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h4 style="color: #856404; margin: 0 0 10px 0;">📋 Próximos pasos:</h4>
						<ul style="margin: 10px 0; padding-left: 20px; color: #856404;">
							<li>Mantente atento a tu correo electrónico para comunicaciones importantes</li>
							<li>Ingresa a tu cuenta en SOFIA Plus para seguir el estado de tu inscripción</li>
							<li>Si no ves tu inscripción, espera hasta 30 minutos para que se ejecute la sincronización</li>
							<li>Prepárate para una experiencia de aprendizaje excepcional</li>
						</ul>
					</div>

					<p style="text-align: center; margin: 30px 0;">
						<a href="http://senasofiaplus.edu.co/sofia-public/" 
						   style="background-color: #28a745; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
							🌐 Acceder a SOFIA Plus
						</a>
					</p>

					<p>Este correo ha sido generado automáticamente.<br/>
					Por favor, no responder.</p>
			</div>
		</div>

			<div style="background-color: #f9f9f9; color: #777777; text-align: center; padding: 15px; font-size: 14px; border-top: 1px solid #eeeeee;">
				<p style="margin: 0;">&copy; %d Betowa - SENA. Todos los derechos reservados.</p>
				<p style="margin: 10px 0 0;">
					<a href="https://portal.senasofiaplus.edu.co/index.php/seguridad/politica-de-confidencialidad" style="color: #555555; text-decoration: none; margin: 0 10px;">Política de Seguridad y Confidencialidad</a>
				</p>
			</div>
		</div>
	</body>`, userName, courseName, formationLevel, enrollmentId, date, currentYear)

	return html, nil
}

func (s *SmtpService) enrollmentErrorTemplate(data map[string]interface{}) (string, error) {
	userName := data["userName"].(string)
	courseName := data["courseName"].(string)
	formationLevel := data["formationLevel"].(string)
	courseId := data["courseId"].(int) // Ahora es el ID del curso
	date := data["date"].(string)
	currentYear := data["currentYear"].(int)
	errorMessage := data["errorMessage"].(string) // Mensaje de error

	html := fmt.Sprintf(`<body style="font-family: 'Arial', sans-serif; background-color: #f4f4f4; margin: 0; padding: 0;">
		<div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
			<div style="text-align: center;">
				<img src="https://betowa.sena.edu.co/assets/banner/banner-email.jpg" 
					alt="Logo SOFIA Plus" 
					style="width: 100%%; height: auto; display: block;"> 
			</div>
			<div style="padding: 30px 20px;">

				<h1 style="color: #d9534f; font-size: 24px; margin-bottom: 20px; text-align: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px;">❌ Error en tu Inscripción</h1>
				<div style="color: #555555; font-size: 16px; line-height: 1.6;">
					<p>Hola <strong>%s</strong>,</p>

					<p>Lamentablemente, tu inscripción al curso no pudo completarse debido al siguiente inconveniente:</p>

					<div style="background-color: #fdecea; border-left: 4px solid #d9534f; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h3 style="color: #d9534f; margin: 0 0 10px 0;">⚠️ Detalles del Error:</h3>
						<p style="margin: 5px 0; font-weight: bold;">%s</p>
					</div>

					<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h4 style="color: #856404; margin: 0 0 10px 0;">📚 Detalles de tu intento de inscripción:</h4>
						<p style="margin: 5px 0;"><strong>Programa:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Nivel de Formación:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Ficha:</strong> #%d</p>
						<p style="margin: 5px 0;"><strong>Fecha de intento:</strong> %s</p>
					</div>

					<p>Te recomendamos verificar la información de tu inscripción y volver a intentarlo. Si el problema persiste, contacta con soporte.</p>

				<p style="text-align: center; margin: 30px 0;">
					<a href="https://betowa.sena.edu.co/" 
						style="background-color: #28a745; 
						color: white; 
						padding: 14px 35px; 
						text-decoration: none; 
						border-radius: 8px; 
						display: inline-block; 
						font-weight: bold; 
						font-size: 16px; 
						box-shadow: 0px 4px 6px rgba(0,0,0,0.1); 
						transition: background-color 0.3s ease;">
						📚 ¡Explora e inscríbete a otro curso en la plataforma!
					</a>
				</p>

					<p>Este correo ha sido generado automáticamente.<br/>
					Por favor, no responder.</p>
			</div>
		</div>

			<div style="background-color: #f9f9f9; color: #777777; text-align: center; padding: 15px; font-size: 14px; border-top: 1px solid #eeeeee;">
				<p style="margin: 0;">&copy; %d Betowa - SENA. Todos los derechos reservados.</p>
				<p style="margin: 10px 0 0;">
					<a href="https://portal.senasofiaplus.edu.co/index.php/seguridad/politica-de-confidencialidad" style="color: #555555; text-decoration: none; margin: 0 10px;">Política de Seguridad y Confidencialidad</a>
				</p>
			</div>
		</div>
	</body>`, 
		userName, errorMessage, courseName, formationLevel, courseId, date, currentYear)

	return html, nil
}


// registrationConfirmationTemplate genera el HTML para confirmación de registro
func (s *SmtpService) registrationConfirmationTemplate(data map[string]interface{}) (string, error) {
	userName := strings.ToUpper(data["userName"].(string))
	documentType := data["documentType"].(string)
	documentNumber := data["documentNumber"].(string)
	date := data["date"].(string)
	currentYear := data["currentYear"].(int)

	html := fmt.Sprintf(`<body style="font-family: 'Arial', sans-serif; background-color: #f4f4f4; margin: 0; padding: 0;">
		<div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
			<div style="text-align: center;">
				<img src="https://betowa.sena.edu.co/assets/banner/banner-email.jpg" 
					alt="Logo SOFIA PLUS" 
					style="width: 220px; height: auto; display: block; margin: 0 auto 20px auto;">
			</div>
			<div style="padding: 30px 20px;">

				<h1 style="color: #333333; font-size: 24px; margin-bottom: 20px; text-align: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px;">Bienvenido a SENA SOFIA PLUS y Betowa – Registro Exitoso</h1>
				<div style="color: #555555; font-size: 16px; line-height: 1.6;">
					<p>Hola <strong>%s</strong>,</p>

					<p>¡Felicitaciones! Tu registro en SENA SOFIA PLUS ha sido completado exitosamente.</p>

					<div style="background-color: #e3f2fd; border-left: 4px solid #2196f3; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h3 style="color: #1976d2; margin: 0 0 10px 0;">👤 Información de tu cuenta:</h3>
						<p style="margin: 5px 0;"><strong>Usuario:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Tipo de documento:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Número de documento:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Fecha de registro:</strong> %s</p>
					</div>

					<div style="background-color: #e3f2fd; border-left: 4px solid #2196f3; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h4 style="color: #1976d2; margin: 0 0 10px 0;">⏱️ Importante - Sincronización:</h4>
						<p style="margin: 5px 0; color: #1976d2; font-weight: bold;">Tu cuenta puede tardar hasta 30 minutos en estar disponible en el sistema académico – administrativo SOFIA Plus una vez recibida esta notificación.</p>
					</div>

					<div style="background-color: #f3e5f5; border-left: 4px solid #9c27b0; padding: 15px; margin: 20px 0; border-radius: 4px;">
						<h4 style="color: #7b1fa2; margin: 0 0 10px 0;">🚀 ¿Qué puedes hacer ahora?</h4>
						<ul style="margin: 10px 0; padding-left: 20px; color: #7b1fa2;">
							<li>Explorar nuestra oferta de programas de formación</li>
						</ul>
					</div>

					<p style="text-align: center; margin: 30px 0;">
						<a href="https://betowa.sena.edu.co" 
						   style="background-color: #2196f3; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
							🚪 Ingresar a Mi Cuenta
						</a>
					</p>

					<p>Este correo ha sido generado automáticamente.<br/>
					Por favor, no responder.</p>
			</div>
		</div>

			<div style="background-color: #f9f9f9; color: #777777; text-align: center; padding: 15px; font-size: 14px; border-top: 1px solid #eeeeee;">
				<p style="margin: 0;">&copy; %d Betowa - SENA. Todos los derechos reservados.</p>
				<p style="margin: 10px 0 0;">
					<a href="https://portal.senasofiaplus.edu.co/index.php/seguridad/politica-de-confidencialidad" style="color: #555555; text-decoration: none; margin: 0 10px;">Política de Seguridad y Confidencialidad</a>
				</p>
			</div>
		</div>
	</body>`, userName, userName, documentType, documentNumber, date, currentYear)

	return html, nil
}
