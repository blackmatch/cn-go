package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// PublishPageData 发布页面数据
type PublishPageData struct {
	PageURL     string `json:"pageURL"`
	Title       string `json:"title"`
	PublishDate string `json:"publishDate"`
}

// MetaData 源数据结构
type MetaData struct {
	PageURL string `json:"pageURL"`
	Code    string `json:"code"`
	Name    string `json:"name"`
}

// GetLatestPageData 获取最新发布的统计页面数据
func GetLatestPageData() PublishPageData {
	// 请求统计局页面数据
	url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 使用 goquery 解析数据
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 取第一条数据（最新发布的）
	firstElement := doc.Find(".center_list .center_list_contlist li a").First()
	pageURL := firstElement.AttrOr("href", "")
	title := firstElement.Find(".cont_tit03").Text()
	updateDate := firstElement.Find(".cont_tit02").Text()
	if !strings.HasPrefix(pageURL, "http") {
		pageURL = "http://www.stats.gov.cn" + pageURL
	}

	return PublishPageData{PageURL: pageURL, Title: title, PublishDate: updateDate}
}

// GetMetaData 获取源数据
func GetMetaData(mData MetaData, s1 string, s2 string) []MetaData {
	time.Sleep(200 * time.Millisecond)
	url := mData.PageURL
	if len(url) == 0 {
		return make([]MetaData, 0)
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 使用 goquery 解析数据
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	mds := make([]MetaData, 0)
	doc.Find(s1).Each(func(i int, s *goquery.Selection) {
		PageURL := ""
		Code := ""
		Name := ""
		s.Find(s2).Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				PageURL = s.Find("a").AttrOr("href", "")
				if len(PageURL) == 0 {
					Code = s.Text()
				} else {
					Code = s.Find("a").Text()
				}
				Code, _ = decodeToGBK(Code)
			} else {
				Name = s.Find("a").Text()
				if len(Name) == 0 {
					Name = s.Text()
				}
				Name, _ = decodeToGBK(Name)
			}
		})

		lastCharIdx := strings.LastIndex(mData.PageURL, "/")
		if lastCharIdx >= 0 && len(PageURL) > 0 {
			tmp := mData.PageURL[0 : lastCharIdx+1]
			PageURL = tmp + PageURL
		}

		mds = append(mds, MetaData{PageURL, Code, Name})
	})

	return mds
}

// GetProvinceData 获取省级数据
func GetProvinceData(pgData PublishPageData) []MetaData {
	// 请求统计局页面数据
	url := pgData.PageURL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 使用 goquery 解析数据
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	prvs := make([]MetaData, 0)
	doc.Find(".provincetable .provincetr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			PageURL := s.Find("a").AttrOr("href", "")
			province, _ := decodeToGBK(s.Find("a").Text())
			PageURL = strings.ReplaceAll(pgData.PageURL, "index.html", PageURL)
			prvs = append(prvs, MetaData{Name: province, PageURL: PageURL})
		})
	})

	return prvs
}

// 获取中文文本（对一些非GBK编码的网站）
func decodeToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}
