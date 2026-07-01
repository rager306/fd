# M051-h1xr44: fd lifecycle warmup auto-recovery

**Vision:** fd-api стартовый warmup сегодня делает 3 попытки с окном ~6с. Если TEI ещё грузит BERT на CPU (~15-20с), все попытки проваливаются с `connection refused`, lifecycle залипает в `degraded` без авто-восстановления — оператор обязан вручную `docker restart fd_api`. Milestone устраняет корневую причину: расширяет стартовое окно и добавляет periodic background recovery, пока TEI снова не станет достижимым. Цель: после `docker compose up/restart` fd-api самостоятельно приходит в `ok` без ручного вмешательства во всех сценариях, где TEI становится здоровым в течение разумного времени.

## Slices

- [x] **S01: Warmup recovery contract** `risk:medium` `depends:[]`
  > After this: После `docker compose down && up` с задержанным TEI fd-api приходит в `/health` 200 без ручного restart

## Boundary Map

Not provided.
