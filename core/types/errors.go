package types

import (
	"fmt"
	"runtime"
)

type ErrWeaponsNotFound struct{}

func (*ErrWeaponsNotFound) Error() string { return "no weaponschedule found" }

type ErrStrStageNotFound struct{ Stage string }

func (err *ErrStrStageNotFound) Error() string { return fmt.Sprintf("no stage found: %s", err.Stage) }

type ErrIntEventNotFound struct{ Event int }

func (err *ErrIntEventNotFound) Error() string { return fmt.Sprintf("no event found: %d", err.Event) }

type ErrStrEventNotFound struct{ Event string }

func (err *ErrStrEventNotFound) Error() string { return fmt.Sprintf("no event found: %s", err.Event) }

type ErrIntTideNotFound struct{ Tide int }

func (err *ErrIntTideNotFound) Error() string { return fmt.Sprintf("no tide found: %d", err.Tide) }

type ErrStrTideNotFound struct{ Tide string }

func (err *ErrStrTideNotFound) Error() string { return fmt.Sprintf("no tide found: %s", err.Tide) }

type StackTrace struct {
	buf       []byte
	stackSize int
}

func NewStackTrace() *StackTrace {
	buf := make([]byte, 1<<16)
	return &StackTrace{buf, runtime.Stack(buf, false)}
}

func (st *StackTrace) Error() string {
	return string(st.buf[0:st.stackSize])
}

type ErrStrWeaponsNotFound struct{ Weapons string }

func (err ErrStrWeaponsNotFound) Error() string {
	return fmt.Sprintf("no weaponschedule found: %s", err.Weapons)
}
