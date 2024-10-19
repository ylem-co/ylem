package resources

import "embed"

//go:embed templates/html
var EmbeddedFileSystem embed.FS
var EmbeddedHtmlTemplates = map[string]string{
	"confirmation_email": "templates/html/confirmation_email.gohtml",
}
