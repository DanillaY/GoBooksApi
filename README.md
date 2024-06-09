![alt text](./Md/logo.png 'text')

----

# GoBooksApi

## Содержание

### Что такое GoBooksApi?
**GoBooksApi** - это api, который позволяет получить доступ к информации в базе данных в формате json.Api Использует GoScrapper для сбора информации на сайтах:

<ol>
	<li>https://book24.ru/</li>
	<li>https://vse-svobodny.com/</li>
	<li>https://www.chitai-gorod.ru/</li>
	<li>https://www.labirint.ru/</li>
	<li>https://www.respublica.ru/</li>
</ol>

Данные о книгах собираются с целью проведения анализа цен в различных цифровых магазинах.

### Как использовать Api?

Для использования api потребуется:

1. Для клонирования проекта с подмодулем требуется написать команду:
	```shell
	git clone https://github.com/DanillaY/GoBooksApi --recurse-submodules
	```
2. Для создания и запуска docker контейнеров требуется написать:
    ```shell
	docker compose up --build
	```
После запуска приложения контейнер GoScrapper начнет собирать данные о книгах в базу данных с выводом соответствующих логов.

Для того, чтобы узнать на каком ip открылись docker контейнеры нужно прописать команду `ipconfig` для windows или `ip addr show` для linux, обычно docker контейнеры начинаются с 172. По умолчанию api будет слушать на 8089 порте, чтобы его поменять потребуется отредактировать файл [dbDokcer.env](./dbDocker.env). 

#### Api Routes

GoBooksApi имеет следующие эндпоинты:

