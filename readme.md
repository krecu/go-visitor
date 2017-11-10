# Visitor

Сервис хранения, систематизации и анализа данных о посетителе (DMP)

  - определения ГЕО информации по IP пользователя (SyPexGeo, MaxMind)
  - определение устройства пользователя по UA (BrowsCap, RegExp)
  - хранение CookieMatching данных пользователя
  - хранение сегментированых данных
  - возможность хранения любого набора данных о пользователе

#### Данные
`id`
Идентификатор пользователя в системе
*пример*:
```json
"id": "63a406058d915def9d3135d9ee34e132"
```
`created`
Дата создания записи в **timestamp**
*пример*:
```json
"created": 1510311263
```
`city`, `country`, `region`
Гео информация, полученная на основании полученного IP
*пример*:
```json
"city": {
    "name": "Appleton",
    "name_ru": "Аплтон",
    "geoname_id": 5244080,
    "mapping_id": 0
},
"country": {
    "name": "United States",
    "name_ru": "США",
    "geoname_id": 225,
    "iso_code": "US",
    "iso_code_3166_1_alpha_3": "USA",
    "mapping_id": 0
},
"region": {
    "name": "Wisconsin",
    "name_ru": "Висконсин",
    "geoname_id": 5279468,
    "iso": "US-WI",
    "mapping_id": 0
},
```
`location`, `postal`
Доп гео информация
*пример*:
```json
"location": {
    "latitude": 44.26193,
    "longitude": 42861.258,
    "time_zone": "America/New_York"
},
"postal": {
    "code": "10013"
},
```
`device`, `browser`, `platform` 
Информация об устройстве
*пример*:
```json
"browser": {
        "name": "Chrome",
        "type": "Browser",
        "version": "18.0",
        "majorver": "",
        "minorver": "",
        "mapping_id": 0
    },
"device": {
    "name": "Nexus 7",
    "type": "Tablet",
    "brand": "Google",
    "mapping_id": 12
},
"platform": {
    "name": "Android",
    "short": "android",
    "version": "4.1",
    "description": "",
    "maker": "",
    "mapping_id": 0
},
```
`personal`, `ip`
Персональная информация о пользователе, если конечно можно :)
*пример*:
```json
"personal": {
        "ua": "Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) Safari/535.19",
        "first_name": "",
        "last_name": "",
        "patronymic": "",
        "age": "",
        "gender": ""
},
"ip": {
        "v4": "165.227.53.107",
        "v6": ""
},
```
`extra`
Дополнительное поле для хранения вычисляемых данных на стороне клиента
*пример*:
```json
"extra": {
    "amberdata": ["121212", "123123"], // пример сегментирования
    "matching": {   // хранение кук
        "dsp_1": "23423rfsdf",
    },
    "campaign": { // пример счетчика просмотров компании
        "111": 2,
    },
    "timeline": {   // пример хранения счетчика временени
        "sessions": [
            {
                "id": "234234sgsvx8tsdyasej",
                "started": 1510311263,
                "seconds": 20 
            }
        ],
        "total": 12,
    }
}
```
`debug`
Отладочная информация
*пример*:
```json
"debug": {
    "GeoProvider": "sypexgeo", // выбранный гео провайдер
    "DeviceProvider": "browscap", // выбранный девайс провайдер
    "TimeGeo": "231.655µs", // затраченное время на вычисления гео
    "TimeDevice": "3.835µs", // затраченное время на вычисление устройства
    "TimeTotal": "239.905µs" // общее время от получения запроса до ответа клиенту
}
```

#### Установка

