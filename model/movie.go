package model

type Movie struct {
    ID            int64  `json:"id"`
    Title         string `json:"title"`
    OriginalTitle string `json:"originalTitle"`
    Year          uint16 `json:"year"`
    Director      string `json:"director"`
    PosterPath    string `json:"posterPath"`
}

type UserMovie struct {
	UserID  uint  `gorm:"primaryKey"`
	MovieID int64 `gorm:"primaryKey"`
	Movie   Movie
	Seen    bool
}
