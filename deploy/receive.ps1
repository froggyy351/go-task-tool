# receive.ps1
# Reassembles tt.exe from Base64 parts delivered via Zoom chat.
#
# How to use:
#   1. Save this file as receive.ps1
#   2. Open PowerShell in the folder where you want tt.exe saved
#   3. Run:  .\receive.ps1
#   4. For each part: copy the Zoom chat message, then press Enter
#   5. When you see "END" in Zoom chat: copy it, press Enter -> auto decode

Write-Host "================================================"
Write-Host " tt.exe receiver"
Write-Host "================================================"
Write-Host ""
Write-Host "For each Zoom message: copy it, then press Enter."
Write-Host "When you receive 'END': copy it, press Enter."
Write-Host ""

$parts   = [System.Collections.Generic.List[string]]::new()
$partNum = 1

while ($true) {
    Write-Host -NoNewline "  [Part $partNum] Copy from Zoom, then press Enter > "
    $null = Read-Host

    $clip = Get-Clipboard -Raw
    if ($null -eq $clip) {
        Write-Host "  Clipboard is empty. Copy the message first, then press Enter."
        continue
    }
    $clip = $clip.Trim()

    if ($clip -eq "END") {
        Write-Host "  Received END. Assembling..."
        break
    }

    $parts.Add($clip)
    Write-Host ("  Part {0} received ({1} chars)" -f $partNum, $clip.Length)
    $partNum++
}

if ($parts.Count -eq 0) {
    Write-Error "No parts received."
    exit 1
}

Write-Host ""
Write-Host "Decoding ($($parts.Count) parts)..."

try {
    $b64   = [string]::Concat($parts)
    $bytes = [System.Convert]::FromBase64String($b64)
    $out   = Join-Path $PWD "tt.exe"
    [System.IO.File]::WriteAllBytes($out, $bytes)

    Write-Host ""
    Write-Host "================================================"
    Write-Host " Saved: $out"
    Write-Host "================================================"
    Write-Host ""
    & $out version
    Write-Host ""
    Write-Host "tt.exe is ready to use."
} catch {
    Write-Error "Decode failed: $_"
    Write-Host "Check that all parts were received in the correct order."
    exit 1
}
