function StartWebSocket() {

    if (window["WebSocket"]) {
        // Connect to websocket
        try {
            socket = new WebSocket("wss://" + document.location.host + "/ws");
            console.log("supports websockets")
    
        } catch (error) {
            console.log("no error")
        }
        
        socket.onopen = () => {
            console.log("Connection open ...");
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            switch (data.type) {
                case 'status':
                    // console.log(`Status: ${data}`);
                    Authenticated(data);
                    break;
                case 'loadUsersResponse':
                    loadUsers(data.users, false);
                    break;
                case 'loadUsrAfterResponse':
                    loadUsers(data.users, true);
                    break;
                case 'getMessages':
                    GetMessages(data.messages, data.Sender, data.Receiver, data.ReceiverID);
                    break;
                case 'SendMessage':
                    AddMessage(data.messages, data.Sender, data.ReceiverID);
                    sendWebSocketMessage({type: 'loadUsrAfterMsg', firstUser: data.Sender, secondUser: data.ReceiverID}); 
                    break;
                case 'IsTyping':
                    Typing(data.Sender, data.isTyping);
                    break;z
                case 'Offline':
                    alert('User offline');
                    break;
                default:
                    console.error('Invalid data type');
                    break;
            }
        };

        socket.onclose = () => {
            console.log("Connection closed");
        };

        socket.onerror = function(event) {
            console.log("Error occurred: " + event.data);
        };
    } else {
        alert("Not supporting websockets");
    }
}

function Typing(Sender,isTyping) {
    const MessagesContainer = document.getElementById('chat-area');
    const ChatBody = document.getElementById('messages');
    const isSender = MessagesContainer.classList.contains(Sender);
    if (!isSender) {
        return;
    } 
    const ChatTyping = document.getElementById('Typing');
    if (isTyping) {
        ChatTyping.classList.remove('Typing');
        ChatBody.scrollTo({
            top: ChatBody.scrollHeight,
            behavior: 'smooth',
        });
    } else {
        ChatTyping.classList.add('Typing');
    }
}

window.sendWebSocketMessage = sendWebSocketMessage;
function sendWebSocketMessage(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.error("Socket not open");
    }
}

function notification(SenderID, Message, Sender) {
    if (Notification.permission === 'default') {
        Notification.requestPermission();
    } else if (Notification.permission === 'granted') {
        new Notification(Sender ,{
            body: Message,
            icon: "../images/pengin.PNG",
            tag: SenderID
        });
    } else {
        alert(`${Sender}: ${Message}`);
    }
}

function loadUsers(Users, newMessage) {
    const UsersContainer = document.getElementById('users-list');

    if (Users === null) {
        console.log('No User found');
        return
    }
    
    // Keep track of existing user IDs, this will be empty when reloading
    const existingUserItems = new Map();
    UsersContainer.querySelectorAll('.user-item').forEach(item => {
        const userId = item.getAttribute('data-user-id');
        existingUserItems.set(userId, item);
    });

    Users.forEach(user => {
        const userId = user.userID.toString();
        let userItem = existingUserItems.get(userId);

        // Create new user item if it doesn't exist
        if (!userItem) {
            userItem = document.createElement('div');
            userItem.className = 'user-item';
            userItem.setAttribute('data-user-id', userId);

            const ChatInfo = document.createElement('div');
            ChatInfo.className = 'chat-info';

            const username = document.createElement('span');
            username.className = 'chat-name';

            const messageTime = document.createElement('span');
            messageTime.className = 'message-Time';

            const lastMessage = document.createElement('div');
            lastMessage.className = 'last-message';

            const status = document.createElement('div');

            ChatInfo.appendChild(username);
            ChatInfo.appendChild(lastMessage);
            ChatInfo.appendChild(messageTime);

            userItem.appendChild(ChatInfo);
            userItem.appendChild(status);

            UsersContainer.appendChild(userItem);
        }

        const ChatInfo = userItem.querySelector('.chat-info');
        const username = ChatInfo.querySelector('.chat-name');
        const messageTime = ChatInfo.querySelector('.message-Time');
        const lastMessage = ChatInfo.querySelector('.last-message');
        const status = userItem.querySelector('.status-icon') || document.createElement('div');

        username.textContent = user.username;
        messageTime.textContent = formatDate(user.timestamp);
        if (formatDate(user.timestamp) !== '') {
            messageTime.textContent = 'last message ' + ShortDate(user.timestamp);
        }

        if (user.sender === "") {
            lastMessage.textContent = MessageLength(user.lastMessage);
        } else if (user.sender !== user.username && user.lastMessage !== "") {
            lastMessage.textContent = MessageLength('You: ' + user.lastMessage);
        } else {
            lastMessage.textContent = MessageLength(UserLength(user.sender) + ': ' + user.lastMessage);
        }

        status.className = `status-icon status-${user.status}`;
        if (!userItem.contains(status)) {
            userItem.appendChild(status);
        }

        userItem.addEventListener('click', () => {
            sendWebSocketMessage({type: 'GetMessages', secondUser: user.userID, Receiver: user.username});
            navigateToPage('Messages');
            const menuItems = document.querySelectorAll('.menu-item');
            menuItems.forEach(i => i.classList.remove('color'));
        });
        
        //* ths shit fuck all my hard work to let the userlist not refresh 
        if (newMessage === true) {
            UsersContainer.appendChild(userItem);
        }
        // Remove from existing map since we've processed this user
        existingUserItems.delete(userId);
    });

    // Remove user items that are no longer in the new user list
    existingUserItems.forEach(item => item.remove());
}