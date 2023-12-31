package resttransport

import (
	"errors"
	"l0wb/models"
	mock_resttransport "l0wb/transport/resttransport/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandler_OrderGet(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name       string
		args       args
		move       func(orderer *mock_resttransport.MockOrderer)
		wantStatus int
	}{
		{
			name: "OK Handler Get",
			args: args{
				uid: "test1",
			},
			move: func(orderer *mock_resttransport.MockOrderer) {
				orderer.EXPECT().Get(gomock.Eq("test1")).Times(1).Return(models.Order{
					OrderUID: "test1",
					Order:    []byte("Hello world"),
				}, nil)
			},
			wantStatus: 200,
		},
		{
			name: "ERR Handler Get: error in orderer get",
			args: args{
				uid: "test1",
			},
			move: func(orderer *mock_resttransport.MockOrderer) {
				orderer.EXPECT().Get(gomock.Eq("test1")).Times(1).Return(models.Order{}, errors.New("some error"))
			},
			wantStatus: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderer := mock_resttransport.NewMockOrderer(ctrl)

			tt.move(orderer)

			h := Handler{
				s:    orderer,
				path: "./mock/",
			}

			recorder := httptest.NewRecorder()
			url := "/?uid=" + tt.args.uid
			request, _ := http.NewRequest(http.MethodGet, url, nil)

			router := h.Init()

			router.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("Handler.Get() error want code %d, got %d", tt.wantStatus, recorder.Code)
			}
		})
	}
}
