// Сверка flang-tui с эталоном digitdisk. Кладётся внутрь пакета
// digitdisk/internal/ui: иначе не дотянуться до незаэкспортированных
// percent, clip, spark, handle. Запускается tools/sverka/run.sh.
//
// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause
package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	fc "flangcolour/flang"
	ff "flangformat/flang"
	fh "flanghistory/flang"
	fs "flangscreen/flang"
	fsc "flangscroll/flang"
	ft "flangtabs/flang"

	rtc "flangcolour/flangrt"
	rtf "flangformat/flangrt"
	rth "flanghistory/flangrt"
	rts "flangscreen/flangrt"
	rtsc "flangscroll/flangrt"
	rtt "flangtabs/flangrt"
)

type tally struct {
	name   string
	inputs int
	diffs  int
	first  string
}

func (t *tally) eq(in string, want, got string) {
	t.inputs++
	if want != got {
		t.diffs++
		if t.first == "" {
			t.first = fmt.Sprintf("вход %s: эталон %q, flang %q", in, want, got)
		}
	}
}

func (t *tally) report(tb testing.TB) {
	if t.first != "" {
		tb.Logf("СВЕРКА %-34s входов %8d расхождений %6d  первое: %s", t.name, t.inputs, t.diffs, t.first)
	} else {
		tb.Logf("СВЕРКА %-34s входов %8d расхождений %6d", t.name, t.inputs, t.diffs)
	}
}

func (t *tally) strict(tb testing.TB) {
	t.report(tb)
	if t.diffs != 0 {
		tb.Errorf("%s: расхождений %d из %d", t.name, t.diffs, t.inputs)
	}
}

func plainTheme() Theme { return Theme{P: Carbon, d: depthNone} }

func texts(ss []string) rts.Value {
	out := make([]rts.Value, len(ss))
	for i, s := range ss {
		out[i] = rts.Text(s)
	}
	return rts.List(out)
}

func nums(xs []float64) rth.Value {
	out := make([]rth.Value, len(xs))
	for i, x := range xs {
		out[i] = rth.Number(x)
	}
	return rth.List(out)
}

// ───────────────────────────── Format ─────────────────────────────

// Ф1 в лоб: %.1f на плотной сетке дробей, где округление ECMAScript и
// округление Go расходятся охотнее всего.
func TestSverkaPercent(t *testing.T) {
	ta := &tally{name: "Проценты / ui.percent"}
	ctx := ff.NewContext()
	check := func(frac float64) {
		want := percent(frac)
		v, err := ff.Procenty(ctx, rtf.Number(frac))
		if err != nil {
			t.Fatalf("flang: %v", err)
		}
		ta.eq(fmt.Sprint(frac), want, v.Str)
	}
	for i := 0; i <= 20000; i++ {
		check(float64(i) / 20000)
	}
	for i := 0; i <= 4000; i++ {
		check(float64(i) / 800)
	}
	for i := -500; i <= 500; i++ {
		check(float64(i) / 1000)
	}
	for _, f := range []float64{0.0005, 0.0015, 0.0025, 0.0035, 0.0045, 0.00125, 0.00375, 0.125, 0.375, 0.625, 0.875} {
		check(f)
	}
	ta.report(t)
}

func TestSverkaPctOf(t *testing.T) {
	ta := &tally{name: "Доля от / ui.pctOf"}
	ctx := ff.NewContext()
	for a := uint64(0); a < 300; a++ {
		for _, b := range []uint64{0, 1, 3, 7, 100, 1024, 999983, 1 << 40} {
			v, err := ff.DolyaOt(ctx, rtf.Number(float64(a)), rtf.Number(float64(b)))
			if err != nil {
				t.Fatalf("flang: %v", err)
			}
			ta.eq(fmt.Sprintf("%d/%d", a, b), fmt.Sprintf("%.17g", pctOf(a, b)), fmt.Sprintf("%.17g", v.Num))
		}
	}
	ta.strict(t)
}

