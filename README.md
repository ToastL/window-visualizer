# Window Visualizer

A real-time, interactive visualization of your windows using a dynamic, flowing grid.

![Visualizer Preview](https://via.placeholder.com/800x450?text=Flowing+Grid+Visualizer)

## Overview

This project tracks the position and dimensions of active Google Chrome windows and visualizes them as disturbances in a flowing, interactive grid. The grid features a modern dark aesthetic with neon gradients and organic, randomized drift effects.

## Tech Stack

- **Backend**: Go (Golang)
    - Uses `CGO` and `AppleScript` (`osascript`) to query window metadata.
    - `Gorilla WebSocket` for real-time communication.
- **Frontend**: Vue 3 + TypeScript + Vite
    - `HTML5 Canvas` for high-performance rendering.
    - Custom physics-like animation loop for the flowing grid effect.

## Prerequisites

- **macOS**: Required for the window tracking logic (uses AppleScript).
- **Go**: 1.25+
- **Node.js**: 18+

## Getting Started

### 1. Start the Backend

The backend listens for window updates and broadcasts them via WebSockets.

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

Server will start on `ws://localhost:8080`.

### 2. Start the Frontend

The frontend visualizes the data.

```bash
cd frontend
npm install
npm run dev
```

Open the provided localhost URL (usually `http://localhost:5173`) in your browser.

## Features

- **Real-time Tracking**: Updates window positions at ~60fps.
- **Flowing Grid**: Grid points drift organically using a randomized sine-wave algorithm.
- **Interactive**: The grid reacts to the presence of Chrome windows.
- **Premium Aesthetic**: Deep slate background, cyan-violet gradients, and glowing effects.
- **Responsive**: Canvas resizes instantly with the window.

## Project Structure

```
.
├── backend/
│   ├── cmd/server/      # Entry point
│   ├── internal/        # Application logic (System, WebSocket, Config)
│   └── go.mod
└── frontend/
    ├── src/
    │   ├── components/  # Vue components (CanvasVisualizer)
    │   ├── composables/ # Logic reuse (useWebSocket)
    │   └── config.ts    # Visual configuration
    └── package.json
```
