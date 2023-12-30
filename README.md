Требуется установленный go

Порядок команд:
1) go build
2) go install

Далее можно пользоваться командами:

1) cloudphoto init (инициализация программы)

2) cloudphoto upload --album ALBUM [--path PHOTOS_DIR] (отправка фотографий в облачное хранилище)

3) cloudphoto download --album ALBUM [--path PHOTOS_DIR] (загрузка фотографий из облачного хранилища)

4) cloudphoto list [--album ALBUM] (просмотр списка альбомов и фотографий в альбоме)

5) cloudphoto delete --album ALBUM [--photo PHOTO] (удаление альбомов и фотографий)

6) cloudphoto mksite (генерация и публикация веб-страниц фотоархива)