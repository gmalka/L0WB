package inmemory

import (
	"l0wb/models"
	"reflect"
	"sync"
	"testing"
)

func Test_inmemory_Add(t *testing.T) {
	type fields struct {
		store map[string][]byte
		m     *sync.Mutex
	}
	type args struct {
		order models.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "OK inmemory cash store ADD",
			fields: fields{
				store: make(map[string][]byte, 10),
				m:     &sync.Mutex{},
			},
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			wantErr: false,
		},
		{
			name: "ERR inmemory cash store ADD: adding an existing value",
			fields: fields{
				store: map[string][]byte{"1": []byte("Bye world")},
				m:     &sync.Mutex{},
			},
			args: args{
				order: models.Order{
					OrderUID: "1",
					Order:    []byte("Hello world"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := inmemory{
				store: tt.fields.store,
				m:     tt.fields.m,
			}
			if err := m.Add(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("inmemory.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inmemory_Get(t *testing.T) {
	type fields struct {
		store map[string][]byte
		m     *sync.Mutex
	}
	type args struct {
		OrderUID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Order
		wantErr bool
	}{
		{
			name: "OK inmemory cash store GET",
			fields: fields{
				store: map[string][]byte{"1": []byte("Bye world")},
				m:     &sync.Mutex{},
			},
			args: args{
				OrderUID: "1",
			},
			want: models.Order{
				OrderUID: "1",
				Order: []byte("Bye world"),
			},
			wantErr: false,
		},
		{
			name: "ERR inmemory cash store GET: cant find order",
			fields: fields{
				store: map[string][]byte{},
				m:     &sync.Mutex{},
			},
			args: args{
				OrderUID: "1",
			},
			want: models.Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := inmemory{
				store: tt.fields.store,
				m:     tt.fields.m,
			}
			got, err := m.Get(tt.args.OrderUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("inmemory.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inmemory.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
