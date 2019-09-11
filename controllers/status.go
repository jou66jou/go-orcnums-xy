package controllers

import (
	"net/http"

	"github.com/otiai10/marmoset"
)

const version = "test0.0.2"

// Status ...
func Status(w http.ResponseWriter, r *http.Request) {
	marmoset.Render(w, true).JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello!",
		"version": version,
	})
}
