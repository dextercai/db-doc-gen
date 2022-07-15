package doc

import (
	"context"
	"fmt"
	"github.com/dextercai/db-doc-gen/model"
	"github.com/dextercai/db-doc-gen/util"
	"io/ioutil"
	"path"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/russross/blackfriday"
)

// createOfflineDoc create offline html、md、pdf、word
func CreateOfflineDoc(docPath string, dbInfo model.DbInfo, tables []model.Table) {
	var (
		docMdArr []string
		docMdStr string
	)
	// 标题
	docMdArr = append(docMdArr, fmt.Sprintf("# %s 数据库文档", dbInfo.DbName))
	docMdArr = append(docMdArr, "*created by github.com/dextercai/db-doc-gen*")
	// 生成基础信息
	docMdArr = append(docMdArr, "### 基础信息")
	docMdArr = append(docMdArr, "| 数据库名称 | 版本 | 字符集 | 排序规则 |")
	docMdArr = append(docMdArr, "| ---- | ---- | ---- | ---- |")
	docMdArr = append(docMdArr, fmt.Sprintf("| %s | %s | %s | %s |", dbInfo.DbName, dbInfo.Version, dbInfo.Charset, dbInfo.Collation))
	docMdArr = append(docMdArr, "")

	// 生成目录
	docMdArr = append(docMdArr, "### 数据库表目录")
	docMdArr = append(docMdArr, "| 序号 | 表名 | 描述 |")
	docMdArr = append(docMdArr, "| ---- | ---- | ---- |")
	for i := range tables {
		docMdArr = append(docMdArr, fmt.Sprintf("| %d | %s | %s |", i+1, tables[i].TableName, tables[i].TableComment))
	}
	// 生成表
	docMdArr = append(docMdArr, "### 数据库表信息")
	for i := range tables {
		docMdArr = append(docMdArr, fmt.Sprintf("#### %s(%s)", tables[i].TableName, tables[i].TableComment))

		docMdArr = append(docMdArr, "##### 列定义")
		docMdArr = append(docMdArr, genTableCol(tables[i].ColList)...)
		docMdArr = append(docMdArr, "")

		docMdArr = append(docMdArr, "##### 索引信息")
		docMdArr = append(docMdArr, genTableIdx(tables[i].IdxList)...)
		docMdArr = append(docMdArr, "")

		docMdArr = append(docMdArr, "####### 数据库定义SQL")
		docMdArr = append(docMdArr, genTableSqlArea()...)
		docMdArr = append(docMdArr, "")

	}
	docMdStr = strings.Join(docMdArr, "\r\n")
	util.WriteToFile(path.Join(docPath, dbInfo.DbName+".md"), docMdStr)
	fmt.Println("md格式输出完成")
	// html
	docMdArr = append([]string{mdCss}, docMdArr...)
	docMdStr = strings.Join(docMdArr, "\r\n")
	htmlPath := path.Join(docPath, dbInfo.DbName+".html")
	convert2Html(docMdStr, htmlPath)
	// pdf
	pdfPath := path.Join(docPath, dbInfo.DbName+".pdf")
	convert2Pdf(htmlPath, pdfPath)
}

// convert2Html md convert to html
func convert2Html(docMdStr, htmlPath string) {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_USE_XHTML
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	extensions := 0
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK
	extensions |= blackfriday.EXTENSION_HEADER_IDS
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_TABLES

	output := blackfriday.Markdown([]byte(docMdStr), renderer, extensions)
	util.WriteToFile(htmlPath, string(output))
	fmt.Println("html generate successfully!")
}

// convert2Pdf md convert to pdf
func convert2Pdf(htmlPath, pdfPath string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var (
		buf []byte
		err error
	)
	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("file:///" + htmlPath),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				Do(ctx)
			return err
		}),
	})
	if err != nil {
		fmt.Println("pdf generate failed! " + err.Error())
	} else {
		err = ioutil.WriteFile(pdfPath, buf, 0644)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
		fmt.Println("pdf generate successfully!")
	}
}
