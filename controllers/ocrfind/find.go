package ocrfind

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type OCRs struct {
	KeyBoard `json:"keyboard"`
}

type KeyBoard struct {
	Nums   []Num    `json:"nums"`
	Engs   []Eng    `json:"engs"`
	Switch []Switch `json:"switch"`
}
type Switch struct {
	Text string `json:"text"`
	Points
}
type Eng struct {
	Char string `json:"char"`
	Points
}
type Num struct {
	N int `json:"num"`
	Points
}
type Points struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}

var (
	RangeY = 10 // 寬容水平座標的差值
	numInt = map[string]int{
		"o": 0, // 英文字母小寫o
		"O": 0, // 英文字母大寫O
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
	}
	charInt = map[rune]int{
		'a': 0,
		'b': 1,
		'c': 2,
		'd': 3,
		'e': 4,
		'f': 5,
		'g': 6,
		'h': 7,
		'i': 8,
		'j': 9,
		'k': 10,
		'l': 11,
		'm': 12,
		'n': 13,
		'o': 14,
		'p': 15,
		'q': 16,
		'r': 17,
		's': 18,
		't': 19,
		'u': 20,
		'v': 21,
		'w': 22,
		'x': 23,
		'y': 24,
		'z': 25,
	}
)

// query與解析hORC的標準內容
func NewOCRs(hORC string) (OCRs, error) {
	var ocrs = OCRs{}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(hORC))
	if err != nil {
		return ocrs, err
	}
	dom.Find("span").Each(func(i int, selection *goquery.Selection) {
		// 0~9數字，涵蓋英文字母o或O
		t := selection.Text()
		num, ok := numInt[t]
		if ok {
			title, exis := selection.Attr("title")
			if exis {
				n := getNum(title, num)
				ocrs.Nums = append(ocrs.Nums, n)
			}
		}

		// 單一字母
		_, ok = charInt[rune(t[0])]
		if ok {
			title, exis := selection.Attr("title")
			if exis {
				e := getEng(title, t)
				ocrs.Engs = append(ocrs.Engs, e)
			}
		}

		// 取得switch鍵
		if t == "123" {
			title, exis := selection.Attr("title")
			if exis {
				sw := getSwitch(title, t)
				ocrs.Switch = append(ocrs.Switch, sw)
			}
		}
	})
	return ocrs, nil
}
