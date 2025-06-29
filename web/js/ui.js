// File: web/js/ui.js

/**
 * Muestra una notificación en la esquina de la pantalla.
 * @param {string} message - El mensaje a mostrar.
 * @param {boolean} isError - Si la notificación es un error.
 */
export function showNotification(message, isError = false) {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `fixed bottom-5 right-5 w-80 p-4 rounded-lg shadow-lg text-white transition-all duration-300 ${isError ? 'bg-red-600' : 'bg-green-500'}`;
    notification.classList.remove('opacity-0', 'translate-y-10');
    setTimeout(() => {
        notification.classList.add('opacity-0', 'translate-y-10');
    }, 4000);
}

/**
 * Muestra una sección específica de la página y oculta las demás.
 * @param {string} sectionId - El ID de la sección a mostrar.
 */
export function showSection(sectionId) {
    document.querySelectorAll('.page-section').forEach(section => {
        section.classList.toggle('active', section.id === sectionId);
    });
}

/**
 * Actualiza la interfaz de autenticación (login/logout, bienvenida).
 * @param {object|null} userState - El estado del usuario (null si no está logueado).
 */
export function updateAuthUI(userState) {
    const authSection = document.getElementById('auth-section');
    const mainNav = document.getElementById('main-nav');
    
    // Limpiar el estado anterior
    authSection.innerHTML = '';
    mainNav.innerHTML = '';

    // Menú de navegación principal
    mainNav.innerHTML += `
        <a href="#catalog" class="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-200 rounded-md">
            <i class="fas fa-film w-6 text-center"></i><span class="ml-3">Catálogo</span>
        </a>
    `;

    if (userState && userState.token) {
        // Usuario logueado
        authSection.innerHTML = `
            <div class="flex items-center">
                <i class="fas fa-user-circle text-3xl text-gray-500"></i>
                <div class="ml-3">
                    <p class="text-sm font-semibold text-gray-800">${userState.username}</p>
                    <a href="#" id="logout-link" class="text-xs text-indigo-600 hover:underline">Cerrar sesión</a>
                </div>
            </div>
        `;
        if (userState.role === 'admin') {
            mainNav.innerHTML += `
                <a href="#upload" class="flex items-center px-4 py-2 mt-2 text-gray-700 hover:bg-gray-200 rounded-md">
                    <i class="fas fa-upload w-6 text-center"></i><span class="ml-3">Subir Video</span>
                </a>
            `;
        }
    } else {
        // Usuario no logueado
        authSection.innerHTML = `
            <a href="#login" class="block w-full text-center bg-indigo-600 text-white font-bold py-2 px-4 rounded-md hover:bg-indigo-700">Iniciar Sesión</a>
            <a href="#register" class="block w-full text-center text-indigo-600 font-bold py-2 px-4 mt-2">Registrarse</a>
        `;
    }
}

/**
 * Renderiza la lista de videos en el catálogo.
 * @param {Array} videos - El array de objetos de video.
 * @param {string} apiBaseUrl - La URL base de la API para construir los enlaces de streaming.
 */
export function renderVideos(videos, apiBaseUrl) {
    const videoList = document.getElementById('videoList');
    videoList.innerHTML = '';
    if (videos && videos.length > 0) {
        videos.forEach(video => {
            const videoElement = document.createElement('div');
            videoElement.className = 'bg-white rounded-lg shadow-md overflow-hidden transform hover:-translate-y-1 transition-transform duration-300';
            videoElement.innerHTML = `
                <video controls class="w-full h-auto bg-black" preload="metadata">
                    <source src="${apiBaseUrl}/stream/${video.file_path}" type="video/mp4">
                    Tu navegador no soporta la etiqueta de video.
                </video>
                <div class="p-4">
                    <h3 class="font-bold text-lg text-gray-800 truncate">${video.title}</h3>
                    <p class="text-sm text-gray-500 mt-1">Categoría: ${video.category}</p>
                </div>
            `;
            videoList.appendChild(videoElement);
        });
    } else {
        videoList.innerHTML = '<p class="text-gray-500 col-span-full text-center">No hay videos disponibles en este momento.</p>';
    }
}
