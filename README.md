##    test_line_dev LINE機器人測試
### Step1 git clone專案
```
git clone https://github.com/gordan0410/test_line_dev.git
```
### Step2 設定參數與環境
1. 前往您的[line developer](https://developers.line.biz/console/), 取得您Line Messaging API的`Channel secret` 及 `Messaging API`
2. 複製並新增至專案內的`config/config.json`
3. 建立mogoDB Docker環境, 請於專案最頂層下:
`docker-compose up -d`
4. 連線至`localhost:9010`確認是否有開啟DB成功

### Step3 運行專案
1. 運行專案, 請於專案最頂層下:
`go run ./test_line_dev/main.go server`
2. 在專案內找到`urls.txt`, 取得`linebot webhook url` 及 `linebot api` 網址
3. 前往您的[line developer](https://developers.line.biz/console/), 於`Messaging API`標籤頁面內, 設定您的Webhook URL為上述產生的`linebot webhook ur`網址, 填寫完後並驗證及開啟`Use webhook`

### Step4 開始測試
* **API 接收訊息並儲存進DB**
1. 開啟您的line, 前往您的[line developer](https://developers.line.biz/console/), 於Messaging API標籤頁面內, 使用QRcode加入line機器人
2. 輸入您的留言並送出
3. 前往`localhost:9010`看是否有該則訊息, ex.
```json
{
    "_id": "ObjectId('63e518efcc8550155d8c0030')",
    "userid": "your user's ID",
    "msg": {
        "id": "17616245141620",
        "text": "your test",
        "emojis": null,
        "mention": null
    }
}
```
* **API 發送訊息給玩家**
1. 使用Postman, 使用POST方法, 網址為上述取得的`linebot api`網址
2. Request Body 格式為`json`, 內容如下:
```json
{
    "user_id":"your user's ID",
    "msg": "msg you want to deliver"
}
```
(註: `user_id`可以在`Basic settings`標籤頁面下找到)

3. 確認您的line是否有收到該則訊息

* **API 取得使用者的歷史訊息紀錄**
1. 使用Postman, 使用GET方法, 網址為上述取得的`linebot api`網址
2. 加入您的query string, ex:
```
?user_id=your_users_ID&content_per_page=10&page=1
```
3. 確認回覆是否為曾經發送過的留言, ex.
```json
{
    "user_id": "your user's ID",
    "messages": [
        {
            "emojis": null,
            "id": "17610153769742",
            "mention": null,
            "text": "好"
        },
        {
            "emojis": null,
            "id": "17610154544859",
            "mention": null,
            "text": "good"
        },
        {
            "id": "17610155119677",
            "keywords": [
                "Depressed"
            ],
            "packageid": "1127953",
            "stickerid": "5223185",
            "stickerresourcetype": "STATIC"
        }
    ],
    "total_page": 5,
    "now_page": 1
}
```

以上, 感謝閱讀