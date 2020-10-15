:: 刪已存在的網路磁碟機 (強制執行加上 /y)
net use Z: /delete /y

:: 建立網路磁碟 (root0217:密碼  Michael:帳號)
net use Z: \\leapsy-nas3\CheckInRecord root0217 /user:Michael /PERSISTENT:NO

:: 在轉換支援中文語系之前 先將英文的年月日儲存到變數中 否則無法提取英文日期
set dateTime=%DATE:~0,4%%DATE:~5,2%%DATE:~8,2%
echo %dateTime%

:: 改變 Windows Command Code Page 成UTF-8 支援中文(支援中文路徑)
chcp 65001

:: 僅上傳今日的檔案到目的地
:: 會自動覆蓋 但建立日期仍然會維持覆蓋之前的日期 不會變成覆蓋之後的日期
copy C:\Users\Fred\Desktop\出勤資料\Rec%dateTime%.csv Z: /y