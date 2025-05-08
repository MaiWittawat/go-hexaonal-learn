// port/auth/authZ.go
package authPort

type JwtAuthZService interface {
	Authorize(email string, roles ...string) (bool, error)
}
