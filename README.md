# Flare

Flare is a syntax highlighting library built using the incremental parser [GPeg](https://github.com/zyedidia/gpeg).

Features:

* Supports fast incremental rehighlighting for use in an editor.
* Can load and use highlighters at runtime.
* Multiple languages are supported, and it is easy to add more (check the `languages/` directory).
* Supports limited context-sensitivity (back references) for handling HERE-doc style blocks.
* Supports including languages within each other.
* Fairly efficient.
* Can format highlighted code as HTML output.

See the `flare` program in `./cmd/flare` for an example that can highlight code files at the terminal
or output highlighted HTML.
