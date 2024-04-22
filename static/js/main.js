var FileData
var isYouke = false


function checkToken() {
    fetch("/", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.code === 200) {
                document.getElementById("upLodaFile").style.display = "block";

                document.getElementById("loginBut").remove();
                userImg = document.getElementById("userImg");
                userImg.style.display = 'block';
                img_ = userImg.querySelector("img");
                img_.src = data["userInfo"]["UserImg"];
            } else {
                document.getElementById("userImg").remove();
                isYouke = true;
                document.getElementById("sx_fenXiang").style.display = "none"
                document.getElementById("garbagebin").style.display = "none"
                let hr = document.getElementsByClassName("fengexian")
                hr[1].style.display = "none"
                hr[2].style.display = "none"
                document.getElementById("settingBut").style.display = "none"
            }
            document.getElementById("load").style.display = "none";

        })
}
function uploadFile() {
    document.getElementById("upLodaFile").addEventListener("click", function () {
        let fileInput = document.getElementById("fileInput");
        fileInput.click();
    });

    document.getElementById("fileInput").addEventListener("change", function () {

        let file = this.files[0];
        if (!file) {
            return;
        }

        const formData = new FormData();
        formData.append("file", file);

        const xhr = new XMLHttpRequest();
        xhr.upload.onprogress = function (event) {
            if (event.lengthComputable) {
                const percentComplete = (event.loaded / event.total) * 100;
                console.log(percentComplete + '% 上传完成');
                // 在这里可以更新UI来显示进度，例如更新进度条s
            }
        };

        xhr.onload = function () {
            if (xhr.status === 200) {
                console.log('上传成功');
                alert("上传成功");
                // console.log(xhr.responseText);
                getFileList();
            } else {
                console.error('上传失败: ' + xhr.statusText);
            }
        };

        // console.log("发送！");
        xhr.open('POST', "/uploadFile", true);
        xhr.send(formData);
    });
}

// 创建文件列表项
function createFileListItem(fileName, upDate, fileSize, fileId) {
    var listItem = document.createElement('div');
    listItem.classList.add('file_list');

    var checkboxDiv = document.createElement('div');
    var checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.id = fileId;
    checkbox.setAttribute('filename', fileName);
    checkbox.fileName = fileName;
    checkboxDiv.appendChild(checkbox);
    listItem.appendChild(checkboxDiv);

    var fileNameDiv = document.createElement('div');
    fileNameDiv.textContent = fileName;
    fileNameDiv.style.overflow = "hidden";
    fileNameDiv.style.textOverflow = "ellipsis";
    fileNameDiv.title = fileName;
    listItem.appendChild(fileNameDiv);

    var upDateDiv = document.createElement('div');
    upDateDiv.textContent = upDate;
    listItem.appendChild(upDateDiv);

    var fileSizeDiv = document.createElement('div');
    fileSizeDiv.textContent = fileSize;
    listItem.appendChild(fileSizeDiv);

    return listItem;
}

// 插入到 #fileList 元素中
function insertFileListItem(fileName, upDate, fileSize, fileId) {
    var fileList = document.getElementById('fileList');
    var listItem = createFileListItem(fileName, upDate, fileSize, fileId);
    fileList.appendChild(listItem);
}

function formatDate(dateTimeString) {
    var date = new Date(dateTimeString);
    var year = date.getFullYear();
    var month = date.getMonth() + 1; // 月份从0开始，需要加1
    var day = date.getDate();

    // 格式化月份和日期，保证是两位数
    month = month < 10 ? '0' + month : month;
    day = day < 10 ? '0' + day : day;

    return year + '-' + month + '-' + day;
}

function setFileSize(a) {
    let byte = a
    let kByte = byte / 1024
    if (kByte < 1) {
        return byte + " B"
    }
    let mByte = kByte / 1024
    if (mByte < 1) {
        return kByte.toFixed(2) + " KB"
    }
    let gByte = mByte / 1024
    if (gByte < 1) {
        return mByte.toFixed(2) + " MB"
    }

    return gByte.toFixed(2) + " GB"

}

