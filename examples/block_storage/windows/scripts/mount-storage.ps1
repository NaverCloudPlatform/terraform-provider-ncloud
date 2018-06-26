# List all disks except the system disk
Get-Disk | Where-Object IsSystem -eq $False
$DiskNumber = (Get-Disk | Where-Object IsSystem -eq $False | select Number).Number

Initialize-Disk -Number $DiskNumber

Get-Disk $DiskNumber | Format-List

New-Partition -DiskNumber $DiskNumber –UseMaximumSize
Add-PartitionAccessPath -DiskNumber $DiskNumber -PartitionNumber 2 –AccessPath "D:\"
Format-Volume -DriveLetter D -FileSystem NTFS -Confirm:$false