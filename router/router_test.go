package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var routeTests = []struct {
	name      string
	routMethod string
	routPatn   string
	reqMethod  string
	reqURL    string
	want      int
}{
	{
		name: "Should respond with a status 200 OK when a route and method match",
		routPatn: "/tasks",
		routMethod:http.MethodGet,
		reqURL:  "/tasks",
		reqMethod:http.MethodGet,
		want:http.StatusOK,
	},
	{
		name: "Should respond with a status 404 Not Found when HTTP method is different",
		routPatn: "/tasks",
		routMethod:http.MethodGet,
		reqURL:  "/tasks",
		reqMethod:http.MethodPost,
		want: http.StatusNotFound,
	},
	{
		name: "Should respond with a status 200 OK when a route match regex and method",
		routPatn: `/tasks/\d`,
		routMethod:http.MethodGet,
		reqURL:  "/tasks/1",
		reqMethod:http.MethodGet,
		want: http.StatusOK,
	},
	{
		name: "Should respond with a status 404 Not Found when route could not be found",
		routPatn: `/tasks\d`,
		routMethod:http.MethodPost,
		reqURL:  "/tasks/a",
		reqMethod:http.MethodPost,
		want: http.StatusNotFound,
	},
}

func TestRoute(t *testing.T) {

	t.Log("Router")

	for _, testcase := range routeTests {
		t.Logf(testcase.name)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(testcase.reqMethod, testcase.reqURL, nil)

		r := Router{}

		r.HandleFunc(testcase.routPatn, testcase.routMethod, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		r.ServeHTTP(rec, req)

		if rec.Code != testcase.want {
			t.Errorf("KO => Got %d wanted %d", rec.Code, testcase.want)
		}
	}
}
