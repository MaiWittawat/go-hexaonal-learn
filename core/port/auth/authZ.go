// port/auth/authZ.go
package port

type JwtAuthZService interface {
	Authorize(email string, roles []string) (bool, error)
}
