package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
)

type PeopleKey string

const (
	PeopleKeyID       PeopleKey = "id"
	PeopleKeyUsername PeopleKey = "username"
)

type Plan string

const (
	PlanNone      Plan = "None"
	PlanArtist    Plan = "Artist"
	PlanArtistPro Plan = "ArtistPro"
)

type Person struct {
	Id         int64
	Username   string
	Name       string
	ImageUrl   string
	Verified   bool
	Plan       Plan
	TrackCount int64
}

type PeopleRepository struct {
	db *sql.DB
}

func NewPeopleRepository(db *sql.DB) *PeopleRepository {
	return &PeopleRepository{db: db}
}

func (pr *PeopleRepository) FindPersonByIndex(ctx context.Context, key PeopleKey, id int64) (person Person, found bool, err error) {
	return pr.queryPersonRow(ctx, "SELECT id, username, name, image_url, verified, plan, track_count FROM people WHERE $1 = $2;", key, id)
}

func (pr *PeopleRepository) Create(ctx context.Context, handle, name, imageUrl string, verified bool, plan Plan, trackCount int64) (person Person, err error) {
	person, _, err = pr.queryPersonRow(
		ctx,
		"select id, username, name, image_url, verified, plan, track_count from new_person($1, $2, $3, $4, $5, $6);",
		handle, name, imageUrl, verified, plan, trackCount,
	)
	return person, err
}

func (pr *PeopleRepository) CreateFollows(ctx context.Context, followerId int64, followeeHandles []string) error {
	_, err := pr.db.ExecContext(ctx, `SELECT new_follows($1, $2);`, followerId, pq.Array(followeeHandles))
	if err != nil {
		slog.Error("failed to create followee", "error", err)
	}
	return err
}

func (pr *PeopleRepository) queryPersonRow(ctx context.Context, query string, args ...any) (person Person, found bool, err error) {
	err = pr.db.QueryRowContext(
		ctx, query, args...,
	).Scan(
		&person.Id,
		&person.Username,
		&person.Name,
		&person.ImageUrl,
		&person.Verified,
		&person.Plan,
		&person.TrackCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return person, false, nil
		}
		return person, false, err
	}
	return person, true, nil
}
