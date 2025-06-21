package v

import pbusers "github.com/nurfianqodar/school-microservices/services/users/pb/users/v1"

var r *pbusers.CreateOneUserRequest

// Validation rules
var (
	ruleCreateOneUserRequest = map[string]string{
		"Email":    "required,email,max=255",
		"Password": "required,min=8",
		"Role":     "required",
	}
	ruleUpdateOnePasswordUserRequest = map[string]string{}
	ruleUpdateOneEmailUserRequest    = map[string]string{}
	ruleUpdateOneRoleUserRequest     = map[string]string{}
)
