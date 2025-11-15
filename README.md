# Memory Analyzer

![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)

Кроссплатформенная утилита для мониторинга использования памяти системы и отдельных процессов в реальном времени. Поддерживает Linux и macOS с удобным интерфейсом терминала.

## Содержание

- [Обзор функциональности](#обзор-функциональности)
- [Требования к системе](#требования-к-системе)
- [Установка](#установка)
- [Использование](#использование)
- [Примеры работы](#примеры-работы)
- [Структура проекта](#структура-проекта)
- [Архитектура](#архитектура)
- [Планы развития](#планы-развития)
- [Лицензия](#лицензия)

## Обзор функциональности

**Memory Analyzer** — это инструмент для мониторинга памяти, предоставляющий следующие возможности:

### Основные возможности

- **Системная статистика памяти**: отображение общей памяти, использованной памяти, доступной памяти и информации о подкачке
- **Список топ процессов**: вывод процессов, потребляющих наибольшее количество памяти, с сортировкой в реальном времени
- **Кроссплатформенность**: единая кодовая база для Linux и macOS
- **Обновление в реальном времени**: автоматическое обновление данных с настраиваемым интервалом (по умолчанию 3 секунды)
- **Красивый интерфейс терминала**: отформатированная таблица с информацией о процессах
- **Процентный расчет**: отображение процента использованной памяти от общего объема

### Поддерживаемые метрики

**Системная память:**
- Total Memory — общий объем оперативной памяти
- Used Memory — использованная память с процентным показателем
- Available Memory — доступная (свободная + кешируемая) память
- Swap Used — используемая подкачка с процентом

**Информация о процессах:**
- PID — идентификатор процесса
- NAME — сокращенное имя процесса (до 15 символов)
- MEMORY — использованная память в удобном формате (B, KB, MB, GB, TB)

## Требования к системе

### Минимальные требования

| Параметр | Требование |
|----------|-----------|
| **ОС** | Linux (2.6.32+) или macOS (10.13+) |
| **Go** | 1.21 или выше |
| **Память** | 50 MB свободной памяти |
| **Процессор** | Любой современный процессор |

### Зависимости по платформам

**Linux:**
- Доступ к `/proc` файловой системе
- Утилита `ps` (обычно предустановлена)

**macOS:**
- Утилиты `ps`, `sysctl`, `vm_stat` (входят в стандартную комплектацию)
- Возможно требуется доступ с правами администратора для некоторых процессов

## Установка

### Предварительно

Убедитесь, что у вас установлен Go 1.21 или выше:

```bash
go version
```

### Вариант 1: Клонирование репозитория

```bash
git clone https://github.com/yourusername/MemoryAnalyzer.git
cd MemoryAnalyzer
```

### Вариант 2: Использование go install

```bash
go install github.com/yourusername/MemoryAnalyzer@latest
```

### Сборка из исходников

```bash
# Загрузка зависимостей
go mod download

# Сборка приложения
go build -o MemoryAnalyzer .

# Или использование Makefile
make build
```

### Установка в систему

```bash
# Скопировать бинарник в системный путь
sudo cp MemoryAnalyzer /usr/local/bin/

# Или использование Makefile
make install
```

## Использование

### Базовый запуск

```bash
./MemoryAnalyzer
```

Или если установлено в систему:

```bash
MemoryAnalyzer
```

### Завершение программы

Нажмите `Ctrl+C` для корректного завершения приложения.

### Примеры использования

**Запуск в фоновом режиме:**

```bash
./MemoryAnalyzer &
```

## Примеры работы

### Вывод программы

```
=== Memory Analyzer ===

System Memory:
Total: 16.00 GB
Used: 8.45 GB (52.8%)
Available: 7.55 GB
Swap Used: 0.00 GB (0.0%)

Top Memory Processes:
PID      NAME            MEMORY
-----------------------------------
2847     Chrome          2.34 GB
1923     Firefox         1.23 GB
1845     VSCode          856.32 MB
2156     Docker          512.00 MB
1234     node            256.00 MB
Updated: 2025-11-15 14:32:45
Press Ctrl+C to exit
```

### Типичные сценарии использования

**Мониторинг утечек памяти:**

```bash
# Запустить и наблюдать за процессом
./MemoryAnalyzer
# Следить за возрастанием памяти конкретного процесса
```

**Поиск прожорливых приложений:**

```bash
# Приложение автоматически показывает топ-10 процессов
./MemoryAnalyzer
# По умолчанию отсортированы по убыванию использования памяти
```

**Анализ полной системной памяти:**

```bash
# Комбинирование с другими утилитами
./MemoryAnalyzer | head -20  # Только системная информация
```

## Структура проекта

```
MemoryAnalyzer/
├── main.go              # Точка входа приложения
├── go.mod               # Модуль Go и зависимости
├── go.sum               # Хеши зависимостей
├── Makefile             # Автоматизация сборки
├── .gitignore           # Правила исключения для Git
├── README.md            # Этот файл
├── LICENSE              # Лицензия проекта (MIT)
│
├── process/             # Пакет для работы с процессами и системной информацией
│   ├── process.go       # Структура ProcessInfo и утилиты для работы с процессами
│   ├── memory.go        # SystemMemoryInfo и функции форматирования памяти
│   ├── dashboard.go     # DisplayDashboard — функция отображения интерфейса
│   └── config.go        # DisplayConfig — конфигурация отображения
│
└── platforms/           # Пакет для кроссплатформенной работы
    ├── memoryReader.go  # MemoryReader — интерфейс для чтения памяти
    ├── linux.go         # LinuxMemoryReader — реализация для Linux
    └── darwin.go        # DarwinMemoryReader — реализация для macOS
```

### Описание компонентов

**main.go**
- Инициализация приложения
- Выбор реализации MemoryReader в зависимости от ОС
- Основной цикл обработки событий
- Обработка сигналов завершения (Ctrl+C)

**process/process.go**
- `ProcessInfo` — информация о процессе (PID, Name, MemoryUsage)
- `getShortProcessName()` — сокращение длинных имен процессов
- `FormatTable()` — форматирование таблицы процессов
- `GetProcessName()` — получение имени процесса по PID

**process/memory.go**
- `SystemMemoryInfo` — структура системной информации о памяти
- `FormatMemorySize()` — преобразование байтов в удобный формат (B, KB, MB, GB, TB)
- `FormatSystemStats()` — форматирование статистики системной памяти

**process/dashboard.go**
- `DisplayDashboard()` — основная функция отображения интерфейса
- Очистка экрана и вывод форматированной информации

**platforms/memoryReader.go**
- `MemoryReader` — интерфейс для кроссплатформенной работы
- Методы: `ReadSystemMemory()`, `GetProcessList()`, `ReadProcessMemory()`

**platforms/linux.go**
- `LinuxMemoryReader` — реализация для Linux
- Чтение из `/proc/[pid]/status` и `/proc/meminfo`
- Парсинг VmRSS для использования памяти процесса

**platforms/darwin.go**
- `DarwinMemoryReader` — реализация для macOS
- Использование команд `ps`, `sysctl`, `vm_stat`
- Парсинг RSS для информации о памяти

## Архитектура

### Паттерны проектирования

Проект использует следующие архитектурные паттерны:

**Strategy Pattern (Стратегия)**

```
MemoryReader (интерфейс)
    ├── LinuxMemoryReader (конкретная стратегия)
    └── DarwinMemoryReader (конкретная стратегия)
```

Позволяет переключаться между реализациями в зависимости от ОС без изменения основного кода.

**Separation of Concerns (Разделение ответственности)**

- `platforms/` — работа с ОС и получение данных
- `process/` — обработка и представление данных
- `main.go` — координация и основной цикл

### Поток данных

```
main()
  ├── Инициализация MemoryReader в зависимости от ОС
  ├── Создание конфигурации (UpdateInterval, TopProcesses)
  └── Основной цикл:
       ├── Получение системной памяти: reader.ReadSystemMemory()
       ├── Получение списка процессов: reader.GetProcessList()
       ├── Получение памяти каждого процесса: reader.ReadProcessMemory(pid)
       └── Отображение: DisplayDashboard()
```

### Особенности реализации

**Linux реализация:**
- Работает через файловую систему `/proc`
- Файлы читаются синхронно для каждого процесса
- VmRSS показывает резидентную память (Resident Set Size)

**macOS реализация:**
- Использует системные команды через `exec.Command()`
- `ps` для списка процессов и RSS памяти
- `sysctl` для общей памяти
- `vm_stat` для информации о страницах памяти и подкачке

## Планы развития

### Краткосрочные (v1.1)

- [ ] Добавить поддержку Windows
- [ ] Реализовать фильтрацию процессов по имени
- [ ] Добавить сохранение истории в CSV
- [ ] Реализовать вывод в JSON формате
- [ ] Добавить конфигурационный файл (.config/MemoryAnalyzer/config.yaml)

### Среднесрочные (v1.2)

- [ ] Web интерфейс через HTML/CSS
- [ ] Гистограмма использования памяти
- [ ] Мониторинг тенденций (увеличение/уменьшение памяти)
- [ ] Уведомления при превышении порога памяти
- [ ] Интеграция с системными логами

### Долгосрочные (v2.0)

- [ ] Поддержка контейнеризации (Docker, Kubernetes)
- [ ] REST API для удаленного мониторинга
- [ ] Интеграция с Prometheus и Grafana
- [ ] Multi-host мониторинг
- [ ] Анализ утечек памяти с профилированием
- [ ] Desktop приложение на основе Go + Gio UI

### Оптимизация и улучшения

- [ ] Кеширование списка процессов (сейчас пересчитывается каждый раз)
- [ ] Параллельное чтение информации о процессах (горутины)
- [ ] Снижение нагрузки на ЦПУ при работе с большим количеством процессов
- [ ] Оптимизация обхода `/proc` на Linux

## Разработка

### Требования для разработчиков

- Go 1.21+
- Git
- Make (для использования Makefile)

### Локальная установка для разработки

```bash
git clone https://github.com/yourusername/MemoryAnalyzer.git
cd MemoryAnalyzer
go mod download
make build
./MemoryAnalyzer
```

### Форматирование кода

```bash
go fmt ./...
gofmt -w .
```

### Проверка кода

```bash
go vet ./...
```

### Сборка для разных платформ

```bash
# Linux (AMD64)
GOOS=linux GOARCH=amd64 go build -o MemoryAnalyzer-linux

# macOS (ARM64 - M1/M2)
GOOS=darwin GOARCH=arm64 go build -o MemoryAnalyzer-darwin-arm64

# macOS (AMD64 - Intel)
GOOS=darwin GOARCH=amd64 go build -o MemoryAnalyzer-darwin-amd64
```

## Решение проблем

### Приложение не запускается на Linux

**Проблема:** Permission denied

**Решение:**
```bash
chmod +x MemoryAnalyzer
./MemoryAnalyzer
```

### Неполная информация о процессах

**Проблема:** Некоторые процессы не отображаются или показывают 0 памяти

**Решение:** 
- Некоторые процессы требуют прав администратора для чтения информации
- Запустите с `sudo` для полной информации
```bash
sudo ./MemoryAnalyzer
```

### Высокая нагрузка на ЦПУ

**Проблема:** Приложение потребляет много CPU

**Решение:** Это может быть вызвано большим количеством процессов. Увеличьте интервал обновления в коде main.go:
```go
UpdateInterval: 5 * time.Second,  // Увеличьте с 3 до 5
```

## Лицензия

Проект распространяется под лицензией **MIT**.

```
MIT License

Copyright (c) 2025 Memory Analyzer Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Сведения о проекте

- **Язык:** Go (Golang)
- **Минимальная версия Go:** 1.21
- **Платформы:** Linux, macOS
- **Статус:** Активная разработка
- **Версия:** 1.0.0

---

