import { reactive, onUnmounted } from 'vue';
import { CONFIG } from '../config';
import type { MessageInfo } from '../types';

export function useWebSocket(id: string) {
    const message = reactive<MessageInfo>({
        id: '',
        x: 0,
        y: 0,
        gridX: 0,
        gridY: 0,
    });

    let ws: WebSocket | null = null;
    let reconnectAttempts = 0;

    const connect = () => {
        if (reconnectAttempts >= CONFIG.MAX_RECONNECT_ATTEMPTS) {
            console.error('Max reconnection attempts reached. Please refresh the page.');
            return;
        }

        try {
            ws = new WebSocket(CONFIG.WS_URL);

            ws.onopen = (event) => {
                console.log('WebSocket connected', event);
                reconnectAttempts = 0;
            };

            ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data) as MessageInfo[];

                    for (const item of data) {
                        if (item.id === id) {
                            message.x = item.x;
                            message.y = item.y;
                            message.gridX = item.gridX;
                            message.gridY = item.gridY;
                        }
                    }
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };

            ws.onclose = () => {
                console.log('WebSocket disconnected, attempting to reconnect...');
                reconnectAttempts++;
                setTimeout(connect, CONFIG.RECONNECT_DELAY);
            };
        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
            reconnectAttempts++;
            setTimeout(connect, CONFIG.RECONNECT_DELAY);
        }
    };

    onUnmounted(() => {
        if (ws) {
            ws.close();
            ws = null;
        }
    });

    return {
        message,
        connect,
    };
}
