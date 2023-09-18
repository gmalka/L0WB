package postgresdb

import (
	"errors"
	"l0wb/models"
	"reflect"
	"regexp"
	"testing"

	_ "github.com/lib/pq"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func Test_noteRepository_Add(t *testing.T) {
	type args struct {
		order models.Order
	}
	tests := []struct {
		name    string
		args    args
		move    func(mock sqlmock.Sqlmock, uid string, data []byte)
		wantErr bool
	}{
		{
			name: "OK postgres store ADD",
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			move: func(mock sqlmock.Sqlmock, uid string, data []byte) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Orders VALUES ($1,$2)")).WithArgs(uid, data).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "ERR postgres store ADD: database error",
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			move: func(mock sqlmock.Sqlmock, uid string, data []byte) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Orders VALUES ($1,$2)")).WithArgs(uid, data).
					WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Errorf("Error while create sqlmock: %v", err)
			}

			r := noteRepository{
				db: db,
			}

			tt.move(mock, tt.args.order.OrderUID, tt.args.order.Order)

			if err := r.Add(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("noteRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_noteRepository_Get(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		args    args
		move    func(mock sqlmock.Sqlmock, uid string)
		want    models.Order
		wantErr bool
	}{
		{
			name: "OK postgres store GET",
			args: args{
				uid: "1",
			},
			move: func(mock sqlmock.Sqlmock, uid string) {
				row := sqlmock.NewRows([]string{"order_uid", "data"}).AddRow("1", "Hello world")

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM Orders WHERE order_uid=$1")).
					WithArgs(uid).
					WillReturnRows(row)
			},
			want: models.Order{
				OrderUID: "1",
				Order:    []byte("Hello world"),
			},
			wantErr: false,
		},
		{
			name: "ERR postgres store GET: database error",
			args: args{
				uid: "1",
			},
			move: func(mock sqlmock.Sqlmock, uid string) {
				row := sqlmock.NewRows([]string{"order_uid", "data"}).AddRow("1", "Hello world")

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM Orders WHERE order_uid=$1")).
					WithArgs(uid).
					WillReturnRows(row).
					WillReturnError(errors.New("some error"))
			},
			want:    models.Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Errorf("Error while create sqlmock: %v", err)
			}

			r := noteRepository{
				db: db,
			}

			tt.move(mock, tt.args.uid)

			got, err := r.Get(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("noteRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("noteRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_noteRepository_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		move    func(mock sqlmock.Sqlmock)
		want    []models.Order
		wantErr bool
	}{
		{
			name: "OK postgres store GET ALL",
			move: func(mock sqlmock.Sqlmock) {
				row := sqlmock.NewRows([]string{"order_uid", "data"}).AddRow("1", "Hello world").AddRow("2", "Good bye")

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM Orders")).
					WillReturnRows(row)
			},
			want: []models.Order{
				{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
				{
					OrderUID: "2",
					Order:    []byte("Good bye"),
				},
			},
			wantErr: false,
		},
		{
			name: "ERR postgres store GET ALL: database get all error",
			move: func(mock sqlmock.Sqlmock) {
				row := sqlmock.NewRows([]string{"order_uid", "data"}).AddRow("1", "Hello world").AddRow("2", "Good bye")

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM Orders")).
					WillReturnRows(row).
					WillReturnError(errors.New("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Errorf("Error while create sqlmock: %v", err)
			}

			r := noteRepository{
				db: db,
			}

			tt.move(mock)

			got, err := r.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("noteRepository.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("noteRepository.GetAll() = %s, want %s", got, tt.want)
			}
		})
	}
}
