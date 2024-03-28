package main

import (
	"bytes"
	"fmt"
	"html/template"
)

func SendEmailWithTemplate(receiverEmail string, subject string, htmlTemplate string, htmlArgs interface{}) error {

	t, err := template.New("email_template").Parse(htmlTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	var tplBuffer bytes.Buffer
	err = t.Execute(&tplBuffer, htmlArgs)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	newHTMLData := tplBuffer.String()

	err = SendEmail(SMTP_CONFIG.Sender, receiverEmail, subject, newHTMLData)
	if err != nil {
		return err
	}

	return nil
}
