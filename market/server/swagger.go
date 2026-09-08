package server

import (
	"github.com/gofiber/fiber/v3"
	"market/api"
)

// swaggerHTML страница Swagger UI, ассеты тянутся с CDN, спека отдаётся нашей ручкой
const swaggerHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <title>Restaraunt</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/swagger/openapi.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`

// setupSwagger вешает Swagger UI и раздачу OpenAPI-спеки
func (s *Server) setupSwagger() {
	s.app.Get("/swagger/openapi.yaml", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/yaml")

		return c.Send(api.Spec)
	})

	s.app.Get("/swagger", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

		return c.SendString(swaggerHTML)
	})
}
