# templUI - Developer Guide

## Overview
templUI is an enterprise-ready UI component library built specifically for Go's `templ` templating library. It combines Go templating with lightweight Vanilla JavaScript and Tailwind CSS to create modern, performant web applications with server-side rendering.

## Key Technologies
- **Go 1.24+** - Backend language and server
- **templ** - Type-safe HTML templating for Go  
- **Tailwind CSS v4.1.3** - Utility-first CSS framework with custom theme system
- **Vanilla JavaScript** - Lightweight client-side interactivity (no framework dependencies)
- **ESBuild** - JavaScript bundling and minification
- **Air** - Live reload for Go development
- **Datastar** - Server-side rendering library for reactive interfaces

## Development Commands

### Primary Development Workflow
```bash
# Start full development environment (recommended)
make dev
# Runs: tailwind-clean, tailwind-watch, templ, server, shiki-highlighter, minify-js-dev, minify-js-components

# Individual development services
make templ           # Watch and compile .templ files
make server          # Start Go server with live reload using Air
make tailwind-watch  # Watch and compile Tailwind CSS
make shiki-highlighter # Start syntax highlighting service
```

### Build & Production
```bash
# Clean and rebuild Tailwind CSS
make tailwind-clean

# Generate static HTML showcases
make build-html      # Compile templ + render showcases

# Validate HTML output
make validate-html   # Build + lint HTML with ESLint

# Generate project assets
make generate-sitemap
make generate-icons
```

### JavaScript & Asset Management
```bash
# Watch and minify core JavaScript bundle
make minify-js-dev

# Watch and minify individual component JavaScript files
make minify-js-components

# Manual linting
npx eslint --fix "out/**/*.html" --ext .html
```

### Docker Development
```bash
# Full stack with Docker Compose
docker-compose up --build

# Individual Docker build
docker build -t templui .
```

## Architecture Overview

### Directory Structure
```
templui/
├── cmd/                    # Command-line applications
│   ├── docs/              # Main documentation server
│   ├── icongen/           # Icon generation utility
│   ├── render-showcases/  # Static HTML generation
│   └── sitemap/          # Sitemap generation
├── internal/              # Private application code
│   ├── components/        # UI component library
│   ├── config/           # Application configuration
│   ├── middleware/       # HTTP middleware
│   ├── ui/              # Page templates and layouts
│   └── utils/           # Shared utilities
├── assets/               # Static assets (CSS, JS, images)
├── static/              # Static files (robots.txt, sitemap.xml)
├── out/                 # Generated HTML showcases
└── shiki/              # Syntax highlighting service
```

### Component Architecture

**Component Structure Pattern:**
Each component follows a consistent structure:
```
internal/components/[component]/
├── [component].templ     # templ template with Go logic
├── [component]_templ.go  # Generated Go code (auto-generated)
├── [component].js        # Vanilla JavaScript (optional)
└── [component].min.js    # Minified JavaScript (auto-generated)
```

**templ Component Pattern:**
```go
// Component props with typed variants
type Props struct {
    ID         string
    Class      string
    Attributes templ.Attributes
    Variant    Variant // Typed enum for styling variants
    Size       Size    // Typed enum for sizes
    // ... component-specific props
}

// Main component template
templ ComponentName(props ...Props) {
    {{ var p Props }}
    if len(props) > 0 {
        {{ p = props[0] }}
    }
    // Component HTML with TwMerge for class management
    <div class={ utils.TwMerge("base-classes", p.variantClasses(), p.Class) }>
        { children... }
    </div>
}

// Variant method for styling logic
func (p Props) variantClasses() string {
    switch p.Variant {
    case VariantPrimary:
        return "bg-primary text-primary-foreground"
    default:
        return "bg-secondary text-secondary-foreground"
    }
}
```

### JavaScript Architecture

**Component JavaScript Pattern:**
- Self-contained IIFE (Immediately Invoked Function Expression)
- DOM-based initialization with `data-*` attributes
- Event delegation for dynamic content
- No external dependencies (pure Vanilla JS)

