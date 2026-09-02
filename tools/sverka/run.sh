#!/bin/sh
# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
#
# Сверка flang-tui с эталоном digitdisk на сетке входов.
#
# ЗАЧЕМ ОТДЕЛЬНЫЙ СКРИПТ. Сверщику нужно дотянуться до незаэкспортированных
# `percent`, `clip`, `spark`, `handle`, `dash`, `cut`, `pct`, `since` — а в Go
# это возможно только ИЗНУТРИ их пакета. Поэтому файлы отсюда кладутся во
# временный клон digitdisk, туда же через `go.mod` подставляются напечатанные
# в Go модули, и `go test` идёт там. В дерево digitdisk ничего не коммитится:
# клон одноразовый и лежит вне обоих репозиториев.
#
#   ./tools/sverka/run.sh              сверка: входы и расхождения по каждому куску
#   ./tools/sverka/run.sh --замер      плюс замер кадра 200×50 с обещаниями и без
#   ./tools/sverka/run.sh --пределы    плюс цена работ в шагах и границы пределов
#
# Нужны: flang на PATH (или FLANG=путь), go, git, сеть для клона digitdisk.
# Клон берётся по DIGITDISK_REF (по умолчанию — закреплённый ниже отпечаток),
# чтобы числа в README сверялись с тем же деревом, на котором их сняли.
set -eu

FLANG="${FLANG:-flang}"
DIGITDISK_REF="${DIGITDISK_REF:-7ea03ed}"
DIGITDISK_URL="${DIGITDISK_URL:-https://github.com/digitable-lol/digitdisk.git}"

ROOT=$(cd "$(dirname "$0")/../.." && pwd)
WORK="${WORK:-${TMPDIR:-/tmp}/flang-tui-sverka}"
BENCH=нет
LIMITS=нет
for arg in "$@"; do
  case "$arg" in
    --замер|--bench) BENCH=да ;;
    --пределы|--limits) LIMITS=да ;;
    *) echo "неизвестный ключ: $arg" >&2; exit 2 ;;
  esac
done

echo "== клон digitdisk $DIGITDISK_REF в $WORK"
rm -rf "$WORK"
mkdir -p "$WORK"
git clone --quiet "$DIGITDISK_URL" "$WORK/digitdisk"
git -C "$WORK/digitdisk" checkout --quiet "$DIGITDISK_REF"
git -C "$WORK/digitdisk" --no-pager log --oneline -1

echo "== печать модулей в Go"
# Имя модуля Go у всех целей одно; здесь оно переписывается, чтобы шесть
# библиотек и эталон ужились в одной сборке. Ровно то же делает `make печать`.
emit_module() {
  # Печать пишет пояснения в stderr; они не беда, поэтому копятся в журнал и
  # показываются только если печать отказала.
  "$FLANG" emit "$ROOT/flang/$1.flang" --target go --out "$WORK/go/$2" >/dev/null 2>>"$WORK/emit.log" \
    || { cat "$WORK/emit.log" >&2; exit 1; }
  grep -rl flangprogram "$WORK/go/$2" | while read -r f; do
    sed -i.bak "s|flangprogram|flang$2|g" "$f" && rm -f "$f.bak"
  done
}
emit_module screen  screen
emit_module colour  colour
emit_module format  format
emit_module history history
emit_module tabs    tabs
emit_module scroll  scroll

# Тот же экран, но без постусловий: цену обещаний иначе не назвать числом.
if [ "$BENCH" = да ] || [ "$LIMITS" = да ]; then
  grep -v '^  обеспечивает ' "$ROOT/flang/screen.flang" > "$WORK/screen-bare.flang"
  "$FLANG" check "$WORK/screen-bare.flang" >/dev/null
  "$FLANG" emit "$WORK/screen-bare.flang" --target go --out "$WORK/go/screenbare" >/dev/null 2>>"$WORK/emit.log" \
    || { cat "$WORK/emit.log" >&2; exit 1; }
  grep -rl flangprogram "$WORK/go/screenbare" | while read -r f; do
    sed -i.bak "s|flangprogram|flangscreenbare|g" "$f" && rm -f "$f.bak"
  done
fi

echo "== подстановка модулей и сверщика в клон"
H="$WORK/digitdisk/host"
{
  echo
  echo "// СВЕРКА (одноразовый клон): библиотеки flang-tui."
  for m in screen colour format history tabs scroll; do
    echo "require flang$m v0.0.0"
    echo "replace flang$m => $WORK/go/$m"
  done
  if [ "$BENCH" = да ] || [ "$LIMITS" = да ]; then
    echo "require flangscreenbare v0.0.0"
    echo "replace flangscreenbare => $WORK/go/screenbare"
  fi
} >> "$H/go.mod"

cp "$ROOT/tools/sverka/report_test.go" "$H/internal/report/zzsverka_test.go"
cp "$ROOT/tools/sverka/ui_test.go"     "$H/internal/ui/zzsverka_test.go"
if [ "$BENCH" = да ] || [ "$LIMITS" = да ]; then
  cp "$ROOT/tools/sverka/bench_test.go" "$H/internal/ui/zzbench_test.go"
fi
if [ "$LIMITS" = да ]; then
  cp "$ROOT/tools/sverka/steps_test.go" "$H/internal/ui/zzsteps_test.go"
fi
( cd "$H" && go mod tidy >/dev/null 2>&1 )

echo
echo "== сверка"
( cd "$H" && go test -count=1 -run TestSverka -v ./internal/report/ ./internal/ui/ ) \
  | grep -E 'СВЕРКА|FAIL|^ok ' | sed 's/^ *zz[a-z_]*\.go:[0-9]*: //'

if [ "$LIMITS" = да ]; then
  echo
  echo "== цена в шагах и границы пределов"
  ( cd "$H" && go test -count=1 -run 'TestCena|TestPredel' -v ./internal/ui/ ) \
    | grep -E 'ШАГОВ|ПОЛОСА|ПРЕДЕЛ|БЕЗ ОБЕЩАНИЙ|FAIL|^ok ' | sed 's/^ *zz[a-z_]*\.go:[0-9]*: //'
fi

if [ "$BENCH" = да ]; then
  echo
  echo "== замер кадра 200×50: эталон, flang с обещаниями, flang без обещаний"
  ( cd "$H" && go test -count=3 -run '^$' -bench BenchmarkKadr -benchmem ./internal/ui/ )
fi

echo
echo "клон остался в $WORK — убрать: rm -rf $WORK"
