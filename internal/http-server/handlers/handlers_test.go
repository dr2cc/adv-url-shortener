package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

func TestPostHandler(t *testing.T) {
	type args struct {
		log      *slog.Logger
		urlSaver URLSaver
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PostHandler(tt.args.log, tt.args.urlSaver); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	type args struct {
		log      *slog.Logger
		urlGeter URLGeter
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHandler(tt.args.log, tt.args.urlGeter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
