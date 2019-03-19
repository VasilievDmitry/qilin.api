package middleware

type RbacRouter interface {
	GetOwner(ctx QilinContext) (string, error)
	GetPermissionsMap() map[string][]string
}