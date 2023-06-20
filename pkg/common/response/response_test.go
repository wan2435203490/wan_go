package r

//func TestSuccess(t *testing.T) {
//	success := Success0("a")
//	fmt.Printf("%#v\n", success)
//	success2 := Success0(1)
//	fmt.Printf("%#v\n", success2)
//	success3 := Success0(true)
//	fmt.Printf("%#v\n", success3)
//}
//
//func TestError(t *testing.T) {
//	err1 := Error0(http.StatusPermanentRedirect, "err1")
//	fmt.Printf("%#v\n", err1)
//	err2 := Error0(http.StatusUnsupportedMediaType, "err2")
//	fmt.Printf("%#v\n", err2)
//	err3 := Error0(http.StatusServiceUnavailable, "")
//	fmt.Printf("%#v\n", err3)
//}

// Success0 discard
//func Success0(data any) *response {
//
//	res := &response{
//		Data: data,
//	}
//
//	res.Message = http.StatusText(http.StatusOK)
//	res.Status = http.StatusOK
//	res.Code = SuccessStatus
//
//	return res
//}
//
//// Error0 discard
//func Error0(httpStatus int, msg string) *response {
//	res := &response{}
//	res.Message = utils.IfThen(len(msg) > 0, msg, http.StatusText(httpStatus)).(string)
//	res.Status = httpStatus
//	res.Code = ErrorStatus
//
//	return res
//}
