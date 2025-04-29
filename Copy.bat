@echo off
setlocal


set "source=D:\Projects\nadil-pos-system-ac\dist"
set "destination=D:\Projects\nadil-pos-system-ac\nadil-pos-backend\static"
set "jsFolder=%destination%\js"


if not exist "%destination%" (
    echo Creating destination folder: %destination%
    mkdir "%destination%"
)


if exist "%jsFolder%" (
    echo Deleting all files in "%jsFolder%"...
    del /Q "%jsFolder%\*.*"
) else (
    echo JS folder does not exist, skipping deletion.
)


echo Copying files from "%source%" to "%destination%"...
xcopy "%source%\*" "%destination%\" /E /H /C /I /Y


if %ERRORLEVEL% equ 0 (
    echo All files have been copied successfully!
) else (
    echo Error occurred during file copying.
)


pause
