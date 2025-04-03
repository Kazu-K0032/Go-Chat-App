function togglePasswordForm() {
  const form = document.querySelector(".l-settings__passwordForm");
  if (form) {
    form.style.display = form.style.display === "none" ? "block" : "none";
  } else {
    // フォームが存在しない場合は、サーバーにリクエストを送信
    window.location.href = "/settings?show_password_form=true";
  }
}
