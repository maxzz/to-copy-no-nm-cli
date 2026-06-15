param(
    [string]$Root = (Split-Path -Parent $PSScriptRoot)
)

Add-Type -AssemblyName System.Drawing

$pngPath = Join-Path $Root "assets\icon-source.png"
$icoPath = Join-Path $Root "assets\icon.ico"

if (-not (Test-Path $pngPath)) {
    throw "Missing icon source: $pngPath"
}

function New-ResizedBitmap([System.Drawing.Image]$source, [int]$size) {
    $bitmap = New-Object System.Drawing.Bitmap $size, $size
    $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
    $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
    $graphics.DrawImage($source, 0, 0, $size, $size)
    $graphics.Dispose()
    return $bitmap
}

$source = [System.Drawing.Image]::FromFile($pngPath)
try {
    $sizes = @(16, 32, 48, 256)
    $bitmaps = New-Object System.Collections.Generic.List[System.Drawing.Bitmap]
    foreach ($size in $sizes) {
        $bitmaps.Add((New-ResizedBitmap $source $size))
    }

    $iconHandle = $bitmaps[$bitmaps.Count - 1].GetHicon()
    $icon = [System.Drawing.Icon]::FromHandle($iconHandle)
    try {
        $stream = [System.IO.File]::Open($icoPath, [System.IO.FileMode]::Create)
        try {
            $icon.Save($stream)
        }
        finally {
            $stream.Close()
        }
    }
    finally {
        $icon.Dispose()
    }

    foreach ($bitmap in $bitmaps) {
        $bitmap.Dispose()
    }
}
finally {
    $source.Dispose()
}

Write-Host "Generated $icoPath"
