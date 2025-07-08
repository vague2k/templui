(function() {
    'use strict';

    // State
    let currentTheme = 'default';
    let generatedCSS = '';

    // Initialize theme system
    function init() {
        // Set up event listeners
        document.addEventListener('click', handleClick);
        
        // Apply initial theme
        applyTheme();
        updateThemeButtons();
    }

    // Handle all click events
    function handleClick(event) {
        const target = event.target;
        
        // Theme selection
        if (target.hasAttribute('data-theme')) {
            event.preventDefault();
            setTheme(target.getAttribute('data-theme'));
            return;
        }
        
        // Check parent elements for data-theme
        const themeButton = target.closest('[data-theme]');
        if (themeButton) {
            event.preventDefault();
            setTheme(themeButton.getAttribute('data-theme'));
            return;
        }
        
        // Handle action buttons
        const actionButton = target.closest('[data-action]');
        if (actionButton) {
            const action = actionButton.getAttribute('data-action');
            
            switch (action) {
                case 'reset-theme':
                    event.preventDefault();
                    setTheme('default');
                    break;
                case 'generate-css':
                    event.preventDefault();
                    generateCSS();
                    break;
                case 'copy-css':
                    event.preventDefault();
                    copyToClipboard();
                    break;
            }
        }
    }

    // Set theme
    function setTheme(theme) {
        currentTheme = theme;
        applyTheme();
        updateThemeButtons();
    }

    // Apply theme to preview
    function applyTheme() {
        const preview = document.getElementById('theme-preview');
        if (!preview) return;
        
        // Remove all theme classes
        preview.className = preview.className.replace(/\btheme-\S+/g, '');
        
        // Add base classes
        preview.className = 'grid gap-4 md:grid-cols-2 lg:grid-cols-10 xl:grid-cols-11';
        
        // Add theme class if not default
        if (currentTheme !== 'default') {
            preview.classList.add(`theme-${currentTheme}`);
        }
    }

    // Update theme button states
    function updateThemeButtons() {
        const buttons = document.querySelectorAll('[data-theme]');
        
        buttons.forEach(button => {
            const isActive = button.getAttribute('data-theme') === currentTheme;
            const dot = button.querySelector('span > span');
            
            if (dot) {
                if (isActive) {
                    dot.setAttribute('data-selected', '');
                } else {
                    dot.removeAttribute('data-selected');
                }
            }
        });
    }

    // Generate CSS by extracting actual computed styles
    function generateCSS() {
        const preview = document.getElementById('theme-preview');
        if (!preview) return;
        
        // Store current dark mode state
        const isDarkMode = document.documentElement.classList.contains('dark');
        
        // Create a proper container that inherits from documentElement
        const tempWrapper = document.createElement('div');
        tempWrapper.style.display = 'none';
        
        // Create inner container for theme
        const tempContainer = document.createElement('div');
        if (currentTheme !== 'default') {
            tempContainer.className = `theme-${currentTheme}`;
        }
        
        tempWrapper.appendChild(tempContainer);
        document.body.appendChild(tempWrapper);
        
        // Extract light mode styles
        document.documentElement.classList.remove('dark');
        const lightStyles = getComputedStyle(tempContainer);
        const lightVars = extractCSSVariables(lightStyles);
        
        // Extract dark mode styles
        document.documentElement.classList.add('dark');
        const darkStyles = getComputedStyle(tempContainer);
        const darkVars = extractCSSVariables(darkStyles);
        
        // Restore original dark mode state
        if (!isDarkMode) {
            document.documentElement.classList.remove('dark');
        }
        
        // Clean up
        document.body.removeChild(tempWrapper);
        
        // Generate CSS string
        const lightCSS = Object.entries(lightVars)
            .map(([key, value]) => `    ${key}: ${value};`)
            .join('\n');
        
        const darkCSS = Object.entries(darkVars)
            .map(([key, value]) => `    ${key}: ${value};`)
            .join('\n');
        
        generatedCSS = `:root {
${lightCSS}
}

.dark {
${darkCSS}
}`;
        
        // Update modal content
        updateModal();
        
        // Show modal (if using a modal system)
        const modal = document.getElementById('css-modal');
        if (modal && window.Modal) {
            window.Modal.show('css-modal');
        }
    }
    
    // Extract CSS variables from computed styles
    function extractCSSVariables(computedStyle) {
        const vars = {};
        const relevantVars = [
            '--background', '--foreground', '--card', '--card-foreground',
            '--popover', '--popover-foreground', '--primary', '--primary-foreground',
            '--secondary', '--secondary-foreground', '--muted', '--muted-foreground',
            '--accent', '--accent-foreground', '--destructive', '--destructive-foreground',
            '--border', '--input', '--ring', '--radius',
            '--chart-1', '--chart-2', '--chart-3', '--chart-4', '--chart-5',
            '--sidebar', '--sidebar-foreground', '--sidebar-primary', '--sidebar-primary-foreground',
            '--sidebar-accent', '--sidebar-accent-foreground', '--sidebar-border', '--sidebar-ring',
            '--surface', '--surface-foreground', '--code', '--code-foreground',
            '--code-highlight', '--code-number', '--selection', '--selection-foreground'
        ];
        
        relevantVars.forEach(varName => {
            const value = computedStyle.getPropertyValue(varName).trim();
            if (value) {
                vars[varName] = value;
            }
        });
        
        return vars;
    }

    // Update modal with generated CSS and theme indicator
    function updateModal() {
        // Update CSS display
        const cssDisplay = document.querySelector('[data-css-display]');
        if (cssDisplay) {
            cssDisplay.textContent = generatedCSS;
        }
        
        // Update theme indicator by getting the actual primary color
        const themeIndicator = document.querySelector('[data-theme-indicator]');
        if (themeIndicator) {
            // Create a temporary element to get the actual primary color
            const temp = document.createElement('div');
            temp.className = currentTheme !== 'default' ? `theme-${currentTheme}` : '';
            temp.style.display = 'none';
            document.body.appendChild(temp);
            
            const primaryColor = getComputedStyle(temp).getPropertyValue('--primary').trim();
            themeIndicator.style.backgroundColor = primaryColor;
            
            document.body.removeChild(temp);
        }
        
        const themeText = document.querySelector('[data-theme-text]');
        if (themeText) {
            themeText.textContent = currentTheme.charAt(0).toUpperCase() + currentTheme.slice(1);
        }
    }

    // Copy to clipboard
    function copyToClipboard() {
        if (!generatedCSS) return;
        
        navigator.clipboard.writeText(generatedCSS)
            .then(() => {
                alert('Copied to clipboard!');
            })
            .catch(err => {
                console.error('Failed to copy:', err);
            });
    }

    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();