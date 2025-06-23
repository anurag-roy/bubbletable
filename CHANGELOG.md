# Changelog

All notable changes to BubbleTable will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Comprehensive package documentation with godoc
- Performance benchmarks for core operations
- Support for custom formatters and renderers
- Multiple key binding presets (Default, Vim, Emacs)
- Built-in themes (Dracula, Monokai, GitHub, Solarized, etc.)
- Headless table core for maximum flexibility
- Struct tag support for automatic column configuration
- Real-time search and filtering
- Multi-state sorting (unsorted → asc → desc → unsorted)
- Efficient pagination for large datasets
- Callback system for user interactions
- Responsive design that adapts to terminal size

### Changed

- Converted from standalone application to reusable library
- Restructured codebase into separate packages (table, renderer, components)
- Improved API design with fluent builder pattern
- Enhanced performance with optimized rendering and data handling

### Fixed

- Memory efficiency improvements for large datasets
- Consistent theme application across all components
- Proper handling of different data types in sorting and formatting

## [0.1.0] - Initial Release

### Added

- Basic table functionality
- Simple rendering capabilities
- Core data structures and operations
