package confread

type ConfigStruct struct {
	Address  string
	Secure   bool
	CertFile string
	KeyFile  string
	Interval int
}
