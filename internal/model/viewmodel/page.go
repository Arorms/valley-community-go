package viewmodel

type Page struct {
	Title  string
	Active string
	User   *CurrentUser
	Unread int
	Flash  string
}

type ErrorPage struct {
	Page
	Message string
}
