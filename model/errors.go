package model

type InvalidStatusValue struct {
	Err error
}

func (i *InvalidStatusValue) Error() string {
	return "Invalid value! Must be either 'enabled' or 'locked'"
}

type AddressTableInUse struct {
	Err error
}

func (p *AddressTableInUse) Error() string {
	return "Address table already in use. Cannot mutate subnet"
}

type PasswordHashMismatch struct {
	Err error
}

func (p *PasswordHashMismatch) Error() string {
	return "Password hashes do not match!"
}
