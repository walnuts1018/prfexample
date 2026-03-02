package handler

import (
	"github.com/walnuts1018/PRFExample/server/usecase"
	"github.com/walnuts1018/PRFExample/server/util/random"
)

type Handler struct {
	u    *usecase.Usecase
	rand random.Random
}

func NewHandler(u *usecase.Usecase, rand random.Random) Handler {
	return Handler{
		u:    u,
		rand: rand,
	}
}
