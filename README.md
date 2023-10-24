Переменные для конфигурирования лежат в config/config.env

Каталог для записи лог-файла должен быть создан или у пользователя, от имени которого запускается сервис, должны быть права на создание каталога, указанного в .env файле 

Запускается с указанием переменной окружения к файлу с конфигурацией. Например, "CONFIG_PATH=/home/novam/go/src/github.com/Raitfolt/juntest/config/config.env"

API

	router.Post("/add", add.New(log, db))
	router.Post("/delete", del.Delete(log, db))
	router.Get("/all", all.List(log, db))
	router.Post("/change", change.Change(log, db))

