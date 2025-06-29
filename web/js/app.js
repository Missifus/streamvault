// Importamos todas las funciones necesarias desde nuestros módulos de API y UI.
// Esto mantiene nuestro código organizado y con responsabilidades claras.
import * as api from './api.js';
import * as ui from './ui.js';

// --- ESTADO GLOBAL DE LA APLICACIÓN ---
// Un objeto que contiene el estado actual de la sesión del usuario.
// Al cargar la página, intenta recuperar los datos desde el localStorage del navegador,
// lo que permite que la sesión del usuario persista entre recargas.
const state = {
    token: localStorage.getItem('authToken') || null,
    role: localStorage.getItem('userRole') || null,
    username: localStorage.getItem('username') || null,
};


// --- FUNCIONES AUXILIARES ---

/**
 * Decodifica un token JWT para extraer la información (payload) del usuario.
 * @param {string} token - El token JWT.
 * @returns {object|null} - El objeto con los datos del usuario o null si hay un error.
 */
function parseJwt(token) {
    try {
        // El token se divide en 3 partes por puntos. El payload está en la segunda parte.
        // atob() decodifica la cadena de base64 a un string JSON.
        return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
        // Si el token es inválido o malformado, devuelve null.
        return null;
    }
}

/**
 * Cierra la sesión del usuario, limpiando el estado y el almacenamiento local.
 */
function handleLogout() {
    // Resetea el estado global de la aplicación.
    state.token = null;
    state.role = null;
    state.username = null;
    // Limpia todos los datos guardados en el almacenamiento del navegador.
    localStorage.clear();
    // Actualiza la interfaz para reflejar que el usuario ha cerrado sesión.
    ui.updateAuthUI(state);
    // Redirige al usuario a la página de login.
    window.location.hash = '#login';
    // Muestra una notificación de éxito.
    ui.showNotification('Has cerrado sesión.');
}


// --- LÓGICA DE EVENTOS Y FORMULARIOS ---

/**
 * Asigna los "escuchadores de eventos" a todos los formularios de la aplicación.
 * Esta función se llama una sola vez cuando la página carga para configurar la interactividad.
 */
function addFormListeners() {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault(); // Previene que la página se recargue.
            const email = document.getElementById('login-email').value;
            const password = document.getElementById('login-password').value;
            try {
                const data = await api.loginUser(email, password); // Llama a la API para hacer login.
                const decodedToken = parseJwt(data.token);
                // Actualiza el estado global con los datos del usuario.
                state.token = data.token;
                state.role = decodedToken.role;
                state.username = decodedToken.username;
                // Guarda los datos en el navegador para persistir la sesión.
                localStorage.setItem('authToken', state.token);
                localStorage.setItem('userRole', state.role);
                localStorage.setItem('username', state.username);
                
                ui.updateAuthUI(state);
                window.location.hash = '#catalog';
                ui.showNotification(`¡Bienvenido, ${state.username}!`);
                e.target.reset(); // Limpia el formulario.
            } catch (error) {
                ui.showNotification(error.message, true); // Muestra errores (ej: "contraseña incorrecta").
            }
        });
    }

    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('register-username').value;
            const email = document.getElementById('register-email').value;
            const password = document.getElementById('register-password').value;
            try {
                const data = await api.registerUser(username, email, password);
                ui.showNotification(data.message);
                window.location.hash = '#login'; // Redirige al login después de un registro exitoso.
                e.target.reset();
            } catch (error) {
                ui.showNotification(error.message, true); // Muestra errores (ej: "email ya en uso").
            }
        });
    }
    
    const uploadForm = document.getElementById('uploadForm');
    if (uploadForm) {
        uploadForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            // Usamos FormData para enviar archivos y texto juntos.
            const formData = new FormData();
            formData.append('title', document.getElementById('title').value);
            formData.append('category', document.getElementById('category').value);
            const videoFile = document.getElementById('videoFile').files[0];
            
            if (!videoFile) {
                ui.showNotification('Por favor, selecciona un archivo de video.', true);
                return;
            }

            formData.append('video', videoFile);
            
            try {
                await api.uploadVideo(formData, state.token);
                ui.showNotification('Video subido exitosamente.');
                e.target.reset();
                loadVideos(); // Recarga la lista de videos para mostrar el nuevo.
                window.location.hash = '#catalog';
            } catch (error) {
                ui.showNotification(error.message, true);
            }
        });
    }

    const authSection = document.getElementById('auth-section');
    if (authSection) {
        // Evento para el botón de Logout (se añade al contenedor padre porque el botón se crea dinámicamente).
        authSection.addEventListener('click', (e) => {
            if (e.target.id === 'logout-link') {
                e.preventDefault();
                handleLogout();
            }
        });
    }
}


// --- LÓGICA DEL PANEL DE ADMINISTRACIÓN ---

