(function() {
    'use strict';

    // Theme color definitions in OKLCH format
    const themeColors = {
        default: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.145 0 0)",
                muted: "oklch(0.975 0 0)",
                "muted-foreground": "oklch(0.556 0 0)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.145 0 0)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.145 0 0)",
                border: "oklch(0.922 0 0)",
                input: "oklch(0.922 0 0)",
                primary: "oklch(0.205 0 0)",
                "primary-foreground": "oklch(0.985 0 0)",
                secondary: "oklch(0.975 0 0)",
                "secondary-foreground": "oklch(0.205 0 0)",
                accent: "oklch(0.975 0 0)",
                "accent-foreground": "oklch(0.205 0 0)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.985 0 0)",
                ring: "oklch(0.708 0 0)",
            },
            dark: {
                background: "oklch(0.145 0 0)",
                foreground: "oklch(0.985 0 0)",
                muted: "oklch(0.231 0 0)",
                "muted-foreground": "oklch(0.687 0 0)",
                popover: "oklch(0.145 0 0)",
                "popover-foreground": "oklch(0.985 0 0)",
                card: "oklch(0.145 0 0)",
                "card-foreground": "oklch(0.985 0 0)",
                border: "oklch(0.231 0 0)",
                input: "oklch(0.231 0 0)",
                primary: "oklch(0.985 0 0)",
                "primary-foreground": "oklch(0.205 0 0)",
                secondary: "oklch(0.231 0 0)",
                "secondary-foreground": "oklch(0.985 0 0)",
                accent: "oklch(0.231 0 0)",
                "accent-foreground": "oklch(0.985 0 0)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.985 0 0)",
                ring: "oklch(0.708 0 0)",
            },
        },
        red: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.145 0 0)",
                muted: "oklch(0.975 0 0)",
                "muted-foreground": "oklch(0.556 0 0)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.145 0 0)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.145 0 0)",
                border: "oklch(0.922 0 0)",
                input: "oklch(0.922 0 0)",
                primary: "oklch(0.6 0.24 27)",
                "primary-foreground": "oklch(0.985 0 0)",
                secondary: "oklch(0.975 0 0)",
                "secondary-foreground": "oklch(0.145 0 0)",
                accent: "oklch(0.975 0 0)",
                "accent-foreground": "oklch(0.145 0 0)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.985 0 0)",
                ring: "oklch(0.6 0.24 27)",
            },
            dark: {
                background: "oklch(0.145 0 0)",
                foreground: "oklch(0.985 0 0)",
                muted: "oklch(0.231 0 0)",
                "muted-foreground": "oklch(0.687 0 0)",
                popover: "oklch(0.145 0 0)",
                "popover-foreground": "oklch(0.985 0 0)",
                card: "oklch(0.145 0 0)",
                "card-foreground": "oklch(0.985 0 0)",
                border: "oklch(0.231 0 0)",
                input: "oklch(0.231 0 0)",
                primary: "oklch(0.6 0.24 27)",
                "primary-foreground": "oklch(0.985 0 0)",
                secondary: "oklch(0.231 0 0)",
                "secondary-foreground": "oklch(0.985 0 0)",
                accent: "oklch(0.231 0 0)",
                "accent-foreground": "oklch(0.985 0 0)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.985 0 0)",
                ring: "oklch(0.6 0.24 27)",
            },
        },
        rose: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.205 0.006 293.334)",
                muted: "oklch(0.975 0.003 293.334)",
                "muted-foreground": "oklch(0.556 0.004 293.334)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.205 0.006 293.334)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.205 0.006 293.334)",
                border: "oklch(0.922 0.003 293.334)",
                input: "oklch(0.922 0.003 293.334)",
                primary: "oklch(0.614 0.215 10.227)",
                "primary-foreground": "oklch(0.985 0.01 10.227)",
                secondary: "oklch(0.975 0.003 293.334)",
                "secondary-foreground": "oklch(0.205 0.006 293.334)",
                accent: "oklch(0.975 0.003 293.334)",
                "accent-foreground": "oklch(0.205 0.006 293.334)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.985 0 0)",
                ring: "oklch(0.614 0.215 10.227)",
            },
            dark: {
                background: "oklch(0.169 0.005 28.668)",
                foreground: "oklch(0.962 0 0)",
                muted: "oklch(0.226 0 0)",
                "muted-foreground": "oklch(0.687 0.003 293.334)",
                popover: "oklch(0.197 0 0)",
                "popover-foreground": "oklch(0.962 0 0)",
                card: "oklch(0.205 0.007 38.753)",
                "card-foreground": "oklch(0.962 0 0)",
                border: "oklch(0.231 0.002 293.334)",
                input: "oklch(0.231 0.002 293.334)",
                primary: "oklch(0.614 0.215 10.227)",
                "primary-foreground": "oklch(0.985 0.01 10.227)",
                secondary: "oklch(0.231 0.002 293.334)",
                "secondary-foreground": "oklch(0.985 0 0)",
                accent: "oklch(0.226 0.003 28.668)",
                "accent-foreground": "oklch(0.985 0 0)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.984 0.009 10.227)",
                ring: "oklch(0.614 0.215 10.227)",
            },
        },
        orange: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.169 0.005 28.668)",
                muted: "oklch(0.977 0.003 81.762)",
                "muted-foreground": "oklch(0.531 0.009 37.67)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.169 0.005 28.668)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.169 0.005 28.668)",
                border: "oklch(0.922 0.006 28.668)",
                input: "oklch(0.922 0.006 28.668)",
                primary: "oklch(0.681 0.224 52.718)",
                "primary-foreground": "oklch(0.983 0.006 81.762)",
                secondary: "oklch(0.977 0.003 81.762)",
                "secondary-foreground": "oklch(0.205 0.007 38.753)",
                accent: "oklch(0.977 0.003 81.762)",
                "accent-foreground": "oklch(0.205 0.007 38.753)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.983 0.006 81.762)",
                ring: "oklch(0.681 0.224 52.718)",
            },
            dark: {
                background: "oklch(0.169 0.005 28.668)",
                foreground: "oklch(0.983 0.006 81.762)",
                muted: "oklch(0.226 0.003 28.668)",
                "muted-foreground": "oklch(0.669 0.005 37.67)",
                popover: "oklch(0.169 0.005 28.668)",
                "popover-foreground": "oklch(0.983 0.006 81.762)",
                card: "oklch(0.169 0.005 28.668)",
                "card-foreground": "oklch(0.983 0.006 81.762)",
                border: "oklch(0.226 0.003 28.668)",
                input: "oklch(0.226 0.003 28.668)",
                primary: "oklch(0.628 0.196 41.428)",
                "primary-foreground": "oklch(0.983 0.006 81.762)",
                secondary: "oklch(0.226 0.003 28.668)",
                "secondary-foreground": "oklch(0.983 0.006 81.762)",
                accent: "oklch(0.226 0.003 28.668)",
                "accent-foreground": "oklch(0.983 0.006 81.762)",
                destructive: "oklch(0.599 0.24 27.214)",
                "destructive-foreground": "oklch(0.983 0.006 81.762)",
                ring: "oklch(0.569 0.157 53.569)",
            },
        },
        green: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.205 0.006 293.334)",
                muted: "oklch(0.975 0.003 293.334)",
                "muted-foreground": "oklch(0.556 0.004 293.334)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.205 0.006 293.334)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.205 0.006 293.334)",
                border: "oklch(0.922 0.003 293.334)",
                input: "oklch(0.922 0.003 293.334)",
                primary: "oklch(0.523 0.179 156.329)",
                "primary-foreground": "oklch(0.985 0.01 10.227)",
                secondary: "oklch(0.975 0.003 293.334)",
                "secondary-foreground": "oklch(0.205 0.006 293.334)",
                accent: "oklch(0.975 0.003 293.334)",
                "accent-foreground": "oklch(0.205 0.006 293.334)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.985 0 0)",
                ring: "oklch(0.523 0.179 156.329)",
            },
            dark: {
                background: "oklch(0.169 0.005 28.668)",
                foreground: "oklch(0.962 0 0)",
                muted: "oklch(0.226 0 0)",
                "muted-foreground": "oklch(0.687 0.003 293.334)",
                popover: "oklch(0.197 0 0)",
                "popover-foreground": "oklch(0.962 0 0)",
                card: "oklch(0.205 0.007 38.753)",
                "card-foreground": "oklch(0.962 0 0)",
                border: "oklch(0.231 0.002 293.334)",
                input: "oklch(0.231 0.002 293.334)",
                primary: "oklch(0.598 0.165 154.756)",
                "primary-foreground": "oklch(0.205 0.081 154.756)",
                secondary: "oklch(0.231 0.002 293.334)",
                "secondary-foreground": "oklch(0.985 0 0)",
                accent: "oklch(0.226 0.003 28.668)",
                "accent-foreground": "oklch(0.985 0 0)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.984 0.009 10.227)",
                ring: "oklch(0.453 0.126 156.329)",
            },
        },
        blue: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.16 0.019 268.87)",
                muted: "oklch(0.976 0.01 231.683)",
                "muted-foreground": "oklch(0.558 0.031 258.723)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.16 0.019 268.87)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.16 0.019 268.87)",
                border: "oklch(0.926 0.016 251.803)",
                input: "oklch(0.926 0.016 251.803)",
                primary: "oklch(0.511 0.27 264.052)",
                "primary-foreground": "oklch(0.984 0.01 231.683)",
                secondary: "oklch(0.976 0.01 231.683)",
                "secondary-foreground": "oklch(0.213 0.019 265.755)",
                accent: "oklch(0.976 0.01 231.683)",
                "accent-foreground": "oklch(0.213 0.019 265.755)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.984 0.01 231.683)",
                ring: "oklch(0.511 0.27 264.052)",
            },
            dark: {
                background: "oklch(0.16 0.019 268.87)",
                foreground: "oklch(0.984 0.01 231.683)",
                muted: "oklch(0.241 0.016 263.317)",
                "muted-foreground": "oklch(0.697 0.017 249.183)",
                popover: "oklch(0.16 0.019 268.87)",
                "popover-foreground": "oklch(0.984 0.01 231.683)",
                card: "oklch(0.16 0.019 268.87)",
                "card-foreground": "oklch(0.984 0.01 231.683)",
                border: "oklch(0.241 0.016 263.317)",
                input: "oklch(0.241 0.016 263.317)",
                primary: "oklch(0.637 0.235 261.226)",
                "primary-foreground": "oklch(0.213 0.019 265.755)",
                secondary: "oklch(0.241 0.016 263.317)",
                "secondary-foreground": "oklch(0.984 0.01 231.683)",
                accent: "oklch(0.241 0.016 263.317)",
                "accent-foreground": "oklch(0.984 0.01 231.683)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.984 0.01 231.683)",
                ring: "oklch(0.53 0.238 269.161)",
            },
        },
        yellow: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.169 0.005 28.668)",
                muted: "oklch(0.977 0.003 81.762)",
                "muted-foreground": "oklch(0.531 0.009 37.67)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.169 0.005 28.668)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.169 0.005 28.668)",
                border: "oklch(0.922 0.006 28.668)",
                input: "oklch(0.922 0.006 28.668)",
                primary: "oklch(0.762 0.185 85.342)",
                "primary-foreground": "oklch(0.238 0.081 38.404)",
                secondary: "oklch(0.977 0.003 81.762)",
                "secondary-foreground": "oklch(0.205 0.007 38.753)",
                accent: "oklch(0.977 0.003 81.762)",
                "accent-foreground": "oklch(0.205 0.007 38.753)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.983 0.006 81.762)",
                ring: "oklch(0.169 0.005 28.668)",
            },
            dark: {
                background: "oklch(0.169 0.005 28.668)",
                foreground: "oklch(0.983 0.006 81.762)",
                muted: "oklch(0.226 0.003 28.668)",
                "muted-foreground": "oklch(0.669 0.005 37.67)",
                popover: "oklch(0.169 0.005 28.668)",
                "popover-foreground": "oklch(0.983 0.006 81.762)",
                card: "oklch(0.169 0.005 28.668)",
                "card-foreground": "oklch(0.983 0.006 81.762)",
                border: "oklch(0.226 0.003 28.668)",
                input: "oklch(0.226 0.003 28.668)",
                primary: "oklch(0.762 0.185 85.342)",
                "primary-foreground": "oklch(0.238 0.081 38.404)",
                secondary: "oklch(0.226 0.003 28.668)",
                "secondary-foreground": "oklch(0.983 0.006 81.762)",
                accent: "oklch(0.226 0.003 28.668)",
                "accent-foreground": "oklch(0.983 0.006 81.762)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.983 0.006 81.762)",
                ring: "oklch(0.485 0.16 54.004)",
            },
        },
        violet: {
            light: {
                background: "oklch(1 0 0)",
                foreground: "oklch(0.169 0.018 276.935)",
                muted: "oklch(0.975 0.005 276.935)",
                "muted-foreground": "oklch(0.551 0.012 276.935)",
                popover: "oklch(1 0 0)",
                "popover-foreground": "oklch(0.169 0.018 276.935)",
                card: "oklch(1 0 0)",
                "card-foreground": "oklch(0.169 0.018 276.935)",
                border: "oklch(0.924 0.005 276.935)",
                input: "oklch(0.924 0.005 276.935)",
                primary: "oklch(0.614 0.253 293.639)",
                "primary-foreground": "oklch(0.984 0.01 231.683)",
                secondary: "oklch(0.975 0.005 276.935)",
                "secondary-foreground": "oklch(0.213 0.017 276.935)",
                accent: "oklch(0.975 0.005 276.935)",
                "accent-foreground": "oklch(0.213 0.017 276.935)",
                destructive: "oklch(0.577 0.237 27.325)",
                "destructive-foreground": "oklch(0.984 0.01 231.683)",
                ring: "oklch(0.614 0.253 293.639)",
            },
            dark: {
                background: "oklch(0.169 0.018 276.935)",
                foreground: "oklch(0.984 0.01 231.683)",
                muted: "oklch(0.238 0.012 265.755)",
                "muted-foreground": "oklch(0.692 0.008 265.755)",
                popover: "oklch(0.169 0.018 276.935)",
                "popover-foreground": "oklch(0.984 0.01 231.683)",
                card: "oklch(0.169 0.018 276.935)",
                "card-foreground": "oklch(0.984 0.01 231.683)",
                border: "oklch(0.238 0.012 265.755)",
                input: "oklch(0.238 0.012 265.755)",
                primary: "oklch(0.571 0.221 294.416)",
                "primary-foreground": "oklch(0.984 0.01 231.683)",
                secondary: "oklch(0.238 0.012 265.755)",
                "secondary-foreground": "oklch(0.984 0.01 231.683)",
                accent: "oklch(0.238 0.012 265.755)",
                "accent-foreground": "oklch(0.984 0.01 231.683)",
                destructive: "oklch(0.433 0.182 27.325)",
                "destructive-foreground": "oklch(0.984 0.01 231.683)",
                ring: "oklch(0.571 0.221 294.416)",
            },
        },
    };

    // Primary colors for theme indicators
    const themePrimaryColors = {
        default: 'oklch(0.205 0 0)',
        red: 'oklch(0.6 0.24 27)',
        rose: 'oklch(0.614 0.215 10.227)',
        orange: 'oklch(0.681 0.224 52.718)',
        green: 'oklch(0.523 0.179 156.329)',
        blue: 'oklch(0.511 0.27 264.052)',
        yellow: 'oklch(0.762 0.185 85.342)',
        violet: 'oklch(0.614 0.253 293.639)'
    };

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
            
            // Remove existing classes
            button.classList.remove('border-foreground/50', 'border-transparent');
            
            // Add appropriate classes
            if (isActive) {
                button.classList.add('border-foreground/50');
            } else {
                button.classList.add('border-transparent');
            }
        });
    }

    // Generate CSS
    function generateCSS() {
        const colors = themeColors[currentTheme];
        if (!colors) return;
        
        const isDark = document.documentElement.classList.contains('dark');
        const mode = isDark ? 'dark' : 'light';
        
        // Generate CSS variables
        const cssVars = Object.entries(colors.light)
            .map(([key, value]) => `    --${key}: ${value};`)
            .join('\n');
        
        const darkCssVars = Object.entries(colors.dark)
            .map(([key, value]) => `    --${key}: ${value};`)
            .join('\n');
        
        generatedCSS = `:root {
${cssVars}
    --radius: 0.5rem;
}
.dark {
${darkCssVars}
    --radius: 0.5rem;
}`;
        
        // Update modal content
        updateModal();
        
        // Show modal (if using a modal system)
        const modal = document.getElementById('css-modal');
        if (modal && window.Modal) {
            window.Modal.show('css-modal');
        }
    }

    // Update modal with generated CSS and theme indicator
    function updateModal() {
        // Update CSS display
        const cssDisplay = document.querySelector('[data-css-display]');
        if (cssDisplay) {
            cssDisplay.textContent = generatedCSS;
        }
        
        // Update theme indicator
        const themeIndicator = document.querySelector('[data-theme-indicator]');
        if (themeIndicator) {
            themeIndicator.style.backgroundColor = themePrimaryColors[currentTheme];
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