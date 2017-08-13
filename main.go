package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"os"
	"strconv"
	"bufio"
)

type VALUER struct{
	url string
	name string
	va string
}

type Valuers []VALUER
var valuers Valuers

func writeCSV(valuers []VALUER, username string ) {
	file, err := os.Open(username + "_valuer_data.csv")
	if err != nil {
		file, _ = os.Create(username + "_valuer_data.csv")
		fmt.Println("create file")
	}
	
	writer := csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
	for i := range valuers {
		//fmt.Println(valuers[i])
		va := valuers[i]
		fmt.Println(va)
		err = writer.Write([]string{va.name, va.url, va.va})
		if err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
	file.Close()
}
func getValuer(url string) error {
	var user_url string
	var valuer_name string
	var valu_num string
	doc, err := goquery.NewDocument(url)
	if err != nil{
		return err
	}
	doc.Find(".valuer_main_box").Each(func(_ int, s *goquery.Selection) {
		var valuer VALUER
		s.Find(".valuer_upper_box").Each(func(_ int, t *goquery.Selection) {
			t.Find("a").Each(func(_ int, u *goquery.Selection) {
				user_url, _ = u.Attr("href")
				valuer.url = user_url
			})
			t.Find("b").Each(func(_ int, u *goquery.Selection) {
				valuer_name = u.Text()
				valuer.name = valuer_name
			})
			t.Find("strong").Each(func(_ int, u *goquery.Selection) {
				valu_num = u.Text()
				valuer.va = valu_num
			})
		})
		//fmt.Println([]string{user_url})
		valuers = append(valuers, valuer)
	})
	return nil
}
func check_num(url string) (int){
	doc, _ := goquery.NewDocument(url)
	//fmt.Println("check num")
	result := 0
	doc.Find(".nv_page_btn").Each(func(_ int, s *goquery.Selection) {
		//fmt.Println("--" + fmt.Sprint(s.Text) + "--")
		result = 1
	})
	//fmt.Println(result)
	return result
}
func main() {
	fmt.Println("input user name")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	username := stdin.Text()
	if check_num("https://valu.is/" + username + "/vaholder") == 0{
		url := "https://valu.is/" + username + "/vaholder"
		err := getValuer(url)
		if err != nil {
			return
		}
	}else{
		for i := 0 ; i < 30; i++ {
			url := "https://valu.is/" + username + "/vaholder?page=" + strconv.Itoa(i)
			//fmt.Println(url)
			err := getValuer(url)
			if err != nil {
				break
			}
		}
	}
	writeCSV(valuers, username)
	fmt.Println("created valuers list")
}