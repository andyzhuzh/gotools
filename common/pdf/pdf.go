package pdf

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)
// https://wkhtmltopdf.org/ 下载wkhtmltox
func MarkdownToHtml(mdFile, htmlFile string) error {
	// "github.com/russross/blackfriday/v2"
	mdFileBytes, err := os.ReadFile(mdFile)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	html := blackfriday.Run(mdFileBytes)
	htmlHeader := `<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
	</head>`
	htmlfooter := `</html>`
	htmlStr := htmlHeader + string(html) + htmlfooter
	//写文件,打开文件，如果不存在则创建
	f, err := os.OpenFile(htmlFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(htmlStr)
	if err != nil {
		return err
	}
	return nil
}

func HtmlToPDF(htmlFile, pdfFile string) error {
	//wkhtmltopdf 下载wkhtmltox
	cmd := exec.Command("wkhtmltopdf", htmlFile, pdfFile)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("HTML 转换为 PDF 失败：%v\n", err)
		return err
	}
	return nil
}

func MarkdownToPDF(mdFile, pdfFile string) {
	filenamewithpath := mdFile                             // /文件名称加扩展名 返回路径最后一个元素，最后目录
	filePath, fileName := filepath.Split(filenamewithpath) //分割路径中的目录与文件
	fileExt := filepath.Ext(filenamewithpath)              //返回路径中的扩展名 如果没有点，返回空
	filenameOnly := strings.TrimSuffix(fileName, fileExt)
	htmlfilename := filenameOnly + ".html"
	htmlfilefullName := filepath.Join(filePath, htmlfilename) //路径连接
	MarkdownToHtml(mdFile, htmlfilefullName)
	HtmlToPDF(htmlfilefullName, pdfFile)
}
