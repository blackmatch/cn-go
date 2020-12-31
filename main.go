package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// NodeData 节点数据
type NodeData struct {
	Code     string     `json:"code"`
	Name     string     `json:"name"`
	Children []NodeData `json:"children"`
}

func main() {
	// fmt.Println("Hello world")
	// pgData := getPageData()
	// // fmt.Println(pgData)
	// prvs := getProvinceData(pgData)
	// // fmt.Println(prvs)
	// // for _, v := range prvs {
	// // 	cts := getCityData(v)
	// // 	fmt.Println(cts)
	// // }
	// fmt.Println(prvs[5])
	// ctData := getCityData(prvs[5])
	// fmt.Println(ctData[5])
	// cntData := getCountyData(ctData[5])
	// fmt.Println(cntData[2])
	// tnData := getTownData(cntData[2])
	// fmt.Println(tnData[2])
	// // vlgData := getVillageData(tnData[2])
	// vlgData := getMetaData(MetaData{PageURL: tnData[2].pageURL, Code: tnData[2].code, Name: tnData[2].name}, ".villagetable .villagetr", "td")
	// fmt.Println(vlgData)
	// bts, _ := json.MarshalIndent(vlgData, "", "  ")
	// fmt.Println(string(bts))

	t1 := time.Now()
	// 获取发布页数据
	latestPageData := GetLatestPageData()
	fmt.Println(latestPageData)
	provinces := GetProvinceData(latestPageData)
	fmt.Println(provinces)
	// 遍历每个省
	nodes := make([]NodeData, 0)
	for _, provData := range provinces {
		fmt.Printf("正在获取 %s 的数据...\n", provData.Name)
		nProv := NodeData{Code: provData.Code, Name: provData.Name}
		// 获取该省所有市级数据
		cities := GetMetaData(provData, ".citytable .citytr", "td")
		for _, ctData := range cities {
			fmt.Printf("正在获取 %s 的数据...\n", ctData.Name)
			nCt := NodeData{Code: ctData.Code, Name: ctData.Name}
			nProv.Children = append(nProv.Children, nCt)
			// 获取该市所有县级数据
			counties := GetMetaData(ctData, ".countytable .countytr", "td")
			for _, cntData := range counties {
				fmt.Printf("正在获取 %s 的数据...\n", cntData.Name)
				nCnt := NodeData{Code: cntData.Code, Name: cntData.Name}
				nCt.Children = append(nCt.Children, nCnt)
				// 获取该县所有镇级数据
				towns := GetMetaData(cntData, ".towntable .towntr", "td")
				for _, townData := range towns {
					fmt.Printf("正在获取 %s 的数据...\n", townData.Name)
					nTown := NodeData{Code: townData.Code, Name: townData.Name}
					nCnt.Children = append(nCnt.Children, nTown)
					// 获取该镇所有乡级数据
					villages := GetMetaData(townData, ".villagetable .villagetr", "td")
					for _, vlgData := range villages {
						nVlg := NodeData{Code: vlgData.Code, Name: vlgData.Name}
						nTown.Children = append(nTown.Children, nVlg)
					}
				}
			}
		}

		nodes = append(nodes, nProv)
		t := time.Now()
		fmt.Printf("\n已耗时：%d\n\n", t.Sub(t1))
		time.Sleep(200 * time.Millisecond)
	}
	bts, _ := json.MarshalIndent(nodes, "", "  ")
	fmt.Println(string(bts))

	t2 := time.Now()
	fmt.Printf("\n共耗时：%d\n", t2.Sub(t1))
}