document.getElementById("removeFileSharing").style.display = "none"
function getFileList() {
    fetch("/queryfilelist", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => response.json())
        .then(data => {
            FileData = JSON.parse(data["data"]);

            if (FileData == null) {
                FileData = {}
            }
            document.querySelector(".active").click()
        });
}

function choiceFile() {
    let file_Check_ALL = document.getElementById("File_Check_ALL")
    file_Check_ALL.addEventListener('change', function () {

        let inputList = document.getElementById("fileList").querySelectorAll("input")
        for (i = 0; i < inputList.length; i++) {
            inputList[i].checked = this.checked;
        }
    });
}

function downLoadFile() {
    document.getElementById("downloadFile").addEventListener("click", function () {
        var inputList = document.getElementById("fileList").querySelectorAll("input")

        for (i = 0; i < inputList.length; i++) {
            (function (input) {
                if (input.checked) {
                    fetch("/downloadfile?fileId=" + input.id)
                        .then(response => response.blob())
                        .then(blob => {
                            var blobUrl = URL.createObjectURL(blob);
                            var link = document.createElement('a');
                            link.href = blobUrl;
                            link.style.display = 'none';
                            console.log("input", input);
                            link.download = input.getAttribute('filename');
                            document.body.appendChild(link);
                            link.click();
                            document.body.removeChild(link);
                        })
                        .catch(error => console.error('文件下载失败:', error));
                }
            })(inputList[i]);
        }
    })
}

function shareFile() {
    document.getElementById("shareFile").addEventListener("click", function () {
        var inputList = document.getElementById("fileList").querySelectorAll("input")

        for (i = 0; i < inputList.length; i++) {
            (function (input) {
                if (input.checked) {
                    fetch("/setFileSharing?fileId=" + input.id, {
                        method: "GET"
                    })
                        .then(response => response.json())
                        .then(data => {
                            let inputList = document.getElementById("fileList").querySelectorAll("input")
                            for (i = 0; i < inputList.length; i++) {
                                inputList[i].checked = false;
                            }

                            getFileList();
                        })
                }
            })(inputList[i]);
        }

        alert("分享成功");
    })
}

function removeFileSharing() {
    document.getElementById("removeFileSharing").addEventListener("click", function () {
        var inputList = document.getElementById("fileList").querySelectorAll("input")
        console.log(inputList);
        for (i = 0; i < inputList.length; i++) {
            (function (input) {
                if (input.checked) {
                    fetch("/removeFileSharing?fileId=" + input.id, {
                        method: "GET"
                    })
                        .then(response => response.json())
                        .then(data => {
                            let inputList = document.getElementById("fileList").querySelectorAll("input")
                            for (i = 0; i < inputList.length; i++) {
                                inputList[i].checked = false;
                            }

                            getFileList();
                        })
                }
            })(inputList[i]);
        }
        alert("取消分享成功");
    })
}

