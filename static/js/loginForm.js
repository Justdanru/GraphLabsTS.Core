const form = document.getElementById("form");

form.addEventListener("submit", function(event) {
    event.preventDefault();
    
    const formData = new FormData(this);
    
    const xhr = new XMLHttpRequest();

    xhr.open("POST", "/api/auth/login", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            //do smth
        }
    };

    xhr.send(JSON.stringify(Object.fromEntries(formData)));
})