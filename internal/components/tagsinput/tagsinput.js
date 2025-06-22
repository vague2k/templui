(function () {
    // IIFE
    function initializeTagsInput(container) {
        const textInput = container.querySelector('.tags-input-text');
        const hiddenInputsContainer = container.querySelector('.tags-input-hidden-container');
        const tagsContainer = container.querySelector('.tags-input-tags-container');
        const name = container.dataset.name;
        const disabled = textInput.hasAttribute('disabled');

        function addTag(value) {
            if (disabled) return;
            const tagValue = value.trim();
            if (!tagValue) return;

            const existingTags = hiddenInputsContainer.querySelectorAll('input[type="hidden"]');
            for (const t of existingTags) {
                if (t.value.toLowerCase() === tagValue.toLowerCase()) {
                    textInput.value = '';
                    return;
                }
            }

            const tagChip = document.createElement('div');
            tagChip.className = "tags-input-chip flex items-center bg-primary text-xs border-transparent gap-1 px-1 " +
                "rounded-md border py-0.5 items-center focus:ring-ring text-primary-foreground font-semibold";
            tagChip.innerHTML = `
                <span>${tagValue}</span>
                <button type="button" class="tags-input-remove hover:text-destructive focus:outline-none"` + (disabled ? "disabled" : "") + `>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            `;

            const hiddenInput = document.createElement('input');
            hiddenInput.type = "hidden";
            hiddenInput.name = name;
            hiddenInput.value = tagValue;

            tagsContainer.appendChild(tagChip);
            hiddenInputsContainer.appendChild(hiddenInput);

            textInput.value = '';
        }

        function removeTag(removeButton) {
            if (disabled) return;
            const tagChip = removeButton.closest('.tags-input-chip');
            if (!tagChip) return;

            const tagValue = tagChip.querySelector('span').textContent.trim();

            const hiddenInput = hiddenInputsContainer.querySelector(`input[type="hidden"][value="${tagValue}"]`);
            if (hiddenInput) {
                hiddenInput.remove();
            }

            tagChip.remove();
        }

        function handleKeyDown(event) {
            if (event.key === 'Enter' || event.key === ',') {
                event.preventDefault();
                addTag(textInput.value);
            } else if (event.key === 'Backspace' && textInput.value === '') {
                event.preventDefault();

                const lastChip = tagsContainer.querySelector('.tags-input-chip:last-child');
                if (lastChip) {
                    const removeButton = lastChip.querySelector('.tags-input-remove');
                    if (removeButton) {
                        removeTag(removeButton);
                    }
                }
            }
        }

        function handleClick(event) {
            if (event.target.closest('.tags-input-remove')) {
                removeTag(event.target.closest('.tags-input-remove'));
            } else {
                textInput.focus();
            }
        }

        // --- Event Listeners ---
        textInput.removeEventListener('keydown', handleKeyDown);
        textInput.addEventListener('keydown', handleKeyDown);
        container.removeEventListener('click', handleClick);
        container.addEventListener('click', handleClick);
    }

    function initAllComponents(root = document) {
        const allTagsInputs = document.querySelectorAll('.tags-input');
        allTagsInputs.forEach(initializeTagsInput);
    }

    const handleHtmxSwap = (event) => {
        const target = event.detail.target || event.detail.elt;
        if (target instanceof Element) {
            requestAnimationFrame(() => initAllComponents(target));
        }
    };

    document.addEventListener("DOMContentLoaded", () => initAllComponents());
    document.body.addEventListener("htmx:afterSwap", handleHtmxSwap);
    document.body.addEventListener("htmx:oobAfterSwap", handleHtmxSwap);
})(); // End of IIFE