package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
	"github.com/yanun0323/goast"
	"github.com/yanun0323/goast/scope"
)

var (
	file   = flag.String("f", "", "ebui go file with functions starting with Preview_")
	helper = flag.Bool("h", false, "show help")
	debug  = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()

	if *helper {
		flag.Usage()
		return
	}

	if *file == "" {
		log.Print("-f is required")
		flag.Usage()
		return
	}

	windowRect := tryGetWindowPosition()

	if err := tryKillPreviousProcess(); err != nil {
		fatal("try kill previous process, err: %+v", err)
		return
	}

	wd, err := findProjectRoot(file)
	if err != nil {
		fatal("find project root, err: %+v", err)
		return
	}

	moduleName, err := findGoModuleName(wd)
	if err != nil {
		fatal("find go module name, err: %+v", err)
		return
	}

	relativeFile, err := filepath.Rel(wd, *file)
	if err != nil {
		fatal("get rel, err: %+v", err)
		return
	}

	ast, err := goast.ParseAst(*file)
	if err != nil {
		fatal("parse ast, err: %+v", err)
		return
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
		fatal("package name not found")
		return
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
		return
	}

	_ = os.RemoveAll(filepath.Join(wd, ".preview/main.go"))

	setWindows := ""
	if windowRect.Dx() > 0 && windowRect.Dy() > 0 {
		setWindows = fmt.Sprintf(`
ebiten.SetWindowPosition(%d, %d)
ebiten.SetWindowSize(%d, %d)
`, windowRect.Min.X, windowRect.Min.Y, windowRect.Dx(), windowRect.Dy())
	}

	mainFn := fmt.Sprintf(`
package main

import (
	preview "%s"

	"github.com/yanun0323/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	%s
	ebiten.SetRunnableOnUnfocused(true)

	view := preview.%s()
	app := ebui.NewApplication(view)
	app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	app.Run("preview")
}
`, importPath, setWindows, fnName)

	mainScope, err := goast.ParseScope(0, []byte(mainFn))
	if err != nil {
		fatal("parse main fn, err: %+v", err)
		return
	}

	newAst, err := goast.NewAst(mainScope...)
	if err != nil {
		fatal("new ast, err: %+v", err)
		return
	}

	_ = os.MkdirAll(filepath.Join(wd, ".preview"), 0755)
	exportFile := filepath.Join(wd, ".preview", "main.go")

	if err := newAst.Save(exportFile, false); err != nil {
		fatal("save main fn, err: %+v", err)
		return
	}

	if err := tryRunPreview(wd, exportFile); err != nil {
		fatal("try run preview, err: %+v", err)
		return
	}
}

func tryGetWindowPosition() (rect image.Rectangle) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("get window position, err: %+v", err)
		}
	}()

	x, y := ebiten.WindowPosition()
	w, h := ebiten.WindowSize()

	return image.Rectangle{
		Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + w, Y: y + h},
	}
}

func findProjectRoot(file *string) (string, error) {
	dir := filepath.Dir(*file)
	for dir != "." && dir != "/" && dir != "" {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		dir = filepath.Join(dir, "..")
	}

	return "", fmt.Errorf("project root not found")
}

func findGoModuleName(wd string) (string, error) {
	goMod, err := os.ReadFile(filepath.Join(wd, "go.mod"))
	if err != nil {
		return "", err
	}

	lines := strings.SplitSeq(string(goMod), "\n")
	for line := range lines {
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

func tryRunPreview(wd string, exportFile string) error {
	log.Printf("wd: %s", wd)
	{
		cmd := exec.Command("pwd")
		cmd.Dir = wd
		output, err := cmd.Output()
		if err != nil {
			return errors.Errorf("pwd, err: %+v", err)
		}
		log.Printf("pwd: %s", string(output))
	}

	exportRelativeFile, err := filepath.Rel(wd, exportFile)
	if err != nil {
		return errors.Errorf("get rel, err: %+v", err)
	}

	cmd := exec.Command("go", "run", "./"+exportRelativeFile)
	cmd.Dir = wd

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			if execErr.ProcessState.ExitCode() == 1 {
				return nil
			}
		}

		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("require go and github.com/yanun0323/ebui/tool/ebui installed")
		}

		return errors.Errorf("run main fn, err: %+v", err)
	}

	return nil
}

func fatal(format string, v ...any) {
	if *debug {
		log.Fatalf(format, v...)
	}
}
