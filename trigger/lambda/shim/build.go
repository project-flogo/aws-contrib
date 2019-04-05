//+build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	fmt.Println("Running build script for the Lambda trigger")

	//Get the dir where build.go is present
	appDir, err := os.Getwd()
	if err != nil {
		fmt.Printf(err.Error())
	}

	var cmd *exec.Cmd

	// Clean up
	fmt.Println("Cleaning up previous executables")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "del", "/q", "handler", "handler.zip")
	} else {
		cmd = exec.Command("rm", "-f", "handler", "handler.zip")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Build an executable for Linux
	fmt.Println("Building a new handler file")
	cmd = exec.Command("go", "build", "-o", "handler")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Dir = filepath.Join(appDir)
	cmd.Env = append(os.Environ(), "GOOS=linux")

	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Zip the executable using the same code as build-lambda-zip
	fmt.Println("Zipping the new handler file")
	inputExe := filepath.Join(cmd.Dir, "handler")
	outputZip := filepath.Join(cmd.Dir, "handler.zip")
	if err := compressExe(outputZip, inputExe); err != nil {
		fmt.Printf("Failed to compress file: %v", err)
	}
}

func writeExe(writer *zip.Writer, pathInZip string, data []byte) error {
	exe, err := writer.CreateHeader(&zip.FileHeader{
		CreatorVersion: 3 << 8,     // indicates Unix
		ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
		Name:           pathInZip,
		Method:         zip.Deflate,
	})
	if err != nil {
		return err
	}

	_, err = exe.Write(data)
	return err
}

func compressExe(outZipPath, exePath string) error {
	zipFile, err := os.Create(outZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	data, err := ioutil.ReadFile(exePath)
	if err != nil {
		return err
	}

	return writeExe(zipWriter, filepath.Base(exePath), data)
}
