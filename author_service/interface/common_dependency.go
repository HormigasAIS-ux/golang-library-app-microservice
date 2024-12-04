package interface_pkg

import (
	ucase "author_service/usecase"
)

type CommonDependency struct {
	AuthorUcase ucase.IAuthorUcase
}