Необходим компилятор [Golang](https://golang.org/) v1.5+.

```sh
$ make install
```

#### Настройка
Конфигурация находится в файле config.yaml
```yaml
app:
  system:
      cpu: 4 // кол-во используемых CPU
      LogFile: "/var/log/log.log" // место логов
      mode: "debug|prod" // режим
      instance: "visitor" // facility значение для graylog
      DebugLevel: 3 // уровень логирования
  graylog:
      host: "192.168.0.48"
      port: 12201
  server:
    grpc:
        listen: ":8081"
    http:
        listen: ":8082"
        ReadTimeout: 5
        WriteTimeout: 5
        MaxHeaderBytes: 32000
        cors:
            AllowedOrigins: ["*"]
            AllowCredentials: true
            AllowedMethods: ["GET", "POST", "PATCH", "DELETE"]
            AllowedHeaders: ["*"]
            MaxAge: 1
            Debug: false
            OptionsPassthrough: false
  cache:
    DefaultExpiration: 10 // ttl протухания записей в кеше
    CleanupInterval: 30 // ttl выталкивания записей из кеша 
  aerospike:
    Hosts:
      - 192.168.0.2:3000
      - 192.168.0.5:3000
    Timeout: 5000,
    ConnectionTimeout: 5000,
    GetTimeout: 10,
    WriteTimeout: 10,
    Ttl: 0 // TTL жизни записи о пользователе
    Set: "users"
    NameSpace: "visitor"
  database:
    SxGeoCity: "SxGeoMax.dat"
    MaxMind: "GeoLite2-City.mmdb"
    BrowsCap: "full_php_browscap.ini"
```

```sh
$ npm install --production
$ NODE_ENV=production node app
```

#### Интеграция
**Visitor** имеет как `HTTP` так и `GRPC` интерфейс взаимодействия
- PHP библиотека https://bitbucket.org/videonow/visitor-client
- GO библиотека находится в этой же репе

На каждый запрос прийдет или сформированные/измененные данные или ошибка (влом писать дальше)

#### Примеры:
**Метод** `post`:
Создание записи о пользователе
*пример*:
```php
<?php
$curl = curl_init();
curl_setopt_array($curl, array(
  CURLOPT_URL => "http://visitor.videonow.ru/",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "POST",
  CURLOPT_POSTFIELDS => "{\n\t\"id\":\"63a406058d915def9d3135d9ee34e132\",\n\t\"ua\":\"Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19\",\n\t\"ip\": \"165.227.53.107\",\n\t\"extra\": {\n\t\t\"test\": 1\n\t}\n}",
  CURLOPT_HTTPHEADER => array(
    "cache-control: no-cache",
    "content-type: application/json",
  ),
));

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if ($err) {
  echo "cURL Error #:" . $err;
} else {
  echo $response;
}
```

**Метод** `get`:
Получение данных о пользователе
*пример*:
```php
<?php
$curl = curl_init();
curl_setopt_array($curl, array(
  CURLOPT_URL => "http://visitor.videonow.ru/63a406058d915def9d3135d9ee34e132",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "GET",
  CURLOPT_HTTPHEADER => array(
    "cache-control: no-cache",
    "content-type: application/json"
  ),
));

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if ($err) {
  echo "cURL Error #:" . $err;
} else {
  echo $response;
}
```

**Метод** `delete`:
Удаление пользователя
*пример*:
```php
<?php
$curl = curl_init();
curl_setopt_array($curl, array(
  CURLOPT_URL => "http://visitor.videonow.ru/63a406058d915def9d3135d9ee34e132",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "REMOVE",
  CURLOPT_HTTPHEADER => array(
    "cache-control: no-cache",
    "content-type: application/json"
  ),
));

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if ($err) {
  echo "cURL Error #:" . $err;
} else {
  echo $response;
}
```

**Метод** `patch`:
Удаление пользователя
*пример*:
```php
<?php

$curl = curl_init();

curl_setopt_array($curl, array(
  CURLOPT_URL => "http://visitor.videonow.ru/63a406058d915def9d3135d9ee34e132",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "PATCH",
  CURLOPT_POSTFIELDS => "{\n\t\"city\":{\n\t\t\"name\": \"test\"\n\t},\n\t\"extra\": {\n\t\t\"newData\": \"data\"\n\t}\n}",
  CURLOPT_HTTPHEADER => array(
    "cache-control: no-cache",
    "content-type: application/json",
  ),
));

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if ($err) {
  echo "cURL Error #:" . $err;
} else {
  echo $response;
}
```












