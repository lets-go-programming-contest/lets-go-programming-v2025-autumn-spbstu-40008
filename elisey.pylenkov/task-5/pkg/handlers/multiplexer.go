package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	// Читаем из всех входных каналов пока они не закроются
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Проверяем все каналы
			for _, ch := range inputs {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case data, ok := <-ch:
					if !ok {
						// Канал закрыт, продолжаем проверять другие
						continue
					}

					if strings.Contains(data, "no multiplexer") {
						continue // Пропускаем данные с этой подстрокой
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					case output <- data:
						// Успешно отправили
					}
				default:
					// Нет данных в этом канале, проверяем следующий
				}
			}
		}
	}
}
