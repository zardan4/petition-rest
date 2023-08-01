package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	petitions "github.com/zardan4/petition-rest"
)

type PetitionPostgres struct {
	db *sqlx.DB
}

func NewPetitionPostgres(db *sqlx.DB) *PetitionPostgres {
	return &PetitionPostgres{
		db: db,
	}
}

func (p *PetitionPostgres) CreatePetition(title, text string, authorId int) (int, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf(`INSERT INTO %s (title, text) VALUES ($1, $2) RETURNING id`, petitionsTable)

	row := tx.QueryRow(query, title, text)

	err = row.Scan(&id)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return 0, err
	}

	query = fmt.Sprintf(`INSERT INTO %s (user_id, petition_id) VALUES ($1, $2)`, usersPetitionsTable)

	_, err = tx.Exec(query, authorId, id)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return 0, err
	}

	return id, tx.Commit()
}

func (s *PetitionPostgres) GetAllPetitions() ([]petitions.Petition, error) {
	var petitions []petitions.Petition

	query := fmt.Sprintf("SELECT * FROM %s", petitionsTable)

	err := s.db.Select(&petitions, query)
	if err != nil {
		return nil, err
	}

	return petitions, nil
}

func (s *PetitionPostgres) GetPetition(petitionId int) (petitions.Petition, error) {
	var petition petitions.Petition

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", petitionsTable)

	err := s.db.Get(&petition, query, petitionId)
	if err != nil {
		return petitions.Petition{}, err
	}

	return petition, nil
}

func (s *PetitionPostgres) DeletePetition(petitionId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s pl USING %s ul WHERE pl.id=ul.petition_id AND ul.user_id=$1 AND ul.petition_id=$2", petitionsTable, usersPetitionsTable)
	res, err := s.db.Exec(query, userId, petitionId)

	rowsChanged, _ := res.RowsAffected()
	if rowsChanged == 0 {
		return errors.New("no rows affected")
	}

	return err
}

func (s *PetitionPostgres) UpdatePetition(petition petitions.UpdatePetitionInput, petitionId, userId int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if petition.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *petition.Title)
		argId++
	}

	if petition.Date != nil {
		setValues = append(setValues, fmt.Sprintf("date=$%d", argId))
		args = append(args, *petition.Date)
		argId++
	}

	if petition.Timeend != nil {
		setValues = append(setValues, fmt.Sprintf("timeend=$%d", argId))
		args = append(args, *petition.Timeend)
		argId++
	}

	if petition.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, *petition.Text)
		argId++
	}

	if petition.Answer != nil {
		setValues = append(setValues, fmt.Sprintf("answer=$%d", argId))
		args = append(args, *petition.Answer)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s pl SET %s FROM %s ul WHERE pl.id = ul.petition_id AND ul.petition_id = $%d AND ul.user_id = $%d", petitionsTable, setQuery, usersPetitionsTable, argId, argId+1)

	args = append(args, petitionId, userId)

	res, err := s.db.Exec(query, args...)

	rowsChanged, _ := res.RowsAffected()
	if rowsChanged == 0 {
		return errors.New("no rows affected")
	}

	return err
}
