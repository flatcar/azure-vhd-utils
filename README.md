<div style="text-align: center">

[![Flatcar OS](https://img.shields.io/badge/Flatcar-Website-blue?logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhLS0gR2VuZXJhdG9yOiBBZG9iZSBJbGx1c3RyYXRvciAyNi4wLjMsIFNWRyBFeHBvcnQgUGx1Zy1JbiAuIFNWRyBWZXJzaW9uOiA2LjAwIEJ1aWxkIDApICAtLT4NCjxzdmcgdmVyc2lvbj0iMS4wIiBpZD0ia2F0bWFuXzEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4Ig0KCSB2aWV3Qm94PSIwIDAgODAwIDYwMCIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgODAwIDYwMDsiIHhtbDpzcGFjZT0icHJlc2VydmUiPg0KPHN0eWxlIHR5cGU9InRleHQvY3NzIj4NCgkuc3Qwe2ZpbGw6IzA5QkFDODt9DQo8L3N0eWxlPg0KPHBhdGggY2xhc3M9InN0MCIgZD0iTTQ0MCwxODIuOGgtMTUuOXYxNS45SDQ0MFYxODIuOHoiLz4NCjxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik00MDAuNSwzMTcuOWgtMzEuOXYxNS45aDMxLjlWMzE3Ljl6Ii8+DQo8cGF0aCBjbGFzcz0ic3QwIiBkPSJNNTQzLjgsMzE3LjlINTEydjE1LjloMzEuOVYzMTcuOXoiLz4NCjxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik02NTUuMiw0MjAuOXYtOTUuNGgtMTUuOXY5NS40aC0xNS45VjI2MmgtMzEuOVYxMzQuOEgyMDkuNFYyNjJoLTMxLjl2MTU5aC0xNS45di05NS40aC0xNnY5NS40aC0xNS45djMxLjINCgloMzEuOXYxNS44aDQ3Ljh2LTE1LjhoMTUuOXYxNS44SDI3M3YtMTUuOGgyNTQuOHYxNS44aDQ3Ljh2LTE1LjhoMTUuOXYxNS44aDQ3Ljh2LTE1LjhoMzEuOXYtMzEuMkg2NTUuMnogTTQ4Ny44LDE1MWg3OS42djMxLjgNCgloLTIzLjZ2NjMuNkg1MTJ2LTYzLjZoLTI0LjJMNDg3LjgsMTUxTDQ4Ny44LDE1MXogTTIzMywyMTQuNlYxNTFoNjMuN3YyMy41aC0zMS45djE1LjhoMzEuOXYyNC4yaC0zMS45djMxLjhIMjMzVjIxNC42eiBNMzA1LDMxNy45DQoJdjE1LjhoLTQ3Ljh2MzEuOEgzMDV2NDcuN2gtOTUuNVYyODYuMUgzMDVMMzA1LDMxNy45eiBNMzEyLjYsMjQ2LjRWMTUxaDMxLjl2NjMuNmgzMS45djMxLjhMMzEyLjYsMjQ2LjRMMzEyLjYsMjQ2LjRMMzEyLjYsMjQ2LjR6DQoJIE00NDguMywzMTcuOXY5NS40aC00Ny44di00Ny43aC0zMS45djQ3LjdoLTQ3LjhWMzAyaDE1Ljl2LTE1LjhoOTUuNVYzMDJoMTUuOUw0NDguMywzMTcuOXogTTQ0MCwyNDYuNHYtMzEuOGgtMTUuOXYzMS44aC0zMS45DQoJdi03OS41aDE1Ljl2LTE1LjhoNDcuOHYxNS44aDE1Ljl2NzkuNUg0NDB6IE01OTEuNiwzMTcuOXY0Ny43aC0xNS45djE1LjhoMTUuOXYzMS44aC00Ny44di0zMS43SDUyOHYtMTUuOGgtMTUuOXY0Ny43aC00Ny44VjI4Ni4xDQoJaDEyNy4zVjMxNy45eiIvPg0KPC9zdmc+DQo=)](https://www.flatcar.org/)
[![Matrix](https://img.shields.io/badge/Matrix-Chat%20with%20us!-green?logo=matrix)](https://app.element.io/#/room/#flatcar:matrix.org)
[![Slack](https://img.shields.io/badge/Slack-Chat%20with%20us!-4A154B?logo=slack)](https://kubernetes.slack.com/archives/C03GQ8B5XNJ)
[![Twitter Follow](https://img.shields.io/twitter/follow/flatcar?style=social)](https://x.com/flatcar)
[![Mastodon Follow](https://img.shields.io/badge/Mastodon-Follow-6364FF?logo=mastodon)](https://hachyderm.io/@flatcar)
[![Bluesky](https://img.shields.io/badge/Bluesky-Follow-0285FF?logo=bluesky)](https://bsky.app/profile/flatcar.org)

</div>


# Azure VHD utilities.

This project provides a Go package to read Virtual Hard Disk (VHD) file, a CLI interface to upload local VHD to Azure storage and to inspect a local VHD.

An implementation of VHD [VHD specification](https://technet.microsoft.com/en-us/virtualization/bb676673.aspx) can be found in the [vhdcore](/vhdcore) package. 


[![Go Report Card](https://goreportcard.com/badge/github.com/flatcar/azure-vhd-utils)](https://goreportcard.com/report/github.com/flatcar/azure-vhd-utils)

# Installation
> Note: You must have Go installed on your machine, at version 1.11 or greater. [https://golang.org/dl/](https://golang.org/dl/) 

    go get github.com/flatcar/azure-vhd-utils

# Features

1. Fast uploads - This tool offers faster uploads by using multiple routines and balancing the load across them.
2. Efficient uploads - This tool will only upload used (non-zero) portions of the disk.
3. Parallelism - This tool can upload segements of the VHD concurrently (user configurable).

# Usage

### Upload local VHD to Azure storage as page blob

```bash
USAGE:
   azure-vhd-utils upload [command options] [arguments...]

OPTIONS:
   --localvhdpath       Path to source VHD in the local machine.
   --stgaccountname     Azure storage account name.
   --stgaccountkey      Azure storage account key.
   --containername      Name of the container holding destination page blob. (Default: vhds)
   --blobname           Name of the destination page blob.
   --parallelism        Number of concurrent goroutines to be used for upload
```

The upload command uploads local VHD to Azure storage as page blob. Once uploaded, you can use Microsoft Azure portal to register an image based on this page blob and use it to create Azure Virtual Machines.

#### Note
When creating a VHD for Microsoft Azure, the size of the VHD must be a whole number in megabytes, otherwise you will see an error similar to the following when you attempt to create image from the uploaded VHD in Azure:

"The VHD http://<mystorageaccount>.blob.core.windows.net/vhds/<vhd-pageblob-name>.vhd has an unsupported virtual size of <number> bytes. The size must be a whole number (in MBs)."

You should ensure the VHD size is even MB before uploading

##### For virtual box:
-------------------
VBoxManage modifyhd <absolute path to file> --resize &lt;size in MB&gt;

##### For Hyper V:
----------------
Resize-VHD -Path <absolute path to file> -SizeBytes 

     http://azure.microsoft.com/blog/2014/05/22/running-freebsd-in-azure/

##### For Qemu:
-------------
qemu-img resize &lt;path-to-raw-file&gt; size

     http://azure.microsoft.com/en-us/documentation/articles/virtual-machines-linux-create-upload-vhd-generic/
 
#### How upload work

Azure requires VHD to be in Fixed Disk format. The command converts Dynamic and Differencing Disk to Fixed Disk during upload process, the conversion will not consume any additional space in local machine.

In case of Fixed Disk, the command detects blocks containing zeros and those will not be uploaded. In case of expandable disks (dynamic and differencing) only the blocks those are marked as non-empty in
the Block Allocation Table (BAT) will be uploaded.

The blocks containing data will be uploaded as chunks of 2 MB pages. Consecutive blocks will be merged to create 2 MB pages if the block size of disk is less than 2 MB. If the block size is greater than 2 MB, 
tool will split them as 2 MB pages.  

With page blob, we can upload multiple pages in parallel to decrease upload time. The command accepts the number of concurrent goroutines to use for upload through parallelism parameter. If the parallelism parameter is not proivded then it default to 8 * number_of_cpus.

### Inspect local VHD

A subset of command are exposed under inspect command for inspecting various segments of VHD in the local machine.

#### Show VHD footer

```bash
USAGE:
   azure-vhd-utils inspect footer [command options] [arguments...]

OPTIONS:
   --path   Path to VHD.
```

#### Show VHD header of an expandable disk

```bash
USAGE:
   azure-vhd-utils inspect header [command options] [arguments...]

OPTIONS:
   --path   Path to VHD.
```

Only expandable disks (dynamic and differencing) VHDs has header.

#### Show Block Allocation Table (BAT) of an expandable disk

```bash
USAGE:
   azure-vhd-utils inspect bat [command options] [arguments...]

OPTIONS:
   --path           Path to VHD.
   --start-range    Start range.
   --end-range      End range.
   --skip-empty     Do not show BAT entries pointing to empty blocks.
```

Only expandable disks (dynamic and differencing) VHDs has BAT.

#### Show block general information

```bash
USAGE:
   azure-vhd-utils inspect block info [command options] [arguments...]

OPTIONS:
   --path   Path to VHD.
```

This command shows the total number blocks, block size and size of block sector

### Show sector bitmap of an expandable disk's block

```bash
USAGE:
   azure-vhd-utils inspect block bitmap [command options] [arguments...]

OPTIONS:
   --path           Path to VHD.
   --block-index    Index of the block.
   
```

# License

This project is published under [MIT License](LICENSE).
