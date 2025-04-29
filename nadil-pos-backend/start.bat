@echo off
chcp 65001 >nul
set "prog1=D:\projects\nadil-pos-system\nadil-pos-backend\nadil-pos-backend.exe"
set "prog2=D:\projects\nadil-pos-system\nadil-pos-backend\nadil-pos-printing-hub\nadil-pos-printing-hub.exe"
set "prog3=D:\projects\nadil-pos-system\nadil-pos-backend\nadil-pos-windows-gui\Nadil-POS-System.exe"


tasklist /FI "IMAGENAME eq nadil-pos-backend.exe" | find /I "nadil-pos-backend.exe" >nul 2>&1
if %errorlevel%==0 (
    echo nadil-pos-backend.exe est déjà en cours d'exécution. Passer...
) else (
    echo Lancement de nadil-pos-backend...
    start "" "%prog1%" 
   
    timeout /t 2 >nul
)

tasklist /FI "IMAGENAME eq nadil-pos-printing-hub.exe" | find /I "nadil-pos-printing-hub" >nul 2>&1

if %errorlevel%==0 (
    echo nadil-pos-printing-hub.exe est déjà en cours d'exécution. Passer...
) else (
    echo Lancement de nadil-pos-printing-hub...
    start "" "%prog2%"
    timeout /t 2 >nul
)


tasklist /FI "IMAGENAME eq Nadil-POS-System.exe" | find /I "Nadil-POS-System" >nul 2>&1
if %errorlevel%==0 (
    echo Nadil-POS-System.exe est déjà en cours d'exécution. Passer...
) else (
    echo Lancement Nadil-POS-System.exe...
    start "" "%prog3%"
)
echo "Cette fenêtre se refermera automatiquement dans 5 secondes
timeout /t 6 >nul
