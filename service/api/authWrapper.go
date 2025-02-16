package api

import "github.com/julienschmidt/httprouter"

// wrapWithAuth applies the authentication middleware to a given httprouter.Handle.
// It ensures that only authenticated users can access the specified handler.
func (rt *_router) wrapWithAuth(handle httprouter.Handle) httprouter.Handle {
	return AuthMiddleware(rt.db, handle)
}
