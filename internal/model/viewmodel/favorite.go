package viewmodel

type FavoritesPage struct {
	Page
	Posts      []PostCard
	PageNo     int
	TotalPages int
	Total      int64
}
