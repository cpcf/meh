# meh - Machine Enhanced Help

## Summary

Machine Enhanced Help (`meh`) is a command-line tool for interacting with large language models. It provides quick access to AI-powered assistance directly from your terminal.

## Features

- **Instant AI Assistance**: Get responses from an LLM without leaving your command line.
- **Context-Aware Queries**: Ask follow-up questions while maintaining conversation context.
- **Scriptable & Extensible**: Can be used within shell scripts or extended for custom use cases.

## Installation

### Prerequisites

## Usage

### Basic Usage
```sh
meh "Explain quantum entanglement in simple terms"
```

### Pipe Usage
```sh
git diff --word-diff=porcelain | meh "Write a commit message for this diff"
```

### Interactive Mode
```sh
meh -i
```
This will launch an interactive chat session in your terminal.

### Using with a File
```sh
meh -f input.txt
```
Processes the contents of `input.txt` as a query.

### Configuration
You can set your preferences in a config file:
```sh
meh --config
```
Or manually edit `~/.config/meh/config.json`:
```json
{
  "api_key": "your-api-key-here",
  "server": "llmserver.local:1234/v1",
}
```

## License
This project is licensed under the MIT License. See `LICENSE` for details.

