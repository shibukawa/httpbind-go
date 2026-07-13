package app

import "net/http"

// Fake DI / cross-package registration helpers — not analyzed as routes.
// There is no static mux.Handle/HandleFunc with a literal pattern in this package.

func fxInvoke(fn any) {}

func otherPackageRegisterRoutes(mux *http.ServeMux) {}

func boot(mux *http.ServeMux) {
	fxInvoke(otherPackageRegisterRoutes)
	otherPackageRegisterRoutes(mux)
}
