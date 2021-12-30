package models

type Points struct{
	Rank int `json:"rank"`
	Score float64 `json:"score"`
	College College `json:"college"`
}

type College struct{
	Name string `json:"name"`
	Id string `json:"id"`
}
