package mailjet

import (
	"fmt"
	"os"

	"github.com/mailjet/mailjet-apiv3-go"
)

func createMailjet(toMail string, templateID int32, variables map[string]interface{}) error {
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

func SendEmailActivation(email, token string) error {
	baseURL := os.Getenv("FRONTEND_BASE_URL")

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%sactivation?key=%s", baseURL, token),
	}

	err := createMailjet(email, 5188578, variables)
	if err != nil {
		return err
	}

	return nil
}

func SendEmailPasswordReset(email, name, token string) error {
	baseURL := os.Getenv("FRONTEND_BASE_URL")

	variables := map[string]interface{}{
		"name": name,
		"link": fmt.Sprintf("%sverification?key=%s", baseURL, token),
	}

	err := createMailjet(email, 5202383, variables)
	if err != nil {
		return err
	}

	return nil
}

func SendEmailVerification(email, token string) error {
	baseURL := os.Getenv("FRONTEND_BASE_URL")

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%sverify/%s", baseURL, token),
	}

	err := createMailjet(email, 5075167, variables)
	if err != nil {
		return err
	}

	return nil
}

func SendConfirmationEmailPasswordChangeSuccess(name, email string) error {
	variables := map[string]interface{}{
		"username": name,
	}

	err := createMailjet(email, 5097353, variables)
	if err != nil {
		return err
	}

	return nil
}

func SendConfirmationEmailUserEmailChangeSuccess(name, oldEmail, newEmail, formattedTime string) error {
	variables := map[string]interface{}{
		"name":      name,
		"oldEmail":  oldEmail,
		"newEmail":  newEmail,
		"updatedAt": formattedTime,
	}

	err := createMailjet(newEmail, 5201222, variables)
	if err != nil {
		return err
	}

	return nil
}
