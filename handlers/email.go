package handlers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/constants"
	"github.com/sirupsen/logrus"
)

func getAWSSVC() *ses.SES {
	region := config.GetCopy().Email.Region
	if sess, err := session.NewSession(&aws.Config{Region: aws.String(region)}); err != nil {
		logrus.Errorln(err.Error())
		return nil
	} else {
		return ses.New(sess) // Create an SES session.
	}
}

func resetPasswordEmailBody(recipient, name, link string, expireTime int) *ses.SendEmailInput {
	interpolatedHtmlBody := fmt.Sprintf(constants.HtmlBody, name, link, expireTime, link)
	interpolatedTextBody := fmt.Sprintf(constants.TextBody, name, link, expireTime)
	sender := config.GetCopy().Email.Sender

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(constants.CharSet),
				Data:    aws.String(constants.Subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(constants.CharSet),
					Data:    aws.String(interpolatedHtmlBody),
				},
				Text: &ses.Content{ // The email body for recipients with non-HTML email clients.
					Charset: aws.String(constants.CharSet),
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
