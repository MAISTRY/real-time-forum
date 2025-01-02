function StartWebSocket() {

    if (window["WebSocket"]) {
        // Connect to websocket
        socket = new WebSocket("wss://" + document.location.host + "/ws");
        console.log("supports websockets")

        socket.onopen = () => {
            console.log("Connection open ...");
        }

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            switch (data.type) {
                case 'loadUsersResponse':
                    loadUsers(data.users);
                    break;
                case 'getMessages':
                    GetMessages(data.messages, data.Sender, data.Receiver, data.ReceiverID);
                    break;
                case 'SendMessage':
                    AddMessage(data.messages, data.Sender);
                    sendWebSocketMessage({type: 'loadUsers'}); // Corrected message type
                    break;
                case 'Offline':
                    alert('User offline');
                    // change this way, noobs only use alert
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
}

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

function loadUsers(Users) {
    const UsersContainer = document.getElementById('users-list');
    
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

        // Remove from existing map since we've processed this user
        existingUserItems.delete(userId);
    });

    // Remove user items that are no longer in the new user list
    existingUserItems.forEach(item => item.remove());
}

function GetMessages(messages, Sender, Receiver, ReceiverID) {

    const MessagesContainer = document.getElementById('chat-area');
    MessagesContainer.innerHTML = '<p style="text-align: center">Loading Messages...</p>';

    const fragments = document.createDocumentFragment();
    const ChatHeader = document.createElement('div');
    ChatHeader.className = 'chat-header';
    
    const Chatimage = document.createElement('img');
    Chatimage.className = 'avatar';
    Chatimage.src = '../images/pengin.PNG';
    
    const HeaderParts = document.createElement('div');

    const ChatName = document.createElement('div');
    ChatName.className = 'chat-name';
    ChatName.textContent = Receiver;
    
    const TypingStatus = document.createElement('div');
    TypingStatus.className = 'Typing';
    TypingStatus.id = `TypingStatus-${ReceiverID}`;

    HeaderParts.appendChild(ChatName);
    HeaderParts.appendChild(TypingStatus);
    ChatHeader.appendChild(Chatimage);
    ChatHeader.appendChild(HeaderParts);

    const ChatBody = document.createElement('div');
    ChatBody.id = 'messages';
    ChatBody.className = 'messages';

    if (messages) {

        const messageChunks = [];
        const totalMessages = messages.length;

        const remainder = totalMessages % 10;
        let startIndex = 0;

        if (remainder !== 0) {
            const firstChunk = messages.slice(0, remainder);
            messageChunks.push(firstChunk);
            startIndex = remainder;
        }

        for (let i = startIndex; i < totalMessages; i += 10) {
            const chunk = messages.slice(i, i + 10);
            messageChunks.push(chunk);
        }

        let ChunkNumber = messageChunks.length-1;
        let test = messages.length;

        ChatBody.addEventListener('scroll', throttle(handleScroll, 200));

        function handleScroll() {
            if (ChatBody.scrollTop === 0 && ChunkNumber != 0) {
                ChunkNumber -= 1;
                loadMoreMessages(ChunkNumber);
            }
        }

        function loadMoreMessages(num) {       

            const previousScrollHeight = ChatBody.scrollHeight;
            const Chunk = document.createDocumentFragment();

            messageChunks[num].forEach(message =>{
                const messageItem = document.createElement('div');
                if (message.FirstUser === Sender) {
                    messageItem.className = 'message sent'; 
                } else {
                    messageItem.className = 'message received';
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

                Chunk.appendChild(messageItem);
            })
            ChatBody.prepend(Chunk);
            ChatBody.scrollTop = ChatBody.scrollHeight - previousScrollHeight;
        }
        loadMoreMessages(ChunkNumber);
    }

    const ChatFooter = document.createElement('div');
    ChatFooter.className = 'input-area';
    
    const inputArea = document.createElement('input');
    inputArea.className = 'message-input';
    inputArea.placeholder = 'Type a message...';
    inputArea.type = 'text';

    const sendButton = document.createElement('button');
    sendButton.className ='send-button';
    sendButton.textContent = 'Send';

    sendButton.addEventListener('click', () => {
        const message = inputArea.value;
        if (message.trim()!== '') {
            sendWebSocketMessage({type: 'SendMessage', message: message, secondUser: ReceiverID, Receiver: Receiver});
            inputArea.value = '';
        }
    });

    inputArea.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            const message = inputArea.value;
            if (message.trim() !== '') {
                sendWebSocketMessage({type: 'SendMessage', message: message, secondUser: ReceiverID, Receiver: Receiver});
                inputArea.value = '';
            }
        }
    });

    ChatFooter.appendChild(inputArea);
    ChatFooter.appendChild(sendButton);

    fragments.appendChild(ChatHeader);
    fragments.appendChild(ChatBody);
    fragments.appendChild(ChatFooter);

    MessagesContainer.innerHTML = '';
    MessagesContainer.appendChild(fragments);

    ChatBody.scrollTo({
        top: ChatBody.scrollHeight,
        behavior: 'smooth',
    });
    
}

function AddMessage(message, Sender) {

    const ChatBody = document.getElementById('messages');

    const messageItem = document.createElement('div');
    if (message.FirstUser === Sender) {
        messageItem.className = 'message sent'; 
    } else {
        messageItem.className = 'message received';

        if (document.visibilityState === 'hidden') {
            console.log('User is not viewing the page.');
            notification(message.FirstUser, message.message, message.Sender);
        } 

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

    ChatBody.scrollTo({
        top: ChatBody.scrollHeight,
        behavior: 'smooth',
    });
}