# Contributing to BubbleTable

Thank you for your interest in contributing to BubbleTable! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/bubbletable.git
   cd bubbletable
   ```
3. **Install dependencies**:
   ```bash
   make setup
   ```

## Development Workflow

### Before Making Changes

1. **Create a new branch** for your feature or fix:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Run tests** to ensure everything works:
   ```bash
   make test
   ```

### Making Changes

1. **Write tests** for new functionality
2. **Update documentation** if needed
3. **Follow Go conventions** and existing code style
4. **Run quality checks**:
   ```bash
   make fmt
   make vet
   make lint
   ```

### Submitting Changes

1. **Commit your changes** with clear messages:

   ```bash
   git add .
   git commit -m "Add feature: description of changes"
   ```

2. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request** on GitHub

## Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and reasonably sized
- Write tests for new functionality

## Testing

- Write unit tests for all new code
- Ensure tests pass: `make test`
- Add benchmarks for performance-critical code

## Documentation

- Update README.md if adding new features
- Add godoc comments for exported functions
- Update API.md for significant API changes
- Include examples for new functionality

## Pull Request Guidelines

### Title and Description

- Use clear, descriptive titles
- Explain what the PR does and why
- Reference any related issues

### Code Quality

- All tests must pass
- Linter checks must pass
- No merge conflicts

### Review Process

- Be responsive to feedback
- Make requested changes promptly
- Keep discussions focused and constructive

## Types of Contributions

### Bug Fixes

- Include reproduction steps
- Add regression tests
- Reference the issue number

### New Features

- Discuss in an issue first
- Follow existing patterns
- Include comprehensive tests
- Update documentation

### Documentation

- Fix typos and improve clarity
- Add missing examples
- Keep documentation up-to-date

### Performance Improvements

- Include benchmarks
- Measure impact quantitatively
- Don't sacrifice readability unnecessarily

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and improve
- Follow the project's values

## Getting Help

- Open an issue for questions
- Check existing issues and PRs
- Ask in discussions for general questions

## Release Process

1. Features and fixes are merged to `main`
2. Releases are tagged with semantic versioning
3. CHANGELOG.md is updated for each release
4. Documentation is updated as needed

Thank you for contributing to BubbleTable!
