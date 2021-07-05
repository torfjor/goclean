package goclean

import (
	"context"
	"fmt"
	pb "goclean/gen"
	"net/http"
	"strings"
)

// Greeting is a greeting.
type Greeting string

// GreeterFunc is a function that greets a person with a random Greeting.
type GreeterFunc func(ctx context.Context, name string) (Greeting, error)

// GreetingStore provides functionality for retrieving greeting templates.
type GreetingStore interface {
	RandomGreetingTemplate(ctx context.Context) (string, error)
}

// NewGreeterFunc returns a configured GreeterFunc. We put all business logic
// here where it's separated from transport and storage related concerns.
func NewGreeterFunc(store GreetingStore) GreeterFunc {
	return func(ctx context.Context, name string) (Greeting, error) {
		if name == "" {
			return "", fmt.Errorf("Cannot greet a person without a name!")
		}

		if len(name) > 128 {
			return "", fmt.Errorf("Cannot greet people with really long names")
		}

		if strings.HasPrefix(name, "King") || strings.HasPrefix(name, "Queen") {
			return "", fmt.Errorf("Cannot greet kings or queens")
		}

		tmpl, err := store.RandomGreetingTemplate(ctx)
		if err != nil {
			return "", err
		}

		return Greeting(fmt.Sprintf(tmpl, name)), nil
	}
}

// ServeHTTP implements http.ServeHTTP.
func (g GreeterFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	greeting, err := g(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(greeting))
}

// Greet implements pb.GreetingServiceServer.
func (g GreeterFunc) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	greeting, err := g(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.GreetResponse{Greeting: string(greeting)}, nil
}
