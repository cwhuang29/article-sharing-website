package utils

const (
	LandingPage       = "/articles/weekly-update"
	LoginPage         = "/login"
	ResetPasswordPath = "/password/reset/"

	CsrfTokenAge             = 6 * 60 * 60 // 6 hours
	LoginMaxAge              = 30 * 86400  // 1 month
	ResetPasswordTokenMaxAge = 60 * 60     // 1 hour

	ResetPasswordMaxRetry = 5

	// Size adjustment for data display on frontend
	// Give slightly more word counts since the CSS style (e.g. text-justify: inter-word) affects layout in varying degrees
	TitleSizeLimit                 = 35.
	SubtitleSizeLimit              = 54.
	OutlineSizeLimit               = 510.
	OutlineSizeLimitWithCoverPhoto = 340.

	// Size limitation for data stored in DB
	TitleBytesLimit    = 255
	SubtitleBytesLimit = 255
	TagsNumLimit       = 5
	TagsBytesLimit     = 20              // Emojis and some Chinese words are 4 bytes long
	FileMaxSize        = 8 * 1000 * 1000 // 8MB

	// The path where the users uploaded files stored (and also the filename prefix)
	UploadImageDir = "public/upload/images/"

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
