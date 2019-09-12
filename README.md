# 目前仍在測試階段
Dockerfile的Tesseract是最新版，辨識結果會有錯誤  
test0.0.3尚在測試新版能否使用  

# 參考
fork origin: https://github.com/otiai10/ocrserver

# 與原版差異
辨識圖片中井字型或同一排數字0~9位置，標準英文鍵盤排列的26個英文字母位置，
`123`、`abc`的手機螢幕鍵盤切換語言位置。

# 條件限制
test0.0.3版本可辨識圖片中井字型1到9的XY位置，但0需要在7到9的下方，  
同水平中不可有重複數字，尚未加入符號的辨識。

# 目前可辨別錯誤
1. `find multiple same or shortage of numbers by y point line`  
圖片中鍵盤缺少或其他位置出現，123、456、789的三種數字組合。
2. `not found number : 0`  
0未被在數字789的底下找到。  

## TODO
1. 數字鍵盤789底下若有多個0
2. 數字鍵盤未確認井字X軸
3. 數字鍵盤未確認多個井字組合
4. 切換鍵未確認必定出現在左下角區域

# Quick Start

## Mac Install tesseract

```sh
brew install tesseract
```

**注意** : test0.0.3版本，tesseract需裝 3.05.01 version
```sh
    git clone https://github.com/tesseract-ocr/tesseract.git
    cd tesseract
    git checkout 2158661
    ./autogen.sh
    ./configure
    make        # take a coffee break
    sudo make install
```

## Go Get & Run

If you have tesseract-ocr and library files on your machine  

```sh
% go get github.com/otiai10/ocrserver/...
% go get github.com/jou66jou/go-orcnums-xy
% cd $GOPATH/github.com/jou66jou/go-orcnums-xy
% go run main.go -p 8080
```  

## or Ready-Made Docker Image

```sh
% docker-compose up -d
```

cf. [docker](https://www.docker.com/products/docker-toolbox)

## Upload Test Image

1. Open http://localhost:8080
2. Upload test image file from `./imagetest` folder 
3. Get json response!


## Result
### Success numbers example

<div align=center><img width="30%" height="30%" src="https://github.com/jou66jou/go-orcnums-xy/blob/master/imagetest/success.jpg" alt="success case"/></div>

### Response

```json
200 OK
-----
{
	"result": {
		"Nums": [
			{
                    "N": 1,	"X1": 102,"Y1": 1190,"X2": 113,"Y2": 122
			},
			{
                    "N": 2, "X1": 316,"Y1": 1190,"X2": 337,"Y2": 122
			},
			{
                    "N": 3, "X1": 533,"Y1": 1190,"X2": 554,"Y2": 122
			},
			{
                    "N": 4, "X1": 97,"Y1": 1362,"X2": 120,"Y2": 139
			},
			{
                    "N": 5, "X1": 316,"Y1": 1362,"X2": 337,"Y2": 139
			},
			{
                    "N": 6, "X1": 533,"Y1": 1362,"X2": 554,"Y2": 139
			},
			{
                    "N": 7, "X1": 99,"Y1": 1533,"X2": 119,"Y2": 156
			},
			{
                    "N": 8, "X1": 316,"Y1": 1532,"X2": 337,"Y2": 156
			},
			{
                    "N": 9, "X1": 533,"Y1": 1532,"X2": 554,"Y2": 156
			},
			{
                    "N": 0, "X1": 316,"Y1": 1704,"X2": 337,"Y2": 173
			}
		]
	},
	"version": "test0.0.1"
}
```