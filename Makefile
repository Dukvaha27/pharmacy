# Makefile

.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  make run       - Обычный запуск"
	@echo "  make dev       - Запуск с hot-reload (перезапуск через 10 сек после сохранения)"
	@echo "  make build     - Сборка бинарника"
	@echo "  make clean     - Очистка временных файлов"

.PHONY: run
run:
	go run cmd/app/main.go

.PHONY: dev
dev:
	@echo "🚀 Запуск с hot-reload (Air)..."
	@echo "📝 Air будет ждать 10 секунд после последнего изменения перед перезапуском"
	@echo "💡 Сохрани файл и подожди 10 секунд - сервер перезапустится"
	@echo ""
	@command -v air >/dev/null 2>&1 || { \
		echo "❌ Air не установлен. Установка..."; \
		go install github.com/air-verse/air@latest; \
		echo "✅ Air установлен"; \
	}
	air

.PHONY: build
build:
	@echo "🔨 Сборка приложения..."
	go build -o bin/app cmd/app/main.go
	@echo "✅ Бинарник создан: bin/app"

.PHONY: clean
clean:
	@echo "🧹 Очистка..."
	rm -rf tmp/ bin/
	@echo "✅ Очистка завершена"