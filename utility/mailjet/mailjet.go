package mailjet

import (
	"os"

	"github.com/mailjet/mailjet-apiv3-go"
)

func CreateMailjet(toMail string, templateID int32, variables map[string]interface{}) error {
	fromMail := os.Getenv("MAILJET_EMAIL")
	fromUsername := os.Getenv("MAILJET_USERNAME")
	publicKey := os.Getenv("MAILJET_PUBLIC_KEY")
	secretKey := os.Getenv("MAILJET_SECRET_KEY")

	mailjetClient := mailjet.NewMailjetClient(publicKey, secretKey)
	messageInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromMail,
				Name:  fromUsername,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: toMail,
				},
			},
			TemplateID:       templateID,
			TemplateLanguage: true,
			Variables:        variables,
		},
	}

	messages := mailjet.MessagesV31{Info: messageInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}

	return nil
}
