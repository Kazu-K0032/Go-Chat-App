// ... existing code ...

// チャット関連のルート
mux.HandleFunc("/api/chat/send", chatController.SendMessage)
mux.HandleFunc("/ws/chat/", service.WebSocketHandler)

// ... existing code ...