function filterFileList() {
    let sx_myFile = document.getElementById("sx_myFile")
    let sx_tuPian = document.getElementById("sx_tuPian")
    let sx_wendang = document.getElementById("sx_wendang")
    let sx_shiPing = document.getElementById("sx_shiPing")
    let sx_yingPing = document.getElementById("sx_yingPing")
    let sx_fenXiang = document.getElementById("sx_fenXiang")
    let garbageBin = document.getElementById("garbagebin")

    sx_myFile.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "block"
        document.getElementById("downloadFile").style.display = "block"
        document.getElementById("shareFile").style.display = "block"
        document.getElementById("removeFile").style.display = "block"
        document.getElementById("removeFileSharing").style.display = "none"
        document.getElementById("huifu").style.display = "none"
        document.getElementById("zd_remove").style.display = "none"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }
        removeActiveStyle()
        document.getElementById('fileList').innerHTML = "";
        sx_myFile.classList.add("active")
        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];
            console.log(a);
            if (a["DeletedStatic"] == 0) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }
    })
    sx_tuPian.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "block"
        document.getElementById("downloadFile").style.display = "block"
        document.getElementById("shareFile").style.display = "block"
        document.getElementById("removeFile").style.display = "block"
        document.getElementById("removeFileSharing").style.display = "none"
        document.getElementById("huifu").style.display = "none"
        document.getElementById("zd_remove").style.display = "none"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }

        removeActiveStyle()
        document.getElementById('fileList').innerHTML = "";
        sx_tuPian.classList.add("active")
        var imageRegex = /\.(jpg|jpeg|png|gif|bmp)$/i;
        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];
            if (imageRegex.test(a["FileName"]) && a["DeletedStatic"] == 0) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }

    })
    sx_wendang.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "block"
        document.getElementById("downloadFile").style.display = "block"
        document.getElementById("shareFile").style.display = "block"
        document.getElementById("removeFile").style.display = "block"
        document.getElementById("removeFileSharing").style.display = "none"
        document.getElementById("huifu").style.display = "none"
        document.getElementById("zd_remove").style.display = "none"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }

        removeActiveStyle()
        document.getElementById('fileList').innerHTML = "";
        sx_wendang.classList.add("active")
        var documentRegex = /\.(doc|docx|pdf|txt|rtf|xlsx|xls|pptx)$/i;
        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];
            if (documentRegex.test(a["FileName"]) && a["DeletedStatic"] == 0) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }
    })
    sx_shiPing.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "block"
        document.getElementById("downloadFile").style.display = "block"
        document.getElementById("shareFile").style.display = "block"
        document.getElementById("removeFile").style.display = "block"
        document.getElementById("removeFileSharing").style.display = "none"
        document.getElementById("huifu").style.display = "none"
        document.getElementById("zd_remove").style.display = "none"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }

        removeActiveStyle()
        document.getElementById('fileList').innerHTML = "";
        sx_shiPing.classList.add("active")
        var videoRegex = /\.(mp4|avi|mov|wmv|mkv|flv)$/i;
        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];
            if (videoRegex.test(a["FileName"]) && a["DeletedStatic"] == 0) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }
    })
    sx_yingPing.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "block"
        document.getElementById("downloadFile").style.display = "block"
        document.getElementById("shareFile").style.display = "block"
        document.getElementById("removeFile").style.display = "block"
        document.getElementById("removeFileSharing").style.display = "none"
        document.getElementById("huifu").style.display = "none"
        document.getElementById("zd_remove").style.display = "none"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }

        removeActiveStyle()
        document.getElementById('fileList').innerHTML = "";
        sx_yingPing.classList.add("active")
        var audioRegex = /\.(mp3|wav|flac|aac|ogg)$/i;
        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];
            if (audioRegex.test(a["FileName"]) && a["DeletedStatic"] == 0) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }
    })
    sx_fenXiang.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "none"
        document.getElementById("downloadFile").style.display = "none"
        document.getElementById("shareFile").style.display = "none"
        document.getElementById("removeFile").style.display = "none"
        document.getElementById("removeFileSharing").style.display = "block"
        document.getElementById("huifu").style.display = "none"
        document.getElementById("zd_remove").style.display = "none"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }

        removeActiveStyle()
        document.getElementById('fileList').innerHTML = "";
        sx_fenXiang.classList.add("active")
        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];
            if (a["IsPublic"] && a["DeletedStatic"] == 0) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }
    })
    garbageBin.addEventListener("click", function () {
        document.getElementById("upLodaFile").style.display = "none"
        document.getElementById("downloadFile").style.display = "none"
        document.getElementById("shareFile").style.display = "none"
        document.getElementById("removeFile").style.display = "none"
        document.getElementById("removeFileSharing").style.display = "none"
        document.getElementById("huifu").style.display = "block"
        document.getElementById("zd_remove").style.display = "block"

        if (isYouke) {
            document.getElementById("upLodaFile").style.display = "none";
            document.getElementById("shareFile").style.display = "none";
            document.getElementById("removeFile").style.display = "none";
        }

        removeActiveStyle();
        document.getElementById('fileList').innerHTML = "";
        garbageBin.classList.add("active")

        let len = 0
        for (i = 0; i < FileData.length; i++) {
            a = FileData[i];

            if (a["DeletedStatic"] == 1) {
                len++;
                insertFileListItem(a["FileName"], formatDate(a["UploadData"]), setFileSize(a["FileSize"]), a["FileIndexId"]);
            }
        }
        console.log("len:", len);
        if (len <= 0) {
            var div = document.createElement('div');
            div.style.textAlign = "center"
            div.className = "file_list"
            div.innerHTML = "没有查到文件捏!"
            document.getElementById('fileList').append(div);
        }
    })


}

