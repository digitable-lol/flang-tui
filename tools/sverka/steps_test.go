// Ф3: во сколько ШАГОВ обходится каждая работа и где кончаются вшитые пределы
// напечатанной программы. Кладётся внутрь пакета digitdisk/internal/ui.
// Запускается tools/sverka/run.sh --пределы.
//
// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause
package ui

import (
	"fmt"
	"testing"

	fh "flanghistory/flang"
	fs "flangscreen/flang"
	fb "flangscreenbare/flang"

	rth "flanghistory/flangrt"
	rtb "flangscreenbare/flangrt"
	rts "flangscreen/flangrt"
)

// цена ищет наименьший предел шагов, при котором работа доходит до конца.
func cena(t *testing.T, name string, run func(steps int) error) int {
	lo, hi := 1, 1
	for {
		if err := run(hi); err == nil {
			break
		}
		hi *= 2
		if hi > 1<<31 {
			t.Fatalf("%s: не уложилось и в 2^31 шагов", name)
		}
	}
	for lo < hi {
		mid := (lo + hi) / 2
		if run(mid) == nil {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	t.Logf("ШАГОВ %-44s %10d  (вшитый предел 1000000: %s)", name, lo,
		map[bool]string{true: "хватает", false: "НЕ ХВАТАЕТ"}[lo <= 1000000])
	return lo
}

func TestCenaVShagah(t *testing.T) {
	th := Theme{P: Carbon, d: depth256}
	reset := th.reset()

	for _, sz := range [][2]int{{80, 24}, {200, 50}, {400, 100}, {800, 200}} {
		w, h := sz[0], sz[1]
		body := texts(benchBody(h))
		bodyB := textsB(benchBody(h))
		cena(t, fmt.Sprintf("Кадр %dx%d с обещаниями", w, h), func(steps int) error {
			ctx := fs.NewContext()
			ctx.MaxSteps = steps
			_, err := fs.Kadr(ctx, body, rts.Number(float64(w)), rts.Number(float64(h)), rts.Number(0), rts.Text("…"), rts.Text(reset))
			return err
		})
		cena(t, fmt.Sprintf("Кадр %dx%d без обещаний", w, h), func(steps int) error {
			ctx := fb.NewContext()
			ctx.MaxSteps = steps
			_, err := fb.Kadr(ctx, bodyB, rtb.Number(float64(w)), rtb.Number(float64(h)), rtb.Number(0), rtb.Text("…"), rtb.Text(reset))
			return err
		})
	}

	for _, n := range []int{60, 240, 1000, 4000} {
		hist := make([]rth.Value, n)
		for i := range hist {
			hist[i] = rth.Number(float64(i%101) / 100)
		}
		list := rth.List(hist)
		cena(t, fmt.Sprintf("График истории из %d замеров, ширина 200", n), func(steps int) error {
			ctx := fh.NewContext()
			ctx.MaxSteps = steps
			_, err := fh.Grafik(ctx, list, rth.Number(200))
			return err
		})
		cena(t, fmt.Sprintf("Дописать замер в историю из %d", n), func(steps int) error {
			ctx := fh.NewContext()
			ctx.MaxSteps = steps
			_, err := fh.DopisatZamer(ctx, list, rth.Number(0.5), rth.Number(float64(n)))
			return err
		})
	}
}

func TestCenaPolosy(t *testing.T) {
	for _, w := range []int{2000, 9997, 9998, 9999, 10000, 20000} {
		ctx := fs.NewContext()
		_, err := fs.Polosa(ctx, rts.Number(0.5), rts.Number(float64(w)), rts.Text("#"), rts.Text("-"))
		if err != nil {
			t.Logf("ПОЛОСА ширина %7d при вшитых пределах: %v", w, err)
		} else {
			t.Logf("ПОЛОСА ширина %7d при вшитых пределах: посчитана", w)
		}
		ctx2 := fs.NewContext()
		ctx2.MaxDepth = 4000000
		ctx2.MaxSteps = 100000000
		_, err2 := fs.Polosa(ctx2, rts.Number(0.5), rts.Number(float64(w)), rts.Text("#"), rts.Text("-"))
		if err2 != nil {
			t.Logf("ПОЛОСА ширина %7d при поднятых пределах: %v", w, err2)
		} else {
			t.Logf("ПОЛОСА ширина %7d при поднятых пределах: посчитана", w)
		}
	}
}
