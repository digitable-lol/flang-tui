// Замер flang-tui против эталона: кадр 200×50 у digitdisk, у flang с
// обещаниями и без. Кладётся внутрь пакета digitdisk/internal/ui.
// Запускается tools/sverka/run.sh --замер.
//
// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause
package ui

import (
	"fmt"
	"strings"
	"testing"

	fs "flangscreen/flang"
	fb "flangscreenbare/flang"

	rtb "flangscreenbare/flangrt"
	rts "flangscreen/flangrt"
)

func benchBody(n int) []string {
	body := make([]string, n)
	for i := range body {
		body[i] = fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m %s", 30+i%200,
			strings.Repeat("ячейка ", 20), fmt.Sprintf("строка %d", i+1))
	}
	return body
}

func textsB(ss []string) rtb.Value {
	out := make([]rtb.Value, len(ss))
	for i, s := range ss {
		out[i] = rtb.Text(s)
	}
	return rtb.List(out)
}

// Эталон: то же, что делает screen.frame на каждой перерисовке — уложить
// раздел в высоту и обрезать каждую строку по ширине.
func goFrame(th Theme, body []string, w, h, scroll int) []string {
	win := refWindow(body, h, scroll)
	out := make([]string, len(win))
	for i, s := range win {
		out[i] = th.clip(s, w)
	}
	return out
}

func BenchmarkKadrGo200x50(b *testing.B) {
	th := Theme{P: Carbon, d: depth256}
	body := benchBody(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = goFrame(th, body, 200, 50, 0)
	}
}

func BenchmarkKadrFlang200x50(b *testing.B) {
	th := Theme{P: Carbon, d: depth256}
	ctx := fs.NewContext()
	body := texts(benchBody(50))
	reset := rts.Text(th.reset())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fs.Kadr(ctx, body, rts.Number(200), rts.Number(50), rts.Number(0), rts.Text("…"), reset); err != nil {
			b.Fatalf("flang: %v", err)
		}
	}
}

func BenchmarkKadrFlangBare200x50(b *testing.B) {
	th := Theme{P: Carbon, d: depth256}
	ctx := fb.NewContext()
	body := textsB(benchBody(50))
	reset := rtb.Text(th.reset())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fb.Kadr(ctx, body, rtb.Number(200), rtb.Number(50), rtb.Number(0), rtb.Text("…"), reset); err != nil {
			b.Fatalf("flang: %v", err)
		}
	}
}

// Ф3: совокупный предел шагов. Кадр 200×50 считается при вшитом пределе
// 1 000 000 или падает на нём? Ответ печатается числом, а не мнением.
func TestPredelShagov(t *testing.T) {
	th := Theme{P: Carbon, d: depth256}
	body := texts(benchBody(50))
	reset := rts.Text(th.reset())
	for _, steps := range []int{1000000, 4000000, 16000000, 64000000} {
		ctx := fs.NewContext()
		ctx.MaxSteps = steps
		_, err := fs.Kadr(ctx, body, rts.Number(200), rts.Number(50), rts.Number(0), rts.Text("…"), reset)
		if err != nil {
			t.Logf("ПРЕДЕЛ %9d шагов: %v", steps, err)
		} else {
			t.Logf("ПРЕДЕЛ %9d шагов: кадр посчитан", steps)
		}
	}
	// Тот же кадр без обещаний: сколько шагов стоит он.
	bodyB := textsB(benchBody(50))
	for _, steps := range []int{1000000, 4000000, 16000000, 64000000} {
		ctx := fb.NewContext()
		ctx.MaxSteps = steps
		_, err := fb.Kadr(ctx, bodyB, rtb.Number(200), rtb.Number(50), rtb.Number(0), rtb.Text("…"), rtb.Text(th.reset()))
		if err != nil {
			t.Logf("БЕЗ ОБЕЩАНИЙ %9d шагов: %v", steps, err)
		} else {
			t.Logf("БЕЗ ОБЕЩАНИЙ %9d шагов: кадр посчитан", steps)
		}
	}
	// И длинная история графика — второй подозреваемый Ф3.
	t.Log("история графика проверяется отдельно в TestPredelIstorii")
}
