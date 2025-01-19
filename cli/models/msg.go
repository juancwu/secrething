package models

type errMsg struct {
	err error
}

func newErrMsg(err error) errMsg {
	return errMsg{err: err}
}

func (e errMsg) Err() error {
	return e.err
}

type navigationMsg struct {
	To page
}

func newNavigationMsg(to page) navigationMsg {
	return navigationMsg{To: to}
}
