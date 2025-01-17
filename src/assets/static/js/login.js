document.querySelector('#submit').addEventListener('click', ()=>{
    const reqData = {
        login:  document.querySelector('#login').value,
        password: document.querySelector('#password').value
    };

    fetch("/login", {
        method: "POST",
        body: JSON.stringify(reqData),
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then((resp) => {
        if (resp.status >= 500) {
            errorMessage("An error occurred, please try later");
        }

        resp.json()
        .then((data) => {
            if (resp.status === 200) {
                loginUser(data);
            }
            else if (resp.status === 400) {
                errorMessage(data.message);
            }
        })
        .catch(error => console.error(error));
    })
    .catch(error => console.error(error));
});

const loginUser = (data) => {
    document.cookie = `access_token=${data.AccessToken}; path=/; SameSite=Strict`;
    document.cookie = `refresh_token=${data.RefreshToken}; path=/; SameSite=Strict`;
    window.location = "/";
};

const errorBox = document.getElementById('errorBox');

const errorMessage = (text) => {
    const msg = `<div id="err-msg" class="alert alert-danger" role="alert">${text}</div>`
    errorBox.innerHTML = msg;

    setTimeout(()=> {
        errorBox.innerHTML = '';
    }, 5000)
};


