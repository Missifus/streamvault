// Configuración básica del reproductor
document.addEventListener('DOMContentLoaded', () => {
    const video = document.getElementById('main-video');
    
    // Intentar cargar el estado previo del reproductor
    const savedTime = localStorage.getItem(`video-time-${video.src}`);
    if (savedTime) {
        video.currentTime = parseFloat(savedTime);
    }
    
    // Guardar el tiempo de reproducción periódicamente
    video.addEventListener('timeupdate', () => {
        localStorage.setItem(`video-time-${video.src}`, video.currentTime);
    });
    
    // Opciones de calidad simuladas (en producción se usaría HLS con múltiples calidades)
    const qualityButton = document.createElement('button');
    qualityButton.textContent = 'Calidad';
    qualityButton.classList.add('quality-btn');
    video.parentNode.insertBefore(qualityButton, video.nextSibling);
    
    qualityButton.addEventListener('click', () => {
        alert('Seleccione calidad:\n- Alta (1080p)\n- Media (720p)\n- Baja (480p)');
    });
});