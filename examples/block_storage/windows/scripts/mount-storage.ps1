# Get the disk that has a raw partition style. (except the system disk)
Get-Disk | Where-Object IsSystem -eq $False
$DiskNumber = (Get-Disk | Where-Object IsSystem -eq $False | select Number).Number

# Initialize the disk.
Initialize-Disk -Number $DiskNumber

Get-Disk $DiskNumber | Format-List

# Partition the disk.
New-Partition -DiskNumber $DiskNumber -UseMaximumSize

Add-PartitionAccessPath -DiskNumber $DiskNumber -PartitionNumber 2 -AccessPath "D:\"

# Format the volume.
Format-Volume -DriveLetter D -FileSystem NTFS -Confirm:$false
