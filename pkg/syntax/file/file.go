package file

import (
	"scion/pkg/syntax/module"
	"scion/pkg/syntax/node"
)

type File struct {
	moduleDefinition node.Node[module.Module]
}
