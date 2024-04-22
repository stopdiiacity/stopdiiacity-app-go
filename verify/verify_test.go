package verify

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	f := func(t *testing.T, url string, expected string) {
		t.Helper()

		var request = VerifyRequest{
			URLs:    []string{url},
			Version: 1,
		}

		var body, marshalErr = json.Marshal(request)
		require.NoError(t, marshalErr)

		var response = Verify(body)
		require.Equal(t, Response{expected, http.StatusOK}, response)
	}

	f(t, "https://jobs.dou.ua/companies/allright/reviews", unsafeMessage)
	f(t, "https://djinni.co/jobs/?company=epam-systems-bb0df", safeMessage)
}
