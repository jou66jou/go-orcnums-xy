# 目前仍在測試階段
# 參考
fork origin: https://github.com/otiai10/ocrserver

# 與原版差異
辨識圖片中井字型0~9的x,y位置。

# 條件限制
test0.0.1版本僅辨識圖片中井字型1到9，且0於7到9的下方  
同水平中不可有重複數字。

# 目前可辨別錯誤
1. `find multiple same or shortage of numbers by y point line`
圖片中鍵盤缺少或其他位置出現，123、456、789的三種數字組合。
2. `not found number : 0`
0未被在數字789的底下找到。  

## TODO
1. 789底下有多個0
2. 確認井字X軸
3. ...

# Quick Start
沒使用fork，因此未修改import路徑，只能先用docker起環境。

## Clone to anywhere

```sh
% git clone github.com/jou66jou/ocrserver/
```

## Ready-Made Docker Image

```sh
% docker run -p 8080:8080 otiai10/ocrserver
# open http://localhost:8080
```

cf. [docker](https://www.docker.com/products/docker-toolbox)
