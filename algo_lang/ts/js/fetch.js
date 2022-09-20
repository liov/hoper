let token = "Bearer.eyJhbGciOiJIUzI1NiJ9.eyJ1aWQiOiIxIiwibG9naW5JcCI6IjEwLjQyLjEuMCIsImxvZ2luTmFtZSI6ImFkbWluIiwibG9naW5UaW1lTWlsbGlzIjoxNTU4NjY5NDE4MjU0LCJleHAiOjE1NTg3NTU4MTh9.uRuK4EvQ5HATumpUt9yywfay_ONxQD6DZawDv8-nDGs"

fetch('http://localhost:8030/api/dubbo/resource/export/res_export_trademark_info', {
    method: 'post', body: JSON.stringify({pageNo:1,pageSize:10}), responseType: 'blob', headers: {
        "Content-Type": "application/json",
        "auth-token": token
    }
})
    .then(res => {
        return res.blob();
    }).then(blob => {


    let a = document.createElement('a');
    a.download = "test.xlsx";
    a.style.display = 'none';
    blob.type = "application/excel";
    let url = window.URL.createObjectURL(blob);
    a.href = url;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
})
