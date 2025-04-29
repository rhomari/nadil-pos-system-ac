Set-ExecutionPolicy Bypass -Scope Process -Force
# Define the JSON URL
$JsonUrl = "http://localhost:2304/updates.json"
$JsonPath = "$env:TEMP\updates.json"

# Download the JSON file
Write-Host "Downloading update list..."
Invoke-WebRequest -Uri $JsonUrl -OutFile $JsonPath

# Read and parse JSON
$UpdateList = Get-Content -Path $JsonPath | ConvertFrom-Json

# Process each update file
foreach ($File in $UpdateList.files) {
    $Url = $File.url
    $Destination = $File.destination
    $ShutdownProcess = $File.shutdown
    $RunAfter = $File.run_after

    # Stop process if specified
    if ($ShutdownProcess) {
        Write-Host "Stopping process: $ShutdownProcess"
        Stop-Process -Name $ShutdownProcess -Force -ErrorAction SilentlyContinue
    }

    # Ensure destination folder exists
    $DestFolder = Split-Path -Path $Destination
    if (!(Test-Path $DestFolder)) {
        Write-Host "Creating folder: $DestFolder"
        New-Item -ItemType Directory -Path $DestFolder -Force
    }

    # Download file
    Write-Host "Downloading: $Url"
    Invoke-WebRequest -Uri $Url -OutFile $Destination

    # Restart process if needed
    if ($RunAfter) {
        Write-Host "Restarting: $RunAfter"
        Start-Process -FilePath $RunAfter
    }
}

Write-Host "Update process completed."
