//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/briandowns/spinner"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var allTargets = [][]string{
	{"aada_mac_x64", "darwin", "amd64"},
	{"aada_mac_arm64", "darwin", "arm64"},
	{"aada_win_x64.exe", "windows", "amd64"},
	{"aada_win_arm.exe", "windows", "arm64"},
	{"aada_linux_x64", "linux", "amd64"},
	{"aada_linux_arm64", "linux", "arm64"},
	{"aada_linux_arm32", "linux", "arm"},
}

var allReleases = [][]string{
	{"aada_mac_x64.zip", "aada_mac_x64", "aada"},
	{"aada_mac_arm64.zip", "aada_mac_arm64", "aada"},
	{"aada_win_x64.zip", "aada_win_x64.exe", "aada.exe"},
	{"aada_win_arm.zip", "aada_win_arm.exe", "aada.exe"},
	{"aada_linux_x64.zip", "aada_linux_x64", "aada"},
	{"aada_linux_arm64.zip", "aada_linux_arm64", "aada"},
	{"aada_linux_arm32.zip", "aada_linux_arm32", "aada"},
}

// Increment the current patch number.
func Patch() error {
	cvs, err := loadVersionInfo()
	if err != nil {
		return err
	}
	cv, err := semver.NewVersion(cvs)
	if err != nil {
		return err
	}
	ncv := cv.IncPatch()
	ioutil.WriteFile("version.info", []byte(ncv.String()), 0644)
	fmt.Println(cv.String(), "->", ncv.String())
	return nil
}

// Builds all supported platform binaries.
func Build() error {
	s := spinner.New(spinner.CharSets[21], 250*time.Millisecond)
	s.Suffix = "loading version info"
	s.Start()
	defer s.Stop()
	cvs, err := loadVersionInfo()
	if err != nil {
		return err
	}
	s.Prefix = cvs + "> "
	s.Stop()
	fmt.Println("building with version", cvs)
	s.Start()
	for _, t := range allTargets {
		fmt.Print(t[0] + " ")
		ok, err := target.Glob(t[0], "*.go", "go.mod", "go.sum", "version.info")
		if err != nil {
			return err
		}
		if ok {
			s.Suffix = "building " + t[1]
			err = buildPlatform(t[1], t[2], t[0])
			if err != nil {
				return err
			}
			s.Stop()
			fmt.Println("built", t[1], "into", t[0])
			s.Start()
		} else {
			fmt.Println("is up to date")
		}
	}
	return nil
}

func Release() error {
	err := Patch()
	if err != nil {
		return err
	}
	err = Package()
	if err != nil {
		return err
	}
	err = Sign()
	if err != nil {
		return err
	}
	return nil
}

// Ensures the build is current and signs Mac binaries via Apple.
func Sign() error {
	err := Build()
	if err != nil {
		return err
	}
	fmt.Println("signing mac binaries (use package for non-mac binaries)")
	err = appleSign("aada_mac_x64")
	if err != nil {
		return err
	}
	err = appleSign("aada_mac_arm64")
	if err != nil {
		return err
	}
	return nil
}

func Package() error {
	err := Build()
	if err != nil {
		return err
	}
	fmt.Println("packaging non-mac binaries (use sign for mac binaries)")
	for _, t := range allReleases {
		// For each release, build a zip file if it doesn't already exist
		ok, err := target.Glob(t[0], t[1])
		if err != nil {
			return err
		}
		if ok {
			zipFile(t[1], t[2], t[0])
		}
	}
	return nil
}

func zipFile(source string, name string, dest string) error {
	fmt.Printf("compressing %s into %s", source, dest)

	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	fw, err := zw.Create(name)
	if err != nil {
		return err
	}
	io.Copy(fw, in)
	zw.Flush()
	zw.Close()

	fmt.Println(" done")
	return nil
}

func buildPlatform(os string, arch string, binary string) error {
	return sh.RunWith(map[string]string{"GOOS": os, "GOARCH": arch},
		"go", "build", "-o", binary)
}

func appleSign(binary string) error {
	fmt.Println("signing", binary)
	err := sh.Copy("aada", binary)
	if err != nil {
		return err
	}
	err = sh.Run("gon", "apple.hcl")
	if err != nil {
		return err
	}
	return os.Rename("aada.zip", binary+".zip")
}

func loadVersionInfo() (string, error) {
	raw, err := ioutil.ReadFile("version.info")
	if err != nil {
		return "", err
	}
	cvs := strings.Trim(string(raw), " \t\r\n")
	return cvs, nil
}
