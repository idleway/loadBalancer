package servicePool

import (
	"bytes"
	"io"
	"loadBalancer/internal/bytesPool"
	cacheModel "loadBalancer/internal/cache"
	"net/http"
	"strconv"
)

func (obj *servicePool) LBFunc(w http.ResponseWriter, req *http.Request) {
	cacheEnabled := obj.cacheEnabled.Load()
	if cacheEnabled {
		reqBodyBuf := bytesPool.Pool.Get().(*bytes.Buffer)
		defer reqBodyBuf.Reset()
		defer bytesPool.Pool.Put(reqBodyBuf)

		if req.Body != nil {
			_, err := io.Copy(reqBodyBuf, req.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBuf.Bytes()))
		}
		cacheKey, err := cacheModel.CalculateKey(cacheModel.CalculateKeyData{
			ReqBody: reqBodyBuf,
			Path:    req.RequestURI,
			Headers: req.Header,
		})
		if err == nil {
			oldResponse, oldResponseExist := obj.cache.Describe(cacheKey)
			if oldResponseExist {
				_, _ = w.Write(oldResponse)
				return
			} else {
				req.Header.Add(cacheModel.CacheKeyCustomHeader, strconv.FormatUint(cacheKey, 10))
			}
		}
	}

	instance := obj.DescribeNextServiceInstance()
	if instance == nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	instance.ServeHTTP(w, req, cacheEnabled)
}
