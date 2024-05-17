function login() {
    var password = "!!8964jss";
    var time_stamp = new Date().getTime()

    sha256Hash(password).then(enc_pwd => sha256Hash(enc_pwd + time_stamp)).then(enc2_pwd =>
        console.log(enc2_pwd, time_stamp));
}

async function sha256Hash(str) {
    const encoder = new TextEncoder();
    const data = encoder.encode(str);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return hashHex;
}

login();