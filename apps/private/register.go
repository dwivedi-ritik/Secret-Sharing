package private

import "net/http"

func PrivateRouter(pathPrefix string) *http.ServeMux {
	privateRouter := http.NewServeMux()

	privateRouter.HandleFunc("GET "+pathPrefix+"getUserInfo", Authorization(GetUserInfoController))
	privateRouter.HandleFunc("GET "+pathPrefix+"getUserMessages", Authorization(GetUserMessagesController))

	privateRouter.HandleFunc("POST "+pathPrefix+"message/encrypt", Authorization(EncryptMessageController))
	privateRouter.HandleFunc("POST "+pathPrefix+"message/decrypt", Authorization(DecryptMessageController))

	return privateRouter
}