func TestSverkaLasted(t *testing.T) {
	ta := &tally{name: "Заняло / ui.lasted"}
	ctx := ff.NewContext()
	for ms := int64(0); ms < 200000; ms += 7 {
		want := lasted(time.Duration(ms) * time.Millisecond)
		v, err := ff.Zanyalo(ctx, rtf.Number(float64(ms)))
		if err != nil {
			t.Fatalf("flang: %v", err)
		}
		ta.eq(fmt.Sprint(ms), want, v.Str)
	}
	ta.report(t)
}

func TestSverkaLevel(t *testing.T) {
	ta := &tally{name: "Уровень доли / ui.level"}
	ctx := fh.NewContext()
	th := plainTheme()
	nameOf := func(s slot) string {
		switch s {
		case th.P.Green:
			return "Спокойно"
		case th.P.Yellow:
			return "Внимание"
		default:
			return "Тревога"
		}
	}
	for i := -100; i <= 1200; i++ {
		f := float64(i) / 1000
		v, err := fh.UrovenDoli(ctx, rth.Number(f))
		if err != nil {
			t.Fatalf("flang: %v", err)
		}
		ta.eq(fmt.Sprint(f), nameOf(th.level(f)), v.Str)
	}
	ta.strict(t)
}

// ───────────────────────────── Screen ─────────────────────────────

func TestSverkaBar(t *testing.T) {
	ta := &tally{name: "Полоса / Theme.bar"}
	ctx := fs.NewContext()
	th := plainTheme()
	for n := 0; n <= 200; n++ {
		for _, f := range []float64{-1, -0.001, 0, 0.001, 0.004, 0.01, 0.125, 0.25, 0.333, 0.4444, 0.5, 0.5001, 0.6667, 0.75, 0.9, 0.99, 0.999, 1, 1.001, 4} {
			v, err := fs.Polosa(ctx, rts.Number(f), rts.Number(float64(n)), rts.Text("█"), rts.Text("─"))
			if err != nil {
				t.Fatalf("flang: %v", err)
			}
			ta.eq(fmt.Sprintf("%v,%d", f, n), th.bar(f, n), v.Str)
		}
	}
	ta.strict(t)
}

func clipCorpus() []string {
	return []string{
		"", "a", "ab", "абвгд", "точка монтирования",
		"\x1b[31mкрасное\x1b[0m", "\x1b[38;2;10;20;30mистинный\x1b[0m",
		"\x1b]0;заголовок\x07текст", "\x1b[?25lскрытый курсор",
		"\x1b[1m\x1b[4mдва подряд\x1b[0m", "хвост\x1b", "обрыв\x1b[31",
		strings.Repeat("длинная строка ", 12),
		"/dev/nvme0n1p2 на / ext4 rw,relatime",
	}
}

func TestSverkaClip(t *testing.T) {
	ta := &tally{name: "Обрезать по ячейкам / Theme.clip"}
	ctx := fs.NewContext()
	th := plainTheme()
	for _, s := range clipCorpus() {
		for n := -3; n <= 60; n++ { // отрицательный предел законен: тип «число»
			v, err := fs.ObrezatPoYacheykam(ctx, rts.Text(s), rts.Number(float64(n)), rts.Text("…"), rts.Text(th.reset()))
			if err != nil {
				t.Fatalf("flang на входе %q,%d: %v", s, n, err)
			}
			ta.eq(fmt.Sprintf("%q,%d", s, n), th.clip(s, n), v.Str)
		}
	}
	ta.strict(t)
}

func TestSverkaPlainWidth(t *testing.T) {
	ta := &tally{name: "Ширина без последовательностей"}
	ctx := fs.NewContext()
	for _, s := range clipCorpus() {
		v, err := fs.ShirinaBezPosledovatelnostey(ctx, rts.Text(s))
		if err != nil {
			t.Fatalf("flang: %v", err)
		}
		ta.eq(fmt.Sprintf("%q", s), fmt.Sprint(plainWidth(s)), fmt.Sprint(int(v.Num)))
	}
	ta.strict(t)
}

