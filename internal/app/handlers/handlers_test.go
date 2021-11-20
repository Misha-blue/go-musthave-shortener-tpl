package handlers

import (
	"net/http"
	"testing"
)

func TestHandleURLRequest(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleURLRequest(tt.args.w, tt.args.r)
		})
	}
}

func Test_getShortenURL(t *testing.T) {
	type args struct {
		storage map[string]string
		e       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getShortenURL(tt.args.storage, tt.args.e); got != tt.want {
				t.Errorf("getShortenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findShornedURL(t *testing.T) {
	type args struct {
		storage map[string]string
		element string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findShornedURL(tt.args.storage, tt.args.element); got != tt.want {
				t.Errorf("findShornedURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		storage map[string]string
		element string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.storage, tt.args.element); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
