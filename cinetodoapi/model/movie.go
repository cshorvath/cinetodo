package model

type Movie struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	Year          uint8  `json:"year"`
	Director      string `json:"director"`
}

type UserMovie struct {
	UserID  int `gorm:"primaryKey"`
	MovieID int `gorm:"primaryKey"`
	Seen    bool
}
