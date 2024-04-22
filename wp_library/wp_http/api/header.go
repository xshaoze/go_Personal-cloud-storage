package wp_api

import (
	wp_sql "DreamVerseCloud/wp_library/wp_sqlite3"
	"DreamVerseCloud/wp_library/wp_token"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

var filePath = "data/file/"

var password = []byte("@?QBJBxq9IV]tj+K}|^EYGJ!n[9$$Yb?T?ND:N*?mIa%]YH.l<f;1dU9 ,Ab5FeNNMD%moiDs{s5e&J,T )|rsRBk65p5d4&i0H?")

func wpMainHeader(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {

		token, err := r.Cookie("token")
		if err != nil {
			log.Println(err)
			data, err := json.Marshal(map[string]interface{}{
				"code": 500,
				"msg":  "游客",
			})
			if err != nil {
				log.Println(err)
			}
			w.Write(data)
			return
		}
		tokenData := wp_token.JwtDecryption(token.Value, password)

		user := wp_sql.QueryUserInfo([]string{
			"userName",
		}, []string{
			tokenData["userName"].(string),
		})

		var reData = make(map[string]interface{})
		if user == nil {
			reData["msg"] = "请重新登录！"
			reData["code"] = 500
		} else {
			if user.UserPasswd == tokenData["userPasswd"].(string) {
				reData["code"] = 200
				reData["msg"] = "验证token成功"
			}

			user.UserPasswd = "凡人之心难以洞悉天机之玄妙。"
			reData["userInfo"] = user
		}
		data, _ := json.Marshal(reData)
		w.Write(data)

	} else {
		w.Write([]byte("Method not alown"))
	}
}

func wpLoginHeader(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var postData map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&postData)
		if err != nil {
			log.Println("wp_api header.go login() 出现错误:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userInfo := wp_sql.QueryUserInfo([]string{
			"userName",
		}, []string{
			fmt.Sprint(postData["userName"]),
		})

		var reData = make(map[string]interface{})
		var data []byte
		if userInfo != nil {
			if userInfo.UserPasswd == fmt.Sprint(postData["userPasswd"]) {
				userInfo.UserPasswd = ""
				reData["code"] = 200
				reData["msg"] = "ok"
				day := 1

				if postData["remberme"] != nil {
					if postData["remberme"].(bool) {
						day = 7
					}
				}

				now := time.Now()
				expiration_time := time.Date(
					now.Year(),
					now.Month(),
					now.Day()+day, 3, 0, 0, 0,
					now.Location(),
				)
				cookie := http.Cookie{
					Name: "token",
					Value: wp_token.JwtEncryption(map[string]interface{}{
						"userName":   postData["userName"],
						"userPasswd": postData["userPasswd"],
					}, password),
					HttpOnly: true,
					Expires:  expiration_time,
					Path:     "/",
				}
				http.SetCookie(w, &cookie)

			} else {
				reData["code"] = 400
				reData["msg"] = "账号或密码错误"
			}
			data, _ = json.Marshal(reData)
		} else {
			reData["code"] = 400
			reData["msg"] = "账户不存在"
			data, _ = json.Marshal(reData)
		}

		w.Write(data)

	} else if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)

	} else {
		w.Write([]byte("Method not alown"))
	}
}

