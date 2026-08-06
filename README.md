<div align="center">

# toofan

**A minimal, lightning-fast typing TUI**  
_Practice with english words or real code snippets. No browser, no account, everything stays local._

<br>

<img src="assets/main.gif" alt="toofan demo" width="750">

</div>

---

## Features

- **Two Modes:** Practice standard English words or real-world code snippets.
- **Curated Lessons:** Hand-written, topic-based code exercises across multiple languages.
- **Dynamic Themes:** Cycle between multiple aesthetic terminal themes (`ctrl+t`).
- **Bot Racing:** Race against 1-5 customizable AI bots that adapt to your speed (`ctrl+b`).
- **Multiplayer Racing:** Compete with others in real-time via the built-in race server (`ctrl+n`).
- **Live Metrics:** Real-time WPM speed and accuracy tracking.
- **Error Review:** See exactly which words you mistyped after every test.
- **Ranks:** Automated progression system based on your typing speed.
- **Offline & Local:** No browser, no account, and zero tracking. All data stays on your machine (except for real-time stats shared during multiplayer races).
- **Data Collection in multiplayer:** While playing multiplayer, to make sure you have a smooth experience, we collect only the username, IP address (required for preventing multiple multiplayer sessions for the same user), room size choice, progress and your WPM. We DO NOT collect history, PBs, config, keystrokes, or anything except the above listed. Our commitment to privacy is absolute.

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

## Installation

⚠️ **Note:** Always take a backup (`ctrl+s`) before updating toofan.

### curl (macOS & Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/vyrx-dev/toofan/master/install.sh | sh
```

### AUR

```bash
paru -S toofan-bin
```

### Go

```bash
go install github.com/vyrx-dev/toofan@latest
```

### Homebrew / Nix / Ubuntu / Fedora

Coming soon.

### Build from Source

If you prefer building manually (requires Go):

```sh
git clone https://github.com/vyrx-dev/toofan.git
cd toofan
go build -o toofan .
mv toofan ~/.local/bin/
```

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

- `config.txt` : Your selected duration, mode, language, and theme
- `results.txt` : Every test result (date, wpm, accuracy, duration, mode)
- `pb.txt` : Your personal bests per mode and duration
</details>

<details>
<summary>Can I backup my data?</summary>

Yes. Press `ctrl+s` to save a backup and `ctrl+r` to restore from one. Backups are saved to `~/.config/toofan/backups/` and can be moved between machines.

</details>

<details>
<summary>Does online mode require account?</summary>

No. Online mode only requires a username that uniquely identifies you across the server. No account or registration is required.

</details>

<details>
<summary>What data is collected in online mode?</summary>

Only IP, room size choice, progress, your WPM, and your username are used during online mode. This data is used to exchange information between clients and log join/leave events. No history, PBs, config, or keystrokes are stored on the server.

</details>

<details>
<summary>Does playing against bots require internet connection?</summary>

No. All bots run locally on your machine. Only multiplayer mode requires internet.

</details>

<details>
<summary>How do I update toofan?</summary>

The update process depends on how you installed it:

**curl (Quick Install):**
Just run the install command again. It will automatically download and replace the old binary.

```bash
curl -fsSL https://raw.githubusercontent.com/vyrx-dev/toofan/master/install.sh | sh
```

**Go:**

```bash
go install github.com/vyrx-dev/toofan@latest
```

**AUR:**
Use your AUR helper to update the package:

```bash
paru -Syu toofan-bin
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

Yes. Everything runs locally and is embedded in the binary. No internet needed.

</details>

<details>
<summary>Want more programming languages?</summary>

We're always looking to add more. If your favorite programming language isn't supported yet, open a PR with a few lesson files and we'll get it in. Check `AGENTS.md` for the file format.

</details>

## Roadmap

- [x] Curl script installation (macOS & Linux)
- [x] Proper documentation for AI and contributors
- [ ] More language support (python, rust, c, typescript, etc.)
- [x] Difficulty levels for english words
- [ ] AUR, Homebrew, Nix packages, Fedora & Ubuntu repos
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