function removeFile() {
    document.getElementById("removeFile").addEventListener("click", function () {
        var inputList = document.getElementById("fileList").querySelectorAll("input")

        if (confirm("是否删除文件！如果想恢复可在回收站恢复！")) {
            for (i = 0; i < inputList.length; i++) {
                (function (input) {
                    if (input.checked) {
                        fetch("/removeFile?fileId=" + input.id, {
                            method: "GET"
                        })
                            .then(response => response.json())
                            .then(data => {
                                let inputList = document.getElementById("fileList").querySelectorAll("input")
                                for (i = 0; i < inputList.length; i++) {
                                    inputList[i].checked = false;
                                }

                                getFileList();
                            })
                    }
                })(inputList[i]);
            }
            alert("删除成功");
        }

    })
    // getFileList();
}

function zd_removeFile() {
    document.getElementById("zd_remove").addEventListener("click", function () {
        var inputList = document.getElementById("fileList").querySelectorAll("input")
        if (confirm("注意！此删除将无法用常规方式恢复，是否删除？")) {
            for (i = 0; i < inputList.length; i++) {
                (function (input) {
                    if (input.checked) {
                        fetch("/zd_removeFile?fileId=" + input.id, {
                            method: "GET"
                        })
                            .then(response => response.json())
                            .then(data => {
                                let inputList = document.getElementById("fileList").querySelectorAll("input")
                                for (i = 0; i < inputList.length; i++) {
                                    inputList[i].checked = false;
                                }

                                getFileList();
                            })
                    }
                })(inputList[i]);
            }
            alert("删除成功");
        }
    })
}

function huiFu() {
    document.getElementById("huifu").addEventListener("click", function () {
        var inputList = document.getElementById("fileList").querySelectorAll("input")

        for (i = 0; i < inputList.length; i++) {
            (function (input) {
                if (input.checked) {
                    fetch("/huifu?fileId=" + input.id, {
                        method: "GET"
                    })
                        .then(response => response.json())
                        .then(data => {
                            let inputList = document.getElementById("fileList").querySelectorAll("input")
                            for (i = 0; i < inputList.length; i++) {
                                inputList[i].checked = false;
                            }

                            getFileList();
                        })
                }
            })(inputList[i]);
        }
        alert("恢复成功");
    })

}

function removeActiveStyle() {
    let active = document.getElementsByClassName("active")
    for (i = 0; i < active.length; i++) {
        active[i].classList.remove("active");
    }
}

function logout() {
    document.getElementById("logout").addEventListener("click", function () {
        fetch("/logout", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            }
        })
        location.reload();
    })
}

function setting() {
    document.getElementById("settingBut").addEventListener("click", function () {
        location.href = "/setOwnMsg"
    })
}

setting();
checkToken();
uploadFile();
getFileList();
choiceFile();
downLoadFile();
shareFile();
filterFileList();
removeFile();
zd_removeFile();
huiFu();
removeFileSharing();
logout();