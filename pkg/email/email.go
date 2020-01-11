package email

import (
	"fmt"

	"github.com/novaladip/geldstroom-api-go/pkg/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailVerification(receiver, token string) error {
	url := fmt.Sprintf("https://geldstroom.cotcapp.my.id/user/verify/email/%v", token)
	htmlContent := fmt.Sprintf(`<strong>To verify your email click <a href="%v">here</a></strong>`, url)
	plainTextContent := fmt.Sprintf("To Verify your email click this link %v", url)
	key := config.ConfigKey.SENDGRID_KEY
	from := mail.NewEmail("Geldstroom", "no-reply@cotcapp.my.id")
	subject := "Verify Email Address"
	to := mail.NewEmail(receiver, receiver)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(key)
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
