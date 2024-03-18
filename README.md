# VK-assigment
## Установка и запуск
```sh 
git clone https://github.com/sivistrukov/vk-assigment.git
cd vk-assigment
make up
```

## Методы

**swagger документация для API расположена на http://localhost:8080/swagger/index.html**

- [post] {{base_url}}/v1/users - создание нового пользователя
- [get] {{base_url}}/v1/actors - получение списка актеров 
- [post] {{base_url}}/v1/actors - добавление нового актера
- [put] {{base_url}}/v1/actors/{id} - обновление данных об актере
- [patch] {{base_url}}/v1/actors/{id} - частичное обновление данных об актере
- [delete] {{base_url}}/v1/actors/{id} - удаление актера
- [get] {{base_url}}/v1/films - получение списка фильмов с поиском и сортировкой
- [post] {{base_url}}/v1/films - добавление нового фильма
- [put] {{base_url}}/v1/films/{id} - обновление данных об фильме
- [patch] {{base_url}}/v1/films/{id} - частичное обновление данных об фильме
- [delete] {{base_url}}/v1/films/{id} - удаление фильма

## База данных

Users:
| id  | username | password(hashed) | is_admin |
| --- | -------- | ---------------- | -------- |
| 1   | admin    | admin            | true     |
| 2   | user     | user             | false    |

Actors:
| id  | fist_name | last_name | middle_name | sex    | birthday   |
| --- | --------- | --------- | ----------- | ------ | ---------- |
| 1   | John      | Doe       | J           | male   | 2006-01-02 |
| 2   | Ryan      | Gosling   | NULL        | male   | 1980-11-12 |
| 3   | Ryan      | Reynolds  | NULL        | male   | 1976-10-23 |
| 4   | Margot    | Robbie    | NULL        | female | 1990-07-02 |

Films:
| id  | title                   | description | release_date | rating |
| --- | ----------------------- | ----------- | ------------ | ------ |
| 1   | Drive                   | ...         | 2011-11-03   | 7      |
| 2   | Blade Runner 2049       | ...         | 2017-10-05   | 8      |
| 3   | Barbie                  | ...         | 2023-07-20   | 7      |
| 4   | The Proposal            | ...         | 2009-06-18   | 8      |
| 5   | The Wolf of Wall Street | ...         | 2013-02-06   | 9      |

Actors and films:
| id  | actor_id | film_id |
| --- | -------- | ------- |
| 1   | 2        | 1       |
| 2   | 2        | 2       |
| 3   | 2        | 3       |
| 4   | 4        | 3       |
| 5   | 3        | 4       |
| 6   | 4        | 5       |
