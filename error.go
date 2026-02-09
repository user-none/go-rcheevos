package rcheevos

/*
#include "rc_error.h"
*/
import "C"

// Result codes
const (
	OK                    = C.RC_OK
	InvalidFuncOperand    = C.RC_INVALID_FUNC_OPERAND
	InvalidMemoryOperand  = C.RC_INVALID_MEMORY_OPERAND
	InvalidConstOperand   = C.RC_INVALID_CONST_OPERAND
	InvalidFPOperand      = C.RC_INVALID_FP_OPERAND
	InvalidConditionType  = C.RC_INVALID_CONDITION_TYPE
	InvalidOperator       = C.RC_INVALID_OPERATOR
	InvalidRequiredHits   = C.RC_INVALID_REQUIRED_HITS
	DuplicatedStart       = C.RC_DUPLICATED_START
	DuplicatedCancel      = C.RC_DUPLICATED_CANCEL
	DuplicatedSubmit      = C.RC_DUPLICATED_SUBMIT
	DuplicatedValue       = C.RC_DUPLICATED_VALUE
	DuplicatedProgress    = C.RC_DUPLICATED_PROGRESS
	MissingStart          = C.RC_MISSING_START
	MissingCancel         = C.RC_MISSING_CANCEL
	MissingSubmit         = C.RC_MISSING_SUBMIT
	MissingValue          = C.RC_MISSING_VALUE
	InvalidLboardField    = C.RC_INVALID_LBOARD_FIELD
	MissingDisplayString  = C.RC_MISSING_DISPLAY_STRING
	OutOfMemory           = C.RC_OUT_OF_MEMORY
	InvalidValueFlag      = C.RC_INVALID_VALUE_FLAG
	MissingValueMeasured  = C.RC_MISSING_VALUE_MEASURED
	MultipleMeasured      = C.RC_MULTIPLE_MEASURED
	InvalidMeasuredTarget = C.RC_INVALID_MEASURED_TARGET
	InvalidComparison     = C.RC_INVALID_COMPARISON
	InvalidState          = C.RC_INVALID_STATE
	InvalidJSON           = C.RC_INVALID_JSON
	APIFailure            = C.RC_API_FAILURE
	LoginRequired         = C.RC_LOGIN_REQUIRED
	NoGameLoaded          = C.RC_NO_GAME_LOADED
	HardcoreDisabled      = C.RC_HARDCORE_DISABLED
	Aborted               = C.RC_ABORTED
	NoResponse            = C.RC_NO_RESPONSE
	AccessDenied          = C.RC_ACCESS_DENIED
	InvalidCredentials    = C.RC_INVALID_CREDENTIALS
	ExpiredToken          = C.RC_EXPIRED_TOKEN
	InsufficientBuffer    = C.RC_INSUFFICIENT_BUFFER
	InvalidVariableName   = C.RC_INVALID_VARIABLE_NAME
	UnknownVariableName   = C.RC_UNKNOWN_VARIABLE_NAME
	NotFound              = C.RC_NOT_FOUND
)

// ErrorString returns the error message for a result code
func ErrorString(result int) string {
	return C.GoString(C.rc_error_str(C.int(result)))
}
