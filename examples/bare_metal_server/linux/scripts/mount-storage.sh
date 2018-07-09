#!/bin/sh

fdisk /dev/xvdb <<EOF
n
p
1


w
EOF

mkfs.xfs /dev/xvdb
mkdir /mnt/a
ls -l /mnt/a
mount /dev/xvdb /mnt/a