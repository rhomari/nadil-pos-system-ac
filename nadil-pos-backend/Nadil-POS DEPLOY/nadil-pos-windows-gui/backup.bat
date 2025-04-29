@echo off
set TIMESTAMP=%date:~-4,4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%
set BACKUP_DIR=c:\backups
set DATABASE=nadilposdb
set USER=postgres
set PGPASSWORD=230404

mkdir "%BACKUP_DIR%" 2>nul

"C:\Program Files\PostgreSQL\16\bin\pg_dump.exe" -U %USER% -h localhost -d %DATABASE% -F c -f "%BACKUP_DIR%\%DATABASE%-%TIMESTAMP%.backup"


echo Backup completed at %TIMESTAMP%