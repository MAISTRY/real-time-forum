document.addEventListener('DOMContentLoaded', () => {
    if (window["WebSocket"]) {
        // Connect to websocket
        conn = new WebSocket("wss://" + document.location.host + "/ws");
        console.log("supports websockets")
    } else {
        alert("Not supporting websockets");
    }
})

function loadUsers(){
    const UsersContainer = document.getElementById('users-list');
}

