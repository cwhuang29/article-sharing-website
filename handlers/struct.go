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

type Tag string

type Article struct {
	Title    string   `form:"title" json:"title" binding:"required"`
	Subtitle string   `form:"subtitle" json:"subtitle" binding:"-"`
	Date     string   `form:"date" json:"date" binding:"required"`
	Authors  []string `form:"authors" json:"authors" binding:"required"`
	Category string   `form:"Category" json:"Category" binding:"required"`
	Tags     []string `form:"Tags" json:"Tags" binding:"required"`
	Content  string   `form:"Content" json:"Content" binding:"required"`
}

// Notice: The field names should be the same as the fields in html template.
// Rewrite by `json:"customize_name"` won't work
type OverviewArticle struct {
	ID       int
	Title    string
	Subtitle string
	Date     string
	Authors  []string
	Category string
	Tags     []string
	Content  string
}

type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
