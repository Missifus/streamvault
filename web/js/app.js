// File: web/js/app.js

// Importamos las funciones desde nuestros módulos
import * as api from './api.js';
import * as ui from './ui.js';

// Estado global de la aplicación
const state = {
    token: localStorage.getItem('authToken') || null,
    role: localStorage.getItem('userRole') || null,
    username: localStorage.getItem('username') || null,
};

function parseJwt(token) {
    try {
        return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
        return null;
    }
}

function handleLogout() {
    state.token = null;
    state.role = null;
    state.username = null;
    localStorage.clear();
    ui.updateAuthUI(state);
    window.location.hash = '#login';
    ui.showNotification('Has cerrado sesión.');
}

function addFormListeners() {
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        try {
            const data = await api.loginUser(email, password);
            const decodedToken = parseJwt(data.token);
            state.token = data.token;
            state.role = decodedToken.role;
            state.username = decodedToken.username;
            localStorage.setItem('authToken', state.token);
            localStorage.setItem('userRole', state.role);
            localStorage.setItem('username', state.username);
            ui.updateAuthUI(state);
            window.location.hash = '#catalog';
            ui.showNotification(`¡Bienvenido, ${state.username}!`);
            e.target.reset();
        } catch (error) { /* el error ya se muestra en la notificación desde api.js */ }
    });

    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('register-username').value;
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;
        try {
            const data = await api.registerUser(username, email, password);
            ui.showNotification(data.message);
            window.location.hash = '#login';
            e.target.reset();
        } catch (error) { /* ... */ }
    });
    
    document.getElementById('uploadForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData();
        formData.append('title', document.getElementById('title').value);
        formData.append('category', document.getElementById('category').value);
        formData.append('video', document.getElementById('videoFile').files[0]);
        try {
            await api.uploadVideo(formData, state.token);
            ui.showNotification('Video subido exitosamente.');
            e.target.reset();
            loadVideos(); // Recargar la lista de videos
            window.location.hash = '#catalog';
        } catch (error) { /* ... */ }
    });

    // Añadir listener para el botón de logout, que se crea dinámicamente
    document.getElementById('auth-section').addEventListener('click', (e) => {
        if (e.target.id === 'logout-link') {
            e.preventDefault();
            handleLogout();
        }
    });
}

async function loadVideos() {
    try {
        const videos = await api.getVideos();
        ui.renderVideos(videos, 'http://localhost:8080');
    } catch (error) {
        document.getElementById('videoList').innerHTML = `<p class="text-red-500 col-span-full text-center">No se pudieron cargar los videos.</p>`;
    }
}

async function handleVerification(token) {
    ui.showSection('verify-section');
    const title = document.getElementById('verify-title');
    const message = document.getElementById('verify-message');
    const icon = document.getElementById('verify-icon');
    try {
        const data = await api.verifyUser(token);
        title.textContent = "¡Cuenta Verificada!";
        message.textContent = data.message;
        icon.className = "fas fa-check-circle text-5xl text-green-500";
    } catch (error) {
        title.textContent = "Error de Verificación";
        message.textContent = error.message;
        icon.className = "fas fa-times-circle text-5xl text-red-500";
    } finally {
        history.pushState("", document.title, window.location.pathname + window.location.search);
    }
}

function router() {
    const urlParams = new URLSearchParams(window.location.hash.split('?')[1]);
    const token = urlParams.get('token');

    if (window.location.hash.startsWith('#verify') && token) {
        handleVerification(token);
    } else {
        const hash = window.location.hash || '#catalog';
        ui.showSection(hash.substring(1) + '-section');
    }
}

// Punto de entrada de la aplicación
function main() {
    window.addEventListener('hashchange', router);
    router();
    addFormListeners();
    ui.updateAuthUI(state);
    loadVideos();
}

// Ejecutar la aplicación cuando el DOM esté listo
document.addEventListener('DOMContentLoaded', main);