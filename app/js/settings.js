function togglePasswordForm() {
  const form = document.querySelector(".l-settings__passwordForm");
  if (form) {
    // フォームの表示状態を切り替え
    form.classList.toggle("is-active");

    // フォームが非表示になった場合、入力をクリア
    if (!form.classList.contains("is-active")) {
      const inputs = form.querySelectorAll("input[type='password']");
      inputs.forEach((input) => (input.value = ""));
    }
  } else {
    // フォームが存在しない場合は、サーバーにリクエストを送信
    window.location.href = "/settings?show_password_form=true";
  }
}
