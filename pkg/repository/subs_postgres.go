package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	petitions "github.com/zardan4/petition-rest"
)

type SubsPostgres struct {
	db *sqlx.DB
}

func NewSubsPostgres(db *sqlx.DB) *SubsPostgres {
	return &SubsPostgres{
		db: db,
	}
}

func (s *SubsPostgres) GetAllSubs(petitionId int) ([]petitions.Sub, error) {
	var subs []petitions.Sub
	query := fmt.Sprintf(`SELECT st.id, st.date, us.name, us.id as user_id
		FROM %s st
		INNER JOIN %s ps ON st.id=ps.sub_id
		INNER JOIN %s up ON ps.petition_id=up.petition_id AND up.petition_id=$1
		INNER JOIN %s us ON us.id=ps.user_id`,
		subsTable, petitionsSubsTable, usersPetitionsTable, usersTable)

	err := s.db.Select(&subs, query, petitionId)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SubsPostgres) CreateSub(petitionId, userId int) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", subsTable)
	row := tx.QueryRow(query)

	var subId int

	err = row.Scan(&subId)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return 0, err
	}

	query = fmt.Sprintf("INSERT INTO %s (sub_id, petition_id, user_id) VALUES ($1, $2, $3)",
		petitionsSubsTable)
	_, err = tx.Exec(query, subId, petitionId, userId)
	if err != nil {
		tx.Rollback() // nolint: errcheck
		return 0, err
	}

	return subId, tx.Commit()
}

func (s *SubsPostgres) DeleteSub(subId, petitionId, userId int) error {
	// query := fmt.Sprintf(`DELETE FROM %s st
	// USING %s ps, %s up
	// WHERE ps.sub_id=st.id AND ps.sub_id=$1
	// AND ps.petition_id=up.petition_id AND ps.petition_id=$2
	// AND ps.user_id=up.user_id AND ps.user_id=$3`,
	// 	subsTable, petitionsSubsTable, usersPetitionsTable)
	query := fmt.Sprintf(`DELETE FROM %s st
	USING %s ps, %s us, %s up
	WHERE st.id=ps.sub_id AND st.id=$1
	AND us.id=ps.user_id AND ps.user_id=$2
	AND ps.petition_id=up.petition_id AND ps.petition_id=$3`,
		subsTable, petitionsSubsTable, usersTable, usersPetitionsTable)

	res, err := s.db.Exec(query, subId, userId, petitionId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (s *SubsPostgres) CheckSignature(petitionId, userId int) (bool, error) {
	subs, err := s.GetAllSubs(petitionId)
	if err != nil {
		return false, err
	}

	for _, sub := range subs {
		subUserIdInt, _ := strconv.Atoi(sub.UserId)
		if subUserIdInt == userId {
			return true, nil
		}
	}

	return false, nil
}
