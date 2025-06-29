// La URL base de nuestra API. Como el frontend y el backend se sirven desde el mismo
// dominio (localhost:8080), podemos usar una ruta relativa. Todas las peticiones
// se dirigirán a /api/..., coincidiendo con la configuración del enrutador de Go.
const API_BASE_URL = '/api';

/**
 * Función genérica y reutilizable para manejar todas las peticiones a la API.
 * Su propósito es centralizar la lógica de 'fetch' y el manejo de errores.
 * @param {string} endpoint - El endpoint específico al que se hará la petición (ej: '/login').
 * @param {object} options - Las opciones para la petición fetch (método, cabeceras, cuerpo, etc.).
 * @returns {Promise<any>} - La respuesta de la API ya convertida a JSON.
 */
async function request(endpoint, options = {}) {
    // Construimos la URL completa para la petición.
    const url = `${API_BASE_URL}${endpoint}`;
    
    try {
        // Ejecutamos la petición al servidor.
        const response = await fetch(url, options);
        
        // Leemos la respuesta como JSON, ya que hemos estandarizado que el backend siempre responda así.
        const data = await response.json();

        // Verificamos si la respuesta del servidor fue exitosa (ej. código 200, 201).
        if (!response.ok) {
            // Si no fue exitosa, usamos el mensaje de error que viene en el JSON del backend.
            // Esto permite mostrar errores específicos como "Email ya en uso".
            const errorMessage = data.message || `Error ${response.status}`;
            throw new Error(errorMessage);
        }
        
        // Si todo fue exitoso, devolvemos los datos.
        return data;
        
    } catch (error) {
        // Capturamos cualquier error que ocurra durante la petición.
        console.error(`API Error en ${endpoint}:`, error);

        // Si el error es de tipo TypeError, generalmente significa que hubo un problema de red
        // (ej. el servidor de Go no está funcionando).
        if (error instanceof TypeError) {
            throw new Error('No se pudo conectar con el servidor. ¿Está funcionando?');
        }
        
        // Re-lanzamos el error para que la función que llamó a 'request' (en app.js)
        // pueda capturarlo y mostrarlo al usuario en una notificación.
        throw error;
    }
}

// --- EXPORTACIÓN DE FUNCIONES ESPECIALIZADAS ---
// Cada una de estas funciones utiliza el 'request' genérico para una tarea específica.
// Esto hace que el código en app.js sea mucho más limpio y fácil de leer.

/**
 * Envía las credenciales para iniciar sesión.
 */
export const loginUser = (email, password) => {
    return request('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    });
};

/**
 * Envía los datos para registrar un nuevo usuario.
 */
export const registerUser = (username, email, password) => {
    return request('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password })
    });
};

/**
 * Obtiene la lista de todos los videos.
 */
export const getVideos = () => {
    return request('/videos');
};

/**
 * Sube un nuevo video.
 * @param {FormData} formData - El objeto FormData que contiene el título, categoría y el archivo.
 * @param {string} token - El token JWT del administrador.
 */
export const uploadVideo = (formData, token) => {
    // Para FormData, no se establece 'Content-Type', el navegador lo hace automáticamente
    // junto con el 'boundary' necesario para la subida de archivos.
    return request('/admin/upload', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: formData
    });
};

/**
 * Obtiene la lista de todos los usuarios (solo para admins).
 */
export const getAdminUsers = (token) => {
    return request('/admin/users', {
        headers: { 'Authorization': `Bearer ${token}` }
    });
};

/**
 * Elimina un usuario (solo para admins).
 */
export const deleteUser = (userId, token) => {
    return request(`/admin/users/${userId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
    });
};

/**
 * Actualiza el rol de un usuario (solo para admins).
 */
export const updateUserRole = (userId, role, token) => {
    return request(`/admin/users/${userId}/role`, {
        method: 'PUT',
        headers: { 
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ role })
    });
};

/**
 * Elimina un video (solo para admins).
 */
export const deleteVideo = (videoId, token) => {
    return request(`/admin/videos/${videoId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
    });
};
