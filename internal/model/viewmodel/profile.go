package viewmodel

type ProfilePage struct {
	Page
	Profile    UserSummary
	Posts      []PostCard
	PageNo     int
	TotalPages int
	Total      int64
	IsMe       bool
}
