@echo off
setlocal enabledelayedexpansion
title Probando Sistema de Autenticación - StreamVault
color 0B

echo.
echo =============================================
echo    PRUEBA DEL SISTEMA DE AUTENTICACIÓN
echo =============================================
echo.

:: ------------------------------------------
:: Configuración de datos de prueba
:: ------------------------------------------
set "test_email=test_user_%random%@streamvault.com"
set "test_password=SecurePass123!"
echo [CONFIG] Usando datos de prueba:
echo   Email:    !test_email!
echo   Password: !test_password!
echo.

:: ------------------------------------------
:: 1. Registrar nuevo usuario
:: ------------------------------------------
echo [PASO 1] Registrando nuevo usuario...
curl -s -X POST http://localhost:8080/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"!test_email!\",\"password\":\"!test_password!\"}"

echo.
echo.
echo Presiona Enter para continuar...
pause >nul

:: ------------------------------------------
:: 2. Verificar email
:: ------------------------------------------
echo.
echo [PASO 2] Verificando email del usuario...
curl -s -X GET "http://localhost:8080/verify?email=!test_email!"

echo.
echo.
echo Presiona Enter para continuar...
pause >nul

:: ------------------------------------------
:: 3. Iniciar sesión y obtener token
:: ------------------------------------------
echo.
echo [PASO 3] Iniciando sesión y obteniendo token JWT...

for /f "tokens=*" %%a in ('curl -s -X POST http://localhost:8080/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"!test_email!\",\"password\":\"!test_password!\"}"') do (
    set "login_response=%%a"
)

:: Procesar respuesta JSON para extraer token
set "token=!login_response:{\"token\":\"=!"
set "token=!token:\"}=!"
echo Token obtenido: !token!

echo.
echo Presiona Enter para continuar...
pause >nul


:: Extraer user_id de la respuesta
set "user_id=!protected_response:*user_id:=!"
set "user_id=!user_id:}=!"
set "user_id=!user_id: =!"
echo User ID: !user_id!

echo.
echo Presiona Enter para continuar...
pause >nul

echo.
echo.
echo =============================================
echo    RESUMEN DEL USUARIO CREADO
echo =============================================
echo    EMAIL:       !test_email!
echo    PASSWORD:    !test_password!
echo    USER ID:     !user_id!
echo    TOKEN:       !token!
echo =============================================
echo.
echo Prueba completada! Presiona Enter para salir...
pause >nul