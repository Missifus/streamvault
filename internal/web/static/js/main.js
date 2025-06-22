// assets/js/main.js

// Configuración
const API_BASE_URL = 'http://localhost:8080/api';

// Manejar autenticación
function handleAuth() {
    // Guardar token
    function saveToken(token) {
        localStorage.setItem('authToken', token);
    }
    
    // Obtener token
    function getToken() {
        return localStorage.getItem('authToken');
    }
    
    // Cerrar sesión
    function logout() {
        localStorage.removeItem('authToken');
        window.location.href = '/index.html';
    }
    
    // Verificar si usuario está autenticado
    function isAuthenticated() {
        return getToken() !== null;
    }
    
    // Validar acceso a páginas protegidas
    function protectPage() {
        if (!isAuthenticated() && !window.location.pathname.includes('index.html')) {
            window.location.href = '/index.html';
        }
    }
    
    // Login
    if (document.getElementById('loginForm')) {
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            
            try {
                const response = await fetch(`${API_BASE_URL}/login`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });
                
                if (!response.ok) throw new Error('Credenciales incorrectas');
                
                const { token } = await response.json();
                saveToken(token);
                window.location.href = '/videos.html';
            } catch (error) {
                showMessage(error.message, 'error');
            }
        });
    }
    
    // Registro
    if (document.getElementById('registerForm')) {
        document.getElementById('registerForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('regEmail').value;
            const password = document.getElementById('regPassword').value;
            const confirmPassword = document.getElementById('confirmPassword').value;
            
            if (password !== confirmPassword) {
                showMessage('Las contraseñas no coinciden', 'error');
                return;
            }
            
            try {
                const response = await fetch(`${API_BASE_URL}/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });
                
                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || 'Error en registro');
                }
                
                showMessage('Registro exitoso. Ahora puedes iniciar sesión.', 'success');
                document.getElementById('showLogin').click();
            } catch (error) {
                showMessage(error.message, 'error');
            }
        });
    }
    
    // Logout
    if (document.getElementById('logoutBtn')) {
        document.getElementById('logoutBtn').addEventListener('click', logout);
    }
    
    // Navegación entre login/registro
    if (document.getElementById('showRegister')) {
        document.getElementById('showRegister').addEventListener('click', (e) => {
            e.preventDefault();
            document.getElementById('loginSection').style.display = 'none';
            document.getElementById('registerSection').style.display = 'block';
        });
    }
    
    if (document.getElementById('showLogin')) {
        document.getElementById('showLogin').addEventListener('click', (e) => {
            e.preventDefault();
            document.getElementById('registerSection').style.display = 'none';
            document.getElementById('loginSection').style.display = 'block';
        });
    }
    
    // Volver atrás
    if (document.getElementById('backBtn')) {
        document.getElementById('backBtn').addEventListener('click', () => {
            window.history.back();
        });
    }
    
    // Proteger páginas
    protectPage();
}

// Manejar videos
function handleVideos() {
    // Cargar lista de videos
    async function loadVideoList() {
        try {
            const token = localStorage.getItem('authToken');
            const response = await fetch(`${API_BASE_URL}/videos`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            
            if (!response.ok) throw new Error('Error al cargar videos');
            
            const videos = await response.json();
            renderVideoList(videos);
        } catch (error) {
            showMessage(error.message, 'error');
        }
    }
    
    // Mostrar lista de videos
    function renderVideoList(videos) {
        const container = document.getElementById('videoList');
        if (!container) return;
        
        container.innerHTML = '';
        
        videos.forEach(video => {
            const videoElement = document.createElement('div');
            videoElement.className = 'video-item';
            videoElement.innerHTML = `
                <div class="video-thumbnail">Miniatura</div>
                <div class="video-info">
                    <h3>${video.title}</h3>
                    <p>${video.duration || '0:00'} | ${video.size || '0MB'}</p>
                </div>
            `;
            
            videoElement.addEventListener('click', () => {
                window.location.href = `/player.html?id=${video.id}`;
            });
            
            container.appendChild(videoElement);
        });
    }
    
    // Cargar y reproducir video
    async function loadVideoPlayer() {
        const urlParams = new URLSearchParams(window.location.search);
        const videoId = urlParams.get('id');
        
        if (!videoId) {
            showMessage('Video no especificado', 'error');
            return;
        }
        
        try {
            const token = localStorage.getItem('authToken');
            const response = await fetch(`${API_BASE_URL}/videos/${videoId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            
            if (!response.ok) throw new Error('Video no encontrado');
            
            const video = await response.json();
            document.getElementById('videoTitle').textContent = video.title;
            
            const videoPlayer = document.getElementById('videoPlayer');
            videoPlayer.src = `${API_BASE_URL}/videos/stream/${videoId}`;
            videoPlayer.load();
        } catch (error) {
            showMessage(error.message, 'error');
        }
    }
    
    // Inicializar según la página
    if (document.getElementById('videoList')) {
        loadVideoList();
    }
    
    if (document.getElementById('videoPlayer')) {
        loadVideoPlayer();
    }
}

// Mostrar mensajes
function showMessage(message, type) {
    const messageDiv = document.getElementById('message');
    if (messageDiv) {
        messageDiv.textContent = message;
        messageDiv.className = `message-${type}`;
        
        setTimeout(() => {
            messageDiv.textContent = '';
            messageDiv.className = '';
        }, 5000);
    } else {
        alert(message);
    }
}

// Inicialización
document.addEventListener('DOMContentLoaded', () => {
    handleAuth();
    handleVideos();
});