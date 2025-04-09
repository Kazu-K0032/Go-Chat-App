document.addEventListener("DOMContentLoaded", function () {
  const messageForm = document.getElementById("messageForm");
  const messageInput = document.getElementById("js-messageInput");
  const messageArea = document.getElementById("js-messageArea");

  // テキストエリアの高さを自動調整する関数
  function adjustTextareaHeight(textarea) {
    textarea.style.height = "auto";
    textarea.style.height = textarea.scrollHeight + "px";
  }

  // テキストエリアの入力時に高さを自動調整
  messageInput.addEventListener("input", function () {
    adjustTextareaHeight(this);
  });

  // Ctrl + Enter で送信
  messageInput.addEventListener("keydown", function (e) {
    if (e.key === "Enter" && e.ctrlKey) {
      e.preventDefault();
      messageForm.dispatchEvent(new Event("submit"));
    }
  });

  messageForm.addEventListener("submit", async function (e) {
    e.preventDefault();

    const formData = new FormData(messageForm);
    const chatID = formData.get("chatID");
    const content = formData.get("content").trim();

    // 空のメッセージは送信しない
    if (!content) return;

    try {
      const response = await fetch(window.location.href, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `chatID=${encodeURIComponent(
          chatID
        )}&content=${encodeURIComponent(content)}`,
      });

      if (response.ok) {
        // メッセージを表示
        const messageData = await response.json();
        const formattedContent = messageData.content
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")
          .replace(/\n/g, "<br>");

        const messageHtml = `
          <div class="l-chatMain__message p-message --sent">
            <div class="l-chatMain__content p-message__content">
              <p class="p-message__text c-txt">${formattedContent}</p>
              <time class="p-message__time c-time">${messageData.time}</time>
            </div>
          </div>
        `;
        messageArea.insertAdjacentHTML("beforeend", messageHtml);

        // 入力欄をクリアして高さをリセット
        messageInput.value = "";
        messageInput.style.height = "auto";

        // メッセージエリアを最下部にスクロール
        messageArea.scrollTop = messageArea.scrollHeight;
      } else {
        console.error("メッセージの送信に失敗しました");
      }
    } catch (error) {
      console.error("エラーが発生しました:", error);
    }
  });

  // 初期表示時にすべてのメッセージエリアの高さを調整
  messageInput.dispatchEvent(new Event("input"));
});
