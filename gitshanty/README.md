# Git Shanty 🏴‍☠️

```
        ⛵
       ___|___
      |       |
      | GIT   |
      |SHANTY |
  ~~~~|_______|~~~~
   ~~~~~~~~~~~~~~~
    ~~~~~~~~~~~~~
     ~~~~~~~~~~~
      ~~~~~~~~~
```

> *"Turn yer git log into glorious sea shanties, ye scallywag!"*

**Git Shanty** be a Go CLI that plunders yer commit history and forges it into magnificent sea shanties. Every commit tells a tale — and every tale deserves a song.

Whether ye shipped a grand new feature or sank the build with a cursed merge, Git Shanty will sing of yer deeds across the seven repos.

---

## ⚓ Installation

```bash
# Clone the treasure map
git clone <repo-url> && cd gitshanty

# Forge yer cutlass
go build -o gitshanty .

# ...or take her for a quick sail
go run .
```

---

## 🗺️ Usage

```bash
# Sing a shanty for yer last commit
gitshanty

# Last 5 commits, one shanty each
gitshanty -n 5

# An epic saga of yer entire repo history
gitshanty -all

# Only Dave's commits (poor Dave)
gitshanty -author "Dave"

# Plunder a specific repo
gitshanty -repo /path/to/repo
```

---

## 🎵 The Shanty Codex

Not all commits be equal, and neither be their songs. Git Shanty reads the soul of each commit and assigns it a shanty style befitting its nature:

| Commit Type | Shanty Style | Description |
|---|---|---|
| ✨ New features | **Celebration Jig** | A lively tune for new treasure added to the hold! |
| 🐛 Bug fixes, bad merges | **Mournful Ballad** | A somber dirge for the bugs we've lost... and the ones we caused. |
| 🔨 Big refactors | **Epic Saga** | A sweeping orchestral tale of code reborn from the deep. |
| 💀 Force pushes, delete sprees | **Mutiny Anthem** | A thunderous roar for those who dared rewrite history! |
| 📝 Regular commits | **Work Song** | A steady rhythm to keep the crew rowing through the day. |
| 📖 Docs, tiny changes | **Lullaby** | A gentle whisper for the smallest of changes, soft as the evening tide. |

---

## 🎸 Fictional Bands of the High Seas

Each shanty be performed by one of these legendary crews:

- **The Merge Conflicts** — Masters of chaos, singers of sorrow
- **Captain Rebase & The Detached HEADs** — They'll rewrite yer history and make it sound good
- **The Fatal Exceptions** — Every show ends in a crash (on purpose)
- **Davy Jones' Cache** — They never forget... unless evicted
- **The Dangling Pointers** — Nobody knows where they'll end up next

---

## 🦜 Example Output

```
🏴‍☠️ Commit: a1b2c3d — "fix: null pointer in auth middleware"
🎵 Style: Mournful Ballad
🎸 Performed by: The Fatal Exceptions

  Oh, the pointer was null and the server did fall,
  The auth middleware answered no call at all.
  We searched through the logs by the lantern's pale light,
  And patched up the hole before end of night.

  Yo ho, yo ho, a developer's life for me! 🏴‍☠️
```

---

## 🏴‍☠️ Fair Winds & Following Seas

Built with mass amounts of ☕, mass amounts of ⚓, and mass amounts of 🦜.

*May yer builds be green, yer merges be clean, and yer shanties be heard across the seven seas.*

**Now hoist the mainsail and `go run .` — adventure awaits!** 🏴‍☠️
