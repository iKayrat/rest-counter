package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

	body := Substr{Text: "AaBbCcDadABCDghj"}

	jsonbody, err := json.Marshal(body)
	require.NoError(t, err)
	log.Println(string(jsonbody), err)

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
