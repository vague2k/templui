(function () {
  function initLabel(label) {
    if (!label || label.hasAttribute("data-initialized")) return;
    label.setAttribute("data-initialized", "true");

    if (
      !label.hasAttribute("for") ||
      !label.hasAttribute("data-disabled-style")
    ) {
      return;
    }

    const forId = label.getAttribute("for");
    const targetElement = forId ? document.getElementById(forId) : null;
    const disabledStyle = label.getAttribute("data-disabled-style");

    if (!disabledStyle) return;

    const classes = disabledStyle.split(" ").filter(Boolean);

    function updateStyle() {
      if (targetElement && targetElement.disabled) {
        label.classList.add(...classes);
      } else {
        label.classList.remove(...classes);
      }
    }

    // Set up mutation observer to detect disabled state changes
    if (targetElement) {
      const observer = new MutationObserver((mutations) => {
        for (const mutation of mutations) {
          if (
            mutation.type === "attributes" &&
            mutation.attributeName === "disabled"
          ) {
            updateStyle();
          }
        }
      });

      observer.observe(targetElement, {
        attributes: true,
        attributeFilter: ["disabled"],
      });
    }

    // Initial style update
    updateStyle();
  }

  function init(root = document) {
    if (
      root instanceof Element &&
      root.matches("label[for][data-disabled-style]")
    ) {
      initLabel(root);
    }
    for (const label of root.querySelectorAll(
      "label[for][data-disabled-style]:not([data-initialized])"
    )) {
      initLabel(label);
    }
  }

  document.addEventListener("DOMContentLoaded", () => init());
})();
