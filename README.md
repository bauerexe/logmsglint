# logmsglint

`logmsglint` — статический анализатор для Go, который проверяет лог-сообщения в `log/slog` и `go.uber.org/zap`.

## Проверяемые правила

1. `logmsg-lowercase` — сообщение должно начинаться со строчной буквы.
2. `logmsg-english` — сообщение должно быть на английском (обнаружение кириллицы считается нарушением).
3. `logmsg-nospecial` — запрещены emoji/спецсимволы и повторяющаяся пунктуация вроде `!!!` и `...`.
4. `logmsg-sensitive` — запрещены потенциально чувствительные ключевые слова (`password`, `token`, `api_key`, `secret`, `bearer` и др.) и кастомные паттерны.

## Конфигурация

По умолчанию линтер пытается прочитать `.logmsglint.yml` в текущей директории запуска.
Также можно явно передать путь через флаг анализатора `-config=/path/to/file.yml`.

Пример конфигурации:

```yaml
rules:
  lowercase: true
  english: true
  nospecial: true
  sensitive: true

sensitive:
  keywords:
    - password
    - token
  patterns:
    - '(?i)card\\s*\\d{4}-\\d{4}-\\d{4}-\\d{4}'
```

- `rules.*` позволяет включать/отключать отдельные правила.
- `sensitive.keywords` переопределяет список стандартных ключевых слов.
- `sensitive.patterns` добавляет ваши regexp-паттерны для поиска чувствительных данных.

## Авто-исправления

Анализатор возвращает `SuggestedFixes` для нарушений:
- `logmsg-lowercase` — приводит первую букву к нижнему регистру.
- `logmsg-english` — удаляет не-ASCII символы.
- `logmsg-nospecial` — убирает control/symbol символы и схлопывает повторяющуюся пунктуацию.
- `logmsg-sensitive` — редактирует найденные секреты в `[REDACTED]`.

## Что считается лог-вызовом

- `slog.Info/Warn/Error/Debug("...")`
- `(*slog.Logger).Info/Warn/Error/Debug("...")`
- `(*zap.Logger).Info/Warn/Error/Debug("...")`
- `(*zap.SugaredLogger).Info/Warn/Error/Debug("...")`

Проверяется только первый аргумент, если это строковый литерал.
Если сообщение не литерал (переменная, `fmt.Sprintf`, конкатенация), вызов пропускается.

## Установка и запуск

Инструкция по установке (скачиванию) CLI:

```bash
go install github.com/bauerex/logmsglint/cmd/logmsglint@latest
```


Запуск из исходников:

```bash
go run ./cmd/logmsglint ./...
```

## Подключение к golangci-lint (module plugin)

Пример `.golangci.yml`:

```yaml
version: "2"

linters-settings:
  custom:
    logmsglint:
      path: ./cmd/golangci-plugin
      description: Checks slog/zap log messages

linters:
  enable:
    - logmsglint
```

> Конкретный способ сборки/подключения может отличаться в зависимости от версии `golangci-lint`; функция `New(any) ([]*analysis.Analyzer, error)` в `cmd/golangci-plugin` подготовлена под модульную схему плагинов.

## Минимальный пример для golangci-lint (plugin)

Минимально заполненный `.golangci.yml`:

```yaml
version: "2"

linters-settings:
  custom:
    logmsglint:
      path: ./cmd/golangci-plugin

linters:
  enable:
    - logmsglint
```

Если хотите сразу передать настройки правил, добавьте рядом файл `.logmsglint.yml`:

```yaml
rules:
  lowercase: true
  english: true
  nospecial: true
  sensitive: true
```

## Как использовать без golangci-lint

Самый простой путь — запуск бинаря `logmsglint`:

```bash
go install github.com/bauerex/logmsglint/cmd/logmsglint@latest
logmsglint ./...
```

С кастомным конфигом:

```bash
logmsglint -config=.logmsglint.yml ./...
```

Локально из исходников (без установки):

```bash
go run ./cmd/logmsglint -config=.logmsglint.yml ./...
```

## CI/CD

В репозитории добавлены шаблоны для автоматического тестирования:
- GitHub Actions: `.github/workflows/ci.yml`
- GitLab CI: `.gitlab-ci.yml`

Оба  запускают `go test ./...`.
