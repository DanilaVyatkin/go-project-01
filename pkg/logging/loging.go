package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type writerHook struct { //Hook на запись
	Writer    []io.Writer
	LogLevels []logrus.Level //Что бы в любой Writer можно было отправлять любое количество уровней логирования
}

func (hook *writerHook) Fire(entry *logrus.Entry) error { //Будет вызываться каждый раз при записи
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line)) //Лайн переводим в массив байт
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level { //Будет возвращать левелы из нашего хука
	return hook.LogLevels
}

//Получаем логер
var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

//Суть этого хука распределить на каждого riter несколько уровней логирования
//kafka (riter) отправляем info, debug
//file (riter) отправляем err, trace
//stdout (riter) отправляем warning, critical

func init() {
	l := logrus.New() //Создаём логгер
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) { //Нужно, что бы показать в каком месте логируем
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s, %d", fileName, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	// CallerPrettyfier передаёт runtime.Frame в котором происходит логирование
	// в нем есть информация о файле в котором происходит логирование (path.Base(frame.File))
	// так же есть информация о текущем лайне на котором происходит логирование (frame.Line)
	// и функция внутри которой мы находимся (frame.Function)
	// CallerPrettyfier возвращает функцию, имя файла и строчку

	err := os.MkdirAll("logs", 0644) // Создаём файл
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640) // Открываем файл
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard) // Что бы логрус не записывал никуда информацию

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
