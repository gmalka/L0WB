package orderservice

import (
	"errors"
	"l0wb/models"
	mock_orderservice "l0wb/services/orderservice/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewOrderService(t *testing.T) {
	tests := []struct {
		name    string
		move    func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher)
		wantErr bool
	}{
		{
			name: "OK Order service create",
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				db.EXPECT().GetAll().Times(1).Return([]models.Order{{OrderUID: "1", Order: []byte("Some text")}, {OrderUID: "2", Order: []byte("Some text")}}, nil)
				cash.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Some text")})).Times(1).Return(nil)
				cash.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "2", Order: []byte("Some text")})).Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ERR Order service create: error on GetAll from store",
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				db.EXPECT().GetAll().Times(1).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "ERR Order service create: error on Add to cash",
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				db.EXPECT().GetAll().Times(1).Return([]models.Order{{OrderUID: "1", Order: []byte("Some text")}, {OrderUID: "2", Order: []byte("Some text")}}, nil)
				cash.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Some text")})).Times(1).Return(errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock_orderservice.NewMockDatabase(ctrl)
			cash := mock_orderservice.NewMockCasher(ctrl)

			tt.move(db, cash)

			_, err := NewOrderService(db, cash)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOrderService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_orderService_Add(t *testing.T) {
	type args struct {
		order models.Order
	}
	tests := []struct {
		name    string
		args    args
		move    func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher)
		wantErr bool
	}{
		{
			name: "OK Add Order Service",
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				cash.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Hello world")})).Times(1).Return(nil)
				db.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Hello world")})).Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ERR Add Order Service: error in store",
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				cash.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Hello world")})).Times(1).Return(nil)
				db.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Hello world")})).Times(1).Return(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "ERR Add Order Service: error in cash",
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				cash.EXPECT().Add(gomock.Eq(models.Order{OrderUID: "1", Order: []byte("Hello world")})).Times(1).Return(errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock_orderservice.NewMockDatabase(ctrl)
			cash := mock_orderservice.NewMockCasher(ctrl)

			tt.move(db, cash)

			o := orderService{
				db:   db,
				cash: cash,
			}

			if err := o.Add(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("orderService.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_orderService_Get(t *testing.T) {
	type args struct {
		OrderUID string
	}
	tests := []struct {
		name    string
		args    args
		move    func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher)
		want    models.Order
		wantErr bool
	}{
		{
			name: "OK Get Order Service: Get from cash",
			args: args{
				OrderUID: "1",
			},
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				cash.EXPECT().Get(gomock.Eq("1")).Times(1).Return(models.Order{
					OrderUID: "1",
					Order:    []byte("Hello there"),
				}, nil)
			},
			want: models.Order{
				OrderUID: "1",
				Order:    []byte("Hello there"),
			},
			wantErr: false,
		},
		{
			name: "OK Get Order Service: Get from db",
			args: args{
				OrderUID: "1",
			},
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				cash.EXPECT().Get(gomock.Eq("1")).Times(1).Return(models.Order{}, errors.New("Some error"))
				db.EXPECT().Get(gomock.Eq("1")).Times(1).Return(models.Order{
					OrderUID: "1",
					Order:    []byte("Hello there"),
				}, nil)
				cash.EXPECT().Add(models.Order{
					OrderUID: "1",
					Order:    []byte("Hello there"),
				}).Times(1).Return(nil)
			},
			want: models.Order{
				OrderUID: "1",
				Order:    []byte("Hello there"),
			},
			wantErr: false,
		},
		{
			name: "ERR Get Order Service: error in cash and store",
			args: args{
				OrderUID: "1",
			},
			move: func(db *mock_orderservice.MockDatabase, cash *mock_orderservice.MockCasher) {
				cash.EXPECT().Get(gomock.Eq("1")).Times(1).Return(models.Order{}, errors.New("Some error"))
				db.EXPECT().Get(gomock.Eq("1")).Times(1).Return(models.Order{}, errors.New("Some error"))
			},
			want: models.Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock_orderservice.NewMockDatabase(ctrl)
			cash := mock_orderservice.NewMockCasher(ctrl)

			tt.move(db, cash)

			o := orderService{
				db:   db,
				cash: cash,
			}

			got, err := o.Get(tt.args.OrderUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
