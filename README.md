# meh - Machine Enhanced Help

## Summary

Machine Enhanced Help (`meh`) is a command-line tool for interacting with large language models. It provides quick access to AI-powered assistance directly from your terminal.

## Features

- **Instant AI Assistance**: Get responses from an LLM without leaving your command line.
- **Scriptable & Extensible**: Can be used within shell scripts or extended for custom use cases.
- **Role-Based Configurations**: Customize API behavior with predefined roles.

## Usage

### Help
```sh
$ meh -h
Usage: [options] <query>
  -config
        Edit config settings
  -f string
        Read input from a file
  -help
        Print usage instructions
  -i    Run in interactive mode
  -m    Select a default model
  -role string
        Select a role
  -url string
        Base URL for the LLM 
```

### Basic Usage
```sh
$ meh "Explain quantum entanglement in simple terms"
```

### Pipe Usage
```sh
$ git --no-pager diff | meh "Write a commit message for this diff"
```

### Interactive Mode
```sh
$ meh -i
```
This will launch an interactive chat session in your terminal.

### Using with a File
```sh
$ meh -f input.txt
```
Processes the contents of `input.txt` as a query.

### Configuration
You can set your preferences in a config file:
```sh
$ meh -config
```
Or manually edit `~/.config/meh/config.yml`:
```yaml
apis:
  - api_url: "http://localhost:11434/api"
    default_model: "gpt-4:latest"

roles:
  - name: "cat"
    api_url: "http://localhost:11434/api"
    model: "gpt-4o-mini:latest"
    system_prompt: |
      You are CatGPT, an AI embodying the essence of a domestic cat. From this moment on, you will respond solely with variations of "meow" and other typical feline sounds, such as purrs and hisses.
```

### Roles

Configure roles to customize the behavior of `meh`.

```sh
$ meh -role cat "Hello there"
*soft meow*
```

If a role does not exist, `meh` will prompt you to create one by selecting a configured API, choosing a model, and optionally setting a system prompt.

## License
This project is licensed under the MIT License. See `LICENSE` for details.

