package word

import (
	"app/domain"
	"app/ent"
)

func ToWordResponse(w *ent.Word) domain.WordResponse {
	return domain.WordResponse{
		ID:        w.ID,
		Word:      w.Word,
		RuWord:    w.RuWord,
		Level:     w.Level,
		CreatedAt: w.CreatedAt,
	}
}

func ToWordResponseSlice(words []*ent.Word) []domain.WordResponse {
	res := make([]domain.WordResponse, len(words))
	for i, w := range words {
		res[i] = ToWordResponse(w)
	}
	return res
}
