function login() {
    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    var time_stamp = new Date().getTime()

    sha256Hash(password).then(enc_pwd => sha256Hash(enc_pwd + time_stamp)).then(enc2_pwd =>
        loginRequest(username, enc2_pwd, time_stamp));
}

async function sha256Hash(str) {
    const encoder = new TextEncoder();
    const data = encoder.encode(str);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return hashHex;
}

async function loginRequest(username, encrypted_pwd, time_stamp) {
    fetch("api/auth/login", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, encrypted_pwd, time_stamp })
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