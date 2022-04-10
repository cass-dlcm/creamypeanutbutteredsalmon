package splatnet

import "fmt"

type ErrWeaponsNotFound []string

func (err *ErrWeaponsNotFound) Error() string {
	return fmt.Sprintf("no weaponschedule found: %s, %s, %s, %s", (*err)[0], (*err)[1], (*err)[2], (*err)[3])
}
