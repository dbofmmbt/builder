package internal

import (
	"go/ast"
	"reflect"
	"slices"
	"strings"

	"github.com/samber/lo"
)

type Visitor struct {
	DesiredStructs []string
	Builders       []Builder
	packageName    string
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	// add case to get package name
	case *ast.File:
		v.packageName = n.Name.Name
	case *ast.TypeSpec:
		switch structType := n.Type.(type) {
		case *ast.StructType:
			if !slices.Contains(v.DesiredStructs, n.Name.Name) {
				// continue searching
				return v
			}

			v.Builders = append(v.Builders, Builder{
				packagee: v.packageName,
				structs: []Struct{
					{
						name: n.Name.Name,
						fields: lo.Map(structType.Fields.List, func(f *ast.Field, _ int) Field {
							tag := reflect.StructTag("")
							if f.Tag != nil {
								tagRaw := f.Tag.Value
								tag = reflect.StructTag(tagRaw[1 : len(tagRaw)-1])
							}

							builderTag := tag.Get("builder")

							return Field{
								name:  f.Names[0].Name,
								typee: f.Type.(*ast.Ident).Name,
								options: Options{
									skip: strings.Contains(builderTag, "skip"),
								},
							}
						}),
					},
				},
			})
		}
	}
	return v
}
