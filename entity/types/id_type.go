package types

import "strconv"


type IdType int

func (id IdType) String() string {
	return strconv.Itoa(int(id))
}
