package net

import (
	"net/http"
	"net/http/httputil"
)

func HttpDumpRequest(request *http.Request) []byte {
	bts, _ := httputil.DumpRequest(request, true)

	return bts
}

func HttpDumpResponse(response *http.Response) []byte {
	bts, _ := httputil.DumpResponse(response, true)

	return bts
}