/**
 * Carga los datos necesarios para el panel de administración (usuarios y videos).
 */
async function loadDashboardData() {
    try {
        // Usamos Promise.all para hacer las dos peticiones a la API en paralelo, mejorando el rendimiento.
        const [users, videos] = await Promise.all([
            api.getAdminUsers(state.token),
            api.getVideos()
        ]);
        // Llama a las funciones de la UI para renderizar las tablas.
        ui.renderAdminUsers(users);
        ui.renderAdminVideos(videos);
    } catch (error) {
        ui.showNotification(error.message, true);
    }
}

/**
 * Configura la navegación por pestañas dentro del panel de administración.
 */
function handleDashboardTabs() {
    const tabs = document.querySelectorAll('#dashboard-tabs a');
    const tabContents = document.querySelectorAll('#dashboard-content > div');
    
    if (tabs.length === 0) return; // Si no hay pestañas, no hace nada.

    tabs.forEach(tab => {
        tab.addEventListener('click', e => {
            e.preventDefault();
            // Actualiza los estilos para resaltar la pestaña activa.
            tabs.forEach(item => item.classList.remove('text-indigo-600', 'border-indigo-600'));
            tab.classList.add('text-indigo-600', 'border-indigo-600');
            // Muestra el contenido de la pestaña seleccionada y oculta las demás.
            const tabId = tab.getAttribute('data-tab');
            tabContents.forEach(content => {
                content.classList.toggle('hidden', content.id !== `${tabId}-tab`);
            });
        });
    });
    tabs[0].click(); // Activa la primera pestaña por defecto.
}

/**
 * Asigna eventos a las tablas del dashboard usando delegación de eventos.
 * Es más eficiente que añadir un listener a cada botón individualmente.
 */
function addDashboardListeners() {
    const dashboardContent = document.getElementById('dashboard-content');
    if (!dashboardContent) return; // Si el elemento no existe, no hace nada.

    dashboardContent.addEventListener('click', e => {
        const action = e.target.getAttribute('data-action');
        const id = e.target.getAttribute('data-id');

        if (!action || !id) return; // Si no se hizo clic en un botón de acción, sale.

        if (action === 'delete-user') {
            ui.showModal('Confirmar Eliminación', `¿Estás seguro de que quieres eliminar al usuario con ID ${id}?`, async () => {
                try {
                    await api.deleteUser(id, state.token);
                    ui.showNotification('Usuario eliminado.');
                    loadDashboardData(); // Recarga los datos para reflejar el cambio.
                } catch (error) {
                    ui.showNotification(error.message, true);
                }
            });
        }

        if (action === 'delete-video') {
            ui.showModal('Confirmar Eliminación', `¿Estás seguro de que quieres eliminar el video con ID ${id}?`, async () => {
                try {
                    await api.deleteVideo(id, state.token);
                    ui.showNotification('Video eliminado.');
                    loadDashboardData(); // Recarga los datos.
                } catch (error) {
                    ui.showNotification(error.message, true);
                }
            });
        }
    });
}


// --- LÓGICA PRINCIPAL Y ENRUTAMIENTO ---

/**
 * Carga la lista de videos desde la API y la muestra en la página.
 */
async function loadVideos() {
    try {
        const videos = await api.getVideos();
        ui.renderVideos(videos);
    } catch (error) {
        document.getElementById('videoList').innerHTML = `<p class="text-red-500 col-span-full text-center">No se pudieron cargar los videos. ${error.message}</p>`;
    }
}

/**
 * El "enrutador" de nuestra aplicación de una sola página (SPA).
 * Lee el hash de la URL (ej: #login) y muestra la sección correspondiente.
 */
function router() {
    const hash = window.location.hash || '#catalog';
    const sectionId = (hash.substring(1) || 'catalog') + '-section';

    // Si la ruta es el dashboard y el usuario es admin, carga los datos.
    if (sectionId === 'dashboard-section' && state.role === 'admin') {
        loadDashboardData();
    }
    
    // Muestra la sección correcta de la página.
    ui.showSection(sectionId);
}

/**
 * La función principal que se ejecuta al cargar la página.
 * Inicializa toda la aplicación.
 */
function main() {
    // Escucha los cambios en el hash de la URL (ej: cuando el usuario hace clic en un enlace).
    window.addEventListener('hashchange', router);
    
    // Configura todos los listeners de los formularios y del dashboard.
    addFormListeners();
    handleDashboardTabs();
    addDashboardListeners();
    
    // Ejecuta las funciones iniciales para configurar el estado de la página.
    router();
    ui.updateAuthUI(state);
    loadVideos();
}

// El punto de entrada de toda la aplicación.
// Se asegura de que el DOM esté completamente cargado antes de ejecutar el script.
document.addEventListener('DOMContentLoaded', main);
