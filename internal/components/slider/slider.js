(function () {
  function initSlider(sliderInput) {
    if (sliderInput.hasAttribute("data-initialized")) return;

    sliderInput.setAttribute("data-initialized", "true");

    const sliderId = sliderInput.id;
    if (!sliderId) return;

    const valueElements = document.querySelectorAll(
      `[data-slider-value][data-slider-value-for="${sliderId}"]`
    );

    function updateValues() {
      valueElements.forEach((el) => {
        el.textContent = sliderInput.value;
      });
    }

    updateValues();
    sliderInput.addEventListener("input", updateValues);
  }

  function initAllComponents(root = document) {
    if (
      root instanceof Element &&
      root.matches('input[type="range"][data-slider-input]')
    ) {
      initSlider(root);
    }
    for (const slider of root.querySelectorAll(
      'input[type="range"][data-slider-input]:not([data-initialized])'
    )) {
      initSlider(slider);
    }
  }

  if (!window.templUI) {
    window.templUI = {};
  }

  window.templUI.slider = {
    initAllComponents: initAllComponents,
  };

  document.addEventListener("DOMContentLoaded", () => initAllComponents());
})();
