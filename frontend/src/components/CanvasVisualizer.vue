<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { CONFIG } from '../config';
import { useWebSocket } from '../composables/useWebSocket';

const id = crypto.randomUUID();
document.title = id;

const { message, connect } = useWebSocket(id);

let resizeHandler: (() => void) | null = null;
let animationFrameId: number | null = null;

const p2_distance = (x1: number, y1: number, x2: number, y2: number) => {
    return Math.sqrt((x1 - x2) ** 2 + (y1 - y2) ** 2);
};

function drawBullet(
    ctx: CanvasRenderingContext2D,
    x1: number,
    y1: number,
    x2: number,
    y2: number,
    scale: number = 2,
): void {
    if (scale < CONFIG.BULLET_MIN_SCALE) {
        return;
    }

    const dx = x2 - x1;
    const dy = y2 - y1;
    const distance = Math.sqrt(dx * dx + dy * dy);

    const offset = CONFIG.BULLET_OFFSET * scale;
    const newX1 = distance > 0 ? x1 + (dx / distance) * offset : x1;
    const newY1 = distance > 0 ? y1 + (dy / distance) * offset : y1;

    const gradient = ctx.createLinearGradient(x1, y1, newX1, newY1);
    gradient.addColorStop(0, CONFIG.COLORS.ACCENT_START);
    gradient.addColorStop(1, CONFIG.COLORS.ACCENT_END);

    ctx.strokeStyle = gradient;
    ctx.globalAlpha = Math.min(1, scale * 1.5);
    ctx.beginPath();
    ctx.moveTo(x1, y1);
    ctx.lineTo(newX1, newY1);
    ctx.lineCap = 'round';
    ctx.lineWidth = CONFIG.BULLET_LINE_WIDTH * scale;
    ctx.stroke();
    ctx.globalAlpha = 1.0;
}

onMounted(() => {
    connect();

    const canvas = document.getElementById('canvas') as HTMLCanvasElement | null;
    if (!canvas) {
        console.error('Canvas element not found');
        return;
    }

    const ctx = canvas.getContext('2d');
    if (!ctx) {
        console.error('Failed to get 2D context from canvas');
        return;
    }

    const canvasRef: HTMLCanvasElement = canvas;
    const ctxRef: CanvasRenderingContext2D = ctx;

    canvasRef.width = window.innerWidth;
    canvasRef.height = window.innerHeight;
    ctxRef.imageSmoothingEnabled = false;

    resizeHandler = () => {
        canvasRef.width = window.innerWidth;
        canvasRef.height = window.innerHeight;
    };
    window.addEventListener('resize', resizeHandler);

    const GRID_SIZE = CONFIG.GRID_SIZE;
    const GRID_SPACING = CONFIG.GRID_SPACING;
    const INTERPOLATION_FACTOR = CONFIG.INTERPOLATION_FACTOR;
    const GRID_WIDTH = GRID_SPACING * GRID_SIZE;
    const MAX_DISTANCE = CONFIG.MAX_DISTANCE;

    let fixedX = 0;
    let fixedY = 0;
    let fixedGridX = 0;
    let fixedGridY = 0;

    let time = 0;

    function loop() {
        time += 0.01;

        fixedX += (message.x - fixedX) * INTERPOLATION_FACTOR;
        fixedY += (message.y - fixedY) * INTERPOLATION_FACTOR;

        fixedGridX += (message.gridX - fixedGridX) * INTERPOLATION_FACTOR;
        fixedGridY += (message.gridY - fixedGridY) * INTERPOLATION_FACTOR;

        ctxRef.fillStyle = CONFIG.COLORS.BACKGROUND;
        ctxRef.fillRect(0, 0, canvasRef.width, canvasRef.height);

        ctxRef.globalCompositeOperation = 'lighter';

        const halfGridWidth = GRID_WIDTH / 2;

        for (let x = 0; x < GRID_SIZE; x++) {
            for (let y = 0; y < GRID_SIZE; y++) {
                const posX = fixedGridX + x * GRID_SPACING;
                const posY = fixedGridY + y * GRID_SPACING;

                const offsetX = Math.floor((halfGridWidth + fixedX - posX) / GRID_WIDTH);
                const offsetY = Math.floor((halfGridWidth + fixedY - posY) / GRID_WIDTH);

                let finalX = offsetX * GRID_WIDTH + posX;
                let finalY = offsetY * GRID_WIDTH + posY;

                const phase = Math.sin(x * 12.9898 + y * 78.233) * 43758.5453;

                const driftX = Math.sin(time * CONFIG.DRIFT.SPEED + phase) * CONFIG.DRIFT.AMPLITUDE;
                const driftY = Math.cos(time * CONFIG.DRIFT.SPEED + phase) * CONFIG.DRIFT.AMPLITUDE;

                finalX += driftX;
                finalY += driftY;

                const distance = p2_distance(fixedX, fixedY, finalX, finalY) / 2 / (GRID_SIZE / 10);

                const scale = Math.max(0, Math.min(1, (MAX_DISTANCE - distance) / (MAX_DISTANCE * 0.5)));

                drawBullet(ctxRef, finalX, finalY, fixedX, fixedY, scale);
            }
        }

        ctxRef.globalCompositeOperation = 'source-over';

        animationFrameId = requestAnimationFrame(loop);
    }

    loop();
});

onUnmounted(() => {
    if (animationFrameId !== null) {
        cancelAnimationFrame(animationFrameId);
    }

    if (resizeHandler) {
        window.removeEventListener('resize', resizeHandler);
    }
});
</script>

<template>
    <canvas id="canvas"></canvas>
</template>

<style scoped>
canvas {
    width: 100vw;
    height: 100vh;
    image-rendering: pixelated;
    image-rendering: crisp-edges;
}
</style>
