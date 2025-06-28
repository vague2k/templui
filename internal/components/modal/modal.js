if (typeof window.modalState === "undefined") {
  window.modalState = {
    openModalId: null,
  };
}

(function () {
  function closeModal(modal, immediate = false) {
    if (!modal || modal.style.display === "none") return;

    const content = modal.querySelector("[data-modal-content]");
    const modalId = modal.id;

    // Update data-open attributes
    modal.setAttribute("data-open", "false");
    updateTriggersState(modalId, false);

    // Apply leaving transitions
    modal.classList.remove("opacity-100");
    modal.classList.add("opacity-0");

    if (content) {
      content.classList.remove("scale-100", "opacity-100");
      content.classList.add("scale-95", "opacity-0");
    }

    function hideModal() {
      modal.style.display = "none";

      if (window.modalState.openModalId === modalId) {
        window.modalState.openModalId = null;
        document.body.style.overflow = "";
      }
    }

    if (immediate) {
      hideModal();
    } else {
      setTimeout(hideModal, 300);
    }
  }

  // Helper function to update trigger states
  function updateTriggersState(modalId, isOpen) {
    const triggers = document.querySelectorAll(
      `[data-modal-trigger="${modalId}"]`
    );
    triggers.forEach((trigger) => {
      trigger.setAttribute("data-open", isOpen.toString());
    });
  }

  function openModal(modal) {
    if (!modal) return;

    // Close any open modal first
    if (window.modalState.openModalId) {
      const openModal = document.getElementById(window.modalState.openModalId);
      if (openModal && openModal !== modal) {
        closeModal(openModal, true);
      }
    }

    const content = modal.querySelector("[data-modal-content]");

    // Display and prepare for animation
    modal.style.display = "flex";

    // Update data-open attributes
    modal.setAttribute("data-open", "true");
    updateTriggersState(modal.id, true);

    // Store as currently open modal
    window.modalState.openModalId = modal.id;
    document.body.style.overflow = "hidden";

    // Force reflow before adding transition classes
    void modal.offsetHeight;

    // Start animations
    modal.classList.remove("opacity-0");
    modal.classList.add("opacity-100");

    if (content) {
      content.classList.remove("scale-95", "opacity-0");
      content.classList.add("scale-100", "opacity-100");

      // Focus first focusable element
      setTimeout(() => {
        const focusable = content.querySelector(
          'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
        );
        if (focusable) focusable.focus();
      }, 50);
    }
  }

  function closeModalById(modalId, immediate = false) {
    const modal = document.getElementById(modalId);
    if (modal) closeModal(modal, immediate);
  }

  function openModalById(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) openModal(modal);
  }

  function handleClickAway(e) {
    const openModalId = window.modalState.openModalId;
    if (!openModalId) return;

    const modal = document.getElementById(openModalId);
    if (!modal || modal.getAttribute("data-disable-click-away") === "true")
      return;

    const content = modal.querySelector("[data-modal-content]");
    const trigger = e.target.closest("[data-modal-trigger]");

    if (
      content &&
      !content.contains(e.target) &&
      (!trigger || trigger.getAttribute("data-modal-trigger") !== openModalId)
    ) {
      closeModal(modal);
    }
  }

  function handleEscKey(e) {
    if (e.key !== "Escape" || !window.modalState.openModalId) return;

    const modal = document.getElementById(window.modalState.openModalId);
    if (modal && modal.getAttribute("data-disable-esc") !== "true") {
      closeModal(modal);
    }
  }

  function initTrigger(trigger) {
    const modalId = trigger.getAttribute("data-modal-trigger");
    if (!modalId) return;

    trigger.addEventListener("click", () => {
      if (
        !trigger.hasAttribute("disabled") &&
        !trigger.classList.contains("opacity-50")
      ) {
        openModalById(modalId);
      }
    });
  }

  function initCloseButton(closeBtn) {
    closeBtn.addEventListener("click", () => {
      const modalId = closeBtn.getAttribute("data-modal-close");
      if (modalId) {
        closeModalById(modalId);
      } else {
        const modal = closeBtn.closest("[data-modal]");
        if (modal && modal.id) {
          closeModal(modal);
        }
      }
    });
  }

  function initAllComponents(root = document) {
    // Initialize modal triggers
    if (root instanceof Element && root.matches("[data-modal-trigger]")) {
      initTrigger(root);
    }
    for (const trigger of root.querySelectorAll("[data-modal-trigger]")) {
      initTrigger(trigger);
    }

    // Initialize modal close buttons
    if (root instanceof Element && root.matches("[data-modal-close]")) {
      initCloseButton(root);
    }
    for (const closeBtn of root.querySelectorAll("[data-modal-close]")) {
      initCloseButton(closeBtn);
    }

    // Initialize modals and their data-open states
    if (root instanceof Element && root.matches("[data-modal]")) {
      initModal(root);
    }
    for (const modal of root.querySelectorAll("[data-modal]")) {
      initModal(modal);
    }
  }

  function initModal(modal) {
    // Initialize trigger states for this modal
    if (modal.id) {
      updateTriggersState(modal.id, false); // Start with closed state
    }

    // Check for modals that should be initially open
    if (modal.getAttribute("data-initial-open") === "true") {
      openModal(modal);
    }
  }

  if (!window.templUI) {
    window.templUI = {};
  }

  window.templUI.modal = {
    initAllComponents: initAllComponents,
  };

  if (typeof window.modalEventsInitialized === "undefined") {
    document.addEventListener("click", handleClickAway);
    document.addEventListener("keydown", handleEscKey);
    window.modalEventsInitialized = true;
  }

  document.addEventListener("DOMContentLoaded", () => initAllComponents());
})();
