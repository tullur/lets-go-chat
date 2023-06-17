package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
)

func TestGreet(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("[GET] /", func() {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			g.Fail(err)
		}

		res := httptest.NewRecorder()
		Greet(res, req)

		g.It("Returns success status code", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})
	})
}
