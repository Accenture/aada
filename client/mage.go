//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"io/ioutil"
	"os"
	"strings"
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
	cvs, err := loadVersionInfo()
	if err != nil {
		return err
	}
	fmt.Println("building with version", cvs)
	for _, t := range allTargets {
		fmt.Print(t[0] + " ")
		ok, err := target.Glob(t[0], "*.go", "go.mod", "go.sum", "version.info")
		if err != nil {
			return err
		}
		if ok {
			err = buildPlatform(t[1], t[2], t[0])
			if err != nil {
				return err
			}
		} else {
			fmt.Println("is up to date")
		}
	}
	return nil
}

// Ensures the build is current and signs Mac binaries via Apple.
func Sign() error {
	err := Build()
	if err != nil {
		return err
	}
	err = appleSign("aada_mac_x64")
	if err != nil {
		return err
	}
	err = appleSign("aada_mac_m1")
	if err != nil {
		return err
	}
	return nil
}

func buildPlatform(os string, arch string, binary string) error {
	err := sh.RunWith(map[string]string{"GOOS": os, "GOARCH": arch},
		"go", "build", "-o", binary)
	fmt.Println("built")
	return err
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
