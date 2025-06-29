// File: web/js/api.js

const API_BASE_URL = 'http://localhost:8080';

/**
 * Función genérica para manejar las peticiones a la API.
 * @param {string} endpoint - El endpoint al que se hará la petición.
 * @param {object} options - Las opciones para la petición fetch.
 * @returns {Promise<any>} - La respuesta de la API.
 */
async function request(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    
    try {
        const response = await fetch(url, options);
        const isJson = response.headers.get('content-type')?.includes('application/json');
        const data = isJson ? await response.json() : await response.text();

        if (!response.ok) {
            const errorMessage = isJson ? (data.error || data.message || JSON.stringify(data)) : data;
            throw new Error(errorMessage);
        }
        return data;
    } catch (error) {
        console.error(`API Error en ${endpoint}:`, error);
        throw error; // Re-lanzamos el error para que el llamador pueda manejarlo.
    }
}

// Exportamos cada función de la API de forma individual.
export const loginUser = (email, password) => {
    return request('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    });
};

export const registerUser = (username, email, password) => {
    return request('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password })
    });
};

export const verifyUser = (token) => {
    return request('/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token })
    });
};

export const getVideos = () => {
    return request('/videos');
};

export const uploadVideo = (formData, token) => {
    return request('/admin/upload', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: formData
    });
};
