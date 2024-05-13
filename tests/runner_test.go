package tests

import (
	"bufio"
	"computer-club-manager/internal/club"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const (
	testFilename     = "input.txt"
	expectedFileName = "expected"
	outFileName      = "out"
)

func findAndRunTests(parent *testing.T, root string) {
	// Every directory with input.txt is considered to be a test.
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		if d.IsDir() {

			parent.Run(d.Name(), func(t *testing.T) {
				findAndRunTests(t, path)
			})

			return filepath.SkipDir
		}

		if d.Name() == testFilename {
			runTest(parent, path)
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		parent.Fail()
	}
}

func readLines(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	return
}

func runTest(t *testing.T, path string) {
	var err error

	dir := filepath.Dir(path)
	os.Stdout, err = os.Create(filepath.Join(dir, outFileName))
	if err != nil {
		t.Fatalf("cannot create '%s' file in %q: %v", outFileName, path, err)
	}

	file, err := os.Open(filepath.Join(dir, testFilename))
	if err != nil {
		t.Fatalf("cannot open '%s' file in %q: %v", testFilename, path, err)
	}
	defer file.Close()

	cc, err := club.NewComputerClub(file)
	if err != nil {
		return
	}

	cc.ProcessEvents()

	expectedLines, err := readLines(filepath.Join(dir, expectedFileName))
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("cannot read expected compilation results file: %v", err)
		}
	}

	outLines, err := readLines(filepath.Join(dir, outFileName))
	if err != nil {
		t.Fatalf("cannot read '%s' file in %q: %v", outFileName, path, err)
	}

	for i := 0; i < len(expectedLines); i++ {
		if expectedLines[i] != outLines[i] {
			t.Errorf("expected: %q,\tgot: %q", expectedLines[i], outLines[i])
		}
	}

}

func Test_Runner(t *testing.T) {
	findAndRunTests(t, ".")
}
