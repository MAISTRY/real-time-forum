document.addEventListener('DOMContentLoaded', () => {
    if (window["WebSocket"]) {
        // Connect to websocket
        socket = new WebSocket("wss://" + document.location.host + "/ws");
        console.log("supports websockets")

        socket.onclose = function(evt) {
            console.log("Connection closed");
        }
    } else {
        alert("Not supporting websockets");
    }
    loadUsers(5);
})

function loadUsers(id){
    const UsersContainer = document.getElementById('users-list');
    
    UsersContainer.innerHTML = '<p style="text-align: center">Loading Users...</p>';
    fetch(`/Data-Users?user=${id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'x-Requested-With': 'XMLHttpRequest'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }
        return response.json();
    })
    .then(Users => {
        const fragments = document.createDocumentFragment();
        
        Users.forEach(user => {
            const userItem = document.createElement('div');
            userItem.className = 'user-item';

            const ChatInfo = document.createElement('div');
            ChatInfo.className = 'chat-info';

            const username = document.createElement('span');
            username.className = 'chat-name';
            username.textContent = user.username;

            const messageTime = document.createElement('span');
            messageTime.className = 'message-Time';
            messageTime.textContent = formatDate(user.timestamp);
        
            
            const lastMessage = document.createElement('div');
            lastMessage.className = 'last-message';
            if (user.sender === ""){
                lastMessage.textContent = MessageLength(user.lastMessage);
            }
            else if (user.sender !== user.username && user.lastMessage !== ""){
                lastMessage.textContent = MessageLength('You: ' + user.lastMessage);
            }
            else {
                lastMessage.textContent = MessageLength(UserLength(user.sender) + ': ' + user.lastMessage);
            }
            
            const status = document.createElement('div');
            status.className = `status-icon status-${user.status}`;

            ChatInfo.appendChild(username);
            ChatInfo.appendChild(messageTime);
            ChatInfo.appendChild(lastMessage);

            userItem.appendChild(ChatInfo);
            userItem.appendChild(status);

            fragments.appendChild(userItem);
        });

        UsersContainer.innerHTML = '';
        UsersContainer.appendChild(fragments);
    })
    .catch(err => {
        console.error(err);
    })
}

