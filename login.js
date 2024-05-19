function login() {
    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    var time_stamp = new Date().getTime()

    sha256Hash(password).then(enc_pwd => sha256Hash(enc_pwd + time_stamp)).then(enc2_pwd =>
        loginRequest(username, enc2_pwd, time_stamp));
}

async function register() {
    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    var password2 = document.getElementById('password2').value;

    if (password != password2) {
        alert("Password must be the same");
        return;
    }

    var encrypted_pwd = await sha256Hash(password);
    const data = { username, encrypted_pwd }
    const urlEncodedData = Object.keys(data).map(key =>
        encodeURIComponent(key) + '=' + encodeURIComponent(data[key])
    ).join('&');

    fetch("api/auth/register", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: urlEncodedData,
    }).then(async response => {
        let json_data = await response.json();
        alert(json_data.message);
        if (response.ok) {
            ChangeToLoginContainer();
        }
    })

}

async function sha256Hash(str) {
    hash = CryptoJS.SHA256(str);
    return hash.toString(CryptoJS.enc.Hex);

    // const encoder = new TextEncoder();
    // const data = encoder.encode(str);
    // const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    // const hashArray = Array.from(new Uint8Array(hashBuffer));
    // const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    // return hashHex;
}

async function loginRequest(username, encrypted_pwd, time_stamp) {
    // 将对象转化为urlencoded格式的字符串
    const data = { username, encrypted_pwd, time_stamp }
    const urlEncodedData = Object.keys(data).map(key =>
        encodeURIComponent(key) + '=' + encodeURIComponent(data[key])
    ).join('&');

    fetch("api/auth/login", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: urlEncodedData,
    }).then(response => response.json())
        .then(data => {
            if (data.token) {
                // 成功获得token, expire_stamp
                setLoginCookie(username, data.token, data.expire_stamp);
            } else {
                alert(`Login Fail: ${data.message}`);
            }
        })
}


async function ChangeToRegisterContainer() {
    document.getElementById('loginContainer').innerHTML = `<input type="text" id="username" placeholder="Username"><br>
            <input type="password" id="password" placeholder="Password"><br>
            <input type="password" id="password2" placeholder="Repeat Password"><br>
            <button onclick="register()">Register</button>
            <button onclick="ChangeToLoginContainer()">To Login</button>`;
}

async function ChangeToLoginContainer() {
    document.getElementById('loginContainer').innerHTML = `<input type="text" id="username" placeholder="Username"><br>
            <input type="password" id="password" placeholder="Password"><br>
            <button onclick="login()">Login</button>
            <button onclick="ChangeToRegisterContainer()">To Register</button>`;
}

async function ChangeToHelloContainer(username) {
    document.getElementById('loginContainer').innerHTML = `<span>Welcome, ${username}!</span>
         <button onclick="logout()">Logout</button>`;
}