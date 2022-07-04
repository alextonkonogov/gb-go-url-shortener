function create() {
    let reqData = {};
    let longURL = document.getElementById("long").value;
    reqData.long = longURL
    let reqDataJSON = JSON.stringify(reqData)

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/create", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var data = JSON.parse(xhr.response);
            document.getElementById("response").innerHTML = `
            <div class="alert alert-success" role="alert">
                <div class="mt-3">
                    <h6>Ваша ссылка</h6>
                    <a href="${data.long}" target="_blank">${data.long}</a>
                </div>
                <div class="mt-3">
                    <h6>Короткая ссылка (ей можно поделиться)</h6>
                    <button type="button" class="btn btn-link" onclick="read(\`/s/${data.short}\`)">/s/${data.short}</button>
                </div>
                <div class="mt-3">
                    <h6>Админская ссылка (оставьте у себя)</h6>
                    <a href="/a/${data.admin}" target="_blank">/a/${data.admin}</a>
                </div>
            </div>`
        }
    };
    try {
        xhr.send(reqDataJSON);
    } catch (err) {
        console.log(err)
    }
}

function read(url) {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", `${url}`, true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var data = JSON.parse(xhr.response);
            document.location.href = data.long;
        }
    };
    try {
        xhr.send();
    } catch (err) {
        console.log(err)
    }
}