// refWindow — четыре строки frame() из screen.go, приведённые дословно.
// Позвать сам frame() нельзя: он рисует раздел по живому снимку, и длину
// содержимого сеткой не задать.
func refWindow(body []string, h, scroll int) []string {
	if scroll > len(body)-1 {
		scroll = len(body) - 1
	}
	if scroll < 0 {
		scroll = 0
	}
	shown := body[scroll:]
	if len(shown) > h {
		shown = shown[:h]
	}
	out := make([]string, 0, h)
	out = append(out, shown...)
	for i := len(shown); i < h; i++ {
		out = append(out, "")
	}
	return out
}

func TestSverkaWindow(t *testing.T) {
	ta := &tally{name: "Уложить раздел / frame"}
	ctx := fs.NewContext()
	for _, n := range []int{0, 1, 2, 3, 5, 9, 24, 50, 137} {
		body := make([]string, n)
		for i := range body {
			body[i] = fmt.Sprintf("строка %d", i+1)
		}
		for h := 1; h <= 60; h++ {
			for _, sc := range []int{-9, -1, 0, 1, 2, 5, 23, 49, 136, 137, 500} {
				v, err := fs.UlozhitRazdel(ctx, texts(body), rts.Number(float64(h)), rts.Number(float64(sc)))
				if err != nil {
					t.Fatalf("flang: %v", err)
				}
				got := make([]string, len(v.List))
				for i, x := range v.List {
					got[i] = x.Str
				}
				ta.eq(fmt.Sprintf("%d,%d,%d", n, h, sc), strings.Join(refWindow(body, h, sc), "\x00"), strings.Join(got, "\x00"))
			}
		}
	}
	ta.strict(t)
}

// ───────────────────────────── History ─────────────────────────────

func TestSverkaPush(t *testing.T) {
	ta := &tally{name: "Дописать замер / push"}
	ctx := fh.NewContext()
	hist := []float64{}
	for i := 0; i < 600; i++ {
		v := float64(i%101) / 100
		if i%17 == 0 {
			v = -1
		}
		hist = push(hist, v)
		fv, err := fh.DopisatZamer(ctx, nums(prevOf(hist, v)), rth.Number(v), rth.Number(float64(histLen)))
		if err != nil {
			t.Fatalf("flang: %v", err)
		}
		got := make([]string, len(fv.List))
		for j, x := range fv.List {
			got[j] = fmt.Sprint(x.Num)
		}
		want := make([]string, len(hist))
		for j, x := range hist {
			want[j] = fmt.Sprint(x)
		}
		ta.eq(fmt.Sprint(i), strings.Join(want, ","), strings.Join(got, ","))
	}
	ta.strict(t)
}

// prevOf восстанавливает историю ДО последнего push: у эталона push чистый,
// и вход у обеих сторон обязан быть один.
func prevOf(after []float64, last float64) []float64 {
	if len(after) == 0 {
		return nil
	}
	if len(after) < histLen {
		return after[:len(after)-1]
	}
	return after[:len(after)-1]
}

