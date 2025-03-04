package controller

import "book-management/iface"

type Controller struct {
	service iface.Service
}

func New(service iface.Service) *Controller {
	return &Controller{
		service: service,
	}
}
