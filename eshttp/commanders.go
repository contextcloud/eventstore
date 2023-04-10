package eshttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/contextcloud/eventstore/es"
)

type CommandNameExtractor func(r *http.Request) (string, error)

func DefaultCommandNameExtractor(r *http.Request) (string, error) {
	if r.URL.Path == "" {
		return "", fmt.Errorf("invalid path for CommandNameExtractor")
	}

	ind := strings.LastIndex(r.URL.Path, "/")
	if ind == -1 {
		return r.URL.Path, nil
	}
	return r.URL.Path[ind+1:], nil
}

func NewCommanders(extractor CommandNameExtractor) func(w http.ResponseWriter, r *http.Request) {
	ext := extractor
	if ext == nil {
		ext = DefaultCommandNameExtractor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		name, err := ext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		unit, err := es.GetUnit(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cmd, err := unit.CreateCommand(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := es.Dispatch(ctx, cmd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
