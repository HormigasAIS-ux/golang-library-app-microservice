package interface_pkg

import (
	ucase "book_service/usecase"
)

type CommonDependency struct {
	BookUcase ucase.IBookUcase
}
