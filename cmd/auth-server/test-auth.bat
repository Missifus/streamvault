@echo off
setlocal enabledelayedexpansion

echo.
echo === REGISTRANDO USUARIO ===
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d "{\"email\":\"test_batch@test.com\",\"password\":\"BatchPass123!\"}"

echo.
echo === VERIFICANDO EMAIL ===
curl -X GET "http://localhost:8080/verify?email=test_batch@test.com"

echo.
echo === INICIANDO SESION ===
for /f "tokens=*" %%a in ('curl -s -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"email\":\"test_batch@test.com\",\"password\":\"BatchPass123!\"}"') do (
    set "response=%%a"
    echo Respuesta: !response!
    
    for /f "delims={},: tokens=2" %%b in ("!response!") do (
        set "token=%%b"
    )
)

:: Limpiar token de comillas extra
set "token=!token:"=!"
echo Token: !token!

echo.
echo === ACCEDIENDO A RUTA PROTEGIDA ===
curl -H "Authorization: Bearer !token!" http://localhost:8080/protected

echo.
echo === TEST COMPLETADO ===
pause