package api

func (a *Api) ValidError(err error) error {
	return err
}

//func (a *Api) ErrorMsg(msg string) error {
//	a.ErrorInternal(msg)
//	return errors.New(msg)
//}
