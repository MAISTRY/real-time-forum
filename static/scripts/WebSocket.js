document.addEventListener('DOMContentLoaded', () => {
    if (window["WebSocket"]) {
        // Connect to websocket
        socket = new WebSocket("wss://" + document.location.host + "/ws");
        console.log("supports websockets")

        socket.onopen = () => {
            console.log("Connection open ...");
            sendWebSocketMessage({
                type: 'loadUsers'
            });
        }

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            switch (data.type) {
                case 'loadUsersResponse':
                    loadUsers(data.users);
                    break;
                case 'getMessagesResponse':
                    GetMessages(data.firstUser, data.secondUser);
                    break;
                case 'sendMessageResponse':
                    sendMessage(data.message, data.Sender, data.Receiver);
                    break;
                default:
                    console.error('Invalid data type');
                    break;
            }
        }

        socket.onclose = () => {
            console.log("Connection closed");
        }

        socket.onerror = function(event) {
            console.log("Error occurred: " + event.data);
        }

    } else {
        alert("Not supporting websockets");
    }
})

function sendWebSocketMessage(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.error("Socket not open");
    }
}

function loadUsers(Users){
    const UsersContainer = document.getElementById('users-list');
    
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
        if (formatDate(user.timestamp) !== '') {
            messageTime.textContent = 'last message ' + formatDate(user.timestamp);
        }
        
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
        ChatInfo.appendChild(lastMessage);
        ChatInfo.appendChild(messageTime);

        userItem.appendChild(ChatInfo);
        userItem.appendChild(status);

        fragments.appendChild(userItem);
    });

    UsersContainer.innerHTML = '<p style="text-align: center">Loading Users...</p>';
    UsersContainer.innerHTML = '';
    UsersContainer.appendChild(fragments);
}

function GetMessages(firstUser, secondUser) {

    const MessagesContainer = document.getElementById('chat-area');
    MessagesContainer.innerHTML = '<p style="text-align: center">Loading Messages...</p>';
    fetch(`/Data-Message?sender=${firstUser}&receiver=${secondUser}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'x-Requested-With': 'XMLHttpRequest'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to fetch messages');
        }
        return response.json();
    })
    .then(messages => {
        const fragments = document.createDocumentFragment();

        const ChatHeader = document.createElement('div');
        ChatHeader.className = 'chat-header';
        
        const Chatimage = document.createElement('img');
        Chatimage.className = 'avatar';
        Chatimage.src = '../images/pengin.PNG';
        
        const ChatName = document.createElement('div');
        ChatName.className = 'chat-name';
        if (secondUser === messages[0].SecondUser){
            ChatName.textContent = messages[0].Receiver;
        } else {
            ChatName.textContent = messages[0].Sender;
        }
        
        ChatHeader.appendChild(Chatimage);
        ChatHeader.appendChild(ChatName);

        const ChatBody = document.createElement('div');
        ChatBody.className = 'messages';
        
        messages.forEach(message =>{

            const messageItem = document.createElement('div');
            if (firstUser === message.FirstUser) {
                messageItem.className = 'message sent'; 
            }else {
                messageItem.className ='message received';
            }

            const messageSender = document.createElement('div');
            messageSender.className = 'message-sender';
            messageSender.textContent = message.Sender + ':';

            const messageContent = document.createElement('div');
            messageContent.className = 'message-content';
            messageContent.textContent = message.message;

            const messageTime = document.createElement('div');
            messageTime.className = 'message-time';
            messageTime.textContent = displayDate(message.timestamp);

            messageItem.appendChild(messageSender);
            messageItem.appendChild(messageContent);
            messageItem.appendChild(messageTime);

            ChatBody.appendChild(messageItem);
        })

        const ChatFooter = document.createElement('div');
        ChatFooter.className = 'input-area';
        
        const inputArea = document.createElement('input');
        inputArea.className = 'message-input';
        inputArea.placeholder = 'Type a message...';
        inputArea.type = 'text';

        const sendButton = document.createElement('button');
        sendButton.className ='send-button';
        sendButton.textContent = 'Send';

        ChatFooter.appendChild(inputArea);
        ChatFooter.appendChild(sendButton);

        fragments.appendChild(ChatHeader);
        fragments.appendChild(ChatBody);
        fragments.appendChild(ChatFooter);

        MessagesContainer.innerHTML = '';
        MessagesContainer.appendChild(fragments);
    })
    .catch(err => {
        console.error(err);
    })
}