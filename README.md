# ppxity

ppxity is a command-line tool that allows you to interact with the [Perplexity API](https://docs.perplexity.ai). It enables you to provide a prompt, along with a set of directories and files, and the tool will compile the content of those files into the prompt and send it to the Perplexity API.
 By default, it uses the claude-3-haiku-20240307 model.
### Why Claude 3 Haiku as default?
Firstly, it works well with the Perplexity API, and secondly the benchmarking results seem to be good.

![img.png](assets/img.png)

### Features
* Input Directories
* Input Files
* Backtrack
* Official Perplexity API support

### Information
This was made for personal use within a few hours, it's far from perfect, it doesn't have a fancy ui yet but I might add that sometime.

I simply wanted to be able to give the AI a bunch of files before asking a question for better responses.

### Requirements
* Perplexity API Key (get one at https://www.perplexity.ai)

### Usage
```bash
Usage:
  ppxity [flags]

Examples:
ppxity -d C:\Users\User\GolandProjects\exampleProject -p "Explain what this project is about."

Flags:
  -k, --api-key string        Perplexity API key (or set PPLX_API_KEY environment variable)
  -D, --debug                 Enable debug mode
  -d, --directories strings   Directories to use for the initial prompt
  -e, --extensions strings    Allowed file extensions to use for the initial prompt (default [go,txt,mod,cs,c,rs,js,ts])
  -f, --files strings         Files to use for the initial prompt
  -h, --help                  help for ppxity
  -m, --model string          Perplexity model to use: e.g. 'claude-3-haiku-20240307' (sonar-small-online, sonar-medium-online, sonar-small-chat, sonar-medium-chat, claude-3-haiku-20240307, codellama-70b-instruct, mistral-7b-instruct, llava-v1.5-7b-wrapper, llava-v1.6-34b, mixtral-8x7b-instruct, mistral-medium, gemma-2b-it, gemma-7b-it, related) (default "claude-3-haiku-20240307")
  -p, --prompt string         Initial prompt for the conversation: e.g. 'Hello, World!'
  -s, --show-initial-prompt   Show the initial prompt
  -t, --timeout int           Timeout in seconds for receiving messages (default 50)
```

### Example use

```bash
git clone https://github.com/0xInception/ppxity
go run main.go -p "What is this project about?" -d /path/to/directory --api-key your_api_key_here -m "llama-3.1-sonar-small-128k-online"
```

or set environment variable:

```bash
export PPLX_API_KEY=your_api_key_here
go run main.go -p "What is this project about?" -d /path/to/directory -m "llama-3.1-sonar-small-128k-online"
```

Alternative:
```bash
ppxity.exe -p "What is this project about?" -d C:\path\to\directory -f C:\path\to\file.go -f C:\path\to\another\file.go --api-key your_api_key_here -m "llama-3.1-sonar-small-128k-online"
```

### Important Note on Model Selection
Model availability depends on your Perplexity API subscription plan. Common models include:
- Free tier: `llama-3.1-sonar-small-128k-online`, `llama-3.1-sonar-large-128k-online`
- Paid tier: Additional models like `claude-sonnet-4-29339`, `gpt-4o`, etc.

Check the [Perplexity API documentation](https://docs.perplexity.ai/getting-started/models) for the complete list of models available to your API key.

### License
MIT
