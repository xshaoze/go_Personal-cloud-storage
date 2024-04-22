package wp_api

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

var port = 8080

func StartServer() {
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(
			"static",
			r.URL.Path[len("/static/"):],
		)
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/", wpMainHeader)
	http.HandleFunc("/login", wpLoginHeader)
	http.HandleFunc("/uploadFile", wpUploadFile)
	http.HandleFunc("/queryfilelist", wpQueryFileListHeader)
	http.HandleFunc("/downloadfile", wpDownloadHandler)
	http.HandleFunc("/setFileSharing", wpSetFileSharingHandler)
	http.HandleFunc("/removeFileSharing", wpRemoveFileSharingHandler)
	http.HandleFunc("/removeFile", wpRemoveFileHandler)
	http.HandleFunc("/zd_removeFile", wp_zd_RemoveFileJHeader)
	http.HandleFunc("/huifu", wpHuiFUFileHeader)
	http.HandleFunc("/logout", deleteCookieHandler)
	http.HandleFunc("/setOwnMsg", setOwnMsgHeader)
	http.HandleFunc("/changeOwnMsg", changeOwnMsgHeader)
	http.HandleFunc("/changeOwnPwd", changeOwnPwdHeader)

	err := http.ListenAndServe(":"+fmt.Sprint(port), nil)
	if err != nil {
		log.Panicln(err)
		return
	}
}
