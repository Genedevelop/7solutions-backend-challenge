package inbound

type IAuthUseCase interface {
	IRegisterUseCase
	ILoginUseCase
}
