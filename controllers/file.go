package controllers

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/otiai10/gosseract"
	"github.com/otiai10/marmoset"
)

var (
	imgexp = regexp.MustCompile("^image")
	numInt = map[string]int{
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
)

type OCRs struct {
	Nums []Num `json="nums"`
}

type Num struct {
	N  int `json="num"`
	X1 int `json="x1"`
	Y1 int `json="y1"`
	X2 int `json="x2"`
	Y2 int `json="y2"`
}

// FileUpload ...
func FileUpload(w http.ResponseWriter, r *http.Request) {
	var errResp = map[string]interface{}{
		"error":   "",
		"version": version,
	}
	render := marmoset.Render(w, true)

	// Get uploaded file
	r.ParseMultipartForm(32 << 20)
	// upload, h, err := r.FormFile("file")
	upload, _, err := r.FormFile("file")
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}
	defer upload.Close()

	// Create physical file
	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	// Make uploaded physical
	if _, err = io.Copy(tempfile, upload); err != nil {
		render.JSON(http.StatusInternalServerError, err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempfile.Name())
	client.Languages = []string{"eng"}
	if langs := r.FormValue("languages"); langs != "" {
		client.Languages = strings.Split(langs, ",")
	}
	if whitelist := r.FormValue("whitelist"); whitelist != "" {
		client.SetWhitelist(whitelist)
	}

	var out string
	out, err = client.HOCRText()
	render.EscapeHTML = false
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}

	// query與解析hORC的標準內容
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(out))

	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}
	var ocrs = OCRs{}
	dom.Find("span").Each(func(i int, selection *goquery.Selection) {
		num, ok := numInt[selection.Text()]
		if ok {
			title, exis := selection.Attr("title")
			if exis {
				n := getN(title, num)
				ocrs.Nums = append(ocrs.Nums, n)
			}
		}
	})

	// 抓出圖片井字型1~9與0的位置
	ocrsKeyboard, err := keepKeyNum(ocrs)
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}

	render.JSON(http.StatusOK, map[string]interface{}{
		"result":  ocrsKeyboard,
		"version": version,
	})
}

// 解析string title (ex: "bbox 102 1190 113 1221; x_wconf 95 1") to struct Num
func getN(title string, num int) Num {
	tmp := strings.Split(title, " ")
	x1, _ := strconv.Atoi(tmp[1])
	y1, _ := strconv.Atoi(tmp[2])
	x2, _ := strconv.Atoi(tmp[3])
	y2, _ := strconv.Atoi(tmp[4][:len(tmp[4])-2]) // 去除最後';'符號
	return Num{num, x1, y1, x2, y2}
}

// 排除非鍵盤上數字
func keepKeyNum(ocrsAll OCRs) (OCRs, error) {
	var ocrsKeyboard = OCRs{}
	var nbox [10][][]int
	var rangeY = 10 // 寬容水平座標的差值

	// 箱子儲存法
	for _, v := range ocrsAll.Nums {
		points := []int{v.X1, v.Y1, v.X2, v.Y2}
		nbox[v.N] = append(nbox[v.N], points)
	}

	// debug: 顯示圖片上所有數字的位置
	// for k, v := range nbox {
	// 	fmt.Println(k, "+", v)
	// }

	// 將水平(Y1 + Y2) +- rangeY 的座標存到pointsY map, 以下方法建立在鍵盤同同一水平上不會重複數字的情況
	var pointsY = map[int][10][]int{} // { "pointsY" : [ yyy1:[x1,y1,x2,y2], ..., yyyN:[x1,y1,x2,y2] ]}

	// 遍歷 1~10 box，數字0最後檢查
	for i := 1; i < 10; i++ {
		// 檢查畫面上數字每個出現位置的水平座標
		for _, nP := range nbox[i] {
			Y := nP[1] + nP[3] // 水平座標(Y1 + Y2)
			noKey := true
			for k, v := range pointsY {
				// 檢查+-rangeY內的水平座標，若有則加入同水平座標的群組
				if k+rangeY > Y && k-rangeY < Y {
					v[i] = nP
					pointsY[k] = v
					noKey = false
					break
				}
			}
			// 若沒發現已有的水平座標則新增 pointsY Map key by X
			if noKey {
				ypV := [10][]int{}
				ypV[i] = nP
				pointsY[Y] = ypV
			}
		}
	}
	combNum := map[string][][]int{
		"123": [][]int{},
		"456": [][]int{},
		"789": [][]int{},
	}
	// 處理每個水平相近的Num, 每當1,4,7為起始, 連續三個數字有紀錄則新增至combNum
	for _, yNum := range pointsY {
		combStr, startIndexs := check3Combo(yNum, 1, 4, 7)
		for i, s := range combStr {
			// 新增連續三個數字為一組進combNum
			combNum[s] = append(combNum[s], yNum[startIndexs[i]], yNum[startIndexs[i]+1], yNum[startIndexs[i]+2])
		}
	}
	// 解析combNum
	isOneCase := true

	for _, v := range combNum {
		if len(v) != 3 { // 若發現任何組合不剛好3個
			isOneCase = false
			break
		}
	}
	if isOneCase { // isOneCase = ture代表各數組在畫面上只有各出現一次（認定畫面上只有鍵盤1~9井字三個數組）
		ocrsKeyboard = CombToOCRs(combNum)
		// 計算4跟1、7跟4的Y1差
		h41 := combNum["456"][0][1] - combNum["123"][0][1]
		h74 := combNum["789"][0][1] - combNum["456"][0][1]
		// 尋找每個0與7的Y1差, 落在h41 or h74 +- rangeY內
		has0 := false
		for _, v := range nbox[0] {
			h07 := v[1] - combNum["789"][0][1]
			if (h07 < h41+rangeY && h07 > h41-rangeY) || (h07 < h74+rangeY && h07 > h74-rangeY) {
				ocrsKeyboard.Nums = append(ocrsKeyboard.Nums, Num{0, v[0], v[1], v[2], v[3]})
				has0 = true
				break
			}
		}
		if !has0 {
			return ocrsKeyboard, errors.New("not found number : 0")
		}
	} else {
		// TODO: 發現水平座標相同的數組有多個或缺少，需再確認數組之間垂直座標
		return ocrsKeyboard, errors.New("find multiple same or shortage of numbers by y point line")
	}
	return ocrsKeyboard, nil
}

// 確認該水平有連續三個數字，回傳數組字串及起始位置slice
func check3Combo(yNum [10][]int, indexs ...int) ([]string, []int) {
	var comboStr = map[int]string{
		1: "123",
		4: "456",
		7: "789",
	}
	var keys []string
	var startIndexs []int
	for _, i := range indexs {
		if len(yNum[i]) > 0 && len(yNum[i+1]) > 0 && len(yNum[i+2]) > 0 {
			keys = append(keys, comboStr[i])
			startIndexs = append(startIndexs, i)
		}
	}
	return keys, startIndexs
}

// 數組Map轉換成OCRs
func CombToOCRs(combNum map[string][][]int) OCRs {
	ocrs := OCRs{}
	apenNums := func(n int, yNum [][]int, nums *[]Num) {
		for i := 0; i < 3; i++ {
			*nums = append(*nums, Num{n + i, yNum[i][0], yNum[i][1], yNum[i][2], yNum[i][3]})
		}
	}
	for k, yNum := range combNum {
		switch k {
		case "123":
			apenNums(1, yNum, &ocrs.Nums)
		case "456":
			apenNums(4, yNum, &ocrs.Nums)
		case "789":
			apenNums(7, yNum, &ocrs.Nums)
		}
	}
	return ocrs
}