*	/getBooks [<span style="color:rgba(9,234,87, 1);">*GET*</span>] - эндпоинт, который принимает различные параметры для фильтрации (см. [query параметры](#query-параметры-для-метода-getBooks) для подробного описания), возвращает объект пагинации и объект "data", внутри которого находится массив книг
*	/getBooksById [<span style="color:rgba(9,234,87, 1);">*GET*</span>] - эндпоинт, который принимает id книги из базы, возвращает объект книги
*	/getProperties [<span style="color:rgba(9,234,87, 1);">*GET*</span>] - эндпоинт, который принимает параметр отдельного поля из базы, возвращает уникальные записи по запрошенному свойству (см. [query параметры](#query-параметры-для-метода-getProperties) для подробного описания)
*	/getMinMaxPrice [<span style="color:rgba(9,234,87, 1);">*GET*</span>] - эндпоинт, который отвечает за получение мапы с тегами минимальной и максимальной цены, не принимает дополнительных параметров.
*	/getBooksByEmail [<span style="color:rgba(9,234,87, 1);">*GET*</span>] - эндпоинт, который отвечает за получение массива книг с тегом books, массив основан на подписках отправленной электронной почты.
*	/deleteBookSubscriber [<span style="color:rgba(255,72,72);">*DELETE*</span>] - эндпоинт, который удаляет указаного подписчика от книги, возвращает сообщение о том был ли удален пользователь.
*	/addNewBookSubscriber [<span style="color:rgba(248,160,76);">*POST*</span>] - эндпоинт, который отвечает за добавление нового пользователя для книги.

Стандартный запрос на получение книг без дополнительных параметров выведет первые 30 записей из базы данных, а также объект пагинации.

#### Возвращаемые модели BookApi

##### Модель book json

В ответ на запрос сервер должен отправить указанное количество моделей book.
**Важно отметить, что не все поля могут быть заполнены!**

Пример одной модели book:

```json
{
	"ID": 662,
	"CurrentPrice": 239,
	"OldPrice": 319,
	"Title": "Камю Альбер: Посторонний",
	"ImgPath": "https://cdn.book24.ru/v2/ASE000000000858196/COVER/cover13d__w410.jpg",
	"PageBookPath": "https://book24.ru/product/postoronniy-5964736/",
	"VendorURL": "https://book24.ru",
	"Vendor": "Book24",
	"Author": "Альбер Камю",
	"Translator": "Галь Нора",
	"ProductionSeries": "Эксклюзивная классика",
	"Category": "Классическая зарубежная литература",
	"Publisher": "АСТ, Neoclassic",
	"ISBN": "978-5-17-137323-8",
	"AgeRestriction": "16+",
	"YearPublish": "2021",
	"PagesQuantity": "128",
	"BookCover": "Мягкий (3)",
	"Format": "115x180  мм",
	"Weight": "0.11  кг",
	"BookAbout": "Посторонний — дебютная работа молодого писателя, своеобразный творческий манифест. Понятие абсолютной свободы — основной постулат этого манифеста. Героя этой повести судят за убийство, которое он совершил по самой глупой из всех возможных причин. И это правда, которую герой не боится бросить в лицо своим судьям, пойти наперекор всему, забыть обо всех условностях и умереть во имя своих убеждений."
},
```
##### Модель pagination json

В api автоматически применяется функция пагинации, вместе с объектами book api отправляет объект пагинации, который описывает состояние текущей страницы.

Пример одной модели book:
```json
"pagination": {
	"Total": 35,
	"PerPage": 30,
	"CurrentPage": 1,
	"LastPage": 2
}
```

#### Описание параметров для имеющихся эндпоинтов

##### Query параметры для метода getBooks

Чтобы легко фильтровать записи из базы данных /getBooks метод принимает различные параметры, которые ограничивают область видимости базы, такие как:

- category - отвечает за категорию книги 
- search - отвечает за фильтрацию по полям title, author и category, находит похожие записи по введенному запросу.
- year - ответчает за год публикации книги, находит книги с определенным годом публикации, год может быть только 1.
- minPrice, maxPrice - отвечает за минимальную и максимальную цену книги, данные параметры чаще всего используются вместе для определения ограничений цены, но также они могут использоваться и отдельно, **важно заметить, что данные параметры сравниваются с текущей ценой (current_price), а не с ценой без скидки (old_price)**.
- stockText - параметр, отвечающий за фльтрацию книг по категории наличия, возможные параметры stockText: "В наличии", "Ожидается" и очень редко "Нет в наличии".
- sortOrder - параметр, отвечающий за порядок сортировки поля, **должен использоваться вместе с полем sortField, который будет описан ниже**, параметр sortOrder может иметь 2 возможных значения:  ASC (от меньшего к большему) и DESC (от большого к меньшему).
- sortField - параметр, отвечающий за название таблицы, для которой будет предназначена сортировка, **должен использоваться вместе с полем sortOrder**, параметр sortFiled может иметь значения названия колонок из базы данных, все возможные колонки можно посмотреть в [скрпите создания базы данных](./books.sql)
- limit - ограничение количества отправляемых моделей book, если требуется убрать ограничение в 30 моделей, то следует отправить запрос с limit равным -1, и тогда будут выведены все имеющиеся книги по данному запросу
- author - параметр, отвечающий за автора книги, также есть возможность вписать несколько авторов **перечисленные авторы должны быть разделены запятой**.
- bookInfoType - параметр, отвечающий за полноту отправляемых данных, изначальное значение bookInfoType равно partial, если значение равно стандартному, тогда только данные из полей id, title, current_price, old_price, img_path, page_book_path, vendor, age_restriction будут загружены и отправлены, иначе все поля будут отправлены.
- vendor - ограничение, отвечающее за название поставщика, поставщик может быть только 1 в параметре.
- pageNum - параметр, отвечающий за пагинацию, данный параметр принимает целочисленное значение страницы, которой вы хотите получить. Стандартным значением pageNum является 1, данный параметр зависит от параметра limit, но если его не указывать то пагинация страниц будет происходит со стандартным значением limit, тоесть 30.  

Пример запроса с применением различных query параметров: <http://172.21.80.1:8089/getBooks?limit=3&author=ницше&title=воля&maxPrice=220>

##### Query параметры для метода getBooksById

Чтобы легко фильтровать записи из базы данных /getBooksById метод принимает параметр id, который возвращает запись книги с указанным id.

Пример запроса с применением различных query параметров: <http://172.21.80.1:8089/getBooksById?id=15>

##### Query параметры для метода getProperties

Чтобы легко фильтровать записи из базы данных /getProperties метод принимает параметр property, который возвращает список всех возможных записей свойства книги.

Пример запроса с применением различных query параметров: <http://172.21.80.1:8089/getProperties?property=category>
##### Query параметры для метода getMinMaxPrice

Данный метод возвращает минимальную и максимальную цену, не принимая никаких параметров.

Пример запроса с применением различных query параметров: <http://172.21.80.1:8089/getMinMaxPrice>

##### Query параметры для метода getBooksByEmail

Данный метод возвращает массив объектов книг в теге books, массив овнован на подписках пользователя, чья электронная почта будет отправлена, принимает параметр userEmail.

Пример запроса с применением различных query параметров: <http://172.21.80.1:8089/getBooksByEmail?userEmail=test@gmail.com>

##### Query параметры для метода addNewBookSubscriber 

Для добавления нового пользователя эндпоинт принимает 2 параметра, а именно bookId, который обозначает ID книги к которой вы хотите подписать пользователя и userEmail, обозначающий почту пользователя, на которое должно прийти оповещение.

Пример запроса с применением различных query параметров: <http://localhost:8089/addNewBookSubscriber?bookId=5391&userEmail=test@gmail.com>

##### Query параметры для метода deleteBookSubscriber

Для удаления подписки пользователя из базы эндпоинт /deleteBookSubscriber принимает идентичные параметры как /addNewBookSubscriber, а именно bookId и userEmail

Пример запроса с применением различных query параметров: <http://localhost:8089/deleteBookSubscriber?bookId=69&userEmail=test@gmail.com>

Также query параметры в GoBooksApi не является регистрозависимыми, это означает, что ответ сервера при запросе <http://172.21.80.1:8089/getBooks?author=Гоголь> не будет отличаться от ответа запроса <http://172.21.80.1:8089/getBooks?author=гОгоЛь>