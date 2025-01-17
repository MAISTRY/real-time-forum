function GetMessages(messages, Sender, Receiver, ReceiverID) {

    const MessagesContainer = document.getElementById('chat-area');
    MessagesContainer.classList = ReceiverID + ' chat-area ' + Sender;
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

    HeaderParts.appendChild(ChatName);
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

    const ChatTyping = document.createElement('div');
    ChatTyping.className = 'message received typing-indicator Typing';
    ChatTyping.id = 'Typing';

    const Dot1 = document.createElement('span');
    Dot1.className = 'dot';
    const Dot2 = document.createElement('span');
    Dot2.className = 'dot';
    const Dot3 = document.createElement('span');
    Dot3.className = 'dot';
    
    ChatTyping.appendChild(Dot1);
    ChatTyping.appendChild(Dot2);
    ChatTyping.appendChild(Dot3);

    const ChatFooter = document.createElement('div');
    ChatFooter.className = 'input-area';
    
    const inputArea = document.createElement('input');
    inputArea.className = 'message-input';
    inputArea.placeholder = 'Type a message...';
    inputArea.type = 'text';

    const sendButton = document.createElement('button');
    sendButton.className ='send-button';
    sendButton.textContent = 'Send';

    let typingTimeout;

    sendButton.addEventListener('click', () => {
        const message = inputArea.value;
        if (message.trim()!== '') {
            sendWebSocketMessage({type: 'SendMessage', message: message, secondUser: ReceiverID, Receiver: Receiver});
            inputArea.value = '';
            sendTypingStatus(false);
            clearTimeout(typingTimeout);
        }
    });

    inputArea.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            const message = inputArea.value;
            if (message.trim() !== '') {
                sendWebSocketMessage({type: 'SendMessage', message: message, secondUser: ReceiverID, Receiver: Receiver});
                inputArea.value = '';
                sendTypingStatus(false);
                clearTimeout(typingTimeout);
            }
        }
    });

    inputArea.addEventListener('input', () => {
        if (inputArea.value.trim() !== '') {
            sendTypingStatus(true);
        }
        clearInterval(typingTimeout);
        typingTimeout = setTimeout(() => { sendTypingStatus(false) }, 5000);
    });

    function sendTypingStatus(isTyping) {
        sendWebSocketMessage({type: 'Typing', isTyping: isTyping, SecondUser: ReceiverID, FirstUser: Sender});
    }

    ChatFooter.appendChild(inputArea);
    ChatFooter.appendChild(sendButton);

    fragments.appendChild(ChatHeader);
    fragments.appendChild(ChatBody);
    fragments.appendChild(ChatTyping);
    fragments.appendChild(ChatFooter);

    MessagesContainer.innerHTML = '';
    MessagesContainer.appendChild(fragments);

    ChatBody.scrollTo({
        top: ChatBody.scrollHeight,
        behavior: 'smooth',
    });
    
}

window.AddMessage = AddMessage;
function AddMessage(message, Sender, ReceiverID) {

    const MessagesContainer = document.getElementById('chat-area');
    const isReceiverID = MessagesContainer.classList.contains(ReceiverID);

    const ChatBody = document.getElementById('messages');

    const messageItem = document.createElement('div');
    if (message.FirstUser === Sender) {
        messageItem.className = 'message sent'; 
    } else {
        if (document.visibilityState === 'hidden' || !isReceiverID) {
            console.log('User is not viewing the page.');
            notification(message.FirstUser, message.message, message.Sender);
            return;
        } 
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

    ChatBody.appendChild(messageItem);

    ChatBody.scrollTo({
        top: ChatBody.scrollHeight,
        behavior: 'smooth',
    });
}