func wpUploadFile(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("token")
	if err != nil {
		// 游客登录
		data, _ := json.Marshal(map[string]interface{}{
			"code": 500,
			"msg":  "好家伙你不许上传!",
		})
		w.Write(data)
		return
	}

	// 解析token数据
	tokenData := wp_token.JwtDecryption(token.Value, password)
	// 接取文件
	err = r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "获取文件失败", http.StatusBadRequest)
		return
	}

	name_md5 := md5.New()
	name_md5.Write([]byte(header.Filename))
	data1 := hex.EncodeToString(name_md5.Sum(nil))
	log.Println("data1", data1)
	tmpFilepath := "data/tmp/" + string(data1)
	tmpfile, err := os.OpenFile(tmpFilepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		// http.Error(w, "文件创建失败", http.StatusInternalServerError)
		log.Println("wpUploadFile():文件创建失败:", err)
		return
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096)
	writer := bufio.NewWriter(tmpfile)

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Println("wpUploadFile():读取文件时出错:", err)
			// http.Error(w, "文件读取失败", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}
		if _, err := writer.Write(buffer[:n]); err != nil {
			log.Println("wpUploadFile():写入文件时出错:", err)
			// http.Error(w, "文件写入失败", http.StatusInternalServerError)
			return
		}
	}

	if err := writer.Flush(); err != nil {
		log.Println("wpUploadFile():刷新缓冲区时出错：", err)
		// http.Error(w, "文件写入失败", http.StatusInternalServerError)
		return
	}

	tmpfile.Close()

	file, err = os.Open(tmpFilepath)
	if err != nil {
		log.Println("wpUploadFile():无法打开文件:", err)
		return
	}

	hash := md5.New()
	reader = bufio.NewReader(file)
	buffer = make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Println("wpUploadFile():读取文件时出错:", err)
			return
		}
		if n == 0 {
			break
		}
		hash.Write(buffer[:n])
	}
	sum := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", sum)

	fileMeta := wp_sql.QueryFileMeta([]string{
		"FileHash_md5",
	}, []string{
		hashString,
	})

	filemeta := wp_sql.FileMeta{
		FileAddr:     filePath + hashString,
		FileHash_md5: hashString,
		FileSize:     header.Size,
	}

	file.Close()
	if fileMeta == nil {
		wp_sql.CreateFileMeta(filemeta)

		err := os.Rename(tmpFilepath, filePath+hashString)
		if err != nil {
			log.Println("wpUploadFile():移动文件失败！:", err)
		}
	} else {
		os.Remove(tmpFilepath)
	}

	fileMeta = wp_sql.QueryFileMeta([]string{
		"FileHash_md5",
	}, []string{
		hashString,
	})
	fileindex := wp_sql.FileIndex{
		FileName:      header.Filename,
		FilePath:      "/",
		UploadData:    time.Now(),
		FileOwnerShip: tokenData["userName"].(string),
		IsPublic:      false,
		DeletedStatic: 0,
		FileMetaId:    fileMeta.FileMetaId,
		FileSize:      fileMeta.FileSize,
	}
	wp_sql.CreateFileIndex(fileindex)
	log.Println(hashString)
	w.Write([]byte(hashString))
}

func wpQueryFileListHeader(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("token")
	if err != nil {
		data1, err := json.Marshal(wp_sql.QueryFileIndexList([]string{"IsPublic"}, []string{"1"}))
		if err != nil {
			log.Panicln(err)
		}
		data, err := json.Marshal(map[string]interface{}{
			"code": 200,
			"msg":  "游客",
			"data": string(data1),
		})
		if err != nil {
			log.Println("wpQueryFileListHeader():转换失败：", err)
		}
		w.Write(data)
		return
	}

	tokenData := wp_token.JwtDecryption(token.Value, password)

	files := wp_sql.QueryFileIndexList([]string{
		"FileOwnerShip",
		// "FilePath",
	}, []string{
		tokenData["userName"].(string),
		// tokenData["FilePath"].(string),
	})

	filess, _ := json.Marshal(files)
	json, err := json.Marshal(map[string]interface{}{
		"code": 200,
		"msg":  "用户:" + tokenData["userName"].(string),
		"data": string(filess),
	})
	if err != nil {
		log.Println("wpQueryFileListHeader():转换失败：", err)
	}
	w.Write(json)

}

