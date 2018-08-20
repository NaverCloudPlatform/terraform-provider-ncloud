#!/bin/sh

fdisk /dev/xvdb <<EOF
n
p
1


w
EOF

mkfs.ext4 /dev/xvdb
mkdir -p /mnt/a
ls -l /mnt/a
mount /dev/xvdb /mnt/a