package ocrfind

import (
	"strconv"
	"strings"
)

func (ocr *OCRs) GetKeyboardSwitch() ([]Switch, error) {
	var ocrsKeyboard = OCRs{}

	return ocrsKeyboard.Switch, nil
}

func getSwitch(title string, text string) Switch {
	tmp := strings.Split(title, " ")
	x1, _ := strconv.Atoi(tmp[1])
	y1, _ := strconv.Atoi(tmp[2])
	x2, _ := strconv.Atoi(tmp[3])
	y2, _ := strconv.Atoi(tmp[4][:len(tmp[4])-1]) // 去除最後';'符號
	return Switch{text, Points{x1, y1, x2, y2}}
}
