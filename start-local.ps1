# Run from Windows PowerShell 5.1 or by double-clicking the companion CMD file.
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$projectRoot = $PSScriptRoot
$backendRoot = Join-Path $projectRoot 'ZZHAN'
$frontendRoot = Join-Path $projectRoot 'ZZHAN_web/user'
$backupRoot = Join-Path $projectRoot '.local-backups'
$logRoot = Join-Path $backupRoot 'local-dev'
$httpClient = $null

function Test-HttpReady {
    param([string]$Url)
    try {
        $response = $httpClient.GetAsync($Url).GetAwaiter().GetResult()
        try { return $response.IsSuccessStatusCode }
        finally { $response.Dispose() }
    }
    catch { return $false }
}

function Test-LocalPort {
    param([int]$Port)
    $probe = New-Object System.Net.Sockets.TcpClient
    try {
        $connection = $probe.ConnectAsync('127.0.0.1', $Port)
        if (-not $connection.Wait(1000)) { return $false }
        return $probe.Connected
    }
    catch { return $false }
    finally { $probe.Dispose() }
}

function Wait-ServiceReady {
    param([string]$Name, [string]$Url, [System.Diagnostics.Process]$Process)
    $deadline = [DateTime]::UtcNow.AddSeconds(30)
    do {
        if (Test-HttpReady $Url) { return }
        $Process.Refresh()
        if ($Process.HasExited) {
            throw "$Name exited with code $($Process.ExitCode). Check logs in: $logRoot"
        }
        Start-Sleep -Milliseconds 500
    } while ([DateTime]::UtcNow -lt $deadline)
    throw "$Name did not become ready within 30 seconds. Check logs in: $logRoot"
}

try {
    New-Item -ItemType Directory -Path $logRoot -Force | Out-Null
    Add-Type -AssemblyName System.Net.Http
    $handler = New-Object System.Net.Http.HttpClientHandler
    $handler.UseProxy = $false
    $httpClient = New-Object System.Net.Http.HttpClient($handler)
    $httpClient.Timeout = [TimeSpan]::FromSeconds(2)

    $backendUrl = 'http://127.0.0.1:9900/api/v1/health'
    if (Test-HttpReady $backendUrl) {
        Write-Host 'Backend is already running; reusing port 9900.'
    }
    else {
        if (Test-LocalPort 9900) {
            throw 'Port 9900 is occupied, but the backend health check failed. Check the existing service first.'
        }
        $configPath = Join-Path $backendRoot 'config.yaml'
        if (-not (Test-Path -LiteralPath $configPath -PathType Leaf)) {
            throw "Backend configuration is missing: $configPath"
        }
        $goExecutable = (Get-Command go.exe -ErrorAction Stop).Source
        $backendExecutable = Join-Path $backupRoot 'blog-server.exe'
        Write-Host 'Building the current backend...'
        $buildProcess = Start-Process -FilePath $goExecutable -ArgumentList @('build', '-o', ('"{0}"' -f $backendExecutable), './cmd/server') -WorkingDirectory $backendRoot -WindowStyle Hidden -Wait -PassThru -RedirectStandardOutput (Join-Path $logRoot 'backend-build.out.log') -RedirectStandardError (Join-Path $logRoot 'backend-build.err.log')
        if ($buildProcess.ExitCode -ne 0) {
            throw "Backend build failed. Check logs in: $logRoot"
        }
        Write-Host 'Starting backend (MySQL must already be running)...'
        $backendProcess = Start-Process -FilePath $backendExecutable -ArgumentList @('-c', ('"{0}"' -f $configPath)) -WorkingDirectory $backendRoot -WindowStyle Hidden -PassThru -RedirectStandardOutput (Join-Path $logRoot 'backend.out.log') -RedirectStandardError (Join-Path $logRoot 'backend.err.log')
        Wait-ServiceReady -Name 'Backend' -Url $backendUrl -Process $backendProcess
    }

    $frontendUrl = 'http://127.0.0.1:3000/'
    if (Test-HttpReady $frontendUrl) {
        Write-Host 'Frontend is already running; reusing port 3000.'
    }
    else {
        if (Test-LocalPort 3000) {
            throw 'Port 3000 is occupied, but the frontend did not respond. Check the existing service first.'
        }
        $nodeExecutable = (Get-Command node.exe -ErrorAction Stop).Source
        $viteScript = Join-Path $frontendRoot 'node_modules/vite/bin/vite.js'
        if (-not (Test-Path -LiteralPath $viteScript -PathType Leaf)) {
            throw "Frontend dependencies are missing. Run npm install in: $frontendRoot"
        }
        Write-Host 'Starting frontend...'
        $frontendProcess = Start-Process -FilePath $nodeExecutable -ArgumentList @(('"{0}"' -f $viteScript), '--host', '127.0.0.1', '--port', '3000', '--strictPort') -WorkingDirectory $frontendRoot -WindowStyle Hidden -PassThru -RedirectStandardOutput (Join-Path $logRoot 'frontend.out.log') -RedirectStandardError (Join-Path $logRoot 'frontend.err.log')
        Wait-ServiceReady -Name 'Frontend' -Url $frontendUrl -Process $frontendProcess
    }

    Write-Host ''
    Write-Host 'Blog is ready: http://127.0.0.1:3000/'
    Write-Host "Logs: $logRoot"
    Write-Host 'The services will keep running after this window closes.'
}
catch {
    Write-Host ''
    Write-Host "Startup failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Logs: $logRoot"
    exit 1
}
finally {
    if ($null -ne $httpClient) { $httpClient.Dispose() }
}
