package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/stopdiiacity/stopdiiacity-app-go/verify"
)

type Empty = struct{}

type CountResponse struct {
	Count uint64 `json:"count"`
}

var (
	_ verify.VerifyRequest
	_ verify.VerifyResponse
)

var (
	count uint64 = 0
)

// VerifyHandler handler
//
// @Accept       json
// @Produce      json
// @Description  Verify
// @Param        verify.VerifyRequest  body      verify.VerifyRequest  true  "request"
// @Success      200                   {object}  verify.VerifyResponse
// @Failure      400  {object}  Empty
// @Router       /verify.json [post]
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	var content, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	atomic.AddUint64(&count, 1)

	var response = verify.Verify(content)
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
}

// CountHandler handler
//
// @Accept       json
// @Produce      json
// @Description  Count
// @Success      200  {object}  CountResponse
// @Failure      400  {object}  Empty
// @Router       /count.json [get]
func CountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	var content, err = json.Marshal(&CountResponse{
		Count: atomic.LoadUint64(&count),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// LinksHandler handler
//
// @Accept       json
// @Produce      json
// @Description  Companies links
// @Success      200  {object}  verify.LinksResponse
// @Failure      400                   {object}  Empty
// @Router       /links.json [get]
func LinksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	var response = &verify.LinksResponse{
		Groups: verify.Prefixes(),
	}

	var content, err = json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
