export const CONFIG = {
    WS_URL: import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws',
    GRID_SIZE: 20,
    GRID_SPACING: 40,
    INTERPOLATION_FACTOR: 1 / 60,
    MAX_DISTANCE: 100,
    BULLET_OFFSET: 5,
    BULLET_MIN_SCALE: 0.05,
    BULLET_LINE_WIDTH: 3,
    COLORS: {
        BACKGROUND: '#0f172a',
        GRID_BASE: '#000',
        ACCENT_START: '#06b6d4',
        ACCENT_END: '#8b5cf6',
    },
    DRIFT: {
        SPEED: 0.001,
        AMPLITUDE: 15,
    },
    RECONNECT_DELAY: 2000,
    MAX_RECONNECT_ATTEMPTS: 10,
} as const;
