package handlers

type UserStatus int

const (
	IsGuest UserStatus = iota
	IsMember
	IsMemberAndVerified
	IsAdmin
	IsAdminAndVerified
)

func (s UserStatus) String() string {
	return [...]string{"guest", "member", "verified member", "admin", "verified admin"}[s]
}

// Notice: The field names should be exact the same as the fields in html template.
// Rewrite by `json:"customize_name"` won't work (but when fetching data by JS, it works)
type Article struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Subtitle   string   `json:"subtitle"`
	Date       string   `json:"date"`
	Authors    []string `json:"authors"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	Outline    string   `json:"outline"`
	CoverPhoto string   `json:"cover_photo"`
	Content    string   `json:"content"`
	AdminOnly  bool     `json:"adminOnly"`
}
