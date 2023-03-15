package delivery

import (
	"encoding/json"
	"loadBalancer/internal/config"
	"net/http"
)

func (obj *adminHandler) describeConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	configBytes, err := json.Marshal(obj.adminUsecase.DescribeConfig(ctx))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(configBytes)
}

func (obj *adminHandler) storeConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newConfig config.ServicePool
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = obj.adminUsecase.StoreConfig(ctx, newConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
