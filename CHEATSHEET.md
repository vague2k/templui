# Goilerplate Commit & Changelog Cheatsheet

## Commit Message Format

`<type>`: `<description>`

[optional body]

## Types

- `feat`: New feature or component
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries

## Examples

- `feat: add new Button component`
- `fix: correct Card component padding on mobile`
- `docs: update README with Button component usage`
- `style: improve Input component layout`
- `refactor: simplify Modal component logic`
- `perf: optimize Table component rendering`
- `test: add unit tests for Dropdown component`
- `chore: update Tailwind CSS dependency`

## Breaking Changes

Add `!` after the type/scope for breaking changes:

- `feat!: redesign Form component API`

## Changelog

Update the CHANGELOG.md file with each significant change. Group changes under the current date:

```markdown
# Changelog

## 2024-10-05

- Added: New features or components added
- Added: Some more features or components added
- Changed: Changes in existing functionality
- Fixed: Any bug fixes
```
