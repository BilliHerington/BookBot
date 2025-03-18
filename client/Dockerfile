# 1. Используем Node.js для сборки фронта
FROM node:18 AS builder

WORKDIR /app

# 2. Копируем package.json и устанавливаем зависимости
COPY package.json package-lock.json ./
RUN npm install --frozen-lockfile

# 3. Копируем код и билдим проект
COPY . .
RUN npm run build

# 4. Используем Nginx для сервинга статики
FROM nginx:1.25 AS production

# 5. Копируем билд React в папку Nginx
COPY --from=builder /app/build /usr/share/nginx/html

# 6. Копируем свой конфиг Nginx в правильное место
COPY nginx.conf /etc/nginx/nginx.conf

# 7. Открываем нужный порт
EXPOSE 3000

# 8. Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]
