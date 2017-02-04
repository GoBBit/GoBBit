package mail

import (
    "log"
    "net/smtp"
    "strings"

    "GoBBit/config"
)

func SMTPConfigured() bool{
    smtpConfig := config.GetInstance().SMTPConfig
    return (smtpConfig.Host != "" && smtpConfig.Port != "" && smtpConfig.User != "")
}

func SendMail(to, subject, message string){
    if !SMTPConfigured(){
        log.Println("SMTP server not configured")
        return
    }

    smtpConfig := config.GetInstance().SMTPConfig

    // Set up authentication information.
    auth := smtp.PlainAuth("", smtpConfig.User, smtpConfig.Pass, smtpConfig.Host)

    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
    toArr := []string{to}
    msg := []byte("To: "+ to +"\r\n" +
        "Subject: "+ subject +"\r\n" +
        "\r\n" +
        ""+ message +"\r\n")

    err := smtp.SendMail(smtpConfig.Host + ":" + smtpConfig.Port, auth, smtpConfig.SenderAddress, toArr, msg)
    if err != nil {
        log.Println(err)
    }
}

func SendUserActivation(to, activationLink, username string){
    if !SMTPConfigured(){
        return
    }

    subject := "Activate your account"
    message := config.GetInstance().UserActivationEmailTemplate
    message = strings.Replace(message, "{username}", username, -1)
    message = strings.Replace(message, "{activationlink}", activationLink, -1)


    SendMail(to, subject, message)
}

