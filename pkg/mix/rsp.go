package mix

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wolfbolin/bolbox/pkg/errors"
	"github.com/wolfbolin/bolbox/pkg/log"
)

type HttpRspWarp struct {
	W http.ResponseWriter
}

func HttpRsp(w http.ResponseWriter) *HttpRspWarp {
	return &HttpRspWarp{
		W: w,
	}
}

func (h *HttpRspWarp) Header(key, val string) *HttpRspWarp {
	h.W.Header().Set(key, val)
	return h
}

func (h *HttpRspWarp) Code(code int) *HttpRspWarp {
	h.W.WriteHeader(code)
	return h
}

func (h *HttpRspWarp) Json(a any) *HttpRspWarp {
	err := json.NewEncoder(h.W).Encode(a)
	if err != nil {
		log.Errorf("Write http response failed: %s", errors.WithStack(err))
		http.Error(h.W, err.Error(), http.StatusInternalServerError)
	}
	return h
}

func (h *HttpRspWarp) Text(format string, a ...any) *HttpRspWarp {
	_, err := fmt.Fprintf(h.W, format, a...)
	if err != nil {
		log.Errorf("Write http response failed: %s", errors.WithStack(err))
		http.Error(h.W, err.Error(), http.StatusInternalServerError)
	}
	return h
}

func (h *HttpRspWarp) BadRequest(e error) *HttpRspWarp {
	log.Warnf("BadRequest: %+v", e)
	h.Code(http.StatusBadRequest)
	h.Text("%s", e.Error())
	return h
}

func (h *HttpRspWarp) ServerError(e error) *HttpRspWarp {
	log.Errorf("InternalServerError: %+v", e)
	h.Code(http.StatusInternalServerError)
	h.Text("%s", e.Error())
	return h
}
