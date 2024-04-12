package repository

import (
	avito "avito_testcase"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestActorPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	authRepo := NewAuthPostgres(sqlxDB)

	type args struct {
		user avito.User
	}
	type mockBehaivior func(args args)

	testTable := []struct {
		name          string
		mockBehaivior mockBehaivior
		args          args
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				user: avito.User{
					Id:       1,
					Username: "username",
					Password: "password",
					Is_admin: true,
				},
			},

			mockBehaivior: func(args args) {

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.user.Id)

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.Username, args.user.Password, args.user.Is_admin).WillReturnRows(rows)
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaivior(testCase.args)

			got, err := authRepo.CreateUser(testCase.args.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.user.Id, got)
			}
		})
	}

}

func TestActorPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	authRepo := NewAuthPostgres(sqlxDB)

	type args struct {
		user avito.User
	}
	type mockBehaivior func(args args)

	testTable := []struct {
		name          string
		mockBehaivior mockBehaivior
		args          args
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				user: avito.User{
					Id:       1,
					Username: "",
					Password: "",
					Is_admin: false,
				},
			},

			mockBehaivior: func(args args) {

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.user.Id)

				mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1 AND password_hash=\\$2").
					WithArgs(args.user.Username, args.user.Password).WillReturnRows(rows)
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaivior(testCase.args)

			got, err := authRepo.GetUser(testCase.args.user.Username, testCase.args.user.Password)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.user, got)
			}
		})
	}

}

func TestActorPostgres_GetUserStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' wasn't expected when opening a stub db connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	authRepo := NewAuthPostgres(sqlxDB)

	type args struct {
		user avito.User
	}
	type mockBehaivior func(args args)

	testTable := []struct {
		name          string
		mockBehaivior mockBehaivior
		args          args
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				user: avito.User{
					Id:       1,
					Username: "",
					Password: "",
					Is_admin: true,
				},
			},

			mockBehaivior: func(args args) {

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.user.Id)

				mock.ExpectQuery("SELECT is_admin FROM users WHERE id=\\$1").
					WithArgs(args.user.Id).WillReturnRows(rows)
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaivior(testCase.args)

			got, err := authRepo.GetUserStatus(testCase.args.user.Id)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.user.Is_admin, got)
			}
		})
	}

}
