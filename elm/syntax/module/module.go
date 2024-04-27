package module

type Module interface {
	module() _Module
}

type _Module struct{}

type Normal struct {
	_Module
	data ModuleData
}

type Port struct {
	_Module
	data ModuleData
}

func (m _Module) module() _Module {
	return m
}

type ModuleData struct {
	moduleName ModuleName
}

func With[R any](
	module Module,
	norm func(*Normal) R,
	port func(*Port) R,
) R {
	switch d := module.(type) {
	case Normal:
		return norm(&d)
	case Port:
		return port(&d)
	}
	panic("unreachable")
}
