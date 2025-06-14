(function () {
  function updateProgressWidth(progressBar) {
    if (!progressBar) return;

    const indicator = progressBar.querySelector("[data-progress-indicator]");
    if (!indicator) return;

    const value = parseFloat(progressBar.getAttribute("aria-valuenow") || "0");
    let max = parseFloat(progressBar.getAttribute("aria-valuemax") || "100");
    if (max <= 0) max = 100;

    let percentage = 0;
    if (max > 0) {
      percentage = (Math.max(0, Math.min(value, max)) / max) * 100;
    }

    indicator.style.width = percentage + "%";
  }

  function initAllComponents(root = document) {
    if (root instanceof Element && root.matches('[role="progressbar"]')) {
      updateProgressWidth(root);
    }
    if (root && typeof root.querySelectorAll === "function") {
      for (const progressBar of root.querySelectorAll('[role="progressbar"]')) {
        updateProgressWidth(progressBar);
      }
    }
  }

  const handleHtmxSwap = (event) => {
    let target;
    if (event.type === "htmx:afterSwap") {
      target = event.detail.elt;
    }
    if (event.type === "htmx:oobAfterSwap") {
      target = event.detail.target;
    }
    if (target instanceof Element) {
      requestAnimationFrame(() => initAllComponents(target));
    }
  };

  document.addEventListener("DOMContentLoaded", () => initAllComponents());
  document.body.addEventListener("htmx:afterSwap", handleHtmxSwap);
  document.body.addEventListener("htmx:oobAfterSwap", handleHtmxSwap);
})();
