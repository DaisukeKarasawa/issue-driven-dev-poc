# Dialectic Lab Frontend

React + TypeScript + Vite frontend for Dialectic Lab.

## Scripts

```bash
npm install
npm run dev
npm run lint
npm run test
npm run build
```

## Runtime notes

- Dev server runs at `http://localhost:5173`.
- `/api` requests are proxied to `http://localhost:8080`.
- The backend must be running to execute graph evaluation from the UI.

## Source layout

- `src/App.tsx`: app orchestration and evaluation flow
- `src/components/*`: editors, parameter panel, result panel, sample loader
- `src/api.ts`: API client + error handling
- `src/types.ts`: shared frontend types

For full project overview, see the repository root `README.md`.
