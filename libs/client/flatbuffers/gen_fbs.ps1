# Stop the script if any command fails. This is a good practice for build scripts.
$ErrorActionPreference = "Stop"

# --- Color Theme ---
$ColorInfo = "Cyan"
$ColorSuccess = "Green"
$ColorError = "Red"
$ColorWarning = "Yellow"
$ColorHighlight = "White"

# --- 1. Verify that the FlatBuffers compiler is installed and in the PATH ---
Write-Host "[~] Checking for FlatBuffers compiler (flatc)..." -ForegroundColor $ColorInfo
$flatcPath = Get-Command flatc -ErrorAction SilentlyContinue
if (-not $flatcPath) {
    Write-Host "[!] Error: flatc not found in your system's PATH." -ForegroundColor $ColorError
    Write-Host "    Please install the FlatBuffers compiler and ensure it's in your PATH."
    exit 1
}
Write-Host "[+] flatc found successfully." -ForegroundColor $ColorSuccess

# --- 2. Define script paths ---
$scriptDir = $PSScriptRoot
$outputDir = Join-Path -Path $scriptDir -ChildPath "generated"
Write-Host "[~] Output directory set to: $outputDir" -ForegroundColor $ColorInfo

# --- 3. Clean up the existing generated folder ---
if (Test-Path $outputDir) {
    Write-Host "[~] Removing existing 'generated' directory..." -ForegroundColor $ColorWarning
    try {
        Remove-Item -Path $outputDir -Recurse -Force
        Write-Host "[+] 'generated' directory removed." -ForegroundColor $ColorSuccess
    } catch {
        Write-Host "[!] Error: Could not remove '$outputDir'. Please check file permissions." -ForegroundColor $ColorError
        Write-Host "    Error details: $_"
        exit 1
    }
}

# --- 4. Find all .fbs schema files ---
Write-Host "[~] Finding FlatBuffers schema (.fbs) files..." -ForegroundColor $ColorInfo
$fbsFiles = (Get-ChildItem -Path $scriptDir -Filter *.fbs -Recurse).FullName
if ($fbsFiles.Count -eq 0) {
    Write-Host "[!] Warning: No .fbs files were found. Nothing to generate." -ForegroundColor $ColorWarning
    exit 0
}
Write-Host "[+] Found $($fbsFiles.Count) schema file(s) to compile." -ForegroundColor $ColorInfo

# --- 5. Compile the schemas using flatc ---
Write-Host "[~] Compiling schemas and generating Go files..." -ForegroundColor $ColorInfo

# Find all subdirectories to use as include paths for flatc
$includePaths = (Get-ChildItem -Path $scriptDir -Directory -Recurse | ForEach-Object { "-I", $_.FullName }) + "-I", $scriptDir

# Using a try/catch block for robust error handling during compilation
try {
    # Use splatting to pass arguments cleanly to the external command
    $flatcArgs = @(
        "--go",
        "-o", $outputDir,
        $includePaths,
        $fbsFiles
    )
    
    # Execute the compiler
    & flatc @flatcArgs

    # Check the exit code of the last native command. A non-zero code means an error occurred.
    if ($LASTEXITCODE -ne 0) {
        # Manually throw a terminating error so the 'catch' block will execute.
        throw "flatc compiler returned a non-zero exit code: $LASTEXITCODE"
    }

    Write-Host ""
    Write-Host "==================================================" -ForegroundColor $ColorSuccess
    Write-Host "  SUCCESS: Go files generated in: $outputDir" -ForegroundColor $ColorHighlight
    Write-Host "==================================================" -ForegroundColor $ColorSuccess
    Write-Host ""

} catch {
    Write-Host ""
    Write-Host "==================================================" -ForegroundColor $ColorError
    Write-Host "  ERROR: flatc compiler failed." -ForegroundColor $ColorHighlight
    Write-Host "  Details: $($_.Exception.Message)" -ForegroundColor $ColorError
    Write-Host "==================================================" -ForegroundColor $ColorError
    Write-Host ""
    exit 1
}

exit 0
