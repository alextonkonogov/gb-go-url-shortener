document.addEventListener("DOMContentLoaded", function(event) {
    Statistics();
});

function Statistics() {
    let reqData = {};
    let str = window.location.href
    reqData.admin = str.split("/")[4];
    console.log(reqData.admin)
    let reqDataJSON = JSON.stringify(reqData)

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/a", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var data = JSON.parse(xhr.response);
            document.getElementById("response").innerHTML = `
            <div class="alert alert-primary" role="alert">
                <div class="mt-3">
                    <h6>Ваша ссылка</h6>
                    <a href="${data.long}" target="_blank">${data.long}</a>
                </div>
                <div class="mt-3">
                    <h6>Короткая ссылка</h6>
                    <a href="/s/${data.short}" target="_blank">/s/${data.short}</a>
                </div>
                <div class="mt-3">
                    <h6>Последний переход по ссылке:</h6>
                    <p>${data.viewed}</p>
                </div>
                <div class="mt-3">
                    <h6>IP:</h6>
                    <p>${data.ip}</p>
                </div>
                <div class="mt-3">
                    <h6>Кол-во переходов:</h6>
                    <p>${data.count}</p>
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




