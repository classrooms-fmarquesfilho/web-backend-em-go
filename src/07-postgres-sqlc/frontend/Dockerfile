# Imagem leve com Node 22 para rodar o dev server do Vite.
#
# Para a aula, usamos o dev server (HMR, sem build). Para produção real,
# faríamos build → adapter-node → servir o binário com Node.
FROM node:22-alpine

WORKDIR /app

# Cache de dependências
COPY package.json package-lock.json ./
RUN npm ci

# Código
COPY . .

EXPOSE 5173

# O `--host 0.0.0.0` está em vite.config.js (server.host),
# necessário para Vite responder a requests vindos de fora do container.
CMD ["npm", "run", "dev"]
