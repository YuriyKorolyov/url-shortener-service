# Авторизация и тестирование в Postman

## Ограничение доступа

- **POST /url** — сохранение URL (без авторизации).
- **GET /url/{alias}** — редирект по alias (без авторизации).
- **DELETE /url/{alias}** — удаление URL (требуется авторизация).

## Как тестировать с Postman

1. Задайте в конфиге или переменной окружения `TOKEN` (например, `my-secret-token`).

2. **Сохранение URL (без токена)**  
   - Method: `POST`  
   - URL: `http://localhost:8082/url`  
   - Body → raw → JSON:
   ```json
   { "url": "https://example.com", "alias": "ex" }
   ```

3. **Редирект (без токена)**  
   - Method: `GET`  
   - URL: `http://localhost:8082/url/ex`

4. **Удаление (с токеном)**  
   - Method: `DELETE`  
   - URL: `http://localhost:8082/url/ex`  
   - Headers:
   - `Authorization`: `Bearer my-secret-token`

5. **Удаление без токена** — ответ `401 Unauthorized`.

6. **Удаление с неверным токеном** — ответ `401 Unauthorized`.
