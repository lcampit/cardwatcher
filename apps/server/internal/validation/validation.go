package validation

import (
	"errors"
	"fmt"
	"strings"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"buf.build/go/protovalidate"
	apperrors "github.com/lcampit/cardwatcher/apps/server/internal/errors"
	errorspb "github.com/lcampit/cardwatcher/gen/go/errors/v1"
	"google.golang.org/protobuf/proto"
)

var pv protovalidate.Validator

func init() {
	var err error
	pv, err = protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize protovalidator: %v", err))
	}
}

// ValidateRequest checks a proto message and returns an AppError with all violations.
func ValidateRequest(req proto.Message, traceID string) error {
	if err := pv.Validate(req); err != nil {
		if validationErrs, ok := errors.AsType[*protovalidate.ValidationError](err); ok {
			var violations []*errorspb.FieldViolation
			for _, violation := range validationErrs.Violations {
				fieldPath := ""
				if violation.Proto.GetField() != nil {
					fieldPath = formatFieldPath(violation.Proto.GetField())
				}

				ruleID := violation.Proto.GetRuleId()
				message := violation.Proto.GetMessage()

				violations = append(violations, &errorspb.FieldViolation{
					Field:       fieldPath,
					Description: message,
					Code:        mapConstraintToCode(ruleID),
				})
			}
			return apperrors.NewValidationFailed(violations, traceID)
		}
		return apperrors.NewInternal("Validation failed", traceID, err)
	}
	return nil
}

// Helper functions
func formatFieldPath(fieldPath *validate.FieldPath) string {
	if fieldPath == nil {
		return ""
	}

	// Build field path from elements
	var parts []string
	for _, element := range fieldPath.GetElements() {
		if element.GetFieldName() != "" {
			parts = append(parts, element.GetFieldName())
		} else if element.GetFieldNumber() != 0 {
			parts = append(parts, fmt.Sprintf("field_%d", element.GetFieldNumber()))
		}
	}

	return strings.Join(parts, ".")
}

func mapConstraintToCode(ruleID string) string {
	switch {
	case strings.Contains(ruleID, "required"):
		return errorspb.AppErrorCode_APP_ERROR_CODE_REQUIRED_FIELD.String()
	default:
		return errorspb.AppErrorCode_APP_ERROR_CODE_INVALID_VALUE.String()
	}
}
