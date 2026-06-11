# make_deploy.ps1
# Builds tt.exe and splits it into Base64 text parts for Zoom chat delivery.
# Run this script from the deploy\ folder.

$projectRoot = Split-Path $PSScriptRoot -Parent
$chunkSize = 100000   # characters per part (~27 parts total)

# --- Build ---
Write-Host "Building tt_release.exe..."
Push-Location $projectRoot
$env:CGO_ENABLED = "0"
go build -ldflags="-s -w" -o tt_release.exe .
if ($LASTEXITCODE -ne 0) { Write-Error "Build failed."; Pop-Location; exit 1 }
Pop-Location
Write-Host "Build OK."

# --- Base64 encode ---
Write-Host "Encoding to Base64..."
$bytes = [System.IO.File]::ReadAllBytes("$projectRoot\tt_release.exe")
$b64   = [System.Convert]::ToBase64String($bytes)
$totalParts = [math]::Ceiling($b64.Length / $chunkSize)

# --- Split ---
Write-Host "Splitting into $totalParts parts..."
for ($i = 0; $i -lt $totalParts; $i++) {
    $start   = $i * $chunkSize
    $len     = [math]::Min($chunkSize, $b64.Length - $start)
    $chunk   = $b64.Substring($start, $len)
    $partNum = "{0:D3}" -f ($i + 1)
    [System.IO.File]::WriteAllText(
        "$PSScriptRoot\part$partNum.txt",
        $chunk,
        [System.Text.Encoding]::ASCII
    )
}

Write-Host ""
Write-Host "========================================="
Write-Host " Done!  $totalParts part files created in deploy\"
Write-Host "========================================="
Write-Host ""
Write-Host "[STEP 1]  Send receive.ps1 first via Zoom chat."
Write-Host "          Customer: save as receive.ps1, then run it."
Write-Host ""
Write-Host "[STEP 2]  Send each part file in order:"
for ($i = 1; $i -le $totalParts; $i++) {
    $n = "{0:D3}" -f $i
    Write-Host ("  [{0}/{1}]  paste part{2}.txt" -f $i, $totalParts, $n)
}
Write-Host ""
Write-Host "[STEP 3]  Send the word:  END"
