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
    
    authSection.innerHTML = '';
    mainNav.innerHTML = `
        <a href="#catalog" class="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-200 rounded-md">
            <i class="fas fa-film w-6 text-center"></i><span class="ml-3">Catálogo</span>
        </a>
    `;

    if (userState && userState.token) {
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
            // CORRECCIÓN: Se ha añadido el contenido correcto para el enlace "Subir Video".
            mainNav.innerHTML += `
                <a href="#upload" class="flex items-center px-4 py-2 mt-2 text-gray-700 hover:bg-gray-200 rounded-md">
                    <i class="fas fa-upload w-6 text-center"></i><span class="ml-3">Subir Video</span>
                </a>
                <a href="#dashboard" class="flex items-center px-4 py-2 mt-2 text-gray-700 hover:bg-gray-200 rounded-md">
                    <i class="fas fa-tachometer-alt w-6 text-center"></i><span class="ml-3">Dashboard</span>
                </a>
            `;
        }
    } else {
        authSection.innerHTML = `
            <a href="#login" class="block w-full text-center bg-indigo-600 text-white font-bold py-2 px-4 rounded-md hover:bg-indigo-700">Iniciar Sesión</a>
            <a href="#register" class="block w-full text-center text-indigo-600 font-bold py-2 px-4 mt-2">Registrarse</a>
        `;
    }
}

/**
 * Renderiza la lista de videos en el catálogo.
 * @param {Array} videos - El array de objetos de video.
 */
export function renderVideos(videos) {
    const videoList = document.getElementById('videoList');
    videoList.innerHTML = '';
    if (videos && videos.length > 0) {
        videos.forEach(video => {
            const videoElement = document.createElement('div');
            videoElement.className = 'bg-white rounded-lg shadow-md overflow-hidden transform hover:-translate-y-1 transition-transform duration-300';
            videoElement.innerHTML = `
                <video controls class="w-full h-auto bg-black" preload="metadata">
                    <source src="/stream/${video.file_path}" type="video/mp4">
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
        videoList.innerHTML = '<p class="text-gray-500 col-span-full text-center">No hay videos disponibles. ¡Sube el primero!</p>';
    }
}

// --- FUNCIONES PARA EL PANEL DE ADMIN ---

export function renderAdminUsers(users) {
    const userList = document.getElementById('admin-user-list');
    userList.innerHTML = '';
    users.forEach(user => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${user.id}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${user.username}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${user.email}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${user.role}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button data-action="delete-user" data-id="${user.id}" class="text-red-600 hover:text-red-900">Eliminar</button>
            </td>
        `;
        userList.appendChild(row);
    });
}

export function renderAdminVideos(videos) {
    const videoList = document.getElementById('admin-video-list');
    videoList.innerHTML = '';
    videos.forEach(video => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${video.id}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${video.title}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${video.category}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button data-action="delete-video" data-id="${video.id}" class="text-red-600 hover:text-red-900">Eliminar</button>
            </td>
        `;
        videoList.appendChild(row);
    });
}

// --- FUNCIONES PARA EL MODAL DE CONFIRMACIÓN ---

let confirmCallback = null;
const modalBackdrop = document.getElementById('modal-backdrop');
const modalConfirm = document.getElementById('modal-confirm');
const modalCancel = document.getElementById('modal-cancel');

export function showModal(title, message, onConfirm) {
    document.getElementById('modal-title').textContent = title;
    document.getElementById('modal-message').textContent = message;
    confirmCallback = onConfirm;
    modalBackdrop.classList.remove('hidden');
    modalBackdrop.classList.add('flex');
}

function hideModal() {
    modalBackdrop.classList.add('hidden');
    modalBackdrop.classList.remove('flex');
    confirmCallback = null;
}

modalConfirm.addEventListener('click', () => {
    if (confirmCallback) {
        confirmCallback();
    }
    hideModal();
});
modalCancel.addEventListener('click', hideModal);
