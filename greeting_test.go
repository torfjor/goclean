package goclean_test

import (
	"context"
	"goclean"
	"strings"
	"testing"
)

type staticGreetingStore struct{}

func (s staticGreetingStore) RandomGreetingTemplate(ctx context.Context) (string, error) {
	return "hello %s", nil
}

func TestGreeterFunc_Authorization(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{"no permissions in context", ctx, true},
		{"wrong permissions in context", context.WithValue(ctx, goclean.ContextKeyPermissions, goclean.None), true},
		{"correct permissions in context", context.WithValue(ctx, goclean.ContextKeyPermissions, goclean.Greet), false},
	}

	greeter := goclean.NewAuthorizedGreetingFunc(goclean.NewGreeterFunc(staticGreetingStore{}), goclean.Greet)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := greeter(tt.ctx, "foo")
			if tt.wantErr && err != nil {
				return
			} else if !tt.wantErr && err != nil {
				t.Fatalf("want nil err, got %v", err)
			} else if tt.wantErr && err == nil {
				t.Fatalf("want non-nil err, got %v", err)
			}
		})
	}
}

func TestGreeterFunc(t *testing.T) {
	type args struct {
		store       goclean.GreetingStore
		nameToGreet string
	}
	tests := []struct {
		name string
		args
		wantErr bool
	}{
		{"succeeds for valid input", args{nameToGreet: "foo", store: staticGreetingStore{}}, false},
		{"fails for empty names", args{nameToGreet: "", store: staticGreetingStore{}}, true},
		{"fails for kings and queens", args{nameToGreet: "King Edward", store: staticGreetingStore{}}, true},
		{"fails for really long names", args{nameToGreet: strings.Repeat("Foo", 100), store: staticGreetingStore{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := goclean.NewGreeterFunc(tt.store)
			_, err := g(context.Background(), tt.nameToGreet)
			if tt.wantErr && err != nil {
				return
			} else if !tt.wantErr && err != nil {
				t.Fatalf("want nil err, got %v", err)
			} else if tt.wantErr && err == nil {
				t.Fatalf("want non-nil err, got %v", err)
			}
		})
	}
}
