package translate_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"go2java/internal/syntaxtree"
	"go2java/internal/translate"
)

func TestExamplesHaveSameGoAndJavaOutput(t *testing.T) {
	examples, err := filepath.Glob(filepath.Join("..", "..", "examples", "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(examples) == 0 {
		t.Fatal("no example files")
	}

	for _, example := range examples {
		example := example
		t.Run(filepath.Base(example), func(t *testing.T) {
			goOut := runCmd(t, filepath.Dir(example), "go", "run", filepath.Base(example))

			tree, err := syntaxtree.LoadTypedFile(example)
			if err != nil {
				t.Fatal(err)
			}
			ir, err := translate.GoToIntermediate(tree)
			if err != nil {
				t.Fatal(err)
			}
			javaCode, err := translate.IntermediateToJava(ir)
			if err != nil {
				t.Fatal(err)
			}

			tmp := t.TempDir()
			mainJava := filepath.Join(tmp, "Main.java")
			if err := os.WriteFile(mainJava, []byte(javaCode), 0o600); err != nil {
				t.Fatal(err)
			}

			runCmd(t, tmp, "javac", "Main.java")
			javaOut := runCmd(t, tmp, "java", "Main")

			if normalize(goOut) != normalize(javaOut) {
				t.Fatalf("output mismatch\ngo: %q\njava: %q", goOut, javaOut)
			}
		})
	}
}

func runCmd(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
	return string(out)
}

func normalize(s string) string { return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n")) }
