package hypervisor

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Installed() bool              { return Installed() }
func (a *Adapter) KVMSupported() bool           { return KVMSupported() }
func (a *Adapter) Install() error               { return Install() }
func (a *Adapter) AddUserToLibvirtGroup() error { return AddUserToLibvirtGroup() }
func (a *Adapter) ConfigureNetwork() error      { return ConfigureNetwork() }
