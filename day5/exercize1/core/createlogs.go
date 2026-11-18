package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//принимаем какой уровень значения мы хотим использовать
/*
например
DEBUG
INFO
warn
error
*/
/*
возвращаем мы сам logger который ьщжем использоать, например logger.DEBUG("")
так же func() error который нужно возвращать когда мы заканчиваем логгирование например defer loggerClose()
и возвращаем ошибку которая у нас могла произайти во время логирования
*/
func NewLogs(logLevel string) (*zap.Logger, func() error, error) {
	//ATOMICLEVEL это специльные тип zap который хранит в себе уровень логгирования
	lvl := zap.NewAtomicLevel()
	// здесь мы преобразуем строку которую передали в функцию например если это был debug то на выходе у нас получится debug
	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		//обработка ошибки
		return nil, nil, fmt.Errorf("unmarshal log level:%w", err)
	}
	/*
		тут мы создаем папку под название logs и будем туда все складывать
		что значит 0755
		0-значит 8 системы исчисления
		\     7         5          5
		    user    user-group    other
			кwx        rwx         rwx
			111        101         101
			read- читать  wrait - менять что либо   x - заходить и работать с файлами(ориентироваться например)
	*/
	if err := os.MkdirAll("features/out/logs", 0755); err != nil {
		// обработка ошибки если не получилось создать
		return nil, nil, fmt.Errorf("mkdir log folder:%w", err)
	}

	// тут мы задаем новый формат времени
	timeStamp := time.Now().UTC().Format("2006-01-02T15-04-05-00000")
	//fmt.Sprintf собирает строку а с помощью filepath.Join обьединяет части пути с учето правил операционноый системы
	// так мы и получаем logs/2026-11-17T15-30-00.12345.log
	// и так когда будет создаваться логгер будет такое описание
	logsFilePath := filepath.Join("features", "out", "logs", fmt.Sprintf("%s.logs", timeStamp))

	//тут открываем файл
	// os.OpenFile() принимает решение - путь к файлу, режим открытия( в нашем случаи(создать файл, открыть для записи)), прова доступа на этот файл,
	/*
		     6         4            4
		    user    user-group    other
			кwx        rwx         rwx
			110       100         100
		r - это читать файл w - это записывать файл x - это создавть файл исполняемым, тоесть деать например скрипт
	*/
	logFile, err := os.OpenFile(logsFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//обрабатываем ошибку если не получилось открыть
		return nil, nil, fmt.Errorf("open log file:%w", err)
	}
	//создаем конфиг для логгирования - можем изменить почти все что хотим
	cfg := zap.NewDevelopmentEncoderConfig()
	// меняем timestamp в логгировании в точности до миллисекунды
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15-04-05-00000")

	//encoder создает нам формат в стринг что бы мы могли читать, тоесть если мы напишешм logger.DEBUG("жопа")
	//то encoder оперделяет как это запись будет выглядить в терминале или в логах
	encoder := zapcore.NewConsoleEncoder(cfg)

	// создаем два подядрка куда мы будем печатать log
	// первый он будет писать в console
	// второй будет писать в наш file
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl), // lvl это ставится уровень логирования котоыре мы передали в аргументах функции
	)
	//ну и так то все, тут мы иницализируем сам логгер
	logger := zap.New(
		// даем ему едро
		core,
		//это комнада значит то откуда пришел сам лог
		zap.AddCaller(),
		//а это, мы показываем что стек вызовов мы будем показывать для Error и выше
		zap.AddStacktrace(zap.ErrorLevel),
	)
	// тут мы создаем свою функцию котрая возращает ошибку - но это не обязательно можно вместо нее вернуть logFile.Close
	cleanup := func() error {
		//мы просим zap сбросить буферезированые записи - Это важно, потому что часть данных может некоторое время находиться в буфере, а не сразу попадать в файл.
		syncErr := logger.Sync()
		//закрываем файл
		closeErr := logFile.Close()
		//позваоляет вернуть сразу нескольок ошибок, поэтому мы их совмещаем и возращаем
		return errors.Join(syncErr, closeErr)
	}
	return logger, cleanup, nil
}
