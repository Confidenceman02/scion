package file

import (
	"scion/elm/syntax/module"
	"scion/elm/syntax/node"
)

type File struct {
	moduleDefinition node.Node[module.Module]
}