**Example Pattern:**
```javascript
(function() {
    // Component state management
    const instances = new Map();
    
    // Initialize component
    function createComponent(element) {
        if (element.hasAttribute("data-initialized")) return;
        element.setAttribute("data-initialized", "true");
        
        // Component logic here
    }
    
    // Auto-initialization
    document.addEventListener("DOMContentLoaded", () => {
        document.querySelectorAll("[data-component]").forEach(createComponent);
    });
})();
```

### CSS Architecture

**Tailwind CSS v4 Configuration:**
- Custom CSS variables for theming
- Built-in dark mode support
- Component-specific utility classes
- TwMerge integration for class conflict resolution

**Theme System:**
```css
:root {
  --background: hsl(0 0% 100%);
  --foreground: hsl(240 10% 3.9%);
  --primary: hsl(240 5.9% 10%);
  /* ... other theme variables */
}

.dark {
  --background: hsl(240 10% 3.9%);
  --foreground: hsl(0 0% 98%);
  /* ... dark mode overrides */
}
```

## Development Patterns & Conventions

### Component Development
1. **Props Structure**: Always use optional props pattern with defaults
2. **Class Management**: Use `utils.TwMerge()` for combining Tailwind classes
3. **Accessibility**: Include proper ARIA attributes and semantic HTML
4. **TypeScript-like Safety**: Leverage Go's type system for component props

### JavaScript Patterns
1. **Self-Contained**: Each component JS file is independent
2. **Data Attributes**: Use `data-*` attributes for configuration
3. **Progressive Enhancement**: Components work without JavaScript
4. **Event Delegation**: Handle dynamic content properly

### Styling Conventions
1. **Utility-First**: Use Tailwind utilities over custom CSS
2. **Component Variants**: Define styling variants as Go enums
3. **Responsive Design**: Mobile-first responsive patterns
4. **Theme Compliance**: Use CSS custom properties for theming

### File Organization
1. **Single Responsibility**: One component per directory
2. **Generated Files**: Never edit `*_templ.go` files (auto-generated)
3. **Consistent Naming**: Use kebab-case for directories, PascalCase for Go types

## Build System

### Asset Pipeline
1. **templ Generate**: Compiles `.templ` files to Go code
2. **Tailwind CSS**: Processes custom CSS with utilities
3. **ESBuild**: Bundles and minifies JavaScript
4. **Air**: Provides live reload for Go development

### Production Build
- Multi-stage Docker build
- Static binary compilation
- Asset minification and compression
- Environment-specific configuration

## Testing & Quality

### Linting & Validation
- **ESLint**: HTML and JavaScript linting
- **html-tailwind**: Tailwind CSS class validation
- **Go fmt**: Standard Go formatting
- **templ fmt**: Template formatting

### HTML Generation
- Showcase components rendered to static HTML
- Automated HTML validation pipeline
- Component integration testing through rendered output

## Environment Configuration

### Development
- `GO_ENV`: Set to "development" for local development
- Live reload enabled for templates and assets
- Detailed error messages and debugging

### Production
- `GO_ENV`: Set to "production" for optimized builds
- Asset bundling and minification
- Static file serving from embedded filesystem

## Dependencies

### Core Dependencies
- `github.com/a-h/templ` - Template engine
- `github.com/Oudwins/tailwind-merge-go` - Tailwind class merging
- `github.com/starfederation/datastar` - Reactive UI library
- `github.com/joho/godotenv` - Environment variable management

### Development Dependencies
- `@html-eslint/eslint-plugin` - HTML linting
- `eslint-plugin-html-tailwind` - Tailwind CSS validation
- `tailwindcss` - CSS framework

This architecture provides a modern, type-safe, and performant foundation for building server-rendered web applications with rich client-side interactivity while maintaining the simplicity and performance benefits of server-side rendering.