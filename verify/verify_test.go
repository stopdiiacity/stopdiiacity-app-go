package verify

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	var (
		urls = []string{
			"https://jobs.dou.ua/companies/allright/reviews",
			"https://djinni.co/jobs/company-epam-systems-bb0df/",
		}
	)

	for _, url := range urls {
		var request = VerifyRequest{
			URLs:    []string{url},
			Version: 1,
		}

		var body, marshalErr = json.Marshal(request)
		require.NoError(t, marshalErr)

		var response = Verify(body)
		require.Equal(t, Response{unsafeMessage, http.StatusOK}, response)
	}
}
