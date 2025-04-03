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
  document.querySelectorAll(".l-section").forEach((content) => {
    content.classList.remove("is-active");
  });
  document.getElementById(targetId).classList.add("is-active");

  // URLのハッシュを更新
  window.location.hash = targetId;
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

  // URLのハッシュから初期タブを設定
  const hash = window.location.hash.slice(1) || "posts";
  const initialTab = document.querySelector(
    `.l-tabs__link[data-tab="${hash}"]`
  );
  if (initialTab) {
    initialTab.classList.add("is-active");
    document.getElementById(hash).classList.add("is-active");
  }
});
