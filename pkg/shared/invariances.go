package shared

import (
	"fmt"
)

type ContractKind string

const (
	KindInvariance    ContractKind = "invariance"
	KindPreCondition  ContractKind = "pre-condition"
	KindPostCondition ContractKind = "post-condition"
	KindUnreachable   ContractKind = "unreachable"
)

// Violation is the panic value raised when a contract check fails.
type Violation struct {
	Kind ContractKind
	Msg  string
}

func (v *Violation) Error() string {
	return fmt.Sprintf("%s broken: %s", v.Kind, v.Msg)
}

func Assert(check bool, msg string) {
	if !check {
		panic(&Violation{KindInvariance, msg})
	}
}

func Assertf(check bool, format string, args ...any) {
	if !check {
		panic(&Violation{KindInvariance, fmt.Sprintf(format, args...)})
	}
}

func PreCondition(check bool, msg string) {
	if !check {
		panic(&Violation{KindPreCondition, msg})
	}
}

func PreConditionf(check bool, format string, args ...any) {
	if !check {
		panic(&Violation{KindPreCondition, fmt.Sprintf(format, args...)})
	}
}

func PostCondition(check bool, msg string) {
	if !check {
		panic(&Violation{KindPostCondition, msg})
	}
}

func PostConditionf(check bool, format string, args ...any) {
	if !check {
		panic(&Violation{KindPostCondition, fmt.Sprintf(format, args...)})
	}
}

func Unreachable(msg string) {
	panic(&Violation{KindUnreachable, msg})
}

func Unreachablef(format string, args ...any) {
	panic(&Violation{KindUnreachable, fmt.Sprintf(format, args...)})
}
