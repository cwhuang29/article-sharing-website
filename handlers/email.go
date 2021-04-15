package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/sirupsen/logrus"
)

const (
	charSet  = "UTF-8"
	subject  = "Reset Password Notification"
	textBody = "Hello %s! " +
		"You are receiving this email because we received a password reset request from your account. " +
		"Copy and paste the following link into your browser to change your pawssword: %s. " +
		"The password reset link will expire in %d minutes. If you didn't request a password reset, no further action is required. " +
		"Please feel free to contact us if you have any further questions."
	htmlBody = `
<div style="width: 68%%; margin-left: auto; margin-right: auto; color: #3D3D3D">
  <div style=" color: #FCD432; text-align: center; font-weight: 700; font-size: 2rem">
    i.news
  </div>
  <p style="font-size: 1.2rem">Hello %s</p>
  <p>You are receiving this email because we received a password reset request from your account. Please click the following button to change your password.</p>
  <div style="display: flex; padding-top:15px; padding-bottom: 15px;">
	<a href="%s" target="_blank" class="button" style="color: #FCFCFC; background-color: #FCD432; padding: 14px 9px; border-radius: 7px; margin: auto; font-weight: 700;">Reset Password</a>
  </div>
  <p>The password reset link will expire in %d minutes. If you didn't request a password reset, no further action is required.</p>
  <p>Please feel free to contact us if you have any further questions.</p>
  <p>Best regards,<br>i.news</p>
  <hr>
  <div style="text-align: center; color: #7E7E7E; font-size: 0.7rem">
  <p>If you're having trouble clicking the "Reset Password" button, copy and paste the following link into your browser: %s</p>
  <p>&copy;&nbsp;2021 i.news All rights reserved.</p>
  </div>
</div>
`
)

func getAWSSVC() *ses.SES {
	region := config.GetConfig().Email.Region
	if sess, err := session.NewSession(&aws.Config{Region: aws.String(region)}); err != nil {
		logrus.Errorln(err.Error())
		return nil
	} else {
		return ses.New(sess) // Create an SES session.
	}
}

func resetPasswordEmailBody(recipient, name, link string, expireTime int) *ses.SendEmailInput {
	interpolatedHtmlBody := fmt.Sprintf(htmlBody, name, link, expireTime, link)
	interpolatedTextBody := fmt.Sprintf(textBody, name, link, expireTime)
	sender := config.GetConfig().Email.Sender

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(interpolatedHtmlBody),
				},
				Text: &ses.Content{ // The email body for recipients with non-HTML email clients.
					Charset: aws.String(charSet),
					Data:    aws.String(interpolatedTextBody),
				},
			},
		},
		Source: aws.String(sender),
		// ConfigurationSetName: aws.String(ConfigurationSet),
	}
}

func logEmailError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			logrus.Error(ses.ErrCodeMessageRejected, aerr.Error())
		case ses.ErrCodeMailFromDomainNotVerifiedException:
			logrus.Error(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		case ses.ErrCodeConfigurationSetDoesNotExistException:
			logrus.Error(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		default:
			logrus.Error(aerr.Error())
		}
	} else {
		logrus.Error(err.Error())
	}
}

// This function is necessary if your account is in Amazon SES sandbox,
// Cause you need to verify every single email address before using them as senders or recipients.
func VerifyRecipientEmail(email string) bool {
	svc := getAWSSVC()

	if _, err := svc.VerifyEmailAddress(&ses.VerifyEmailAddressInput{EmailAddress: aws.String(email)}); err != nil {
		logEmailError(err)
		return false
	}
	logrus.Info("Verification sent to address: " + email)
	return true
}

func SendResetPasswordEmail(recipient, name, link string, expireMins int) bool {
	svc := getAWSSVC()
	input := resetPasswordEmailBody(recipient, name, link, expireMins)

	if result, err := svc.SendEmail(input); err != nil {
		logEmailError(err)
		return false
	} else {
		logrus.Info("Email Sent to address: " + recipient)
		logrus.Info(result)
		return true
	}
}
