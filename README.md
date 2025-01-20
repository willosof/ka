# KA: Kill All, but Better

## 🚀 Introduction
KA (Kill All) is your ultimate sidekick for managing runaway processes on macOS. If you’ve ever been frustrated by `killall`'s inability to handle applications with numerous sub-processes, KA is here to make your life easier. With an intuitive interface and powerful killing powers, KA turns the chaotic mess of sub-process wrangling into a smooth, interactive killing experience.

🎯 **Why?** Because macOS's `killall` can be frustratingly ineffective when working with complex applications that spawn endless subprocesses. KA steps in to make sure nothing gets left behind.

## ✨ Features
- **Smart Process Killing**: Target processes by name with surgical precision—no stragglers left behind.
- **Interactive Interface**: Use a sleek, scrollable, multi-select interface to choose which processes to terminate.
- **Signal Flexibility**: Supports sending custom signals (default is `SIGTERM`, but feel free to get creative).
- **Safety First**: No more accidental self-termination; KA excludes itself from the kill list.
- **Batch Killing**: Mass-murder processes with a single command or confirm individually. Your call.
- **Highlighting for Clarity**: See your target process names highlighted in glorious ANSI colors.

## 🛠 Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/willosof/ka.git
   cd ka
   ```

2. **Build the binary:**
   ```bash
   go build -o ka
   ```

3. **Move it to your path:**
   ```bash
   mv ka /usr/local/bin/
   ```

4. **Profit:**
   ```bash
   ka 
   ```

## 🧑‍💻 Usage
```bash
Usage: ka [options] process_name

Options:
  -s SIGNAL   Signal to send (e.g., -s 9 for SIGKILL)
  -y          Assume yes; kill all matching processes without confirmation
```

### Examples
1. Kill all `node` processes interactively:
   ```bash
   ka node
   ```

2. Kill all `chrome` processes immediately with `SIGKILL`:
   ```bash
   ka -s 9 -y chrome
   ```

3. Pretend you're the Terminator and hunt down processes by name. 🌟

## 🔧 Roadmap
- **Windows Support**: Because macOS users shouldn’t have all the fun.
- **Better ANSI Styling**: Upgrade from our current `survey` library to the shiny new Charm ecosystem for even slicker visuals.
- **Enhanced Filtering**: Add support for regex-based process filtering.

## 🤝 Contributing
We’d love your help to make KA even better! Here’s how you can get involved:

1. Fork the repo 🍴
2. Make your changes ✍️
3. Submit a pull request 🛠️

Please ensure your code is well-tested (on macOS) and adheres to Go best practices. We’re especially keen on contributions to:
- Add Windows/Linux support
- Improve our ANSI output and interface
- Squash bugs we haven’t noticed yet

## 🚨 Limitations
- **MacOS Only**: Currently, KA is tested and functional only on macOS. If you're on Windows or Linux, feel free to help us expand compatibility!
- **Dependencies**: KA relies on `pgrep` and `ps`, which are standard on macOS. If these are missing, things might get wobbly.

## 🎉 Credits
Made with ❤️ by developers who just wanted processes to die properly. Inspired by the chaos of modern macOS development.

## 📄 License
MIT License. Do whatever you want, but don't blame us if you kill the wrong thing. 😉

