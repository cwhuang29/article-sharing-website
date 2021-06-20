package constants

const (
	ResetPasswordMaxRetry = 5

	Subject  = "Reset Password Notification"
	CharSet  = "UTF-8"
	TextBody = "Hello %s! " +
		"You are receiving this email because we received a password reset request from your account. " +
		"Copy and paste the following link into your browser to change your pawssword: %s. " +
		"The password reset link will expire in %d minutes. If you didn't request a password reset, no further action is required. " +
		"Please feel free to contact us if you have any further questions."
	HtmlBody = `
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
