package entity

import "mysqldump-slice/entity/types"

type Primary[T types.IdInterface] struct {
	keys map[string][]T
}
