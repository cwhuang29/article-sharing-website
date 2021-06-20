package constants

const (
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
)
