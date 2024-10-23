package standard

import "time"

type SwaggerSuccessResponse[D any] struct {
	Status string `json:"status" example:"success"`
	Data   D      `json:"data"`
}

type SwaggerMetadata struct {
	Timestamp time.Time `json:"timestamp" example:"2024-10-23T13:52:29.459752901+07:00"`
	RequestID string    `json:"request_id" example:"IIJORWQpIvDMzzNf"`
}

type SwaggerBadRequestErrorResponse struct {
	Status           string          `json:"status" example:"error"`
	Error            string          `json:"error" example:"invalid_request"`
	ErrorDescription string          `json:"error_description" example:"invalid field: require string but got int"`
	Metadata         SwaggerMetadata `json:"metadata"`
}

type SwaggerNotFoundErrorResponse struct {
	Status           string          `json:"status" example:"error"`
	Error            string          `json:"error" example:"not_found"`
	ErrorDescription string          `json:"error_description" example:"not found user with id 323979471029873"`
	Metadata         SwaggerMetadata `json:"metadata"`
}

type SwaggerDuplicatedErrorResponse struct {
	Status           string          `json:"status" example:"error"`
	Error            string          `json:"error" example:"duplicated"`
	ErrorDescription string          `json:"error_description" example:"something has already existed"`
	Metadata         SwaggerMetadata `json:"metadata"`
}

type SwaggerInvalidCredentialsErrorResponse struct {
	Status           string          `json:"status" example:"error"`
	Error            string          `json:"error" example:"invalid_credentials"`
	ErrorDescription string          `json:"error_description" example:"username or password is invalid"`
	Metadata         SwaggerMetadata `json:"metadata"`
}

type SwaggerUnauthorizedErrorResponse struct {
	Status           string          `json:"status" example:"error"`
	Error            string          `json:"error" example:"forbidden"`
	ErrorDescription string          `json:"error_description" example:"require authentication"`
	Metadata         SwaggerMetadata `json:"metadata"`
}

type SwaggerForbiddenErrorResponse struct {
	Status           string          `json:"status" example:"error"`
	Error            string          `json:"error" example:"forbidden"`
	ErrorDescription string          `json:"error_description" example:"not enough permission to access"`
	Metadata         SwaggerMetadata `json:"metadata"`
}
