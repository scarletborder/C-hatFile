function updateResults(results) {
    var resultsList = document.getElementById('results');
    resultsList.innerHTML = '';  // 清空现有结果

    results.forEach(function (result) {
        var li = document.createElement('li');
        li.innerHTML = `Title: ${result.title}, Tags: ${result.tags.join(', ')}  <a onclick="DownloadRecord('${result.url}', '${result.title}')">DOWNLOAD</a>`;
        resultsList.appendChild(li);
    });
}

function DownloadRecord(url, filename) {
    var token = getCookie("token")
    fetch(url, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    }).then(response => response.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.download = filename;
            a.href = url;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        })
        .catch(err => console.error('Download failed:', err));
}


function startSearch() {
    var token = getCookie("token")
    var title = document.getElementById('searchTitle').value;
    var tag = document.getElementById('searchTag').value;

    if (title == "" && tag == "") {
        alert("Need at least one argument");
        return;
    }
    let url = `/api/search/search?`

    if (title != "") {
        url += `title=${title}`
        if (tag != "") {
            url += `&tags=${tag}`
        }
    } else if (tag != "") {
        url += `tags=${tag}`
    }




    fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    }).then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    }).then(data => {
        var results = data.results;
        if (!results) {
            alert("No results");
            return;
        }

        // 模拟搜索结果
        // for (let i = 0; i < 123; i++) {
        //     results.push({
        //         title: title + " Result " + (i + 1),
        //         url: "http://example.com",
        //         tags: [tag]
        //     });
        // }

        localStorage.setItem('searchResults', JSON.stringify(results)); // 在localStorage中存储results
        performSearch();
        // console.log('Success:', data);
    }).catch(error => {
        alert('Error:', error);
    });
}


function performSearch(page = 1) {
    var results = JSON.parse(localStorage.getItem('searchResults'));
    const resultsPerPage = 10;
    const totalPages = Math.ceil(results.length / resultsPerPage);
    const pageResults = results.slice((page - 1) * resultsPerPage, page * resultsPerPage);

    updateResults(pageResults);
    renderPagination(page, totalPages);
}

function renderPagination(currentPage, totalPages) {
    var paginationDiv = document.getElementById('pagination');
    var totalPagesInfo = document.getElementById('totalPagesInfo'); // 获取显示总页数的div

    paginationDiv.innerHTML = ''; // 清空现有的分页按钮
    totalPagesInfo.textContent = `Total pages: ${totalPages}`; // 更新总页数信息

    let startPage = Math.max(1, currentPage - 3);
    let endPage = Math.min(totalPages, currentPage + 3);

    for (let i = startPage; i <= endPage; i++) {
        let pageLink = document.createElement('button');
        pageLink.innerText = i;
        pageLink.onclick = function () { performSearch(i); };
        if (currentPage === i) {
            pageLink.disabled = true;
        }
        paginationDiv.appendChild(pageLink);
    }
}