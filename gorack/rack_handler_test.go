package gorack

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testBody = "\x00OMG, test body\nhalps"

func TestRackHandler(t *testing.T) {

	// package variable
	gorackRunner = "../ruby/gorack"

	ts := httptest.NewServer(NewRackHandler("../ruby/config_test.ru"))
	defer ts.Close()

	body := strings.NewReader(testBody)

	res, err := http.Post(ts.URL, "text/plain", body)
	if err != nil {
		t.Error(err)
	}

	got, err := ioutil.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		t.Error(err)
	}

	exp := testBody

	if exp != string(got) {
		t.Errorf("\nGot:%s\nExp:%s", got, exp)
	}
}

func TestIpcEcho(t *testing.T) {
	t.SkipNow()
}
