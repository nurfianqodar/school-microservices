package httperr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nurfianqodar/school-microservices/utils/httpres"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidRequestBody = New(http.StatusBadRequest, "invalid request body")
	ErrInternalServer     = New(http.StatusInternalServerError, "internal server error")
)

type HTTPErr interface {
	error
	Send(w http.ResponseWriter)
}

type httperr struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

func (e *httperr) Error() string {
	return e.Message
}

func New(code int, message string, detail ...any) HTTPErr {
	var finalDetail any
	if len(detail) == 0 {
		finalDetail = nil
	} else if len(detail) == 1 {
		finalDetail = detail[0]
	} else {
		finalDetail = detail
	}
	return &httperr{
		Code:    code,
		Message: message,
		Detail:  finalDetail,
	}
}

func (e *httperr) Send(w http.ResponseWriter) {
	w.WriteHeader(e.Code)
	if err := json.NewEncoder(w).Encode(httpres.New(false, e)); err != nil {
		log.Printf("error: unable to write error response. %s", err.Error())
		fmt.Fprintln(w, "unable to write response")
	}
}

func ConvertGRPCErrorToHTTPErr(err error) HTTPErr {

	var httpCode int
	var finalDetail any

	st := status.Convert(err)
	message := st.Message()
	detail := make([]any, 0, len(st.Details()))

	// Convert code
	switch st.Code() {
	case codes.AlreadyExists:
		httpCode = http.StatusConflict
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
	case codes.Aborted:
		httpCode = http.StatusConflict
	default:
		httpCode = http.StatusInternalServerError
	}

	for _, d := range st.Details() {
		switch info := d.(type) {
		case *epb.BadRequest:
			detail = append(detail, info.GetFieldViolations())
		case *epb.PreconditionFailure:
			detail = append(detail, info.GetViolations())
		case *epb.QuotaFailure:
			detail = append(detail, info.GetViolations())
		}
	}

	if len(detail) == 0 {
		finalDetail = nil
	} else if len(detail) == 1 {
		finalDetail = detail[0]
	} else {
		finalDetail = detail
	}

	return &httperr{
		Code:    httpCode,
		Message: message,
		Detail:  finalDetail,
	}

}
