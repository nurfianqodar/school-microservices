package errs

import (
	"log"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertValidationError(e validator.ValidationErrors, trans ut.Translator) error {
	st := status.New(codes.InvalidArgument, "invalid input data")
	fieldViolations := make([]*epb.BadRequest_FieldViolation, 0, len(e))
	for i, fieldError := range e {
		fieldViolations[i] = &epb.BadRequest_FieldViolation{
			Field:       fieldError.Field(),
			Description: fieldError.Translate(trans),
			Reason:      fieldError.Translate(trans),
		}
	}
	ds, err := st.WithDetails(&epb.BadRequest{
		FieldViolations: fieldViolations,
	})

	if err != nil {
		log.Printf("error: failed to create error detail. %s", err.Error())
		return status.Error(codes.Internal, "internal server error")
	}

	return ds.Err()
}
