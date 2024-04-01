package main_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestSimple(t *testing.T) {
	doTest(t, "testdata/simple")
}

func TestContext(t *testing.T) {
	doTest(t, "testdata/withcontext", "-c")
}

func doTest(t *testing.T, path string, ff ...string) {
	var snapsConfigured = snaps.WithConfig(snaps.Dir(path))
	args := append(append([]string{"run", "."}, ff...), path)
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("error executing cmd: %s", err.Error())
	}

	t.Log("running test")
	b, err := os.ReadFile(path + "/bus/generated.go")
	if err != nil {
		t.Fatal(err)
	}
	content := string(b)
	t.Log("content is ok")
	snapsConfigured.MatchSnapshot(t, content)
}
