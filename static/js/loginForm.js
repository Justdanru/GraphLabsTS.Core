const form = document.getElementById("form");
const loginInput = document.getElementById("loginInput")
const passwordInput = document.getElementById("passwordInput")
const loginWarning = document.getElementById("loginWarning")
const passwordWarning = document.getElementById("passwordWarning")
const systemWarning = document.getElementById("systemWarning")
const loginFormatWarning = document.getElementById("loginFormatWarning")
const passwordFormatWarning = document.getElementById("passwordFormatWarning")

const loginRegex = /^[a-zA-Z0-9]{5,20}$/;
const passwordRegex = /^[a-zA-Z0-9<>%!.,'\/@#$*()]{7,30}$/

function userNotFoundWarning() {
    loginInput.classList.remove("border-secondary")
    loginInput.classList.add("border-danger")
    loginWarning.classList.remove("visually-hidden")
}

function loginFormatWarningFunc() {
    loginInput.classList.remove("border-secondary")
    loginInput.classList.add("border-danger")
    loginFormatWarning.classList.remove("visually-hidden")
}

function passwordFormatWarningFunc() {
    passwordInput.classList.remove("border-secondary")
    passwordInput.classList.add("border-danger")
    passwordFormatWarning.classList.remove("visually-hidden")
}

function wrongPasswordWarning() {
    passwordInput.classList.remove("border-secondary")
    passwordInput.classList.add("border-danger")
    passwordWarning.classList.remove("visually-hidden")
}

function systemErrorWarning() {
    passwordWarning.classList.remove("visually-hidden")
}

function dropWarnings() {
    loginInput.classList.remove("border-danger")
    loginInput.classList.add("border-secondary")
    passwordInput.classList.remove("border-danger")
    passwordInput.classList.add("border-secondary")
    passwordWarning.classList.add("visually-hidden")
    loginWarning.classList.add("visually-hidden")
}

form.addEventListener("submit", function(event) {
    event.preventDefault();
    
    const formData = new FormData(this);
    
    const xhr = new XMLHttpRequest();

    xhr.open("POST", "/api/auth/login", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 400) {
            answer = JSON.parse(this.responseText)
            if (answer.error == "user not found") {
                userNotFoundWarning();
            } else if (answer.error == "wrong password") {
                wrongPasswordWarning();
            } else if (answer.error == "wrong login format") {
                loginFormatWarningFunc();
            } else if (answer.error == "wrong password format") {
                passwordFormatWarningFunc();
            }
        }
        if (this.readyState == 4 && this.status == 500) {
            systemErrorWarning();
        }
    };

    xhr.send(JSON.stringify(Object.fromEntries(formData)));
});

loginInput.addEventListener("input", function() {
    if (loginRegex.test(loginInput.value)) {
        loginInput.classList.remove("border-danger")
        loginInput.classList.add("border-secondary")
        loginWarning.classList.add("visually-hidden")
        loginFormatWarning.classList.add("visually-hidden")
    } else {
        loginFormatWarningFunc();
    }
});

passwordInput.addEventListener("input", function() {
    if (passwordRegex.test(passwordInput.value)) {
        passwordInput.classList.remove("border-danger")
        passwordInput.classList.add("border-secondary")
        passwordWarning.classList.add("visually-hidden")
        passwordFormatWarning.classList.add("visually-hidden")
    } else {
        passwordFormatWarningFunc();
    }
});