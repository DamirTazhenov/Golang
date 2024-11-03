# Используем Node.js для сборки React-приложения
FROM node:14 AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json для установки зависимостей
COPY front-end/package.json front-end/package-lock.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем весь исходный код для сборки фронтенда
COPY front-end ./

# Собираем приложение
RUN npm run build

# Используем nginx для раздачи статических файлов
FROM nginx:alpine

# Копируем скомпилированные статические файлы в папку nginx
COPY --from=build /app/build /usr/share/nginx/html

# Экспонируем порт 80 для nginx
EXPOSE 80

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]
