<div align="center">

# toofan

**A minimal, lightning-fast typing TUI**  
_Practice with english words or real code snippets. No browser, no account, everything stays local._

<br>

<img src="assets/main.gif" alt="toofan demo" width="750">

</div>

---

> **This is the `toofan-online` branch.** It adds multiplayer racing and AI bots.
> Still early, so expect a few rough edges. Found a bug? [Open an issue](https://github.com/vyrx-dev/toofan/issues).
> Want the stable offline-only version? Head to the [`master` branch](https://github.com/vyrx-dev/toofan/tree/master).

## Features

- **Two Modes:** Practice standard English words or real-world code snippets.
- **Curated Lessons:** Hand-written, topic-based code exercises across multiple languages.
- **Dynamic Themes:** Cycle between multiple aesthetic terminal themes (`ctrl+t`).
- **Bot Racing:** Race against 1-5 AI bots that adapt to your typing speed (`ctrl+b`).
- **Multiplayer Racing:** Compete with others in real-time online races (`ctrl+n`).
- **Live Metrics:** Real-time WPM speed and accuracy tracking.
- **Error Review:** See exactly which words you mistyped after every test.
- **Ranks:** Automated progression system based on your typing speed.
- **Offline-first:** Solo mode and bots work with zero internet. Only multiplayer needs a connection.

<p align="center">
  <img src="assets/code-snippets-grid.png" width="48%" title="Real Code Snippets" alt="Real Code Snippets" />
  <img src="assets/lession-grid.png" width="48%" title="Curated Topics & Lessons" alt="Curated Topics & Lessons" />
  <img src="assets/languages-grid.png" width="48%" title="Multiple Languages Supported" alt="Multiple Languages Supported" />
  <img src="assets/theme-grid.png" width="48%" title="Dynamic Built-in Themes" alt="Dynamic Built-in Themes" />
</p>

## Profile Dashboard

A personal overview of your typing speed history, personal bests across durations, and a daily activity map to keep you consistent. Press `ctrl+p` to open.

<div align="center">
<img src="assets/profile-new.png" width="95%">
</div>

## Installation (toofan-online)

> The `curl` and AUR install methods below are specific to the online build.
> For the stable offline-only version, see the [`master` branch](https://github.com/vyrx-dev/toofan/tree/master).

### curl (macOS & Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/vyrx-dev/toofan/toofan-online/install-online.sh | sh
```

### AUR

```bash
paru -S toofan-online-bin
```

### Go

```bash
go install github.com/vyrx-dev/toofan@v1.0.1-online
```

### Build from Source

```bash
git clone -b toofan-online https://github.com/vyrx-dev/toofan.git
cd toofan
go build -o toofan .
mv toofan ~/.local/bin/
```

### GitHub Release

Download the binary for your platform from the [releases page](https://github.com/vyrx-dev/toofan/releases/tag/v1.0.1-online).

## FAQ

<details>
<summary>How are stats calculated?</summary>

```text
raw      = total_chars / 5 / elapsed_minutes
wpm      = (total_chars - uncorrected_errors) / 5 / elapsed_minutes
accuracy = (total_chars - all_mistakes) / total_chars × 100
```

- **wpm** - your net speed. Every 5 characters count as one "word". Uncorrected mistakes are subtracted.
- **accuracy** - counts every wrong keystroke, even if you corrected it with backspace.
- **raw** - your gross speed before any penalty.
- **errors** - press `e` on the results page to see exactly which words you mistyped.
</details>

<details>
<summary>Where are my files stored?</summary>

Everything lives in `~/.config/toofan/` as plain text files:

- `config.txt` : Your selected duration, mode, language, theme, and multiplayer username
- `results.txt` : Every test result (date, wpm, accuracy, duration, mode)
- `pb.txt` : Your personal bests per mode and duration
</details>

<details>
<summary>Can I backup my data?</summary>

Yes. Press `ctrl+s` to save a backup and `ctrl+r` to restore from one. Backups are saved to `~/.config/toofan/backups/` and can be moved between machines.

</details>

<details>
<summary>Does multiplayer need an account?</summary>

No. You just pick a username when you start a multiplayer race. No sign-up, no email, no password. The username is stored locally in your config file so you don't have to type it every time.

</details>

<details>
<summary>What data is shared during multiplayer?</summary>

Only what's needed for the race to work:

- **Your username** so other players can see who they're racing
- **Your IP address** to prevent duplicate sessions (not stored after the race)
- **Your typing progress and WPM** so the race bars update live
- **Room settings** like mode, language, difficulty, duration

That's it. Your typing history, personal bests, config, and keystrokes never leave your machine. None of this is saved on the server after a race ends.

</details>

<details>
<summary>Is multiplayer safe to use?</summary>

The server is a simple Go HTTP server. No database, no persistent storage, no analytics. Race data only lives in memory and goes away when the room closes or the server restarts.

The server code is fully open source in the `race_server/` directory if you want to look through it or run your own. If you'd rather stay offline, just don't press `ctrl+n`. Solo mode and bots work without any internet.

</details>

<details>
<summary>Can I run my own multiplayer server?</summary>

Yes. The server is a standalone Go binary in `race_server/`:

```bash
cd race_server
go build -o toofan-server .
./toofan-server --port 8525
```

Then set your server URL in toofan. No external dependencies needed, just Go standard library.

</details>

<details>
<summary>Do bots need internet?</summary>

No. Bots run entirely on your machine. They match your recent average speed with some variation so it feels natural. No internet needed.

</details>

<details>
<summary>How do I update toofan-online?</summary>

**Go:**

```bash
go install github.com/vyrx-dev/toofan@v1.0.1-online
```

**Build from source:**

```bash
cd toofan && git pull && go build -o toofan . && mv toofan ~/.local/bin/
```

</details>

<details>
<summary>How do I uninstall Toofan?</summary>

If you installed via the `curl` Quick Install, simply delete the binary and the configuration folder:

```bash
rm ~/.local/bin/toofan
rm -rf ~/.config/toofan
```

_(If you built it from source and moved it globally, run `sudo rm /usr/local/bin/toofan` instead)._

</details>

<details>
<summary>Does it work offline?</summary>

Yes. Solo mode and bots are fully offline. Everything is baked into the binary.
Only multiplayer (`ctrl+n`) needs a connection. If you never press `ctrl+n`, it works exactly like the stable version.

</details>

<details>
<summary>Want more programming languages?</summary>

We're always looking to add more. If your favorite programming language isn't supported yet, open a PR with a few lesson files and we'll get it in. Check `AGENTS.md` for the file format.

</details>

## Roadmap

- [x] Curl script installation (macOS & Linux)
- [x] Proper documentation for AI and contributors
- [x] More language support (Go, Rust, Python, C, C++, C#, TypeScript, Java)
- [x] Difficulty levels for english words
- [x] Bot racing
- [x] Multiplayer racing
- [ ] Homebrew, Nix, Fedora & Ubuntu packages
- [x] Fix top pane alignment to match bottom panes in profile

## Contributors

- Huge thanks to [@aaravmaloo](https://github.com/aaravmaloo) for creating online mode and implementing AI bots.

(As a contributor, to get your name to mentions, please create a PR along with the contributions you made.)
## Contributing

- New snippets : Drop a file in `internal/lang/data/<language>/lessons/` and rebuild
- New languages : Just a folder with lesson files
- New themes : One Go file with a color palette
- Bug fixes and UX improvements

If you're using an AI coding assistant, read [`AGENTS.md`](AGENTS.md) first.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) : TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) : Terminal styling

---

<a href="https://www.star-history.com/#vyrx-dev/toofan&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=vyrx-dev/toofan&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=vyrx-dev/toofan&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=vyrx-dev/toofan&type=date&legend=top-left" />
 </picture>
</a>
