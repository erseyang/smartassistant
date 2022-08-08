# 在X86主机上运行智汀家庭云

本教程，我们将描述如何在 X86 主机上通过 U 盘引导运行智汀家庭云，以及如何将智汀家庭云烧写到主机硬盘

**重要:** 您需要具备基本的**linux**操作经验，包括运行shell脚本，运行基本的linux命令，若对以上存疑请提前了解后再阅读本教程。

## 环境准备

* X86主机，支持efi引导，2GB以上内存，16GB以上硬盘，网口已联网
* 8GB以上空间的 U 盘

## 烧录U盘

访问 https://github.com/zhiting-tech/smartassistant/releases ，下载对应版本的 openwrt 镜像, 
如：openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img.gz ，解压后写入 U 盘

Linux系统下：

``` shell
gzip -d -k openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img.gz
sudo dd if=./openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img of=/dev/sdb bs=4MiB && sync
# 可选：创建分区，拷贝 openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img 文件到新分区，用于后续烧写到硬盘
```

Windows系统下：

先使用7z等工具解压 openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img.gz，
然后使用 [USB Image Tool](https://www.alexpage.de/usb-image-tool/) 等工具将 img 文件写入 U 盘。

## U盘启动智汀家庭云

1. 将烧录好的 U 盘插入，将网线连接任一网口，开机
2. 设置 BIOS 为上电启动，U盘为第一启动项，引导模式为efi，并保存
3. 重启后选择 **smartassistant**, 进入系统
4. 等待约 2 分钟后，执行 docker ps 指令，查看相关容器是否启动
5. 局域网中手机运行智汀APP，扫描并添加

## 写入硬盘

此方法会烧录一个全新的智汀家庭云到硬盘，U盘系统数据将不会迁移；开始前请确认 U 盘中已手动增加分区，并且拷贝好 openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img 文件

1. /etc/init.d/smartassistant stop
2. /etc/init.d/dockerd stop
3. umount /mnt/data
4. 以上操作可参考 “[停止智汀家庭云相关功能](#停止智汀家庭云相关功能)” 永久停用
5. mkdir /mnt/disk
6. mount /dev/sdb3 /mnt/disk
7. ~/init-efi.sh /mnt/disk/openwrt-21.02.3-x86-64-generic-ext4-combined-efi.img /dev/sda ，输入y确认
8. 安装成功后提示重启并将U盘拔出，至此openwrt系统安装完毕
9. 系统启动后等待约 2 分钟，局域网中手机运行智汀APP，扫描并添加

## 停止智汀家庭云相关功能

以下操作将永久停用智汀家庭云相关功能，适合 U 盘烧录场景，减少每次操作步骤

``` shell
/etc/init.d/smartassistant disable
/etc/init.d/make_shared disable
/etc/init.d/dockerd disable

vim /etc/config/fstab
# disable /mnt/data mount
# disable swap
```