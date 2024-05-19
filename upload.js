document.getElementById('uploadForm').addEventListener('submit', function (event) {
    event.preventDefault(); // 阻止表单的默认提交行为

    const fileInput = document.getElementById('fileInput');
    const textInput = document.getElementById("textInput");
    const file = fileInput.files[0];
    const text = textInput.value;


    if (!text) {
        alert("at least one tag");
        return;
    }

    if (!file) {
        alert('Choose one file first');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);
    formData.append("tags", text);
    var token = getCookie("token")
    fetch('/api/upload/file', {
        method: 'POST',
        body: formData,
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
            alert(data.message);
        })
        .catch(error => {
            console.error('Error:', error);
            alert('fail');
        });
});