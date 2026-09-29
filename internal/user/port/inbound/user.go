package inbound

type IUserUseCase interface {
	IGetUserUseCase
	IListUsersUseCase
	IUpdateUserUseCase
	IDeleteUserUseCase
	ICountUsersUseCase
}