func wpDownloadHandler(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("fileId")
	fileIndex := wp_sql.QueryFileIndex([]string{"id"}, []string{queryParam})
	fileMeta := wp_sql.QueryFileMeta([]string{"FileMetaId"}, []string{fmt.Sprint(fileIndex.FileMetaId)})
	file, err := os.Open(fileMeta.FileAddr) // 替换为你要提供的文件的路径
	if err != nil {
		http.Error(w, "文件未找到", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", `attachment; filename=`+fileIndex.FileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(fileIndex.FileSize))

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "无法复制文件内容到响应体", http.StatusInternalServerError)
		return
	}
}

func wpSetFileSharingHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	data, _ := json.Marshal(map[string]interface{}{
		"code": 500,
		"msg":  "游客",
	})
	if err != nil {
		w.Write(data)
		return
	}
	tokenData := wp_token.JwtDecryption(token.Value, password)
	user := wp_sql.QueryUserInfo([]string{
		"userName",
	}, []string{
		tokenData["userName"].(string),
	})
	var reData = make(map[string]interface{})
	if user != nil {
		if user.UserPasswd == tokenData["userPasswd"].(string) {
			reData["code"] = 200
			reData["msg"] = "验证token成功"
		} else {
			w.Write(data)
			return
		}
	} else {
		w.Write(data)
		return
	}

	queryParam := r.URL.Query().Get("fileId")

	wp_sql.SetFileSharing(queryParam)

	log.Println(queryParam)
	w.Write([]byte(queryParam))
}

func wpRemoveFileSharingHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	data, _ := json.Marshal(map[string]interface{}{
		"code": 500,
		"msg":  "游客",
	})
	if err != nil {
		w.Write(data)
		return
	}
	tokenData := wp_token.JwtDecryption(token.Value, password)
	user := wp_sql.QueryUserInfo([]string{
		"userName",
	}, []string{
		tokenData["userName"].(string),
	})
	var reData = make(map[string]interface{})
	if user != nil {
		if user.UserPasswd == tokenData["userPasswd"].(string) {
			reData["code"] = 200
			reData["msg"] = "验证token成功"
		} else {
			w.Write(data)
			return
		}
	} else {
		w.Write(data)
		return
	}

	queryParam := r.URL.Query().Get("fileId")

	a := wp_sql.RemoveFileSharing(queryParam)

	log.Println(queryParam)

	mapp := make(map[string]interface{})

	if a {
		mapp["code"] = 200
		mapp["msg"] = "删除成功"
	} else {
		mapp["code"] = 500
		mapp["msg"] = "删除失败"
	}

	redata, _ := json.Marshal(mapp)
	w.Write([]byte(redata))
}
func wpRemoveFileHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	data, _ := json.Marshal(map[string]interface{}{
		"code": 500,
		"msg":  "游客",
	})
	if err != nil {
		w.Write(data)
		return
	}
	tokenData := wp_token.JwtDecryption(token.Value, password)
	user := wp_sql.QueryUserInfo([]string{
		"userName",
	}, []string{
		tokenData["userName"].(string),
	})
	var reData = make(map[string]interface{})
	if user != nil {
		if user.UserPasswd == tokenData["userPasswd"].(string) {
			reData["code"] = 200
			reData["msg"] = "验证token成功"
		} else {
			w.Write(data)
			return
		}
	} else {
		w.Write(data)
		return
	}

	queryParam := r.URL.Query().Get("fileId")

	wp_sql.RemoveFile(queryParam)

	log.Println(queryParam)
	w.Write([]byte(queryParam))
}

func wp_zd_RemoveFileJHeader(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	data, _ := json.Marshal(map[string]interface{}{
		"code": 500,
		"msg":  "游客",
	})
	if err != nil {
		w.Write(data)
		return
	}
	tokenData := wp_token.JwtDecryption(token.Value, password)
	user := wp_sql.QueryUserInfo([]string{
		"userName",
	}, []string{
		tokenData["userName"].(string),
	})
	var reData = make(map[string]interface{})
	if user != nil {
		if user.UserPasswd == tokenData["userPasswd"].(string) {
			reData["code"] = 200
			reData["msg"] = "验证token成功"
		} else {
			w.Write(data)
			return
		}
	} else {
		w.Write(data)
		return
	}

	queryParam := r.URL.Query().Get("fileId")

	wp_sql.ZdRemoveFile(queryParam)

	log.Println(queryParam)
	w.Write([]byte(queryParam))
}

