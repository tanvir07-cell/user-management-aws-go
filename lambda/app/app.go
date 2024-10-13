package app

import "usermanagement.tanvirrifat.io/api"

type App struct{
	ApiHandler api.ApiHandler
}

func NewApp() App{
	return App{
		ApiHandler: api.NewApiHandler(),
	}
}