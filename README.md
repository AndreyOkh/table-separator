# table-separator

Приложение предназначено для разделения таблицы по колонке. Принимается таблица в формате ODS, на выходе множество таблиц в формате CSV.

### Принимаемые параметры
| Ключ | Значение                                                            | Параметр по умолчанию                                                                     |
|------|---------------------------------------------------------------------|-------------------------------------------------------------------------------------------|
| `-f` | Путь к файлу                                                        | -                                                                                         |
| `-o` | Директория куда сохранять таблицы                                   | По умолчанию в каталоге запуска создается директория с именем "files_ГГГГ-ММ-ДД_ЧЧ-мм-СС" |
| `-c` | Колонка по которой необходимо фильтровать. Нумерация начинается с 0 | По умолчанию таблица фильтруется по 5 колонке (F)                                         |
