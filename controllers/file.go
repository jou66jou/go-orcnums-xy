package controllers

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jou66jou/go-orcnums-xy/controllers/ocrfind"
	"github.com/otiai10/gosseract"
	"github.com/otiai10/marmoset"
)

var (
	imgexp = regexp.MustCompile("^image")
)

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
	// if whitelist := r.FormValue("whitelist"); whitelist != "qwertyuiopasdfghjklzxcvbnm" {
	// 	client.SetWhitelist(whitelist)
	// }

	// 設定三種特定辨認名單並取得個別hOCR內容
	var outNum, outEng, outSwitch string
	client.SetWhitelist(ocrfind.GetNumWiteList())
	outNum, err = client.HOCRText()
	render.EscapeHTML = false
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}
	client.SetWhitelist(ocrfind.GetEngWiteList())
	outEng, err = client.HOCRText()
	render.EscapeHTML = false
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}
	client.SetWhitelist(ocrfind.GetSwitchWiteList())
	outSwitch, err = client.HOCRText()
	render.EscapeHTML = false
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}
	// 取得OCRs
	ocrs, err := ocrfind.NewOCRs(outNum, outEng, outSwitch)
	if err != nil {
		errResp["error"] = err.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}

	// 抓出圖片中井字型1~9與0的位置
	nums, nErr := ocrs.GetKeyboardNum()
	// 抓出圖片中標準鍵盤26個字母位置
	engs, eErr := ocrs.GetKeyboardEng()
	if nErr != nil && eErr != nil {
		errResp["error"] = nErr.Error() + "; " + eErr.Error()
		render.JSON(http.StatusBadRequest, errResp)
		return
	}

	var ocrsKeyboard ocrfind.OCRs
	ocrsKeyboard.Nums = nums
	ocrsKeyboard.Engs = engs
	// TODO:未確認左下角位置
	ocrsKeyboard.Switch = ocrs.Switch

	render.JSON(http.StatusOK, map[string]interface{}{
		"result":  ocrsKeyboard,
		"version": version,
	})
}
