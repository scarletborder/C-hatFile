async function GetSelfView() {
    var token = getCookie("token");

    fetch("/api/manager/self", {
        method: 'Get',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    }).then(resp => {
        if (!resp.ok) {
            alert(resp.text);
            return;
        }
        return resp.json();
    }).then(data => {
        let llist = document.getElementById("self_view");
        let res = `<p>In addition, you have ${data["filenum"]} files</p>`
        let info_b = "";
        for (const value of data["files"]) {
            let info = `<li>"${value["name"]}" tags: ${value["tags"]} upload_time: ${value["upload_time"]} <a onclick="DeleteFile(${value["id"]})">DELETE</a></li>`
            info_b += info;
        }
        res = res + "<ul>\n" + info_b + "\n</ul>";
        llist.innerHTML = res;
    })
}

async function DeleteFile(file_id) {
    var token = getCookie("token");

    const params = {
        fid: file_id
    }
    const querystr = new URLSearchParams(params).toString();
    const url = `/api/manager/del?${querystr}`;
    fetch(url, {
        method: "Delete",
        headers: {
            'Authorization': `Bearer ${token}`
        }
    }).then(resp => {
        if (!resp.ok) {
            alert(resp.text);
            return;
        }
        return resp.json();
    }).then(data => {
        alert(data["message"]);
        GetSelfView();
    })
}