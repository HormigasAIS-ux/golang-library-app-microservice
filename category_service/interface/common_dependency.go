package interface_pkg

import (
	ucase "category_service/usecase"
)

type CommonDependency struct {
	CategoryUcase ucase.ICategoryUcase
}
