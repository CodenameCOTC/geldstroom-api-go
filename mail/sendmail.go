package sendmail

import (
	"fmt"

	"github.com/novaladip/geldstroom-api-go/config"
	"github.com/novaladip/geldstroom-api-go/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Send(receiver, body string) {
	key := config.GetKey()

	from := mail.NewEmail("Geldstroom", "no-reply@geldstroom.my.id")
	subject := "New Login"
	to := mail.NewEmail(receiver, receiver)
	message := mail.NewSingleEmail(from, subject, to, body, body)
	client := sendgrid.NewSendClient(key.SENDGRID_KEY)
	response, err := client.Send(message)
	if err != nil {
		logger.ErrorLog.Println(err)
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
