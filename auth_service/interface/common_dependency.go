package interface_pkg

import ucase "auth_service/usecase"

type CommonDependency struct {
	AuthUcase ucase.IAuthUcase
	UserUcase ucase.IUserUcase
}
