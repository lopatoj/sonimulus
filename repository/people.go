package repository

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
)

type Plan int

const (
	PlanNone Plan = iota
	PlanArtist
	PlanArtistPro
)

type Person struct {
	Id         int64
	Handle     string
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

func (pr *PeopleRepository) GetPersonById(id int64) (*Person, error) {
	return pr.queryPersonRow("SELECT id, handle, name, image_url, verified, plan, track_count FROM people WHERE id = $1;", id)
}

func (pr *PeopleRepository) GetPersonByHandle(handle string) (*Person, error) {
	return pr.queryPersonRow("SELECT id, handle, name, image_url, verified, plan, track_count FROM people WHERE handle = $1;", handle)
}

func (pr *PeopleRepository) CreatePerson(handle, name, imageUrl string, verified bool, plan Plan, trackCount int64) (*Person, error) {
	return pr.queryPersonRow(
		"select id, handle, name, image_url, verified, plan, track_count from new_person($1, $2, $3, $4, $5, $6);",
		handle, name, imageUrl, verified, planToStr(plan), trackCount,
	)
}

func (pr *PeopleRepository) CreateFollows(followerId int64, followeeHandle []string) error {
	_, err := pr.db.Exec(`SELECT new_follows($1, $2);`, followerId, pq.Array(followeeHandle))
	if err != nil {
		slog.Error("failed to create followee", "error", err)
	}
	return err
}

func (pr *PeopleRepository) queryPersonRow(query string, args ...any) (*Person, error) {
	person := &Person{}
	var plan string
	err := pr.db.QueryRow(
		query, args...,
	).Scan(
		&person.Id,
		&person.Handle,
		&person.Name,
		&person.ImageUrl,
		&person.Verified,
		&plan,
		&person.TrackCount,
	)
	person.Plan = strToPlan(plan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return person, nil
}

func strToPlan(plan string) Plan {
	switch plan {
	case "Artist":
		return PlanArtist
	case "ArtistPro":
		return PlanArtistPro
	default:
		return PlanNone
	}
}

func planToStr(plan Plan) string {
	switch plan {
	case PlanArtist:
		return "Artist"
	case PlanArtistPro:
		return "ArtistPro"
	default:
		return "None"
	}
}
