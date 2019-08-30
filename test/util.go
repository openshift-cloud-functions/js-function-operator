package test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func AssertGetRequest(t *testing.T, url string, expectedStatusCode int, expectedBody []byte) {
	res, err := http.Get(url)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	assert.Equal(t, expectedStatusCode, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	assert.Equal(t, expectedBody, b)
}
