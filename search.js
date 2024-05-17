function updateResults(results) {
    var resultsList = document.getElementById('results');
    resultsList.innerHTML = '';  // 清空现有结果

    results.forEach(function (result) {
        var li = document.createElement('li');
        li.textContent = `Title: ${result.title}, URL: ${result.url}, Tags: ${result.tags.join(', ')}`;
        resultsList.appendChild(li);
    });
}


function startSearch() {
    var title = document.getElementById('searchTitle').value;
    var tag = document.getElementById('searchTag').value;

    // 模拟搜索结果
    var results = [];
    for (let i = 0; i < 123; i++) {
        results.push({
            title: title + " Result " + (i + 1),
            url: "http://example.com",
            tags: [tag]
        });
    }

    localStorage.setItem('searchResults', JSON.stringify(results)); // 在localStorage中存储results
    performSearch();
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