<!-- SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru> -->
<!-- SPDX-License-Identifier: BSD-2-Clause -->

# flang-tui — terminal screen layout in flang

[По-русски](README.ru.md)

A library of pure, total functions: **data plus width and height → finished
lines**. No orders, no reads, no writes. Raw terminal mode, key input, the
output buffer and SIGWINCH stay with the host; numbers and strings come in,
strings go out.

It is *emitted* to **Go** and to **C**, so it travels inside somebody else's
binary — no runtime, no Node, no interpreter.

The reference behaviour is [digitdisk](https://github.com/digitable-lol/digitdisk),
`host/internal/{report,ui,sysinfo}`. Every piece **was** diffed against it over
a grid of 494,467 inputs, once, and the reference's answers at the boundaries
now live here as examples. The diff itself is over; that is said in full below,
with the numbers and with the commit the differ can be fetched from.

**There is no second language in this tree.** Not Go, not Python, not
JavaScript, and no `Makefile` either: the entry point is `./ярлык` — 112 lines
of `sh`, 56 of them code — over a list that is itself a flang program.

The source is written in flang, whose surface is Russian. File names are
English, module names are English, the prose inside is Russian — that is the
language's own convention, and this repository follows it. Two names are
Cyrillic on purpose, `ярлык` and `ярлыки.flang`, taken from the language's own
tree: `./ярлык проверка` is what a person types, not an internal name.

---

## Six modules

| File | Module | What it does | Reference in digitdisk |
|---|---|---|---|
| `flang/screen.flang` | «Screen» | fill bar, clipping by *printing* cells, fitting a section into a height, the whole frame | `ui/widgets.go` `bar`/`clip`/`plainWidth`, `ui/screen.go` `frame` |
| `flang/colour.flang` | «Colour» | truecolour, the xterm-256 cube, the sixteen basics | `ui/theme.go` `seq`, `cube256` |
| `flang/format.flang` | «Format» | sizes, shares, percentages, durations, grouped counters, the em-dash placeholder, truncation | `report/report.go`, `report/places.go`, `ui/widgets.go`, `ui/screen.go`, `sysinfo/sysinfo.go` |
| `flang/history.flang` | «History» | running-graph history and its row of glyphs | `ui/screen.go` `push`, `ui/widgets.go` `spark` |
| `flang/tabs.flang` | «Tabs» | tabs: state plus keypress → state | `ui/screen.go` `handle` |
| `flang/scroll.flang` | «Scroll» | scrolling: offset, viewport height, content length, keypress → offset | `ui/screen.go` `handle` + `frame` |

`tools/licensing.flang` is the licence guard: a plan for `flang io`, the method
borrowed from digitdisk. There is no Python, no JavaScript and **no Go** in this
tree, and that is one of its four checks. `.go` moved into the guard's list of
foreign extensions on 2 September 2026, when the last Go left; while the differ
was here, it had to be an exception.

`tools/sverka/` is **gone**, and with it the last hand-written Go in the tree.
It was a differ — a script plus Go files dropped into a disposable clone of
digitdisk — and it did its job: 494,869 inputs, 0 divergences. What it guarded
now lives here, as examples: see *The diff: what it proved, and where it went*.

## Getting it

The flang compiler is a single binary that needs only `cc`:

```sh
brew install flang        # or: asdf plugin add flang
# or from a clone of the language:  make -C bootstrap
```

Then, and there is no `make` here any more:

```sh
./ярлык             # the list of shortcuts, read out of ярлыки.flang
./ярлык проверка    # flang check + flang test: modules, guard, shortcuts
./ярлык ведомость   # the proof ledger, printed in full
./ярлык печать      # emit to Go and to C, build both
./ярлык пределы     # F3: where the emitted program's depth limit binds
./ярлык лицензии    # the licence guard
```

`ярлык` is 112 lines of `sh`, 56 of them code, and it holds **no list**: it asks
the binary for the command string and runs it. The list itself is
`ярлыки.flang` — a *program*, type-checked, with examples, whose
`«Сколько ярлыков»` carries the count as a postcondition rather than as prose.
The trick is the language's own (`ярлык` + `ярлыки.flang` in its tree), and the
reason to copy it is the same: the `Makefile` that used to be here had no
dependency and no timestamp rule in it — only `.PHONY` and a list of command
strings, which is exactly what `make` was being kept for.

One thing the Makefile did carry has been dropped on purpose: the list of
module names (`МОДУЛИ = screen colour format history tabs scroll`). It could
drift from the tree in silence — a seventh module would simply not be checked
and not be emitted, and everything would stay green. The shortcuts say
`flang/*.flang` instead; there is nothing left to drift.

Emission puts each module in its own `out-go/<name>` and `out-c/<name>`. All
eight emit targets name the Go module `flangprogram`, so after emission the
shortcut rewrites it to `flang<name>`: six libraries in one build would
otherwise collide. That is the *only* edit made to emitted code, and `sed`
makes it, not a hand. (`make` still runs — inside `out-c/<name>`, over a
Makefile the *compiler* printed. That one is generated, gitignored, and nobody
hand-edits it.)

---

## The diff: what it proved, and where it went

**It is over, and it is not pretending otherwise.** The diff against digitdisk
was a one-off argument — *this was rewritten correctly* — not a standing check.
It never ran in CI, not once: the workflow of this repository has never invoked
`tools/sverka/run.sh`. Keeping a second language in the tree to make a point
that had already been made was the whole of what `tools/sverka/` was doing.

**Last run: 2 September 2026.** flang 0.6.2 (`bacde89`), digitdisk pinned at
`7ea03ed` (0.5.0), Go 1.26.5, Linux. Output, piece by piece:

| Piece | Reference | Inputs | Divergences |
|---|---|---:|---:|
| Байты (bytes) | `report.Bytes` | 10,922 | 0 |
| Процент целого (percent of) | `report.pct` | 3,600 | 0 |
| Прочерк (dash) | `report.dash` | 13 | 0 |
| Обрезать (cut) | `report.cut` | 225 | 0 |
| Прошло (since) | `report/places.since` | 10,811 | 0 |
| Длительность (duration) | `sysinfo.HumanDuration` | 230,782 | 0 |
| Проценты (percent) | `ui.percent` | 25,014 | 0 |
| Доля от (share of) | `ui.pctOf` | 2,400 | 0 |
| Заняло (lasted) | `ui.lasted` | 28,572 | 0 |
| Уровень доли (level) | `ui.Theme.level` | 1,301 | 0 |
| Полоса (bar) | `ui.Theme.bar` | 4,020 | 0 |
| Обрезать по ячейкам (clip) | `ui.Theme.clip` | 896 | 0 |
| Ширина без последовательностей | `ui.plainWidth` | 14 | 0 |
| Уложить раздел (window) | `ui.screen.frame` | 5,940 | 0 |
| Дописать замер (push) | `ui.push` | 600 | 0 |
| График (spark) | `ui.Theme.spark` | 837 | 0 |
| Переключить (tab switch + scroll reset) | `ui.screen.handle` | 300 | 0 |
| Прокрутить (scroll) | `ui.screen.handle` + `frame` | 2,835 | 0 |
| Высота тела (body height) | `ui.screen.bodyHeight` | 201 | 0 |
| Куб 256 (cube) | `ui.cube256` | 140,608 | 0 |
| Цвет (colour seq) | `ui.Theme.seq` | 24,576 | 0 |
| **Total** | | **494,467** | **0** |

(A further 402 inputs — `detectDepth`, `UsableTERM`, `PaletteByName` — belong to
the sibling repository
[flang-env](https://github.com/digitable-lol/flang-env). 494,869 across both.)

**How to re-run it, in full, today.** The differ was not deleted from history;
it was deleted from the tree. It is one `git show` away, at the commit where it
last lived:

```sh
git show 6bb627cc45966bfc24cc2680e1eec2196fb2a43d:tools/sverka/run.sh > /tmp/run.sh
git show 6bb627cc45966bfc24cc2680e1eec2196fb2a43d:tools/sverka/ui_test.go     > /tmp/ui_test.go
git show 6bb627cc45966bfc24cc2680e1eec2196fb2a43d:tools/sverka/report_test.go > /tmp/report_test.go
# tools/sverka/{bench,steps}_test.go are there too, for --замер and --пределы
```

Whoever changes the behaviour of a diffed function should do exactly that, and
say the new numbers. That is the honest replacement for a check that was never
running anyway.

### What came back into the tree instead: 362 examples, taken from the reference

The diff itself cannot be written in flang — flang has no way to call Go, and
`пример` compares against a *written-down value*, not against another program.
What **can** be carried across is the reference's answers, and they were:

* the inputs were picked by hand at the boundaries the grid had been walking —
  unit steps of `Байты` (1,048,575 → `1024.0 КиБ`, the reproduced quirk), the
  rounding ties `Ф1` found (`0.0005` → `0.1%`, `1050 ms` → `1.1 с`), the
  clamps of `Прокрутить`, every string in the `clip` corpus, the whole grid
  wherever it was small enough to carry entire (`dash`, `plainWidth`,
  `UsableTERM`, `PaletteByName`);
* the **expected values were produced by digitdisk**, on the pinned `7ea03ed`,
  by a generator run once against that clone — not typed out by hand, and not
  read off flang;
* they were then put into the modules and run: **362 examples, 362 passed, 0
  failed, first try.** That run is itself a diff — a small, permanent one.

**And it is smaller by three orders of magnitude, which is the point of saying
it in numbers.** 494,869 inputs against a live reference, run by hand and never
in CI, became **362 inputs pinned in the source and run on every push**, plus
the 249 examples that were already there. 0.07% of the grid, ∞× the frequency.
An example pins a value; it does not notice if digitdisk changes. Nothing here
claims otherwise.
---

## The ledger: what actually carries each claim

Three words, and they never blur. **«доказано»** (proved) — about *all* inputs.
**«сетка N»** (grid of N) — computed on N of the author's values; that is not a
proof. **«объявлено, не доказано»** (declared, not proved) — neither theorem
nor examples; the runtime computes it on whatever inputs arrive. The numbers
are `flang check --proof` output.

| Module | Functions | All total | Claims | proved | grid | declared, not proved | Examples |
|---|---:|---|---:|---:|---:|---:|---:|
| Screen | 19 | yes | 25 | 4 | 15 | 6 | 101 |
| Colour | 9 | yes | 5 | 1 | 4 | 0 | 53 |
| Format | 33 | yes | 14 | 4 | 9 | 1 | 200 |
| History | 9 | yes | 8 | 3 | 5 | 0 | 51 |
| Tabs | 5 | yes | 8 | 0 | 6 | 2 | 37 |
| Scroll | 6 | yes | 7 | 0 | 7 | 0 | 66 |
| **Total** | **81** | **81 of 81** | **67** | **12** | **46** | **9** | **508** |

`tools/licensing.flang` is 22 functions, all total, 30 examples; `ярлыки.flang`
is 5 functions, all total, 6 examples. No ledger is printed for the guard
because it declares a `план`, whose laws the binary does not judge — and it says
so itself.

Examples went from 187 to 508 in this change: 321 of them are the reference's
answers, carried in when the differ left (the section above). Claims did not
move — 67, of which 12 proved. Adding examples does not prove anything, and the
`сетка N` column grew for exactly that reason: a bigger grid is still a grid.

Twelve proved claims is 18% of sixty-seven. The rest is grid and declared. That
is written down as it is; nobody here will call a grid a proof.

---

## Emission: eight targets, two verified

| Target | Modules emitted | Check |
|---|---:|---|
| Go | 6 of 6 | `gofmt -l` empty, `go vet ./...` clean, `go build ./...` clean |
| C | 6 of 6 | `cc -std=c99 -Wall -Wextra -Werror -pedantic -O2 -flto`, **0 warnings about our code** |

Run on `cc (Ubuntu 15.2.0-16ubuntu1) 15.2.0`, Linux, flang 0.6.2, `./ярлык всё`
exiting 0. A caveat, so that "zero warnings" is not a lie: the linker prints its
own housekeeping line three times per run — `lto-wrapper: warning: using serial
compilation of 2 LTRANS jobs`. That is about how GCC parallelises LTO, not about
our code; there is not one diagnostic about the sources themselves, and
`-Werror` would not have let one through. Each module also emits a runner,
`flang_cli`; the request is JSON on stdin:

```sh
echo '{"fn":"Байты","args":[{"n":"2576980378"}]}' | out-c/format/flang_cli
{"ok":true,"value":{"s":"2.4 ГиБ"}}
```

---

## Speed, and why a production build must be emitted without postconditions

A 200×50 frame — fifty coloured lines of 200 cells, which is what digitdisk
repaints on every tick, every keypress and every SIGWINCH:

| What | Nanoseconds per frame | Against Go |
|---|---:|---:|
| Go (`clip` plus window slicing) | 125,903 … 143,641 | 1× |
| flang emitted to Go, **with postconditions** | 53,548,575 … 56,106,170 | ≈400× |
| flang emitted to Go, **without postconditions** | 10,085,322 … 10,746,541 | ≈75× |

Three runs each, `go test -bench -benchmem`, AMD EPYC 7742, 1 September 2026.
Hence the number this measurement existed for: **postconditions cost 5.3×**
(55.9 ms against 10.1 ms). Not the 7× the earlier reconnaissance guessed, but
not 1.5× either.

**This one could not be carried into the language, and it is not pretending
it was.** A benchmark needs a stopwatch and a rival implementation; flang has
neither — `пример` compares a value, not a duration, and there is nothing in
the language that would call Go. `tools/sverka/bench_test.go` is therefore the
one file here whose job is simply *finished*: the table above is a dated
measurement, not a running check, and re-taking it means fetching that file out
of `6bb627cc45966bfc24cc2680e1eec2196fb2a43d` (see the diff section) and
running `run.sh --замер` again. Anyone who changes `«Обрезать по ячейкам»` or
`«Кадр»` and cares about the 5.3× should do exactly that.

So: **postconditions are written, and a working build is emitted without
them.** 10 ms a frame is 99 frames a second; 56 ms is 18 frames a second and a
lag you can feel with your hand. A claim has to be written down — otherwise
nobody ever checks it — but recomputing it on every repaint in production buys
nothing: checking lives in `flang check` and `flang test`, not in the hot loop.

The emit flag that drops postconditions is being written in the language right
now. Until it lands, the number above was obtained honestly: the postconditions
were stripped from a copy of the source with one line
(`grep -v '^  обеспечивает '`), the copy was checked by the same `flang check`
and emitted by the same `flang emit`.

---

## Falsifiers: what was tested and what came out

### F1 — rounding. **FIRED. Found, fixed, re-measured.**

`к строке` prints by ECMAScript rules, while Go's `%.1f` rounds the *exact*
binary value. The naive `«Округлить к чётному» от (значение умножить на 10)`
diverged from the reference:

* percentages — **331 divergences out of 25,014**; the first: input `0.0005`,
  reference `0.1%`, ours `0.0%`;
* measurement duration — **115 out of 28,572**; the first: 1050 ms, reference
  `1.1 с`, ours `1.0 с`.

The cause is named with a number, not guessed: `0.05 * 10` in double precision
is *exactly* `0.5`, the tie goes to even — whereas exact `0.05` is
0.050000000000000002775…, and it crosses the midpoint. The cure is Dekker
splitting: the exact error of a product is computed with additions and
multiplications alone (`«Погрешность произведения»`), and the sign of that
error decides the tie. After the fix: **0 divergences on the same 53,586
inputs**.

The bound of the fix is named too: the 2²⁷+1 multiplier overflows double
precision around 10²⁹². No share, size or duration lives there.

**And one divergence from the reference that is *reproduced*, not fixed.**
`report.Bytes` picks its unit by integer division but prints a rounded
fraction, so 1,048,575 bytes come out as `1024.0 КиБ`, not `1.0 МиБ`. The same
happens here, and there is an example pinning it.

### F2 — "tabs and scrolling are not pure". **DID NOT FIRE.**

Both are pure; the arguments they were missing are named and passed in:

* tabs need *one* number — how many tabs there are. Width, height and SIGWINCH
  have nothing to do with choosing a tab;
* scrolling needs *two* — body height and content length. The host measures
  both. Body height is its own function, `«Высота тела» от строк и обвязка`
  (the reference's chrome is 5 rows), and it is diffed: 201 inputs, 0
  divergences.

### F3 — the cumulative step limit. **DID NOT FIRE — but a different limit does.**

The built-in step limit is 1,000,000. A frame does not come close: cost grows
linearly in the number of *lines*, not cells.

| Work | Steps | Is 1,000,000 enough |
|---|---:|---|
| Frame 80×24 | 25 | yes |
| Frame 200×50 | 51 | yes |
| Frame 400×100 | 101 | yes |
| Frame 800×200 | 201 | yes |
| Graph of 4000 samples, width 200 | 3,802 | yes |
| Push a sample into a history of 4000 | 2 | yes |

Those step counts came out of `tools/sverka/steps_test.go`, and they are the
half of it that could **not** be carried over: bisecting a step budget needs
`ctx.MaxSteps` on a list-taking function, and `flang run --args` takes a flat
object of scalars only — `«Кадр»` wants a list of lines. The table stays as a
dated measurement (2 September 2026, the numbers reproduced unchanged).

What does bind is the **call-depth limit of 10,000**, and it binds exactly
where recursion walks cells: a bar of width 9,999 computes, a bar of width
10,000 does not, and the binary says so in words:

```
функция «Клетки полосы» превысила предел глубины вызовов (10000) на глубине 10001
```

**That half did carry over, and without Go.** `./ярлык пределы` emits «Screen»
to C, builds it, and asks the *emitted* runner for a bar of 9,999 cells and one
of 10,000: the first must compute, the second must answer
`FLANG_RECURSION_LIMIT`, or the shortcut exits 1. It runs in CI on every push.
The number it pins is the emitted program's `FL_MAX_DEPTH`, which is what the
paragraph above is about — the binary's own limit is 20,000 and binds two cells
earlier (9,998 computes, 9,999 does not, under `flang run --max-depth 10000`).
Two runtimes, two frames of accounting; the difference is one call, and it is
named here rather than glossed.

What a host should do, as a number: nobody draws a 10,000-cell bar in a
terminal, and if they must, the limit is raised at emit time
(`flang emit --max-depth N`) or on the emitted program's context
(`ctx.MaxDepth`). With the limit raised, the same widths compute.

---

## What is not here, said out loud

* **The reference has no digit grouping and no decimal comma.** Go's `fmt` is
  locale-blind: digitdisk writes `2.4 ГиБ` with a dot, always, and never groups
  digits. So the functions come in two ranks and are never mixed: `«Байты»`,
  `«Проценты»`, `«Длительность»`, `«Прошло»`, `«Заняло»` match the reference
  exactly; `«Байты знаком»`, `«Проценты знаком»`, `«Разрядами»`,
  `«Байты десятично»` are ours, and there is nothing to diff them against. That
  is why they are absent from the diff table above: only what has a reference
  appears there.
* **digitdisk has no Home/End on tabs** — those sequences are not decoded at
  all. The `«Первая»` and `«Последняя»` variants are our addition, with no
  reference. Scroll's `«В начало»` does have one (`g`); `«В конец»` does not.
* **Surrogate pairs.** `длина` counts code points, and so does Go's `[]rune`;
  above the basic plane they agree, but no diff was run on such strings.
* **Character width** is one cell per code point, exactly as in the reference.
  Neither there nor here is there CJK double width or emoji clustering.

---

## Where the "emits to Go" boundary comes from

Everything under `flang/` is total functions with no orders at all, so it emits
to all eight targets. The sibling repository
[flang-env](https://github.com/digitable-lol/flang-env) is built the same way
but also carries a `план` that reads the environment itself; such a plan emits
to exactly one target and is executed by `flang io`. Its README quotes the
compiler's refusal verbatim and works the boundary out.

## Licence

BSD-2-Clause (`LICENSE`; a Russian translation in `LICENSE-RU.md`). The tree is
written from scratch; behaviour is diffed against digitdisk through its open
source, but not one line of anyone else's code is here.
