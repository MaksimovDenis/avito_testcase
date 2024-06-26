# Avtio 🚀  

Веб-сервер "Сервис баннеров"  

# Запуск  

1. Сколнировать репозиторий:
```bash   
git clone https://github.com/MaksimovDenis/avito_testcase.git
```
2. Перейти в директорию проекта (если Вы не в ней):  
```bash    
cd avito_testcase 
```
3. Из дериктории проекта выполнить команды:  
```bash      
make build
```
```bash    
make up 
```   
4. Swagger документация доступна по ссылке:  
```bash   
http://localhost:8000/swagger/index.html
```  
5. Если вам необходимо перезапустить проект:  
```bash      
make chmod
```
```bash      
make restart
```
6. Остальные функции прописаны в Makefile:  
```bash      
make help
```  

# Реализация  
- Подход с чистой архитектурой  
- Язык программирование: Golang 1.21.6  
- БД - PostgreSQL latest  
- КЭШ - Redis latest  
- Для реализации http сервера использовались стандартные библиотеки go.
- Для деплоя зависимостей и самого сервиса использовался Docker и Docker-compose  
- Структура БД и схема кеширования, которая применялась в проекте:    
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/164db992-db6c-404b-a74a-37cd5682c829)  
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/c46ef32f-0576-48bd-b0cf-ba7904e1d350)  

# Ход решения  
- Основные условия ТЗ - выполнены (в конце этого файла будет приведены примеры запросов)   

- Был изменён подход к авторизации, вместо явной передачи токена в качестве параметра, добавил полную систему аутентификации и авторизации. При регистрации указывается роль (admin или обычный пользователь). В результате пользователь получает JWT токен, далее пользователю доступен функционал сервера, в зависимости от роли.   
- Из дополнительных условий - были реализованы тесты для остальных сценариев.
- Проведено нагрузочное тестирование в Postman по сценарию получения баннеров в количестве = 5000.  
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/517b2ac0-74ea-445d-b989-0ac88b8d6c34)  

# Возникшие проблемы и что следует улучшить  
- Проблему скорости ответа удалась решить за счёт кеширование через Redis.
- Метод получения всех списков баннеров, с учётом вложенного списка tagID вероятно можно написать более оптимально. 
- Стоит улучшить систему валидации полей, основные сценарии учтены, но не все. 
- Необходимо улучшить систему обработки всех возможных ошибок.
- Не все дополнительные пункты задания были выполнены.   

# Примеры запросов  
- Авторизация и аутентификация:  
 ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/74ce7129-8129-4898-91b8-8d4ba3304a1c)  
 ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/6b048a32-ff88-464f-bd8c-0049e80d4c19)

 - Получение баннеров по фиче и тегу:  
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/4f075dcf-dd69-4485-8865-304802f7ee0d)

Если необходима последняя версия баннера, то получать её пользователь будет напрямую из БД, а не из кеша.
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/91bf8e8e-0583-458a-8c12-397539f0590b)

Если баннер не активен, то он будет доступен пользователям только со статусом админинстратор. 

 - Получение списка баннеров по фиче и/или тегу:  
 ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/b54e7f92-6ad4-46b5-8e97-850bcff97b1f)  
 ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/9b9f8bf2-79a6-492b-b57f-bbfc7ab26763)  
 ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/c33d3980-e7f1-4aa2-b1ea-894fd0c8bc0a)  
 ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/6c945435-26df-492b-b02e-1312b8413fbe)  

  - Создание баннера:
  ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/572f9a41-f72c-49da-aac4-deb7067220f3)  
  Т.к. баннер определяется фичей и тегом, то создать такой же не допускается  
  ![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/930a6922-b6ee-46e8-8dd3-982bf84ac2df)  
  
  - Обновление баннера:  
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/601c42b6-9b83-4e2e-8719-3cd7fc6574d9)  
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/8e2f3818-8458-4e6d-91cd-2123fabbc4a1)  
  
  - Удаление баннера:
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/35de93ca-27e4-440a-8e70-8d39cb3cb796)  
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/0023987d-8879-4249-ab2c-b8f3b5ab00e0)  

# Метрики сервиса
 - Связка prometheus + grafana, паттерн 4 golden signals.
![image](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/e9908c1f-2e6e-4e2d-b223-bc53b4f6bfbf)






