package ports

type Hypervisor interface {
	Installed() bool
	KVMSupported() bool
	Install() error
	AddUserToLibvirtGroup() error
	ConfigureNetwork() error
}
