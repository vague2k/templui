(function () {
  function handleDropdownItemClick(event) {
    const item = event.currentTarget;

    // Check if this item should prevent dropdown from closing
    if (item.dataset.preventClose === "true") {
      return; // Don't close the dropdown
    }

    const popoverContent = item.closest("[data-popover-id]");
    if (popoverContent) {
      const popoverId = popoverContent.dataset.popoverId;
      if (window.closePopover) {
        window.closePopover(popoverId, true);
      } else {
        console.warn("popover.Script's closePopover function not found.");
        document.body.click(); // Fallback
      }
    }
  }

  function init(root = document) {
    // Select items with 'data-dropdown-item' but not 'data-dropdown-submenu-trigger'
    const items = root.querySelectorAll(
      "[data-dropdown-item]:not([data-dropdown-submenu-trigger])"
    );
    items.forEach((item) => {
      item.removeEventListener("click", handleDropdownItemClick);
      item.addEventListener("click", handleDropdownItemClick);
    });
  }

  window.templUI = window.templUI || {};
  window.templUI.dropdown = { init: init };

  document.addEventListener("DOMContentLoaded", () => init());
})();
