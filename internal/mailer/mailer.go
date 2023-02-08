package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	mg     *mailgun.MailgunImpl
	sender string
}

func New(yourDomain, privateAPIKey, sender string) Mailer {

	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	return Mailer{
		mg:     mg,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	body := htmlBody.String()

	subject := "Lebena Api Registration"

	message := m.mg.NewMessage(m.sender, subject, "", recipient)

	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i := 1; i <= 3; i++ {
		resp, id, err := m.mg.Send(ctx, message)

		if nil == err {
			fmt.Printf("ID: %s Resp: %s\n", id, resp)
			return nil
		}

		time.Sleep(2000 * time.Millisecond)
	}

	return err

}
