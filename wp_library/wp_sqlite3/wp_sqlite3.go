package wp_sql

import (
	"DreamVerseCloud/wp_library/wp_tool"
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type FileIndex struct {
	FileIndexId   int
	FileName      string
	FilePath      string // 显示在前端的文件路径
	UploadData    time.Time
	FileOwnerShip string
	IsPublic      bool
	DeletedStatic int8
	FileMetaId    int64
	FileSize      int64
}

type FileMeta struct {
	FileMetaId   int64
	FileHash_md5 string
	FileAddr     string // 在后端的文件路径
	FileSize     int64
}

// UserData 表示从userData表中检索的用户信息
type UserData struct {
	UserName   string
	UserPasswd string
	UserImg    string
	Sex        int
	EMail      string
}

var FileMetas map[string]FileMeta

var dbConn *sql.DB

func init() {
	wp_tool.CreateDirRecursively("data")
	wp_tool.CreateDirRecursively("data/file")
	wp_tool.CreateDirRecursively("data/tmp")
	if QueryUserInfo([]string{"userName"}, []string{"admin"}) == nil {
		var userData = UserData{
			UserName:   "admin",
			UserPasswd: "123456",
			UserImg:    "/static/userData/default.jpg",
			Sex:        0,
			EMail:      "admin@admin.com",
		}
		CreateUserInfo(&userData)
	}
}

func NewDB() *sql.DB {
	sqlit3Path := "data/data.db"
	conn, err := sql.Open("sqlite3", sqlit3Path)
	if err != nil {
		log.Fatal("打开数据库失败", err)
	}

	// 初始化文件表
	conn.Exec(`
	CREATE TABLE IF NOT EXISTS "main"."fileMeta" (
		"FileMetaId" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"FileHash_md5" text,
		"FileAddr" text,
		"FileSize" INTEGER
	  );

	  CREATE TABLE "main"."fileIndex" (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"FileName" TEXT,
		"FilePath" TEXT,
		"UploadData" DATE,
		"FileOwnerShip" TEXT,
		"IsPublic" INTEGER,
		"DeletedStatic" INTEGER,
		"FileMetaId" INTEGER,
		"FileSize" INTEGER
	  );

	`)
	// 初始化用户表
	conn.Exec(`
	CREATE TABLE IF NOT EXISTS "main"."userData" (
		"userName" text NOT NULL,
		"userPasswd" text NOT NULL,
		"userImg" TEXT,
		"sex" integer,
		"eMail" TEXT,
		PRIMARY KEY ("userName"),
		CONSTRAINT "userName" UNIQUE ("userName" ASC)
	  );
	`)
	return conn
}

func CreateFileMeta(fileMeta FileMeta) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}
	stmt, err := dbConn.Prepare(`
	INSERT INTO fileMeta 
	(FileHash_md5, FileAddr, FileSize) 
	VALUES 
	(?, ?, ?)
	`)
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = stmt.Exec(
		fileMeta.FileHash_md5,
		fileMeta.FileAddr,
		fileMeta.FileSize,
	)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CreateFileIndex(data FileIndex) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}
	stmt, err := dbConn.Prepare(`
	INSERT INTO fileIndex 
	(FileName, FilePath, UploadData, FileOwnerShip, IsPublic, DeletedStatic, FileMetaId, FileSize) 
	VALUES 
	(?, ?, ?, ?, ?, ?, ?,?)
	`)
	if err != nil {
		log.Println("CreateFileIndex:1:", err)
		return false
	}
	_, err = stmt.Exec(
		data.FileName,
		data.FilePath,
		data.UploadData,
		data.FileOwnerShip,
		data.IsPublic,
		data.DeletedStatic,
		data.FileMetaId,
		data.FileSize,
	)
	if err != nil {
		log.Println("CreateFileIndex:2:", err)
		return false
	}
	return true
}

func QueryFileIndexList(keys []string, values []string) []*FileIndex {
	if len(keys) != len(values) {
		log.Println("key and value 应该等长")
		return nil
	}
	if dbConn == nil {
		dbConn = NewDB()
	}
	var query strings.Builder
	query.WriteString(`SELECT * FROM fileIndex WHERE `)
	args := make([]interface{}, len(values))
	for i, key := range keys {
		if i > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString(key)
		query.WriteString(" = ?")
		args[i] = values[i]
	}

	rows, err := dbConn.QueryContext(context.Background(), query.String(), args...)
	if err != nil {
		log.Println("QueryFileIndexList：查询出错：", err)
		return nil
	}
	defer rows.Close()

	var fileIndexes []*FileIndex
	for rows.Next() {
		var fileIndex FileIndex
		err := rows.Scan(
			&fileIndex.FileIndexId,
			&fileIndex.FileName,
			&fileIndex.FilePath,
			&fileIndex.UploadData,
			&fileIndex.FileOwnerShip,
			&fileIndex.IsPublic,
			&fileIndex.DeletedStatic,
			&fileIndex.FileMetaId,
			&fileIndex.FileSize,
		)
		if err != nil {
			log.Println("QueryFileIndexList：扫描结果集出错：", err)
			return nil
		}
		fileIndexes = append(fileIndexes, &fileIndex)
	}

	if err := rows.Err(); err != nil {
		log.Println("QueryFileIndexList：扫描结果集出错：", err)
		return nil
	}

	return fileIndexes
}

