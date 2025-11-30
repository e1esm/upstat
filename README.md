<div align="center" width="100%">
    <img src="./docs/assets/upstat.png" width="128" alt="" />
</div>

# Upstat

> Простой и удобный самодостаточный инструмент мониторинга статуса

![](./docs/assets/dashboard.png)

## 💻 Живая демонстрация

Попробуйте сами.

Демо-сервер (расположение: Сингапур): [https://demo.upstat.com](https://upstat.chamanbudhathoki.com.np/)

Имя пользователя: `demo`
Пароль: `demodemo`

## ⭐ Возможности

Пока что функций не много, но вот текущий список:

- Мониторинг доступности HTTP(s)
- Графики статуса и задержек
- Уведомления через Discord
- 60-секундные интервалы проверки
- Красивый, реактивный и быстрый интерфейс
- Несколько страниц статуса
- Привязка страниц статуса к конкретным доменам
- График пинга
- Информация о сертификатах
- PWA (Progressive Web App)
- Поддержка баз данных Sqlite и Postgres

И десятки более мелких функций, которые будут добавлены.

## 🔧 Установка

### 🐳 Docker

Для Sqlite

```bash
curl https://raw.githubusercontent.com/chamanbravo/upstat/main/docker-compose-sqlite.yml -o docker-compose.yml
docker compose up
```

Для Postgres

```bash
curl -O https://raw.githubusercontent.com/chamanbravo/upstat/main/docker-compose.yml
docker compose up
```

Upstat теперь работает по адресу http://localhost:3000

> [!ВАЖНО]
> Не забудьте изменить значения переменных окружения перед развертыванием.

### 💪🏻 Без Docker

Требования:

- Node.js 14 / 16 / 18 / 20.4
- npm 9
- Golang 1.21+
- Postgres (опционально)

```shell
cp .sample.env .env
```

```shell
air
cd web && npm run dev
```

## Технологический стек

- React
- Shadcn
- Golang
- Postgres/Sqlite

## 🙌 Участие в разработке

Я приветствую contributions! Вклад в развитие - это то, что делает сообщество открытого исходного кода таким удивительным местом для обучения, вдохновения и творчества. Любой ваш вклад **чрезвычайно ценится**.

Если у вас есть предложение, как улучшить проект, пожалуйста, сделайте форк репозитория, внесите изменения и создайте pull request. Вы также можете просто открыть issue с тегом "enhancement".
Не забудьте поставить звезду проекту! Заранее спасибо!

1. Сделайте форк проекта
2. Создайте ветку для вашей функции (`git checkout -b feature/AmazingFeature`)
3. Зафиксируйте изменения (`git commit -m 'Add some AmazingFeature'`)
4. Отправьте изменения в ветку (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request


## 🖼 Дополнительные скриншоты

Создание монитора

<img src="./docs/assets/create.png" width="512" alt="" />

Страница монитора

<img src="./docs/assets/chart.png" width="512" alt="" />

Страница настроек

<img src="./docs/assets/settings.png" width="512" alt="" />

Уведомления

<img src="./docs/assets/notifications.png" width="512" alt="" />

<img src="./docs/assets/discord_notification.png" width="512" alt="" />