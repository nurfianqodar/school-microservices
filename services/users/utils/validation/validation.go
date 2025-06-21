package v

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTrans "github.com/go-playground/validator/v10/translations/id"
	pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"
)

var (
	Validate *validator.Validate
	Trans    ut.Translator
)

func init() {

	// Init translator
	idLoc := id.New()
	enLoc := en.New()
	uni := ut.New(idLoc, idLoc, enLoc)

	Trans, _ = uni.GetTranslator("id")
	Validate = validator.New()

	idTrans.RegisterDefaultTranslations(Validate, Trans)

	// Register validation rules for gRPC requests
	Validate.RegisterStructValidationMapRules(ruleCreateOneUserRequest, pbusers.CreateOneUserRequest{})
	Validate.RegisterStructValidationMapRules(ruleUpdateOneEmailUserRequest, pbusers.UpdateOneEmailUserRequest{})
	Validate.RegisterStructValidationMapRules(ruleUpdateOnePasswordUserRequest, pbusers.UpdateOnePasswordUserRequest{})
	Validate.RegisterStructValidationMapRules(ruleUpdateOneRoleUserRequest, pbusers.UpdateOneRoleUserRequest{})
}
