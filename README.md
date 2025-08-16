# process-chronicle  

a process timer record for windows

## build  
1. install golang compiler
2. clone this project and into project root directory
3. use `go mod tidy` to add required package
4. follow [this](https://docs.fyne.io/) to build the environment
5. use `go run .\main.go` to run project  

資料儲存
* filter:
    * prefix(route): 要過濾的程式(前綴)
* register:
    * alias: 程式名稱
    * path: 程式路徑
    * totalTime: 總開啟時間
    * lastOpened: 最後開啟日期