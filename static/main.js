var socket = null;
let o = document.getElementById('output');
let userField = document.getElementById('username');
let messageField = document.getElementById('message');

window.onbeforeunload = function () {
    let jsonData = {};
    jsonData['action'] = 'left';
    socket.send(JSON.stringify(jsonData));
}
document.addEventListener('DOMContentLoaded', () => {
    socket = new WebSocket('ws://localhost:8080/ws');
    socket.onopen = () => {
        console.log("Connected to websocket");
    };
    socket.onclose = () => {
        console.log("Disconnected from websocket");
    };
    socket.onerror = (err) => {
        console.log("Error: ", err);
    };
    socket.onmessage = msg => {
        let data = JSON.parse(msg.data);
        switch (data.action) {
            case "list_users":
                let ul = document.getElementById("online_users");
                while (ul.firstChild) ul.removeChild(ul.firstChild);
                if (data.connected_users.length > 0) {
                    data.connected_users.forEach(user => {
                        let li = document.createElement("li");
                        li.classList.add("list-group-item");
                        li.innerText = "🟢 " + user;
                        ul.appendChild(li);
                    });
                }
                break;
            case "broadcast":
                o.innerHTML = data.message + "<br>" + o.innerHTML;
                break;
        }
    };

    let userInput = document.getElementById('username');
    userInput.addEventListener('change', () => {
        let jsonDate = {};
        jsonDate['action'] = "username";
        jsonDate['username'] = userInput.value;
        socket.send(JSON.stringify(jsonDate));
    });

    document.getElementById('message').addEventListener('keydown', (event) => {
        if (event.code === "Enter") {
            if (!socket) {
                console.log("No socket connection");
                return false;
            }
            event.preventDefault();
            event.stopPropagation();
            if ((userField.value === '') || (messageField.value === '')) {
                errorMessage('Please enter your name and message');
                return false;
            } else {
                sentMessage();
            }

        }
    });

    document.getElementById('sendBtn').addEventListener('click', () => {
        if ((userField.value === '') || (messageField.value === '')) {
            errorMessage('Please enter your name and message');
            return false;
        } else {
            sentMessage();
        }
    });
});

// function send message
function sentMessage() {
    let jsonData = {};
    jsonData['action'] = 'broadcast';
    jsonData['username'] = document.getElementById('username').value;
    jsonData['message'] = document.getElementById('message').value;

    socket.send(JSON.stringify(jsonData));
    document.getElementById('message').value = '';
}

function errorMessage(msg) {
    notie.alert({
        type: 'error',
        text: msg,
    })
}