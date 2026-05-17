function getValue(id){
    return document.getElementById(id).value;
}

function get(id){
    return document.getElementById(id);
}

function createElement(element){
    return document.createElement(element);
}

function showMessage(text, isError = false) {
    const msg = get("message");
    msg.textContent = text;
    msg.style.color = isError ? "red" : "green";
    setTimeout(() => msg.textContent = "", 3000);
}