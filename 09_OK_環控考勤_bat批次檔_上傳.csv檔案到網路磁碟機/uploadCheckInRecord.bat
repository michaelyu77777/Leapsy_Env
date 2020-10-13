:: 刪除上次連線
net use Y: /delete

:: 連上網路磁碟
:: root0217:密碼  Michael:帳號
net use Y: \\leapsy-nas3\APP\member\MichaelYu\LeapsyCheckInRecordBackup root0217 /user:Michael /PERSISTENT:NO

:: 僅上傳今日的檔案到目的地
copy D:\Users\fish0\Documents\考勤系統\Amber-匯出設定與匯出檔案\Rec%DATE:~0,4%%DATE:~5,2%%DATE:~8,2%.csv Y: