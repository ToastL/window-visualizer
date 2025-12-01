<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { CONFIG } from '../config';
import { useWebSocket } from '../composables/useWebSocket';
import type { Vector2 } from '../types';

const id = crypto.randomUUID();
document.title = id;

const { message, connect } = useWebSocket(id);

let resizeHandler: (() => void) | null = null;
let animationFrameId: number | null = null;

const p2_distance = (p1: Vector2, p2: Vector2) => {
    return Math.sqrt((p1.x - p2.x) ** 2 + (p1.y - p2.y) ** 2);
};

function drawBullet(
    ctx: CanvasRenderingContext2D,
    p1: Vector2,
    p2: Vector2,
    scale: number = 2,
): void {
    if (scale < CONFIG.BULLET_MIN_SCALE) {
        return;
    }

    const dx = p2.x - p1.x;
    const dy = p2.y - p1.y;
    const distance = Math.sqrt(dx * dx + dy * dy);

    const offset = CONFIG.BULLET_OFFSET * scale;
    const newPos: Vector2 = {
        x: distance > 0 ? p1.x + (dx / distance) * offset : p1.x,
        y: distance > 0 ? p1.y + (dy / distance) * offset : p1.y,
    };

    const gradient = ctx.createLinearGradient(p1.x, p1.y, newPos.x, newPos.y);
    gradient.addColorStop(0, CONFIG.COLORS.ACCENT_START);
    gradient.addColorStop(1, CONFIG.COLORS.ACCENT_END);

    ctx.strokeStyle = gradient;
    ctx.globalAlpha = Math.min(1, scale * 1.5);
    ctx.beginPath();
    ctx.moveTo(p1.x, p1.y);
    ctx.lineTo(newPos.x, newPos.y);
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

    let fixedPosition: Vector2 = { x: 0, y: 0 };
    let fixedGridPosition: Vector2 = { x: 0, y: 0 };

    let time = 0;

    function loop() {
        time += 0.01;

        fixedPosition.x += (message.x - fixedPosition.x) * INTERPOLATION_FACTOR;
        fixedPosition.y += (message.y - fixedPosition.y) * INTERPOLATION_FACTOR;

        fixedGridPosition.x += (message.gridX - fixedGridPosition.x) * INTERPOLATION_FACTOR;
        fixedGridPosition.y += (message.gridY - fixedGridPosition.y) * INTERPOLATION_FACTOR;

        ctxRef.fillStyle = CONFIG.COLORS.BACKGROUND;
        ctxRef.fillRect(0, 0, canvasRef.width, canvasRef.height);

        ctxRef.globalCompositeOperation = 'lighter';

        const halfGridWidth = GRID_WIDTH / 2;

        for (let x = 0; x < GRID_SIZE; x++) {
            for (let y = 0; y < GRID_SIZE; y++) {
                const pos: Vector2 = {
                    x: fixedGridPosition.x + x * GRID_SPACING,
                    y: fixedGridPosition.y + y * GRID_SPACING,
                };

                const offset: Vector2 = {
                    x: Math.floor((halfGridWidth + fixedPosition.x - pos.x) / GRID_WIDTH),
                    y: Math.floor((halfGridWidth + fixedPosition.y - pos.y) / GRID_WIDTH)
                };

                let finalPos: Vector2 = {
                    x: offset.x * GRID_WIDTH + pos.x,
                    y: offset.y * GRID_WIDTH + pos.y,
                };

                const distance = p2_distance(fixedPosition, finalPos) / 2 / (GRID_SIZE / 10);

                const scale = Math.max(0, Math.min(1, (MAX_DISTANCE - distance) / (MAX_DISTANCE * 0.5)));

                drawBullet(ctxRef, finalPos, fixedPosition, scale);
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

<style>
html, body, #app {
    margin: 0;
    padding: 0;

    overflow: hidden;
}

canvas {
    image-rendering: pixelated;
    image-rendering: crisp-edges;
}
</style>
