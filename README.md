### Что было изменено

1. Два канала: 
   * users - в него отправляются новые юзеры после генерации
   * save - канал, в который записываются результаты записи в файл

2. Генерация отдельного юзера в горутине - свободынй воркер заберет себе нового юзера

Время работы:
 - без 'worker pool': 112 сек
 - пять воркеров: 20 сек
