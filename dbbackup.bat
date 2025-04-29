@echo off
setlocal enabledelayedexpansion


set "BACKUP_DIR=%~dp0"


set "PG_BIN=C:\Program Files\PostgreSQL\16\bin"  
set "PG_USER=postgres"
set "PG_DB=nadilposdb"
set "PG_HOST=localhost"
set "PG_PORT=5432"


for /f "tokens=2 delims==" %%I in ('wmic OS Get localdatetime /value') do set datetime=%%I
set "TIMESTAMP=!datetime:~0,4!-!datetime:~4,2!-!datetime:~6,2!_!datetime:~8,2!-!datetime:~10,2!-!datetime:~12,2!"

set "BACKUP_FILE=%BACKUP_DIR%%PG_DB%_backup_%TIMESTAMP%.backup"

"%PG_BIN%\pg_dump" -U %PG_USER% -h %PG_HOST% -p %PG_PORT% -F c -b -v -f "%BACKUP_FILE%" %PG_DB%


if exist "%BACKUP_FILE%" (
    echo Backup completed successfully: %BACKUP_FILE%
) else (
    echo Backup failed!
)

pause
