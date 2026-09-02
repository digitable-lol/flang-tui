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
`host/internal/{report,ui,sysinfo}`. Every piece is diffed against it over a
grid of inputs; the numbers below are the output of a run, not a promise.

The source is written in flang, whose surface is Russian. File names are
English, module names are English, the prose inside is Russian — that is the
language's own convention, and this repository follows it.

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
borrowed from digitdisk. There is no Python and no JavaScript in this tree, and
that is one of its four checks.

`tools/sverka/` is the differ: a script plus Go files that get dropped into a
disposable clone of digitdisk. It is the only hand-written Go in the tree, and
it is here for exactly one reason — so that anybody, not only their author, can
re-check the numbers in this README.

## Getting it

The flang compiler is a single binary that needs only `cc`:

```sh
brew install flang        # or: asdf plugin add flang
# or from a clone of the language:  make -C bootstrap
```

Then:

```sh
make проверка   # flang check + flang test, module by module
make печать     # emit to Go and to C, build both
make лицензии   # the licence guard
```

Emission puts each module in its own `out-go/<name>` and `out-c/<name>`. All
eight emit targets name the Go module `flangprogram`, so after emission `make`
rewrites it to `flang<name>`: six libraries in one build would otherwise
collide. That is the *only* edit made to emitted code, and `sed` makes it, not
a hand.

---

## Diff against digitdisk: 494,719 inputs, 0 divergences

**Reproduced by one command:**

```sh
./tools/sverka/run.sh              # inputs and divergences, piece by piece
./tools/sverka/run.sh --пределы    # plus step costs and where the limits bind
./tools/sverka/run.sh --замер      # plus the frame benchmark, with and without postconditions
```

The script clones digitdisk itself at the pinned `7ea03ed` (0.5.0), emits the
modules to Go, drops them into its `go.mod` and puts the differ *inside* the
`internal/report` and `internal/ui` packages — there is no other way to reach
the unexported `percent`, `clip`, `spark`, `handle`. The clone is disposable and
lives outside both repositories: nothing is committed into digitdisk's tree.
Every row below is `go test -v` output from that run.

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
| Переключить (tab switch) | `ui.screen.handle` | 150 | 0 |
| Прокрутить (scroll) | `ui.screen.handle` + `frame` | 2,835 | 0 |
| Высота тела (body height) | `ui.screen.bodyHeight` | 201 | 0 |
| Куб 256 (cube) | `ui.cube256` | 140,608 | 0 |
| Цвет (colour seq) | `ui.Theme.seq` | 24,576 | 0 |
| **Total** | | **494,317** | **0** |

(A further 402 diff inputs — `detectDepth`, `UsableTERM`, `PaletteByName` —
belong to the sibling repository
[flang-env](https://github.com/digitable-lol/flang-env), and so do their
numbers. The total across both repositories is 494,719.)

**Two caveats, without which the table would be lying.**

1. `Полоса` was diffed from width 0 up: the argument's type is `нат`
   (non-negative), a negative width is *outside the contract*, and emitted code
   does not type-check its arguments — the compiler says so itself at emit
   time. The reference returns an empty string for a negative width; a host
   calling emitted code must guard the width itself.
2. `Уложить раздел` and `Прокрутить` were diffed against the four lines of
   `frame()` quoted *verbatim* inside the differ, not against `frame()` itself:
   it renders a section from a live snapshot, and the content length cannot be
   driven from a grid. The key step, however, was diffed against the real
   `handle` — that one can be called.

---

## The ledger: what actually carries each claim

Three words, and they never blur. **«доказано»** (proved) — about *all* inputs.
**«сетка N»** (grid of N) — computed on N of the author's values; that is not a
proof. **«объявлено, не доказано»** (declared, not proved) — neither theorem
nor examples; the runtime computes it on whatever inputs arrive. The numbers
are `flang check --proof` output.

| Module | Functions | All total | Claims | proved | grid | declared, not proved | Examples |
|---|---:|---|---:|---:|---:|---:|---:|
| Screen | 19 | yes | 25 | 4 | 15 | 6 | 43 |
| Colour | 9 | yes | 5 | 1 | 4 | 0 | 17 |
| Format | 33 | yes | 14 | 4 | 9 | 1 | 73 |
| History | 9 | yes | 8 | 3 | 5 | 0 | 23 |
| Tabs | 5 | yes | 8 | 0 | 6 | 2 | 13 |
| Scroll | 6 | yes | 7 | 0 | 7 | 0 | 18 |
| **Total** | **81** | **81 of 81** | **67** | **12** | **46** | **9** | **187** |

`tools/licensing.flang` is 21 functions, all total, 27 examples; no ledger is
printed for it because it declares a `план`, whose laws the binary does not
judge — and it says so itself.

Twelve proved claims is 18% of sixty-seven. The rest is grid and declared. That
is written down as it is; nobody here will call a grid a proof.

---

## Emission: eight targets, two verified

| Target | Modules emitted | Check |
|---|---:|---|
| Go | 6 of 6 | `gofmt -l` empty, `go vet ./...` clean, `go build ./...` clean |
| C | 6 of 6 | `cc -std=c99 -Wall -Wextra -Werror -pedantic -O2 -flto`, **0 warnings about our code** |

Run on `cc (Ubuntu 15.2.0-16ubuntu1) 15.2.0`, Linux, flang 0.6.2, `make all`
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

Three runs each, `go test -bench -benchmem`, AMD EPYC 7742. Hence the number
this measurement existed for: **postconditions cost 5.3×** (55.9 ms against
10.1 ms). Not the 7× the earlier reconnaissance guessed, but not 1.5× either.

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

What does bind is the **call-depth limit of 10,000**, and it binds exactly
where recursion walks cells: a bar of width 9,999 computes, a bar of width
10,000 does not, and the binary says so in words:

```
функция «Клетки полосы» превысила предел глубины вызовов (10000) на глубине 10001
```

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
  `«Байты десятично»` are ours, and there is nothing to diff them against. The
  ledger lists them as undiffed.
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
