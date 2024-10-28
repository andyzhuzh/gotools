package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// import "archive/zip"  "IO"
func ZipFile(zipFileFullName string, sourceFileFullName string) error {
	// zipFile(`D:\data\test.zip`, `D:\data\PivotStyleLight16.xlsx`)
	//"archive/zip" "io" "os" "path/filepath"
	// 1.创建ZIP文件
	zipFileObj, err := os.Create(zipFileFullName)
	if err != nil {
		return err
		// panic(err)
	}
	// 2.创建一个Writer对象
	zipWriter := zip.NewWriter(zipFileObj)

	// 3.向ZIP对象添加文件
	filename := filepath.Base(sourceFileFullName) //文件名称加扩展名
	w, err := zipWriter.Create(filename)
	if err != nil {
		return err
		// panic(err)
	}
	// 4.打开待压缩文件
	f, err := os.Open(sourceFileFullName)
	if err != nil {
		return err
		// panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return err
		// panic(err)
	}
	zipWriter.Close()
	return nil
}

// import "archive/zip"  "IO" "path/filepath"
func UnzipFile(zipFileName string, destFilePath string) error {
	//"archive/zip" "io" "os" "path/filepath"
	// unZipFile(`D:\data\test.zip`, `D:\data\test`)
	// 1.打开ZIP文件
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
		// panic(err)
	}
	defer zipFile.Close()

	// 2.遍历ZIP文件
	for _, file := range zipFile.File {
		fileNameWithPath := file.Name
		filefullPath := filepath.Join(destFilePath, filepath.Dir(fileNameWithPath))
		filefullName := filepath.Join(destFilePath, fileNameWithPath)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(filefullPath, os.ModePerm)
			continue
		}
		// 3.创建文件夹
		if err := os.MkdirAll(filefullPath, os.ModePerm); err != nil {
			return err
		}
		// 4.解压到目标文件
		destFile, err := os.OpenFile(filefullName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		f, err := file.Open()
		if err != nil {
			return err
		}
		// 5.写入文件
		if _, err := io.Copy(destFile, f); err != nil {
			return err
		}
		destFile.Close()
		f.Close()

	}

	return nil
}
