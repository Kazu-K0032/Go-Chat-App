document.addEventListener('DOMContentLoaded', function() {
    const sendButton = document.getElementById('js-sendButton');
    const messageText = document.getElementById('js-messageText');
    const messageList = document.getElementById('js-messageList');

    sendButton.addEventListener('click', function() {
        const text = messageText.value.trim();
        if (text === '') return;

        const messageElement = document.createElement('div');
        messageElement.className = 'p-messageBoard__item';
        messageElement.textContent = text;
        
        messageList.appendChild(messageElement);
        messageText.value = '';
    });
});
