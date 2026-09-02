// Сверка flang-tui с эталоном digitdisk. Кладётся внутрь пакета
// digitdisk/internal/report: иначе не дотянуться до незаэкспортированных
// dash, cut, pct, since. Запускается tools/sverka/run.sh.
//
// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause
package report

import (
	"fmt"
	"testing"
	"time"

	"digitdisk/internal/sysinfo"

	fmtf "flangformat/flang"
	rt "flangformat/flangrt"
)

func num(v float64) rt.Value { return rt.Number(v) }
func str(s string) rt.Value  { return rt.Text(s) }

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
		tb.Logf("СВЕРКА %-28s входов %8d расхождений %6d  первое: %s", t.name, t.inputs, t.diffs, t.first)
	} else {
		tb.Logf("СВЕРКА %-28s входов %8d расхождений %6d", t.name, t.inputs, t.diffs)
	}
	if t.diffs != 0 {
		tb.Errorf("%s: расхождений %d из %d", t.name, t.diffs, t.inputs)
	}
}

func call1(tb testing.TB, f func(*rt.Ctx, rt.Value) (rt.Value, error), a rt.Value) rt.Value {
	v, err := f(fmtf.NewContext(), a)
	if err != nil {
		tb.Fatalf("flang: %v", err)
	}
	return v
}

func call2(tb testing.TB, f func(*rt.Ctx, rt.Value, rt.Value) (rt.Value, error), a, b rt.Value) rt.Value {
	v, err := f(fmtf.NewContext(), a, b)
	if err != nil {
		tb.Fatalf("flang: %v", err)
	}
	return v
}

// Сетка размеров: степени двойки и десятки, их окрестности, круглые числа.
func byteGrid() []int64 {
	var g []int64
	for e := 0; e <= 61; e++ {
		p := int64(1) << uint(e)
		for _, d := range []int64{-2, -1, 0, 1, 2} {
			if p+d >= 0 {
				g = append(g, p+d, -(p + d))
			}
		}
	}
	for _, p := range []int64{1, 10, 100, 1000, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e12, 1e15} {
		for _, d := range []int64{-1, 0, 1, 3, 7, 512, 1023, 1025} {
			g = append(g, p*7+d, p+d)
		}
	}
	for n := int64(0); n < 4096; n++ {
		g = append(g, n)
	}
	for n := int64(0); n < 3000; n++ {
		g = append(g, 1048576+n*337, 1073741824+n*1048573)
	}
	return g
}

func TestSverkaBytes(t *testing.T) {
	ta := &tally{name: "Байты / report.Bytes"}
	for _, n := range byteGrid() {
		ta.eq(fmt.Sprint(n), Bytes(n), call1(t, fmtf.Bayty, num(float64(n))).Str)
	}
	ta.report(t)
}

func TestSverkaPct(t *testing.T) {
	ta := &tally{name: "Процент целого / pct"}
	for a := uint64(0); a < 400; a++ {
		for _, b := range []uint64{0, 1, 3, 7, 100, 1024, 4096, 999983, 1 << 40} {
			want := fmt.Sprintf("%.1f", pct(a, b))
			gv := call2(t, fmtf.ProcentCelogo, num(float64(a)), num(float64(b)))
			ta.eq(fmt.Sprintf("%d/%d", a, b), want, fmt.Sprintf("%.1f", gv.Num))
		}
	}
	ta.report(t)
}

func TestSverkaDash(t *testing.T) {
	ta := &tally{name: "Прочерк / dash"}
	cases := []string{"", " ", "  ", "\t", "\n", " \t\n\v\f\r ", "ext4", " ext4 ", "—", " ", "", "0", "  0  "}
	for _, s := range cases {
		ta.eq(fmt.Sprintf("%q", s), dash(s), call1(t, fmtf.Procherk, str(s)).Str)
	}
	ta.report(t)
}

func TestSverkaCut(t *testing.T) {
	ta := &tally{name: "Обрезать / cut"}
	texts := []string{"", "a", "ab", "abc", "абвгд", "точка монтирования", "/dev/nvme0n1p2", "ёжик хвост", "…"}
	for _, s := range texts {
		for n := 0; n <= 24; n++ {
			ta.eq(fmt.Sprintf("%q,%d", s, n), cut(s, n), call2(t, fmtf.Obrezat, str(s), num(float64(n))).Str)
		}
	}
	ta.report(t)
}

func TestSverkaSince(t *testing.T) {
	ta := &tally{name: "Прошло / report.since"}
	for s := int64(0); s < 400000; s += 37 {
		want := since(time.Duration(s) * time.Second)
		ta.eq(fmt.Sprint(s), want, call1(t, fmtf.Proshlo, num(float64(s))).Str)
	}
	ta.report(t)
}

func TestSverkaHumanDuration(t *testing.T) {
	ta := &tally{name: "Длительность / HumanDuration"}
	for s := int64(-100); s < 3000000; s += 13 {
		ta.eq(fmt.Sprint(s), sysinfo.HumanDuration(float64(s)), call1(t, fmtf.Dlitelnost, num(float64(s))).Str)
	}
	for _, s := range []float64{0.5, 1.9, 59.999, 86399.5, 1e9} {
		ta.eq(fmt.Sprint(s), sysinfo.HumanDuration(s), call1(t, fmtf.Dlitelnost, num(s)).Str)
	}
	ta.report(t)
}
