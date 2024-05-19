// Login function to validate and set cookie


// function setLoginCookie(username, token) {
//     var expires = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000); // 7 days from now
//     document.cookie = `token=${token}; expires=${expires.toUTCString()}; path=/; username=${username};`;
//     localStorage.setItem('tokenExpires', expires.getTime()); // Store expires as timestamp

//     updateLoginUI(username);
// }


function setLoginCookie(username, token, expire_stamp) {
    var expires_date = new Date(expire_stamp); // 设置cookie有效期为7天
    var expiresStr = `expires=${expires_date.toUTCString()}`; // 将过期时间转换为GMT格式的字符串
    var pathStr = "path=/"; // 设置cookie的路径，确保在整个网站中有效
    document.cookie = `token=${token};`; // 组合字符串设置cookie
    document.cookie = `${expiresStr};`
    document.cookie = `${pathStr};`
    document.cookie = `username=${username};`

    localStorage.setItem('tokenExpires', expire_stamp); // 在localStorage中存储过期时间的时间戳

    updateLoginUI(username); // 更新页面上的用户界面
}


function updateLoginUI(username) {
    ChangeToHelloContainer(username);
}
function changeContent(page) {
    const event = new CustomEvent('navChange', { detail: page });
    document.getElementById('header').dispatchEvent(event);
}

document.addEventListener('DOMContentLoaded', function () {
    function waitForElement() {
        if (document.getElementById('loginContainer')) {
            checkLoginStatus();
        } else {
            setTimeout(waitForElement, 500); // 每500毫秒检查一次
        }
    }

    waitForElement();
});

function checkLoginStatus() {
    var token = getCookie('token');
    var expiresTimestamp = localStorage.getItem('tokenExpires');
    var now = Date.now();

    if (token && expiresTimestamp > now) {
        var username = getCookie("username"); // Assume you store/get this somewhere
        updateLoginUI(username);
    } else {
        logout();
    }
}

function getCookie(name) {
    var value = `; ${document.cookie}`;
    var parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

function logout() {
    var keys = document.cookie.match(/[^ =;]+(?=\=)/g);
    if (keys) {
        for (var i = keys.length; i--;)
            document.cookie = keys[i] + '=0;expires=' + new Date(0).toUTCString()
    }
    localStorage.removeItem('tokenExpires'); // Clear the stored expires time
    localStorage.removeItem("username");

    fetch('header.html')
        .then(response => response.text())
        .then(html => {
            document.getElementById('header-container').innerHTML = html;
            // 初始化header中可能需要的任何JavaScript
            // initializeHeader();
        })
        .catch(error => {
            console.error('Failed to load the header:', error);
        });
}