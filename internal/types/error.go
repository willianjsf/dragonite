package types

// MatrixErrorCode contém os códigos de erros do Matrix
type MatrixErrorCode string

const (
	M_FORBIDDEN    MatrixErrorCode = "M_FORBIDDEN"
	M_BAD_JSON     MatrixErrorCode = "M_BAD_JSON"
	M_NOT_FOUND    MatrixErrorCode = "M_NOT_FOUND"
	M_UNAUTHORIZED MatrixErrorCode = "M_UNAUTHORIZED"
	M_USER_IN_USE  MatrixErrorCode = "M_USER_IN_USE"
)

// ErrorResponse é uma estrutura com o formato de erros padrão do Matrix
type ErrorResponse struct {
	ErrCode MatrixErrorCode `json:"errcode"`
	Message string          `json:"error"`
}

func NewErrorResponse(code MatrixErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		ErrCode: code,
		Message: message,
	}
}
