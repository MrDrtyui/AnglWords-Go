package domain

import "time"

type WordCreateDto struct {
	Word string `json:"word"`
}

type WordResponse struct {
	ID        int       `json:"id"`
	Word      string    `json:"word"`
	RuWord    string    `json:"ruWord"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"createdAt"`
}
