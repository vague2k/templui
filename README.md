# Goilerplate

Modern UI Components for Go & Templ

<img src="./assets/img/gopher.svg" alt="Goilerplate Logo" width="200"/>

## About

Goilerplate is a growing library of modern UI components designed specifically for Go and Templ. It leverages Alpine.js for enhanced interactivity and Tailwind CSS for effortless styling. Whether you're building a small website or a large web application, Goilerplate provides the tools you need to create sleek, responsive interfaces with ease.

## Features

- **Go-native Implementation**: Optimized for Go developers, seamlessly integrating with Go backends.
- **Templ-first Design**: Leverages the full power of Templ for type-safe, high-performance templating.
- **Server-Side Rendering (SSR) Focus**: Excellent performance and SEO benefits out of the box.
- **Alpine.js Integration**: Enhanced client-side interactivity when needed.
- **Tailwind CSS Styling**: Modern, utility-first styling that's highly customizable.
- **Accessible Components**: Built with accessibility in mind, following WCAG guidelines.
- **TypeSafe**: Utilizing Go's type system for robust, error-resistant development.

## Getting Started

There are two main ways to use Goilerplate in your projects:

1. **Use as a Package Library**

   Install Goilerplate as a Go package:

   ```
   go get github.com/axzilla/goilerplate
   ```

   Then import and use components in your Templ files:

   ```go
   import "github.com/axzilla/goilerplate/pkg/components"

   // In your Templ files
   @components.Button(components.ButtonProps{Text: "Click me"})
   ```

2. **Copy Components to Your Codebase**

   Visit our [components documentation](https://goilerplate.com/docs/components) to find and copy the components you need directly into your project.

You can also mix and match these approaches based on your project needs.

For a quick start, check out our [Goilerplate quickstart](https://github.com/axzilla/goilerplate-quickstart) template, which provides a pre-configured setup using Goilerplate as a package library.

For detailed setup instructions and examples, visit our [how to use guide](https://goilerplate.com/docs/how-to-use).

## Components

Explore our growing list of components in the [components documentation](https://goilerplate.com/docs/components). Each component comes with usage examples and code snippets.

## Inspiration

Goilerplate draws inspiration from several popular UI libraries and frameworks:

- [shadcn/ui](https://ui.shadcn.com/)
- [Pines UI](https://devdojo.com/pines)
- [daisyUI](https://daisyui.com/)

We're exploring whether to make Goilerplate a strict port of shadcn/ui or to create a unique blend of various inspirations. The project is still evolving, and community feedback will play a crucial role in shaping its direction.

## Current Status

- **Heavy Development**: The project is under active development. Expect frequent updates and potential breaking changes until we reach a stable version.
- Actively growing component library
- Regular updates and bug fixes
- Continuous development based on community feedback

## Changelog

See [releases](https://github.com/axzilla/goilerplate/releases) for a detailed list of changes in each version.

## Contributing

We welcome contributions from the community! Whether it's adding new components, improving existing ones, or enhancing documentation, your input is valuable. Please check our [contributing guidelines](CONTRIBUTING.md) for more information on how to get involved.

## Feedback

Your feedback is crucial in shaping the future of Goilerplate. If you have suggestions, feature requests, or encounter any issues, please [open an issue](https://github.com/axzilla/goilerplate/issues) on our GitHub repository or reach out to us through our website.

## License

Goilerplate is open-source software licensed under the [MIT license](LICENSE).

## Support

For support, questions, or discussions, please [open an issue](https://github.com/axzilla/goilerplate/issues) on our GitHub repository.

---

Built with ❤️ by the Go community, for the Go community.

[![Product Hunt](https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=494295&theme=light)](https://www.producthunt.com/posts/goilerplate-1?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-goilerplate-1)
