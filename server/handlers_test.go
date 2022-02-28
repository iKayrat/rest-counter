package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
)

func TestAddCounter(t *testing.T) {
	server := NewTestServer(t)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	w := httptest.NewRecorder()

	expected := 5

	req, err := http.NewRequest("POST", fmt.Sprint(ts.URL, "/rest/counter/add/", expected), nil)
	require.NoError(t, err)

	server.router.ServeHTTP(w, req)

	var actual map[string]int
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err)

	require.Equal(t, 0+expected, actual["counter"])
	require.Equal(t, 200, w.Code)
}

func TestSubCounter(t *testing.T) {
	server := NewTestServer(t)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	w := httptest.NewRecorder()

	expected := 5

	req, err := http.NewRequest("POST", fmt.Sprint(ts.URL, "/rest/counter/sub/", expected), nil)
	require.NoError(t, err)

	server.router.ServeHTTP(w, req)

	var actual map[string]int
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err)

	require.Equal(t, 0-expected, actual["counter"])
	require.Equal(t, 200, w.Code)
}

func TestGetVal(t *testing.T) {
	server := NewTestServer(t)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", ts.URL+"/rest/counter/val", nil)
	require.NoError(t, err)

	server.router.ServeHTTP(w, req)

	expected := 50
	set, err := server.store.Client.Set("counter", expected, 0).Result()
	require.NoError(t, err)
	require.Equal(t, 200, w.Code)
	require.Equal(t, set, "OK")

	res, err := server.store.Client.Get("counter").Int()
	if err == redis.Nil {
		require.Equal(t, 200, w.Code)
		require.EqualError(t, err, "redis: nil")
	} else if err != nil {
		require.Equal(t, 400, w.Code)
		require.EqualError(t, err, "counter doesn't exist")
	}
	require.NoError(t, err)

	require.Equal(t, res, expected)
	require.Equal(t, 200, w.Code)
}

func TestGetSubst(t *testing.T) {
	testserver := NewTestServer(t)

	ts := httptest.NewServer(testserver.router)
	defer ts.Close()

	w := httptest.NewRecorder()

	body := Body{Text: "AaBbCcDadABCDghj"}

	jsonbody, err := json.Marshal(body)
	require.NoError(t, err)
	// log.Println(string(jsonbody), err)

	req, err := http.NewRequest("POST", ts.URL+"/rest/substr/find", bytes.NewBuffer(jsonbody))
	require.NoError(t, err)

	testserver.router.ServeHTTP(w, req)

	actual := RespBody{}
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err)

	expected := RespBody{"adABCDghj", 9}
	require.Equal(t, 200, w.Code)
	require.Equal(t, expected, actual)

}

func TestFindStr(t *testing.T) {
	testserver := NewTestServer(t)

	testCase := []string{"rest", "counter", "find"}

	expected := [][]string{
		{"/rest/counter/add/:i", "/rest/counter/sub/:i", "/rest/substr/find", "/rest/self/find/:str", "/rest/hash/calc", "/rest/email/check/", "/rest/inn/check/", "/rest/counter/val", "/rest/hash/result/:id"},
		{"/rest/counter/add/:i", "/rest/counter/sub/:i", "/rest/counter/val"},
		{"/rest/substr/find", "/rest/self/find/:str"},
	}

	actual := make([]string, 0)

	for i, test := range testCase {
		ts := httptest.NewServer(testserver.router)
		defer ts.Close()

		w := httptest.NewRecorder()

		// log.Println("test:#", i)
		req, err := http.NewRequest("POST", fmt.Sprint(ts.URL, "/rest/self/find/", test), nil)
		if err != nil {
			break
		}
		require.NoError(t, err)

		testserver.router.ServeHTTP(w, req)

		err = json.Unmarshal(w.Body.Bytes(), &actual)
		// log.Println("err", err)
		require.NoError(t, err)

		require.Equal(t, expected[i], actual)
	}
}

func TestCheckEmail(t *testing.T) {
	testserver := NewTestServer(t)

	testCase := BodyEmail{
		Emails: []string{
			"text@text.com",
			"texttext.com",
			"test@test.ru",
			"ttc@",
			"_dfs32452*.&%#^~!@#$%^&*()(*&@gmai.com",
			"safe@life.com",
		},
	}

	type ExpectedBody struct {
		Valid_emails []string `json:"valid_emails"`
	}

	expected := ExpectedBody{}

	expected.Valid_emails = []string{"text@text.com", "test@test.ru", "safe@life.com"}

	jsonbody, err := json.Marshal(testCase)
	require.NoError(t, err)

	// actual := make([]string, 0)
	actual := ExpectedBody{}

	ts := httptest.NewServer(testserver.router)
	defer ts.Close()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("POST", fmt.Sprint(ts.URL, "/rest/email/check/"), bytes.NewBuffer(jsonbody))
	if err != nil {
		t.Error(err)
	}
	require.NoError(t, err)

	testserver.router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
	// }
}
