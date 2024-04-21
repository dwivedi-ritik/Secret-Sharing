package auth

import "net/http"

func UserRouter(pathPrefix string) *http.ServeMux {
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("POST "+pathPrefix+"login", LoginUserController)
	userRouter.HandleFunc("POST "+pathPrefix+"signup", AddUserController)
	return userRouter
}
