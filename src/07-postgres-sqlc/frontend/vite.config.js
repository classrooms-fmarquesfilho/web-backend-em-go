import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// API_TARGET permite trocar o destino do proxy:
//   - desenvolvimento local:    http://localhost:8080  (default)
//   - dentro de docker-compose: http://backend:8080
export default defineConfig({
  plugins: [sveltekit()],
  server: {
    host: '0.0.0.0',          // necessário dentro de container
    proxy: {
      // Tudo em /api/* vai para a API Go,
      // removendo o prefixo /api antes de enviar.
      '/api': {
        target: process.env.API_TARGET || 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, '')
      }
    }
  }
});
