# M044-9vahk2: Upgrade OpenAPI contract to OAS 3.2.0

**Vision:** Перевести fd `/openapi.json` и `/docs` с OpenAPI 3.1.0 на OAS 3.2.0 отдельным follow-up после M041, сохранив совместимость runtime API и обновив verifier/evidence под новый schema contract.

## Success Criteria

- `/openapi.json` returns an OAS 3.2.0 document.
- `/docs` renders or clearly supports the OAS 3.2.0 schema.
- Final contract verifier checks OAS 3.2.0, not 3.1.0.
- No runtime API behavior regresses from M041 acceptance.
- Mandatory Go gates pass after final changes.

## Slices

- [ ] **S01: OAS 3.2 delta and validation strategy** `risk:medium` `depends:[]`
  > After this: После slice есть короткий research artifact: какие поля/семантика меняются для нашего spec, каким инструментом валидировать 3.2.0, и какие текущие 3.1 assumptions надо заменить.

- [ ] **S02: Emit and verify OAS 3.2.0 spec** `risk:medium` `depends:[S01]`
  > After this: `GET /openapi.json` returns `openapi: 3.2.0`; unit tests and verifier assertions expect 3.2.0; `/docs` still renders.

- [ ] **S03: Final OAS 3.2 acceptance and closure** `risk:low` `depends:[S02]`
  > After this: Full fd contract verification passes with OAS 3.2.0, plus mandatory Go gates and milestone artifacts.

## Boundary Map

| In scope | Out of scope |
|---|---|
| OpenAPI document version/shape upgrade to OAS 3.2.0 | Changing embedding runtime, cache, auth semantics, or performance targets |
| Verifier/test/docs adjustments for OAS 3.2.0 | Rewriting M041 acceptance claims |
| Validation tooling evidence | Implementing unrelated API endpoints |
