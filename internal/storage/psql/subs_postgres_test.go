package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func Test_CreateSub(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}

	defer db.Close()

	r := NewSubsPostgres(db)

	type args struct {
		petitionId int
		userId     int
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				petitionId: 1,
				userId:     1,
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", subsTable)).WillReturnRows(rows)

				mock.ExpectExec(fmt.Sprintf(`INSERT INTO %s \(sub_id, petition_id, user_id\) VALUES \(\$1, \$2, \$3\)`, petitionsSubsTable)).
					WithArgs(id, args.petitionId, args.userId).
					WillReturnResult(sqlmock.NewResult(int64(id), 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty field petitionId",
			args: args{
				userId: 1,
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("no petitionid provided"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", subsTable)).WillReturnRows(rows)

				mock.ExpectExec(fmt.Sprintf(`INSERT INTO %s \(sub_id, petition_id, user_id\) VALUES \(\$1, \$2, \$3\)`, petitionsSubsTable)).
					WithArgs(id, args.petitionId, args.userId).
					WillReturnResult(sqlmock.NewResult(int64(id), 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty field userId",
			args: args{
				petitionId: 1,
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("no userid provided"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", subsTable)).WillReturnRows(rows)

				mock.ExpectExec(fmt.Sprintf(`INSERT INTO %s \(sub_id, petition_id, user_id\) VALUES \(\$1, \$2, \$3\)`, petitionsSubsTable)).
					WithArgs(id, args.petitionId, args.userId).
					WillReturnResult(sqlmock.NewResult(int64(id), 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd insert error",
			args: args{
				petitionId: 1,
				userId:     1,
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", subsTable)).WillReturnRows(rows)

				mock.ExpectExec(fmt.Sprintf(`INSERT INTO %s \(sub_id, petition_id, user_id\) VALUES \(\$1, \$2, \$3\)`, petitionsSubsTable)).
					WithArgs(id, args.petitionId, args.userId).
					WillReturnError(errors.New("error while handling 2nd insert"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	// act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			res, err := r.CreateSub(testCase.args.petitionId, testCase.args.userId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, res)
			}
		})
	}
}
