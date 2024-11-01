package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/dbofmmbt/builder/example"
	"github.com/dbofmmbt/builder/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		srcFile := args[0]
		targetStructs := args[1:]

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, srcFile, nil, parser.Mode(0))
		if err != nil {
			panic(err)
		}
		v := &internal.Visitor{
			DesiredStructs: targetStructs,
		}
		ast.Walk(v, node)

		if v.Builders == nil {
			fmt.Println("no builder struct found")
			return
		}

		g := internal.NewGenerator()

		builderFileName := fmt.Sprint(strings.TrimSuffix(srcFile, path.Ext(srcFile)), ".builder.go")

		builderFile, err := os.OpenFile(builderFileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		g.Generate(v.Builders, builderFile)
	},
}

func Execute() {
	rootCmd.AddCommand(&cobra.Command{
		Use: "example",
		Run: func(cmd *cobra.Command, args []string) {
			example.Example()

		},
	})
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