func wpHuiFUFileHeader(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	data, _ := json.Marshal(map[string]interface{}{
		"code": 500,
		"msg":  "游客",
	})
	if err != nil {
		w.Write(data)
		return
	}
	tokenData := wp_token.JwtDecryption(token.Value, password)
	user := wp_sql.QueryUserInfo([]string{
		"userName",
	}, []string{
		tokenData["userName"].(string),
	})
	var reData = make(map[string]interface{})
	if user != nil {
		if user.UserPasswd == tokenData["userPasswd"].(string) {
			reData["code"] = 200
			reData["msg"] = "验证token成功"
		} else {
			w.Write(data)
			return
		}
	} else {
		w.Write(data)
		return
	}

	queryParam := r.URL.Query().Get("fileId")

	wp_sql.HuiFU(queryParam)

	log.Println(queryParam)
	w.Write([]byte(queryParam))
}

func deleteCookieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		w.Write([]byte("Cookie 'token' has been deleted."))
	}
}

func setOwnMsgHeader(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/setMessage.html"))
		tmpl.Execute(w, nil)
	}
}

func changeOwnPwdHeader(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a, _ := io.ReadAll(r.Body)
		var data1 = make(map[string]interface{})
		json.Unmarshal(a, &data1)
		log.Println(data1)

		token, err := r.Cookie("token")
		if err != nil {
			log.Println(err)
			data, _ := json.Marshal(map[string]interface{}{
				"code": 500,
				"msg":  "没有权限",
			})
			w.Write(data)
			return
		}

		tokenData := wp_token.JwtDecryption(token.Value, password)

		userInfo := wp_sql.QueryUserInfo([]string{
			"userName",
		}, []string{
			fmt.Sprint(tokenData["userName"]),
		})

		var reData = make(map[string]interface{})

		if userInfo != nil {
			if userInfo.UserPasswd == fmt.Sprint(data1["oldPwd"]) {
				userInfo.UserPasswd = ""
				reData["code"] = 200
				reData["msg"] = "ok"
				wp_sql.ChangeAdminPassword(tokenData["userName"].(string), data1["pwd"].(string))
			} else {
				reData["code"] = 500
				reData["msg"] = "账户或密码错误"
			}
		} else {
			reData["code"] = 500
			reData["msg"] = "账户或密码错误"
		}
		data11, err := json.Marshal(reData)
		if err != nil {
			log.Panicln(err)
		}

		w.Write(data11)
	}
}

func changeOwnMsgHeader(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 最大 10MB 的文件大小限制
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var key []string
	var value []string

	// 检查是否有文件被上传
	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
	} else {
		// 获取上传的文件
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		localFile, err := os.Create("static/userData/" + handler.Filename)
		if err != nil {
			log.Panicln(err)
		}
		defer localFile.Close()
		// 将上传的文件内容写入到本地文件中
		_, err = io.Copy(localFile, file)
		if err != nil {
			log.Panicln(err)
		}
		key = append(key, "userImg")
		value = append(value, "/static/userData/"+handler.Filename)
	}

	log.Println("r.MultipartForm.Value:", r.MultipartForm.Value)

	email := r.MultipartForm.Value["email"]
	sex := r.MultipartForm.Value["sex"]

	log.Println()

	if len(email) != 0 {
		key = append(key, "eMail")
		value = append(value, email[0])
	}
	if len(sex) != 0 {
		key = append(key, "sex")
		value = append(value, sex[0])
	}

	wp_sql.UpdateUserInfo(key, value)

	aaa := map[string]interface{}{
		"code": 200,
	}
	data, _ := json.Marshal(aaa)
	w.Write(data)
}
