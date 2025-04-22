package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/yanun0323/goast"
	"github.com/yanun0323/goast/scope"
)

var (
	file = flag.String("f", "", "file")
)

func main() {
	flag.Parse()

	// XXX: Remove me
	if len(*file) == 0 {
		flag.Set("f", "/Users/Shared/Project/personal/go/ebui/ui_text.go")
	}

	if err := tryKillPreviousProcess(); err != nil {
		log.Fatalf("try kill previous process, err: %+v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("get wd, err: %+v", err)
	}

	moduleName, err := findGoModuleName(wd)
	if err != nil {
		log.Fatalf("find go module name, err: %+v", err)
	}

	relativeFile, err := filepath.Rel(wd, *file)
	if err != nil {
		log.Fatalf("get rel, err: %+v", err)
	}

	ast, err := goast.ParseAst(*file)
	if err != nil {
		log.Fatalf("parse ast, err: %+v", err)
	}

	var pkgName string
	ast.IterScope(func(s goast.Scope) bool {
		if s.Kind() == scope.Package {
			s.Node().IterNext(func(n *goast.Node) bool {
				switch n.Text() {
				case "package", "\t", "\n", " ":
					return true
				default:
					pkgName = n.Text()
					return false
				}
			})
		}
		return pkgName == ""
	})

	if pkgName == "" {
		log.Fatalf("package name not found")
	}

	importPath := moduleName
	if strings.Contains(relativeFile, "/") {
		importPath = fmt.Sprintf("%s/%s", moduleName, relativeFile)
		spImportPath := strings.Split(importPath, "/")
		spImportPath[len(spImportPath)-1] = pkgName
		importPath = strings.Join(spImportPath, "/")
	}

	var fnName string
	ast.IterScope(func(s goast.Scope) bool {
		name, ok := s.GetFuncName()
		if ok && strings.HasPrefix(name, "Preview_") {
			fnName = name
		}
		return fnName == ""
	})

	if fnName == "" {
		log.Fatalf("Preview function not found. The function name must start with Preview_")
	}

	_ = os.RemoveAll(filepath.Join(wd, ".preview"))
	_ = os.MkdirAll(filepath.Join(wd, ".preview"), 0755)

	mainFn := fmt.Sprintf(`
package main

import (
	"github.com/yanun0323/ebui"
	preview "%s"
)

func main() {
	view := preview.%s()
	app := ebui.NewApplication(view)
	app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	app.Run("preview")
}
`, importPath, fnName)

	mainScope, err := goast.ParseScope(0, []byte(mainFn))
	if err != nil {
		log.Fatalf("parse main fn, err: %+v", err)
	}

	newAst, err := goast.NewAst(mainScope...)
	if err != nil {
		log.Fatalf("new ast, err: %+v", err)
	}

	exportFile := filepath.Join(wd, ".preview", "main.go")

	if err := newAst.Save(exportFile, false); err != nil {
		log.Fatalf("save main fn, err: %+v", err)
	}

	cmd := exec.Command("go", "run", exportFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("run main fn, err: %+v", err)
	}

	// view := ebui.Preview_Text()
	// app := ebui.NewApplication(view)
	// app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	// app.Run("preview")
}

func findGoModuleName(wd string) (string, error) {
	goMod, err := os.ReadFile(filepath.Join(wd, "go.mod"))
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(goMod), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module name not found")
}

func tryKillPreviousProcess() error {
	cmd := exec.Command("ps", "aux")
	grepCmd := exec.Command("grep", ".*go-build/.*/main$")

	psOutput, err := cmd.Output()
	if err != nil {
		return errors.Errorf("ps aux, err: %+v", err)
	}

	grepCmd.Stdin = bytes.NewReader(psOutput)
	grepOutput, err := grepCmd.Output()
	if err != nil {
		return nil
	}

	fmt.Println("找到的程序：")
	fmt.Println(string(grepOutput))

	lines := strings.SplitSeq(string(grepOutput), "\n")
	for line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid := fields[1]
		fmt.Printf("找到 PID: %s\n", pid)

		cmd := exec.Command("kill", pid)
		if err := cmd.Run(); err != nil {
			return errors.Errorf("kill, err: %+v", err)
		}

		fmt.Printf("已成功終止 PID %s 的程序\n", pid)
	}

	return nil
}
