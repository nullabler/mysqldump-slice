package types

import "fmt"


type UidType string

func (str UidType) String() string {
	return fmt.Sprintf("'%s'", string(str))
}
