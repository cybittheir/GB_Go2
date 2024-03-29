package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
Задание 2. Graceful shutdown

Цель задания:
Научиться правильно останавливать приложения.

Что нужно сделать:
-В работе часто возникает потребность правильно останавливать приложения.
Например, когда наш сервер обслуживает соединения, а нам хочется,
чтобы все текущие соединения были обработаны и лишь потом произошло выключение сервиса.
Для этого существует паттерн graceful shutdown.
-Напишите приложение, которое выводит квадраты натуральных чисел на экран,
а после получения сигнала ^С обрабатывает этот сигнал, пишет «выхожу из программы» и выходит.

Советы и рекомендации:
Для реализации данного паттерна воспользуйтесь каналами и оператором select с default-кейсом.
*/

func main() {

	fmt.Println("-= Start =-")

	x := make(chan os.Signal, 1)
	signal.Notify(x, syscall.SIGINT, syscall.SIGTERM)

	i := 1
	for {

		select {
		case sig := <-x:
			fmt.Println("Graceful shutdown", sig)
			return
		default:
			fmt.Println("i ^ 2 =", i*i)
			time.Sleep(200 * time.Millisecond)
			i++
		}

	}

}
