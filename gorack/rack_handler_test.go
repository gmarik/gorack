package gorack

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRackHandler(t *testing.T) {

	// package variable
	gorackRunner = "../ruby/gorack"

	ts := httptest.NewServer(NewRackHandler("../ruby/config_test.ru"))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}

	got, err := ioutil.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		t.Error(err)
	}

	exp := "hellozzzz"

	if exp != string(got) {
		t.Errorf("\nGot:%s\nExp:%s", got, exp)
	}
}
