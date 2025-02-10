# `go-clack/core`

Clack contains low-level primitives for implementing your own command-line applications.

Currently exposes `Prompt` as well as:

- `TextPrompt`
- `PasswordPrompt`
- `PathPrompt`
- `ConfirmPrompt`
- `SelectPrompt`
- `MultiSelectPrompt`
- `GroupMultiSelectPrompt`
- `SelectPathPrompt`
- `MultiSelectPathPrompt`
- `SelectKeyPrompt`

Each `Prompt` accepts a `Render` function.

```go
p := core.NewTextPrompt(core.TextPromptParams{
  Render: func(p *core.TextPrompt) string {
    return fmt.Sprintf("What's your name?\n%s", p.ValueWithCursor)
  },
})

name, err := p.Run()
if (err != nil) {
  // Handle prompt's cancellation
  os.Exit(0)
}
```

## Get Started

To start using `go-clack/core`, follow these steps:

### 1. Install the Package

First, add the `go-clack/core` package to your Go project:

```bash
go get github.com/orochaa/go-clack/core
```

### 2. Create a Prompt

To create and run a simple text prompt, you can use the following code:

```go
// main.go
package main

import (
  "fmt"
  "os"
  "github.com/orochaa/go-clack/core"
)

func main() {
  p := core.NewTextPrompt(core.TextPromptParams{
    Render: func(p *core.TextPrompt) string {
    return fmt.Sprintf("What's your name?\n%s", p.ValueWithCursor())
    },
  })

  name, err := p.Run()
  if err != nil {
      fmt.Println("Prompt was canceled.")
      os.Exit(0)
  }

  fmt.Printf("Hello, %s!\n", name)
}
```

### 3. Run Your Application

Compile and run your application:

```bash
go run main.go
```

This will present a text prompt asking for the user's name. The input will be captured and printed back as a greeting.

## Explore More

The `go-clack/core` package provides various other prompts for different types of user inputs, such as password inputs, file path selections, confirmations, and more. Explore the available prompts and customize their behavior using the Render function, or create a brand new one extending the core `Prompt`.
