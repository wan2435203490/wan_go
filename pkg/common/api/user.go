package api

import "wan_go/sdk/pkg/jwtauth/user"

func (a *Api) IsAdmin() bool {
	return user.IsAdmin(a.Context)
}
