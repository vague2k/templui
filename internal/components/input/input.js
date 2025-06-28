(function () {
  function initPasswordToggle(button) {
    if (button.hasAttribute("data-initialized")) {
      return;
    }

    button.setAttribute("data-initialized", "true");

    button.addEventListener("click", function () {
      const inputId = button.getAttribute("data-toggle-password");
      const input = document.getElementById(inputId);
      if (input) {
        const iconOpen = button.querySelector(".icon-open");
        const iconClosed = button.querySelector(".icon-closed");

        if (input.type === "password") {
          input.type = "text";
          iconOpen.classList.add("hidden");
          iconClosed.classList.remove("hidden");
        } else {
          input.type = "password";
          iconOpen.classList.remove("hidden");
          iconClosed.classList.add("hidden");
        }
      }
    });
  }

  function init(root = document) {
    const buttons = root.querySelectorAll(
      "[data-toggle-password]:not([data-initialized])"
    );
    buttons.forEach((button) => {
      initPasswordToggle(button);
    });
  }

  window.templUI = window.templUI || {};
  window.templUI.input = { init: init };

  document.addEventListener("DOMContentLoaded", () => init());
})();