func TestSverkaSpark(t *testing.T) {
	ta := &tally{name: "График / Theme.spark"}
	ctx := fh.NewContext()
	th := plainTheme()
	histories := [][]float64{
		{},
		{0},
		{1},
		{-1},
		{-1, -1, -1},
		{0, 0.5, 1},
		{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		{1.5, -0.5, 0.75, 0.7499999, 0.142857, 0.1428572},
	}
	long := make([]float64, 300)
	for i := range long {
		long[i] = float64(i%101) / 100
	}
	histories = append(histories, long)
	for _, h := range histories {
		for n := -2; n <= 90; n++ {
			v, err := fh.Grafik(ctx, nums(h), rth.Number(float64(n)))
			if err != nil {
				t.Fatalf("flang: %v", err)
			}
			ta.eq(fmt.Sprintf("len=%d,n=%d", len(h), n), th.spark(h, n), v.Str)
		}
	}
	ta.strict(t)
}

// ───────────────────────────── Tabs ─────────────────────────────

func TestSverkaTabs(t *testing.T) {
	ta := &tally{name: "Переключить / screen.handle"}
	ctx := ft.NewContext()
	type kase struct {
		k    key
		mk   func() rtt.Value
		name string
	}
	kases := []kase{
		{key{kind: keyRight}, ft.VariantSleduyuschaya, "вправо"},
		{key{kind: keyTab}, ft.VariantSleduyuschaya, "tab"},
		{key{kind: keyLeft}, ft.VariantPredyduschaya, "влево"},
		{key{kind: keyShiftTab}, ft.VariantPredyduschaya, "shift-tab"},
		{key{kind: keyDown}, ft.VariantMimo, "вниз"},
		{key{kind: keyPgUp}, ft.VariantMimo, "страницей вверх"},
	}
	for d := 1; d <= 9; d++ {
		r := rune('0' + d)
		kases = append(kases, kase{key{kind: keyRune, r: r}, func(n int) func() rtt.Value {
			return func() rtt.Value { return ft.VariantNomerom(rtt.Number(float64(n))) }
		}(d), fmt.Sprintf("цифра %d", d)})
	}
	for tab := 0; tab < len(sections); tab++ {
		for _, c := range kases {
			s := &screen{tab: tab, rows: 40}
			s.handle(c.k, nil)
			v, err := ft.Pereklyuchit(ctx, ft.SozdatVkladki(rtt.Number(float64(tab)), rtt.Number(float64(len(sections)))), c.mk())
			if err != nil {
				t.Fatalf("flang: %v", err)
			}
			got := -1
			for _, f := range v.Fields {
				if f.Name == "открыта" {
					got = int(f.Value.Num)
				}
			}
			ta.eq(fmt.Sprintf("tab=%d,%s", tab, c.name), fmt.Sprint(s.tab), fmt.Sprint(got))
		}
	}
	ta.strict(t)
}

// ───────────────────────────── Scroll ─────────────────────────────

func TestSverkaScroll(t *testing.T) {
	ta := &tally{name: "Прокрутить / handle+frame"}
	ctx := fsc.NewContext()
	type kase struct {
		k    key
		mk   func() rtsc.Value
		name string
	}
	kases := []kase{
		{key{kind: keyDown}, fsc.VariantStrokoyVniz, "вниз"},
		{key{kind: keyUp}, fsc.VariantStrokoyVverh, "вверх"},
		{key{kind: keyPgDn}, fsc.VariantStraniceyVniz, "страницей вниз"},
		{key{kind: keyPgUp}, fsc.VariantStraniceyVverh, "страницей вверх"},
		{key{kind: keyRune, r: 'j'}, fsc.VariantStrokoyVniz, "j"},
		{key{kind: keyRune, r: 'k'}, fsc.VariantStrokoyVverh, "k"},
		{key{kind: keyRune, r: 'g'}, fsc.VariantVNachalo, "g"},
	}
	for _, rows := range []int{8, 10, 24, 40, 60} {
		h := rows - 5
		if h < 1 {
			h = 1
		}
		for _, length := range []int{0, 1, 2, 3, 5, 9, 24, 50, 137} {
			for _, sc := range []int{0, 1, 2, 5, 23, 49, 136, 137, 500} {
				for _, c := range kases {
					s := &screen{rows: rows, scroll: sc}
					s.handle(c.k, nil)
					// Верхний край режет frame(); его четыре строки приведены
					// дословно, потому что позвать frame() сеткой нельзя.
					if s.scroll > length-1 {
						s.scroll = length - 1
					}
					if s.scroll < 0 {
						s.scroll = 0
					}
					v, err := fsc.Prokrutit(ctx, rtsc.Number(float64(sc)), rtsc.Number(float64(h)), rtsc.Number(float64(length)), c.mk())
					if err != nil {
						t.Fatalf("flang: %v", err)
					}
					ta.eq(fmt.Sprintf("rows=%d,len=%d,sc=%d,%s", rows, length, sc, c.name), fmt.Sprint(s.scroll), fmt.Sprint(int(v.Num)))
				}
			}
		}
	}
	ta.strict(t)
}

func TestSverkaBodyHeight(t *testing.T) {
	ta := &tally{name: "Высота тела / bodyHeight"}
	ctx := fsc.NewContext()
	for rows := 0; rows <= 200; rows++ {
		s := &screen{rows: rows}
		v, err := fsc.VysotaTela(ctx, rtsc.Number(float64(rows)), rtsc.Number(5))
		if err != nil {
			t.Fatalf("flang: %v", err)
		}
		ta.eq(fmt.Sprint(rows), fmt.Sprint(s.bodyHeight()), fmt.Sprint(int(v.Num)))
	}
	ta.strict(t)
}

// ───────────────────────────── Colour ─────────────────────────────

func TestSverkaCube256(t *testing.T) {
	ta := &tally{name: "Куб 256 / cube256"}
	ctx := fc.NewContext()
	for r := 0; r < 256; r += 5 {
		for g := 0; g < 256; g += 5 {
			for b := 0; b < 256; b += 5 {
				v, err := fc.Kub256(ctx, rtc.Number(float64(r)), rtc.Number(float64(g)), rtc.Number(float64(b)))
				if err != nil {
					t.Fatalf("flang: %v", err)
				}
				ta.eq(fmt.Sprintf("%d,%d,%d", r, g, b), fmt.Sprint(cube256(RGB{uint8(r), uint8(g), uint8(b)})), fmt.Sprint(int(v.Num)))
			}
		}
	}
	ta.strict(t)
}

func TestSverkaSeq(t *testing.T) {
	ta := &tally{name: "Цвет / Theme.seq"}
	ctx := fc.NewContext()
	for r := 0; r < 256; r += 17 {
		for g := 0; g < 256; g += 17 {
			for b := 0; b < 256; b += 17 {
				sl := slot{c: RGB{uint8(r), uint8(g), uint8(b)}, ansi: 32}
				for _, bg := range []bool{false, true} {
					tt := Theme{P: Carbon, d: depthTrue}
					v, err := fc.CvetIstinnyy(ctx, rtc.Flag(bg), rtc.Number(float64(r)), rtc.Number(float64(g)), rtc.Number(float64(b)))
					if err != nil {
						t.Fatalf("flang: %v", err)
					}
					ta.eq(fmt.Sprintf("true %d,%d,%d,%v", r, g, b, bg), tt.seq(sl, bg), v.Str)

					t2 := Theme{P: Carbon, d: depth256}
					v2, err := fc.CvetIzPalitry(ctx, rtc.Flag(bg), rtc.Number(float64(r)), rtc.Number(float64(g)), rtc.Number(float64(b)))
					if err != nil {
						t.Fatalf("flang: %v", err)
					}
					ta.eq(fmt.Sprintf("256 %d,%d,%d,%v", r, g, b, bg), t2.seq(sl, bg), v2.Str)

					t3 := Theme{P: Carbon, d: depth16}
					v3, err := fc.CvetIzShestnadcati(ctx, rtc.Flag(bg), rtc.Number(32))
					if err != nil {
						t.Fatalf("flang: %v", err)
					}
					ta.eq(fmt.Sprintf("16 %v", bg), t3.seq(sl, bg), v3.Str)
				}
			}
		}
	}
	ta.strict(t)
}

