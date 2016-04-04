/*
   Send Mail in Golang via SMTP
 */
package routes

import (
    "bytes"
    "net/smtp"
    "github.com/ahermida/dartboardAPI/api/Config"
    "text/template"
)

//Struct to handle template
type TemplateData struct {
    From    string
    To      string
    Subject string
    Body    string
    Link    string
}

//send mail function
func setupEmail(to string, token string) error {
  var err error
  var doc bytes.Buffer

  //email template
  const emailTemplate = `From: {{.From}}
Subject: {{.Subject}}
To: {{.To}}


{{.Body}}
http://localhost:8080/auth/{{.Link}}
Sincerely,
{{.From}}
`
  //Send to "to" -- user's email
  context := &TemplateData{
       From: "Albert Hermida",
         To: to,
    Subject: "Authorize your account for my Project",
       Body: "Hey, Thanks for signing up! Click on the link to authorize your account:",
       Link: token, //such a small string, not much reason to implement bytes string cc
  }

  //make new template for email
  t := template.New("emailTemplate")

  //parse it
  t, err = t.Parse(emailTemplate)

  //if there's an error parsing our template, return it
  if err != nil {
      return err
  }

  //if there's an error executing the template, return it
  exErr := t.Execute(&doc, context)
  if exErr != nil {
      return exErr
  }

  //basic smtp auth
  auth := smtp.PlainAuth("",
    config.Email.Username,
    config.Email.Password,
    config.Email.EmailServer,
  )

  //send mail via gmail smtp
  errSending := smtp.SendMail("smtp.gmail.com:587",
    auth,
    config.Email.Username,
    []string{to},
    doc.Bytes())

  //if there's an error sending, let ourselves know
  if errSending != nil {
    return errSending
  }

  //return nil error if everything went as planned
  return nil
}

//send mail function
func recoverEmail(to string, token string) error {
  var err error
  var doc bytes.Buffer

  //email template
  const emailTemplate = `From: {{.From}}
Subject: {{.Subject}}
To: {{.To}}


{{.Body}}
http://localhost:8080/recovery/{{.Link}}
Sincerely,
{{.From}}
`
  //Send to "to" -- user's email
  context := &TemplateData{
       From: "Albert Hermida",
         To: to,
    Subject: "Authorize your account for my Project",
       Body: "Click on the link to reset your password:",
       Link: token, //such a small string, not much reason to implement bytes string cc
  }

  //make new template for email
  t := template.New("emailTemplate")

  //parse it
  t, err = t.Parse(emailTemplate)

  //if there's an error parsing our template, return it
  if err != nil {
      return err
  }

  //if there's an error executing the template, return it
  exErr := t.Execute(&doc, context)
  if exErr != nil {
      return exErr
  }

  //basic smtp auth
  auth := smtp.PlainAuth("",
    config.Email.Username,
    config.Email.Password,
    config.Email.EmailServer,
  )

  //send mail via gmail smtp
  errSending := smtp.SendMail("smtp.gmail.com:587",
    auth,
    config.Email.Username,
    []string{to},
    doc.Bytes())

  //if there's an error sending, let ourselves know
  if errSending != nil {
    return errSending
  }

  //return nil error if everything went as planned
  return nil
}
