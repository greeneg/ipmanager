package model

type InvalidStatusValue struct {
	Err error
}

func (i *InvalidStatusValue) Error() string {
	return "Invalid value! Must be either 'enabled' or 'locked'"
}

type PasswordHashMismatch struct {
	Err error
}

func (p *PasswordHashMismatch) Error() string {
	return "Password hashes do not match!"
}
