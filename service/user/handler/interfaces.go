package handler

import (
	pbusers "github.com/MamangRust/monolith-ecommerce-pb/users"
)

type UserQueryHandler interface {
	pbusers.UserQueryServiceServer
}

type UserCommandHandler interface {
	pbusers.UserCommandServiceServer
}
