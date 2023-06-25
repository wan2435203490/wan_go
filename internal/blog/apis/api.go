package apis

// 会有并发bug
// var (
//
//	a blogApi
//
// )
//type blogApi struct {
//	api.Api
//}

//func MakeContext(a *api.Api) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if a.MakeContext(c) != nil {
//			c.Abort()
//			return
//		}
//
//		c.Next()
//	}
//}
