// アイコン画像の変更を処理する関数
function handleIconChange(event) {
  const file = event.target.files[0];
  if (file) {
    // 画像のプレビューを表示
    const reader = new FileReader();
    reader.onload = function (e) {
      const img = document.getElementById("profile-icon");
      img.src = e.target.result;

      // フォームを自動送信
      event.target.form.submit();
    };
    reader.readAsDataURL(file);
  }
}

// タブ切り替えの処理
function handleTabClick(event) {
  event.preventDefault();
  const targetId = event.target.getAttribute("data-tab");

  // タブのアクティブ状態を切り替え
  document.querySelectorAll(".l-tabs__link").forEach((tab) => {
    tab.classList.remove("is-active");
  });
  event.target.classList.add("is-active");

  // コンテンツの表示を切り替え
  document.querySelectorAll(".p-post").forEach((content) => {
    content.style.display = "none";
  });
  document.getElementById(targetId).style.display = "block";
}

// ページ読み込み時の処理
document.addEventListener("DOMContentLoaded", function () {
  // 保存された画像を復元
  const savedIcon = localStorage.getItem("selectedIcon");
  if (savedIcon) {
    const img = document.getElementById("profile-icon");
    img.src = savedIcon;
  }

  // タブクリックイベントの設定
  document.querySelectorAll(".l-tabs__link").forEach((tab) => {
    tab.addEventListener("click", handleTabClick);
  });

  // 初期表示のタブを設定
  const initialTab = document.querySelector(".l-tabs__link.is-active");
  if (initialTab) {
    const targetId = initialTab.getAttribute("data-tab");
    document.getElementById(targetId).style.display = "block";
  }
});