func QueryFileIndex(keys []string, values []string) *FileIndex {
	if len(keys) != len(values) {
		log.Println("key and value 应该等长")
		return nil
	}
	if dbConn == nil {
		dbConn = NewDB()
	}
	var query strings.Builder
	query.WriteString(`SELECT * FROM fileIndex WHERE `)
	args := make([]interface{}, len(values))
	for i, key := range keys {
		if i > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString(key)
		query.WriteString(" = ?")
		args[i] = values[i]
	}

	row := dbConn.QueryRowContext(context.Background(), query.String(), args...)
	var fileIndex FileIndex
	err := row.Scan(
		&fileIndex.FileIndexId,
		&fileIndex.FileName,
		&fileIndex.FilePath,
		&fileIndex.UploadData,
		&fileIndex.FileOwnerShip,
		&fileIndex.IsPublic,
		&fileIndex.DeletedStatic,
		&fileIndex.FileMetaId,
		&fileIndex.FileSize,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到记录
			return nil
		}
		// log.Println("QueryFileIndex：查询出错：", err)
		return nil
	}
	return &fileIndex
}

func QueryFileMeta(keys []string, values []string) *FileMeta {
	if len(keys) != len(values) {
		log.Println("key and value 应该等长")
		return nil
	}
	if dbConn == nil {
		dbConn = NewDB()
	}

	var query strings.Builder
	query.WriteString(`SELECT * FROM fileMeta WHERE `)
	args := make([]interface{}, len(values))
	for i, key := range keys {
		if i > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString(key)
		query.WriteString(" = ?")
		args[i] = values[i]
	}

	row := dbConn.QueryRowContext(context.Background(), query.String(), args...)
	var fileMeta FileMeta
	err := row.Scan(
		&fileMeta.FileMetaId,
		&fileMeta.FileHash_md5,
		&fileMeta.FileAddr,
		&fileMeta.FileSize,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到记录
			return nil
		}
		// log.Println("QueryFileIndex：查询出错", err)
		return nil
	}
	return &fileMeta
}

func QueryUserInfo(keys []string, values []string) *UserData {
	if len(keys) != len(values) && len(keys) != 0 {
		log.Println("len(keys):", len(keys), "len(values):", len(values))
		log.Println("key and value 应该等长,且长度不能为0")
		return nil
	}
	if dbConn == nil {
		dbConn = NewDB()
	}

	var query strings.Builder
	query.WriteString(`SELECT * FROM userData WHERE `)
	args := make([]interface{}, len(values))
	for i, key := range keys {
		if i > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString(key)
		query.WriteString(" = ?")
		args[i] = values[i]
	}

	row := dbConn.QueryRowContext(context.Background(), query.String(), args...)
	var user UserData
	err := row.Scan(&user.UserName, &user.UserPasswd, &user.UserImg, &user.Sex, &user.EMail)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到记录
			return nil
		}
		log.Println("QueryUserInfo：查询出错：", err)
		return nil
	}

	return &user
}

func CreateUserInfo(userData *UserData) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}
	stmt, err := dbConn.Prepare(`
	INSERT INTO userData 
	(userName, userPasswd, userImg, sex, eMail) 
	VALUES 
	(?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = stmt.Exec(
		userData.UserName,
		userData.UserPasswd,
		userData.UserImg,
		userData.Sex,
		userData.EMail,
	)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func SetFileSharing(fileId string) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}

	_, err := dbConn.Exec("UPDATE fileIndex SET IsPublic = 1 WHERE id = ?", fileId)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func RemoveFileSharing(fileId string) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}

	_, err := dbConn.Exec("UPDATE fileIndex SET IsPublic = 0 WHERE id = ?", fileId)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func RemoveFile(fileId string) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}

	_, err := dbConn.Exec("UPDATE fileIndex SET DeletedStatic = 1 WHERE id = ?", fileId)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func ZdRemoveFile(fileId string) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}

	_, err := dbConn.Exec("UPDATE fileIndex SET DeletedStatic = 2 WHERE id = ?", fileId)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func HuiFU(fileId string) bool {
	if dbConn == nil {
		dbConn = NewDB()
	}

	_, err := dbConn.Exec("UPDATE fileIndex SET DeletedStatic = 0 WHERE id = ?", fileId)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func ChangeAdminPassword(username string, newPassword string) bool {
	// 执行 SQL 更新语句
	_, err := dbConn.ExecContext(context.Background(), `
		UPDATE "userData" SET "userPasswd" = ? WHERE "userName" = ?
	`, newPassword, username)
	if err != nil {
		log.Println("Error updating password:", err)
		return false
	}
	return true
}

func UpdateUserInfo(keys []string, values []string) bool {
	if len(keys) != len(values) {
		return false
	}

	var setValues []string
	for i := range keys {
		setValues = append(setValues, keys[i]+"='"+values[i]+"'")
	}
	setClause := strings.Join(setValues, ", ")

	query := "UPDATE userData SET " + setClause

	_, err := dbConn.Exec(query)
	if err != nil {
		return err == nil
	}

	return true
}
