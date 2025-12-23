package word

import (
	"app/ent"
	"app/ent/user"
	"app/ent/word"
	"app/internal/db"
	"context"
	"log"
)

type Repository interface {
	GetWord(ctx context.Context, wordString string) (*ent.Word, error)
	GetAllWords(ctx context.Context) ([]*ent.Word, error)
	GetMyWords(ctx context.Context, userId int) ([]*ent.Word, error)
	CreateWord(ctx context.Context, userId int, wordString string, levelString string, wordRu string) (*ent.Word, error)
	RemoveMyWord(ctx context.Context, userId int, wordId int) error
}

type PostgresRepository struct {
	Db *db.Db
}

func NewPostgresRepository(db *db.Db) Repository {
	return &PostgresRepository{Db: db}
}

func (r *PostgresRepository) GetWord(ctx context.Context, wordString string) (*ent.Word, error) {
	return r.Db.Client.Word.Query().Where(word.WordEQ(wordString)).Only(ctx)
}

func (r *PostgresRepository) GetAllWords(ctx context.Context) ([]*ent.Word, error) {
	return r.Db.Client.Word.Query().All(ctx)
}

func (r *PostgresRepository) GetMyWords(ctx context.Context, userId int) ([]*ent.Word, error) {
	log.Println(userId)
	return r.Db.Client.Word.Query().Where(word.HasUserWith(user.IDEQ(userId))).All(ctx)
}

func (r *PostgresRepository) CreateWord(ctx context.Context, userId int, wordString string, levelString string, wordRu string) (*ent.Word, error) {
	existingWord, err := r.Db.Client.Word.
		Query().
		Where(word.WordEQ(wordString)).
		Only(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if existingWord != nil {
		err = existingWord.Update().
			AddUserIDs(userId).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return existingWord, nil
	}

	newWord, err := r.Db.Client.Word.
		Create().
		SetWord(wordString).
		SetRuWord(wordRu).
		SetLevel(levelString).
		AddUserIDs(userId).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return newWord, nil
}

func (r *PostgresRepository) RemoveMyWord(ctx context.Context, userId int, wordId int) error {
	return r.Db.Client.User.UpdateOneID(userId).RemoveWordIDs(wordId).Exec(ctx)
}
