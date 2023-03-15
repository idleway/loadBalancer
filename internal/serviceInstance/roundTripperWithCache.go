package serviceInstance

import (
	"bytes"
	"io"
	"loadBalancer/internal/bytesPool"
	cacheModel "loadBalancer/internal/cache"
	"net/http"
	"strconv"
)

type transportWithCache struct {
	http.RoundTripper
	cacheChan chan cacheModel.StoreInCache
}

func (t *transportWithCache) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	cacheKey, err := strconv.ParseUint(req.Header.Get(cacheModel.CacheKeyCustomHeader), 10, 64)
	if err != nil {
		return resp, nil
	}

	freeBuffers := true
	respBodyBuf := bytesPool.Pool.Get().(*bytes.Buffer)
	defer func() {
		if freeBuffers {
			respBodyBuf.Reset()
			bytesPool.Pool.Put(respBodyBuf)
		}
	}()

	if resp.Body != nil {
		_, err = io.Copy(respBodyBuf, resp.Body)
		if err != nil {
			return
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBuf.Bytes()))
	}

	freeBuffers = false
	go func() {
		t.cacheChan <- cacheModel.StoreInCache{
			Key:      cacheKey,
			RespBody: respBodyBuf,
		}
	}()

	return resp, nil
}
