package ocrfind

import (
	"errors"
	"strconv"
	"strings"
)

func (ocr *OCRs) GetKeyboardEng() ([]Eng, error) {
	var engs []Eng
	var cbox [26][][]int
	// 箱子儲存法
	for _, v := range ocr.Engs {
		points := []int{v.X1, v.Y1, v.X2, v.Y2}
		cbox[charInt[rune(v.Char[0])]] = append(cbox[charInt[rune(v.Char[0])]], points)
	}

	// debug: 顯示圖片上所有字母的位置
	// for k, v := range cbox {
	// 	fmt.Println(string(rune(k+97)), "+", v)
	// }

	var pointsY = map[int][26][]int{} // { "pointsY" : [ a:[x1,y1,x2,y2], ..., z:[x1,y1,x2,y2] ]}

	// 歸類同水平的英文字母，假設每個水平英文字母未重複
	for i := 0; i < 26; i++ {
		// 檢查畫面上數字每個出現位置的水平座標
		for _, cP := range cbox[i] {
			Y := cP[1] + cP[3] // 水平座標(Y1 + Y2)
			noKey := true
			for k, v := range pointsY {
				// 檢查+-RangeY內的水平座標，若有則加入同水平座標的群組
				if k+RangeY > Y && k-RangeY < Y {
					v[i] = cP
					pointsY[k] = v
					noKey = false
					break
				}
			}
			// 若沒發現已有的水平座標則新增 pointsY Map key by X
			if noKey {
				ypV := [26][]int{}
				ypV[i] = cP
				pointsY[Y] = ypV
			}
		}
	}

	// 檢查手機鍵盤關鍵三組同水平英文字母組合
	var comboKeys = []string{
		"qwertyuiop", "asdfghjkl", "zxcvbnm",
	}

	for i := 0; i < 26; i++ {
		if len(cbox[i]) == 0 { // 發生未檢測到某字母
			return nil, errors.New("not found char : " + string(rune(i+97)))
		}
	}

	// 紀錄符合的組合數
	count := 0
	// 檢查同一水平的字母集群
	for _, clus := range pointsY {
		// 檢查是否為連續組合
		for _, ks := range comboKeys {
			hasCombo := true
			for _, k := range ks {
				if len(clus[charInt[k]]) == 0 { // 連續中斷
					hasCombo = false
					break
				}
			}
			if hasCombo { // 該集群為連續組合
				count++
				for _, k := range ks { // 不在前面for做是考慮減少append
					ps := clus[charInt[k]]
					engs = append(engs, Eng{string(k), Points{ps[0], ps[1], ps[2], ps[3]}})
				}
			}
		}
	}
	// 若不為三組則判斷錯誤
	if count != 3 {
		return nil, errors.New("find multiple same or shorage of combo chars by y point line")
	}

	return engs, nil
}

func getEng(title string, c string) Eng {
	tmp := strings.Split(title, " ")
	x1, _ := strconv.Atoi(tmp[1])
	y1, _ := strconv.Atoi(tmp[2])
	x2, _ := strconv.Atoi(tmp[3])
	y2, _ := strconv.Atoi(tmp[4][:len(tmp[4])-1]) // 去除最後';'符號
	return Eng{c, Points{x1, y1, x2, y2}}
}
