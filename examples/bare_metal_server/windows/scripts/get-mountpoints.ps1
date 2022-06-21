# Name getmntptdtls.ps1
# This will output specific details for mountpoints

$TotalGB = @{Name="Capacity(GB)";expression={[math]::round(($_.Capacity/ 1073741824),2)}}
$FreeGB = @{Name="FreeSpace(GB)";expression={[math]::round(($_.FreeSpace / 1073741824),2)}}
$FreePerc = @{Name="Free(%)";expression={[math]::round(((($_.FreeSpace / 1073741824)/($_.Capacity / 1073741824)) * 100),0)}}

function get-mountpoints {
    $volumes = Get-WmiObject win32_volume -Filter "DriveType='3'"
    $volumes | Select Name, Label, DriveLetter, FileSystem, $TotalGB, $FreeGB, $FreePerc | Format-Table -AutoSize
}

get-mountpoints